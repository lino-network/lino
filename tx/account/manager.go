package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/account/model"
	"github.com/lino-network/lino/types"
	"github.com/tendermint/go-crypto"
)

// how many days the stake will increase to the maximum
var CoinDays int64 = 8

// the minimal requirement to open an account
var OpenBankFee types.Coin = types.NewCoin(5 * types.Decimals)

// linoaccount encapsulates all basic struct
type AccountManager struct {
	accountStorage *model.AccountStorage `json:"account_manager"`
}

// NewLinoAccount return the account pointer
func NewAccountManager(key sdk.StoreKey) *AccountManager {
	return &AccountManager{
		accountStorage: model.NewAccountStorage(key),
	}
}

// check if account exist
func (accManager *AccountManager) IsAccountExist(ctx sdk.Context, accKey types.AccountKey) bool {
	accountInfo, _ := accManager.accountStorage.GetInfo(ctx, accKey)
	return accountInfo != nil
}

// Implements types.AccountManager.
func (accManager *AccountManager) CreateAccount(
	ctx sdk.Context, accKey types.AccountKey, pubkey crypto.PubKey, registerFee types.Coin) sdk.Error {
	if accManager.IsAccountExist(ctx, accKey) {
		return ErrAccountAlreadyExists(accKey)
	}
	bank, err := accManager.accountStorage.GetBankFromAddress(ctx, pubkey.Address())
	if err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}
	if bank.Username != "" {
		return ErrBankAlreadyRegistered()
	}

	if !bank.Balance.IsGTE(registerFee) {
		return ErrRegisterFeeInsufficient()
	}

	accountInfo := &model.AccountInfo{
		Username: accKey,
		Created:  ctx.BlockHeight(),
		PostKey:  pubkey,
		OwnerKey: pubkey,
		Address:  pubkey.Address(),
	}
	if err := accManager.accountStorage.SetInfo(ctx, accKey, accountInfo); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}

	bank.Username = accKey
	if err := accManager.accountStorage.SetBankFromAddress(ctx, pubkey.Address(), bank); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}

	accountMeta := &model.AccountMeta{
		LastActivity:   ctx.BlockHeight(),
		ActivityBurden: types.DefaultActivityBurden,
	}
	if err := accManager.accountStorage.SetMeta(ctx, accKey, accountMeta); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}
	reward := &model.Reward{types.NewCoin(0), types.NewCoin(0), types.NewCoin(0)}
	if err := accManager.accountStorage.SetReward(ctx, accKey, reward); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}
	return nil
}

// use coin to present stake to prevent overflow
func (accManager *AccountManager) GetStake(ctx sdk.Context, accKey types.AccountKey) (types.Coin, sdk.Error) {
	bank, err := accManager.accountStorage.GetBankFromAccountKey(ctx, accKey)
	if err != nil {
		return types.NewCoin(0), ErrGetStake(accKey).TraceCause(err, "")
	}
	pendingStakeQueue, err := accManager.accountStorage.GetPendingStakeQueue(ctx, bank.Address)
	if err != nil {
		return types.NewCoin(0), err
	}

	accManager.removeExpiredTXFromPendingStakeQueue(ctx, bank, pendingStakeQueue)

	stake := bank.Stake
	for _, pendingStake := range pendingStakeQueue.PendingStakeList {
		stake = stake.Plus(
			types.RatToCoin(
				sdk.NewRat(ctx.BlockHeader().Time-pendingStake.StartTime,
					pendingStake.EndTime-pendingStake.StartTime).
					Mul(pendingStake.Coin.ToRat())))
	}
	if err := accManager.accountStorage.SetPendingStakeQueue(ctx, bank.Address, pendingStakeQueue); err != nil {
		return types.NewCoin(0), err
	}

	if err := accManager.accountStorage.SetBankFromAddress(ctx, bank.Address, bank); err != nil {
		return types.NewCoin(0), err
	}
	return stake, nil

}

func (accManager *AccountManager) AddCoinToAddress(
	ctx sdk.Context, address sdk.Address, coin types.Coin) (err sdk.Error) {
	if coin.IsZero() {
		return nil
	}
	bank, _ := accManager.accountStorage.GetBankFromAddress(ctx, address)
	if bank == nil {
		bank = &model.AccountBank{
			Address: address,
			Balance: coin,
		}
		if err := accManager.accountStorage.SetPendingStakeQueue(ctx, address, &model.PendingStakeQueue{}); err != nil {
			return err
		}
	} else {
		bank.Balance = bank.Balance.Plus(coin)
	}
	pendingStake := model.PendingStake{
		StartTime: ctx.BlockHeader().Time,
		EndTime:   ctx.BlockHeader().Time + CoinDays*24*3600,
		Coin:      coin,
	}
	if err := accManager.addPendingStakeToQueue(ctx, address, pendingStake); err != nil {
		return ErrAddCoinToAddress(address).TraceCause(err, "")
	}
	if err := accManager.accountStorage.SetBankFromAddress(ctx, bank.Address, bank); err != nil {
		return ErrAddCoinToAddress(address).TraceCause(err, "")
	}
	return nil
}

func (accManager *AccountManager) AddCoin(
	ctx sdk.Context, accKey types.AccountKey, coin types.Coin) (err sdk.Error) {
	address, err := accManager.GetBankAddress(ctx, accKey)
	if err != nil {
		return ErrAddCoinToAccount(accKey).TraceCause(err, "")
	}
	if err := accManager.AddCoinToAddress(ctx, address, coin); err != nil {
		return ErrAddCoinToAccount(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager *AccountManager) MinusCoin(
	ctx sdk.Context, accKey types.AccountKey, coin types.Coin) (err sdk.Error) {
	accountBank, err := accManager.accountStorage.GetBankFromAccountKey(ctx, accKey)
	if err != nil {
		return ErrMinusCoinToAccount(accKey).TraceCause(err, "")
	}
	if !accountBank.Balance.IsGTE(coin) {
		return ErrAccountCoinNotEnough()
	}
	pendingStakeQueue, err := accManager.accountStorage.GetPendingStakeQueue(ctx, accountBank.Address)
	if err != nil {
		return err
	}
	accountBank.Balance = accountBank.Balance.Minus(coin)

	// update pending stake queue, remove expired transaction
	accManager.removeExpiredTXFromPendingStakeQueue(ctx, accountBank, pendingStakeQueue)

	for len(pendingStakeQueue.PendingStakeList) > 0 {
		pendingStake := pendingStakeQueue.PendingStakeList[0]
		unstakeCoin := types.RatToCoin(
			sdk.NewRat(pendingStake.EndTime-ctx.BlockHeader().Time,
				pendingStake.EndTime-pendingStake.StartTime).
				Mul(pendingStake.Coin.ToRat()))
		if coin.IsGTE(unstakeCoin) {
			// if withdraw coin is larger than unstakeCoin, withdraw all unstake coin and pop from queue
			accountBank.Stake = accountBank.Stake.Plus(pendingStakeQueue.PendingStakeList[0].Coin.Minus(unstakeCoin))
			coin = coin.Minus(unstakeCoin)
			pendingStakeQueue.PendingStakeList = pendingStakeQueue.PendingStakeList[1:]
		} else {
			// new end time = current time + (unstake - withdraw)/unstake*(end time - current time)
			pendingStakeQueue.PendingStakeList[0].EndTime = ctx.BlockHeader().Time +
				unstakeCoin.Minus(coin).ToRat().Quo(unstakeCoin.ToRat()).
					Mul(sdk.NewRat(pendingStakeQueue.PendingStakeList[0].EndTime-ctx.BlockHeader().Time)).Evaluate()
			pendingStakeQueue.PendingStakeList[0].Coin = pendingStakeQueue.PendingStakeList[0].Coin.Minus(coin)
			coin = types.NewCoin(0)
			break
		}
	}
	if coin.IsPositive() {
		accountBank.Stake = accountBank.Balance
	}
	if err := accManager.accountStorage.SetPendingStakeQueue(ctx, accountBank.Address, pendingStakeQueue); err != nil {
		return err
	}

	if err := accManager.accountStorage.SetBankFromAddress(ctx, accountBank.Address, accountBank); err != nil {
		return ErrMinusCoinToAccount(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager *AccountManager) GetBankAddress(ctx sdk.Context, accKey types.AccountKey) (sdk.Address, sdk.Error) {
	accountInfo, err := accManager.accountStorage.GetInfo(ctx, accKey)
	if err != nil {
		return nil, ErrGetBankAddress(accKey).TraceCause(err, "")
	}
	return accountInfo.Address, nil
}

func (accManager *AccountManager) GetOwnerKey(ctx sdk.Context, accKey types.AccountKey) (*crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.accountStorage.GetInfo(ctx, accKey)
	if err != nil {
		return nil, ErrGetOwnerKey(accKey).TraceCause(err, "")
	}
	return &accountInfo.OwnerKey, nil
}

func (accManager *AccountManager) GetPostKey(ctx sdk.Context, accKey types.AccountKey) (*crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.accountStorage.GetInfo(ctx, accKey)
	if err != nil {
		return nil, ErrGetPostKey(accKey).TraceCause(err, "")
	}
	return &accountInfo.PostKey, nil
}

func (accManager *AccountManager) GetBankBalance(ctx sdk.Context, accKey types.AccountKey) (types.Coin, sdk.Error) {
	accountBank, err := accManager.accountStorage.GetBankFromAccountKey(ctx, accKey)
	if err != nil {
		return types.Coin{}, ErrGetBankBalance(accKey).TraceCause(err, "")
	}
	return accountBank.Balance, nil
}

func (accManager *AccountManager) GetSequence(ctx sdk.Context, accKey types.AccountKey) (int64, sdk.Error) {
	accountMeta, err := accManager.accountStorage.GetMeta(ctx, accKey)
	if err != nil {
		return 0, ErrGetSequence(accKey).TraceCause(err, "")
	}
	return accountMeta.Sequence, nil
}

func (accManager *AccountManager) IncreaseSequenceByOne(ctx sdk.Context, accKey types.AccountKey) sdk.Error {
	accountMeta, err := accManager.accountStorage.GetMeta(ctx, accKey)
	if err != nil {
		return ErrGetSequence(accKey).TraceCause(err, "")
	}
	accountMeta.Sequence += 1
	if err := accManager.accountStorage.SetMeta(ctx, accKey, accountMeta); err != nil {
		return ErrIncreaseSequenceByOne(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager *AccountManager) AddIncomeAndReward(
	ctx sdk.Context, accKey types.AccountKey, originIncome, actualReward types.Coin) sdk.Error {
	reward, err := accManager.accountStorage.GetReward(ctx, accKey)
	if err != nil {
		return ErrAddIncomeAndReward(accKey).TraceCause(err, "")
	}
	// 1% of total income
	reward.OriginalIncome = reward.OriginalIncome.Plus(originIncome)
	reward.ActualReward = reward.ActualReward.Plus(actualReward)
	reward.UnclaimReward = reward.UnclaimReward.Plus(actualReward)
	if err := accManager.accountStorage.SetReward(ctx, accKey, reward); err != nil {
		return ErrAddIncomeAndReward(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager *AccountManager) ClaimReward(ctx sdk.Context, accKey types.AccountKey) sdk.Error {
	reward, err := accManager.accountStorage.GetReward(ctx, accKey)
	if err != nil {
		return ErrClaimReward(accKey).TraceCause(err, "")
	}
	if err := accManager.AddCoin(ctx, accKey, reward.UnclaimReward); err != nil {
		return ErrClaimReward(accKey).TraceCause(err, "")
	}
	reward.UnclaimReward = types.NewCoin(0)
	if err := accManager.accountStorage.SetReward(ctx, accKey, reward); err != nil {
		return ErrClaimReward(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager *AccountManager) UpdateLastActivity(
	ctx sdk.Context, accKey types.AccountKey) sdk.Error {
	accountMeta, err := accManager.accountStorage.GetMeta(ctx, accKey)
	if err != nil {
		return ErrUpdateLastActivity(accKey).TraceCause(err, "")
	}
	accountMeta.LastActivity = ctx.BlockHeight()
	if err := accManager.accountStorage.SetMeta(ctx, accKey, accountMeta); err != nil {
		return ErrUpdateLastActivity(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager *AccountManager) IsMyFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) bool {
	return accManager.accountStorage.IsMyFollower(ctx, me, follower)
}

func (accManager *AccountManager) IsMyFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) bool {
	return accManager.accountStorage.IsMyFollowing(ctx, me, following)
}

func (accManager *AccountManager) SetFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) sdk.Error {
	if accManager.accountStorage.IsMyFollower(ctx, me, follower) {
		return nil
	}
	meta := model.FollowerMeta{
		CreatedAt:    ctx.BlockHeight(),
		FollowerName: follower,
	}
	accManager.accountStorage.SetFollowerMeta(ctx, me, meta)
	return nil
}

func (accManager *AccountManager) SetFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) sdk.Error {
	if accManager.accountStorage.IsMyFollowing(ctx, me, following) {
		return nil
	}
	meta := model.FollowingMeta{
		CreatedAt:     ctx.BlockHeight(),
		FollowingName: following,
	}
	accManager.accountStorage.SetFollowingMeta(ctx, me, meta)
	return nil
}

func (accManager *AccountManager) RemoveFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) sdk.Error {
	if !accManager.accountStorage.IsMyFollower(ctx, me, follower) {
		return nil
	}
	accManager.accountStorage.RemoveFollowerMeta(ctx, me, follower)
	return nil
}

func (accManager *AccountManager) RemoveFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) sdk.Error {
	if !accManager.accountStorage.IsMyFollowing(ctx, me, following) {
		return nil
	}
	accManager.accountStorage.RemoveFollowingMeta(ctx, me, following)
	return nil
}

func (accManager *AccountManager) addPendingStakeToQueue(
	ctx sdk.Context, address sdk.Address, pendingStake model.PendingStake) sdk.Error {
	pendingStakeQueue, err := accManager.accountStorage.GetPendingStakeQueue(ctx, address)
	if err != nil {
		return err
	}
	pendingStakeQueue.PendingStakeList = append(pendingStakeQueue.PendingStakeList, pendingStake)
	return accManager.accountStorage.SetPendingStakeQueue(ctx, address, pendingStakeQueue)
}

func (accManager *AccountManager) removeExpiredTXFromPendingStakeQueue(
	ctx sdk.Context, bank *model.AccountBank, pendingStakeQueue *model.PendingStakeQueue) {
	for len(pendingStakeQueue.PendingStakeList) > 0 {
		if pendingStakeQueue.PendingStakeList[0].EndTime < ctx.BlockHeader().Time {
			bank.Stake = bank.Stake.Plus(pendingStakeQueue.PendingStakeList[0].Coin)
			pendingStakeQueue.PendingStakeList = pendingStakeQueue.PendingStakeList[1:]
		} else {
			break
		}
	}
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}
