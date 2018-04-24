package account

import (
	"reflect"

	"github.com/lino-network/lino/tx/account/model"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-crypto"
)

var (
	// how many days the stake will increase to the maximum
	CoinDays         int64 = 8
	TotalCoinDaysSec int64 = CoinDays * 24 * 3600

	// maximum transaction cpacity cost
	CapacityUsagePerTransaction = types.NewCoin(1 * types.Decimals)

	// transaction cpacity recover period
	TransactionCapacityRecoverPeriod int64 = 24 * 3600 * 8
)

// linoaccount encapsulates all basic struct
type AccountManager struct {
	accountStorage model.AccountStorage `json:"account_manager"`
}

// NewLinoAccount return the account pointer
func NewAccountManager(key sdk.StoreKey) AccountManager {
	return AccountManager{
		accountStorage: model.NewAccountStorage(key),
	}
}

// check if account exist
func (accManager AccountManager) IsAccountExist(
	ctx sdk.Context, accKey types.AccountKey) bool {
	accountInfo, _ := accManager.accountStorage.GetInfo(ctx, accKey)
	return accountInfo != nil
}

// Implements types.AccountManager.
func (accManager AccountManager) CreateAccount(
	ctx sdk.Context, accKey types.AccountKey, pubkey crypto.PubKey,
	registerFee types.Coin) sdk.Error {
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
		Created:  ctx.BlockHeader().Time,
		PostKey:  pubkey,
		OwnerKey: pubkey,
		Address:  pubkey.Address(),
	}
	if err := accManager.accountStorage.SetInfo(ctx, accKey, accountInfo); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}

	bank.Username = accKey
	if err := accManager.accountStorage.SetBankFromAddress(
		ctx, pubkey.Address(), bank); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}

	accountMeta := &model.AccountMeta{
		LastActivity:        ctx.BlockHeader().Time,
		TransactionCapacity: types.NewCoin(0),
	}
	if err := accManager.accountStorage.SetMeta(ctx, accKey, accountMeta); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}
	reward :=
		&model.Reward{types.NewCoin(0), types.NewCoin(0), types.NewCoin(0), types.NewCoin(0)}
	if err := accManager.accountStorage.SetReward(ctx, accKey, reward); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}

	if err := accManager.accountStorage.SetGrantKeyList(ctx, accKey, &model.GrantKeyList{}); err != nil {
		return err
	}
	return nil
}

// use coin to present stake to prevent overflow
func (accManager AccountManager) GetStake(
	ctx sdk.Context, accKey types.AccountKey) (types.Coin, sdk.Error) {
	bank, err := accManager.accountStorage.GetBankFromAccountKey(ctx, accKey)
	if err != nil {
		return types.NewCoin(0), ErrGetStake(accKey).TraceCause(err, "")
	}
	pendingStakeQueue, err := accManager.accountStorage.GetPendingStakeQueue(ctx, bank.Address)
	if err != nil {
		return types.NewCoin(0), err
	}

	accManager.updateTXFromPendingStakeQueue(ctx, bank, pendingStakeQueue)

	stake := bank.Stake
	if err := accManager.accountStorage.SetPendingStakeQueue(
		ctx, bank.Address, pendingStakeQueue); err != nil {
		return types.NewCoin(0), err
	}

	if err := accManager.accountStorage.SetBankFromAddress(ctx, bank.Address, bank); err != nil {
		return types.NewCoin(0), err
	}
	return stake.Plus(types.RatToCoin(pendingStakeQueue.StakeCoinInQueue)), nil

}

func (accManager AccountManager) AddCoinToAddress(
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
		if err := accManager.accountStorage.SetPendingStakeQueue(
			ctx, address, &model.PendingStakeQueue{}); err != nil {
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
	if err := accManager.addPendingStakeToQueue(ctx, address, bank, pendingStake); err != nil {
		return ErrAddCoinToAddress(address).TraceCause(err, "")
	}
	if err := accManager.accountStorage.SetBankFromAddress(ctx, bank.Address, bank); err != nil {
		return ErrAddCoinToAddress(address).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) AddCoin(
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

func (accManager AccountManager) MinusCoin(
	ctx sdk.Context, accKey types.AccountKey, coin types.Coin) (err sdk.Error) {
	accountBank, err := accManager.accountStorage.GetBankFromAccountKey(ctx, accKey)
	if err != nil {
		return ErrMinusCoinToAccount(accKey).TraceCause(err, "")
	}
	if !accountBank.Balance.IsGTE(coin) {
		return ErrAccountCoinNotEnough()
	}
	pendingStakeQueue, err :=
		accManager.accountStorage.GetPendingStakeQueue(ctx, accountBank.Address)
	if err != nil {
		return err
	}
	accountBank.Balance = accountBank.Balance.Minus(coin)

	// update pending stake queue, remove expired transaction
	accManager.updateTXFromPendingStakeQueue(ctx, accountBank, pendingStakeQueue)

	for len(pendingStakeQueue.PendingStakeList) > 0 {
		lengthOfQueue := len(pendingStakeQueue.PendingStakeList)
		pendingStake := pendingStakeQueue.PendingStakeList[lengthOfQueue-1]
		if coin.IsGTE(pendingStake.Coin) {
			// if withdraw money more than last pending transaction, remove last transaction
			coin = coin.Minus(pendingStake.Coin)
			pendingStakeQueue.StakeCoinInQueue =
				pendingStakeQueue.StakeCoinInQueue.Sub(sdk.NewRat(
					pendingStakeQueue.LastUpdateTime-pendingStake.StartTime,
					TotalCoinDaysSec).Mul(pendingStake.Coin.ToRat()))
			pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Minus(pendingStake.Coin)
			pendingStakeQueue.PendingStakeList =
				pendingStakeQueue.PendingStakeList[:lengthOfQueue-1]
		} else {
			// otherwise try to cut last pending transaction
			pendingStakeQueue.StakeCoinInQueue =
				pendingStakeQueue.StakeCoinInQueue.Sub(sdk.NewRat(
					pendingStakeQueue.LastUpdateTime-pendingStake.StartTime,
					TotalCoinDaysSec).Mul(coin.ToRat()))
			pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Minus(coin)
			pendingStakeQueue.PendingStakeList[lengthOfQueue-1].Coin =
				pendingStakeQueue.PendingStakeList[lengthOfQueue-1].Coin.Minus(coin)
			coin = types.NewCoin(0)
			break
		}
	}
	if coin.IsPositive() {
		accountBank.Stake = accountBank.Balance
	}
	if err := accManager.accountStorage.SetPendingStakeQueue(
		ctx, accountBank.Address, pendingStakeQueue); err != nil {
		return err
	}

	if err := accManager.accountStorage.SetBankFromAddress(
		ctx, accountBank.Address, accountBank); err != nil {
		return ErrMinusCoinToAccount(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) GetBankAddress(
	ctx sdk.Context, accKey types.AccountKey) (sdk.Address, sdk.Error) {
	accountInfo, err := accManager.accountStorage.GetInfo(ctx, accKey)
	if err != nil {
		return nil, ErrGetBankAddress(accKey).TraceCause(err, "")
	}
	return accountInfo.Address, nil
}

func (accManager AccountManager) GetOwnerKey(
	ctx sdk.Context, accKey types.AccountKey) (*crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.accountStorage.GetInfo(ctx, accKey)
	if err != nil {
		return nil, ErrGetOwnerKey(accKey).TraceCause(err, "")
	}
	return &accountInfo.OwnerKey, nil
}

func (accManager AccountManager) GetPostKey(
	ctx sdk.Context, accKey types.AccountKey) (*crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.accountStorage.GetInfo(ctx, accKey)
	if err != nil {
		return nil, ErrGetPostKey(accKey).TraceCause(err, "")
	}
	return &accountInfo.PostKey, nil
}

func (accManager AccountManager) GetBankBalance(
	ctx sdk.Context, accKey types.AccountKey) (types.Coin, sdk.Error) {
	accountBank, err := accManager.accountStorage.GetBankFromAccountKey(ctx, accKey)
	if err != nil {
		return types.Coin{}, ErrGetBankBalance(accKey).TraceCause(err, "")
	}
	return accountBank.Balance, nil
}

func (accManager AccountManager) GetSequence(
	ctx sdk.Context, accKey types.AccountKey) (int64, sdk.Error) {
	accountMeta, err := accManager.accountStorage.GetMeta(ctx, accKey)
	if err != nil {
		return 0, ErrGetSequence(accKey).TraceCause(err, "")
	}
	return accountMeta.Sequence, nil
}

func (accManager AccountManager) IncreaseSequenceByOne(
	ctx sdk.Context, accKey types.AccountKey) sdk.Error {
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

func (accManager AccountManager) AddIncomeAndReward(
	ctx sdk.Context, accKey types.AccountKey,
	originIncome, friction, actualReward types.Coin) sdk.Error {
	reward, err := accManager.accountStorage.GetReward(ctx, accKey)
	if err != nil {
		return ErrAddIncomeAndReward(accKey).TraceCause(err, "")
	}
	reward.OriginalIncome = reward.OriginalIncome.Plus(originIncome)
	reward.FrictionIncome = reward.FrictionIncome.Plus(friction)
	reward.ActualReward = reward.ActualReward.Plus(actualReward)
	reward.UnclaimReward = reward.UnclaimReward.Plus(actualReward)
	if err := accManager.accountStorage.SetReward(ctx, accKey, reward); err != nil {
		return ErrAddIncomeAndReward(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) ClaimReward(ctx sdk.Context, accKey types.AccountKey) sdk.Error {
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

func (accManager AccountManager) IsMyFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) bool {
	return accManager.accountStorage.IsMyFollower(ctx, me, follower)
}

func (accManager AccountManager) IsMyFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) bool {
	return accManager.accountStorage.IsMyFollowing(ctx, me, following)
}

func (accManager AccountManager) SetFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) sdk.Error {
	if accManager.accountStorage.IsMyFollower(ctx, me, follower) {
		return nil
	}
	meta := model.FollowerMeta{
		CreatedAt:    ctx.BlockHeader().Time,
		FollowerName: follower,
	}
	accManager.accountStorage.SetFollowerMeta(ctx, me, meta)
	return nil
}

func (accManager AccountManager) SetFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) sdk.Error {
	if accManager.accountStorage.IsMyFollowing(ctx, me, following) {
		return nil
	}
	meta := model.FollowingMeta{
		CreatedAt:     ctx.BlockHeader().Time,
		FollowingName: following,
	}
	accManager.accountStorage.SetFollowingMeta(ctx, me, meta)
	return nil
}

func (accManager AccountManager) RemoveFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) sdk.Error {
	if !accManager.accountStorage.IsMyFollower(ctx, me, follower) {
		return nil
	}
	accManager.accountStorage.RemoveFollowerMeta(ctx, me, follower)
	return nil
}

func (accManager AccountManager) RemoveFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) sdk.Error {
	if !accManager.accountStorage.IsMyFollowing(ctx, me, following) {
		return nil
	}
	accManager.accountStorage.RemoveFollowingMeta(ctx, me, following)
	return nil
}

func (accManager AccountManager) CheckUserTPSCapacity(
	ctx sdk.Context, me types.AccountKey, tpsCapacityRatio sdk.Rat) sdk.Error {
	accountMeta, err := accManager.accountStorage.GetMeta(ctx, me)
	if err != nil {
		return ErrCheckUserTPSCapacity(me).TraceCause(err, "")
	}
	stake, err := accManager.GetStake(ctx, me)
	if err != nil {
		return ErrCheckUserTPSCapacity(me).TraceCause(err, "")
	}
	if accountMeta.TransactionCapacity.IsGTE(stake) {
		accountMeta.TransactionCapacity = stake
	} else {
		incrementRatio := sdk.NewRat(
			ctx.BlockHeader().Time-accountMeta.LastActivity,
			TransactionCapacityRecoverPeriod)
		if incrementRatio.GT(sdk.OneRat) {
			incrementRatio = sdk.OneRat
		}
		accountMeta.TransactionCapacity =
			accountMeta.TransactionCapacity.Plus(types.RatToCoin(
				stake.Minus(accountMeta.TransactionCapacity).ToRat().Mul(incrementRatio)))
	}
	currentTxCost := types.RatToCoin(CapacityUsagePerTransaction.ToRat().Mul(tpsCapacityRatio))
	if currentTxCost.IsGT(accountMeta.TransactionCapacity) {
		return ErrAccountTPSCapacityNotEnough(me)
	}
	accountMeta.TransactionCapacity = accountMeta.TransactionCapacity.Minus(currentTxCost)
	accountMeta.LastActivity = ctx.BlockHeader().Time
	if err := accManager.accountStorage.SetMeta(ctx, me, accountMeta); err != nil {
		return ErrIncreaseSequenceByOne(me).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) UpdateDonationRelationship(
	ctx sdk.Context, me, other types.AccountKey) sdk.Error {
	relationship, err := accManager.accountStorage.GetRelationship(ctx, me, other)
	if err != nil {
		return err
	}
	if relationship == nil {
		relationship = &model.Relationship{0}
	}
	relationship.DonationTimes += 1
	if err := accManager.accountStorage.SetRelationship(ctx, me, other, relationship); err != nil {
		return err
	}
	return nil
}

func (accManager AccountManager) GrantPubKeyToUser(
	ctx sdk.Context, me types.AccountKey, grantUser types.AccountKey,
	validityPeriod int64, grantLevel int64) sdk.Error {
	pubKey, err := accManager.GetPostKey(ctx, grantUser)
	if err != nil {
		return err
	}

	grantKeyList, err := accManager.accountStorage.GetGrantKeyList(ctx, me)
	if err != nil {
		return err
	}

	idx := 0
	for idx < len(grantKeyList.GrantPubKeyList) {
		if grantKeyList.GrantPubKeyList[idx].Expire < ctx.BlockHeader().Time ||
			grantKeyList.GrantPubKeyList[idx].Username == grantUser {
			grantKeyList.GrantPubKeyList = append(
				grantKeyList.GrantPubKeyList[:idx], grantKeyList.GrantPubKeyList[idx+1:]...)
			continue
		}
	}
	newGrantPubKey := model.GrantPubKey{
		Username: grantUser,
		PubKey:   *pubKey,
		Expire:   ctx.BlockHeader().Time + validityPeriod,
	}
	grantKeyList.GrantPubKeyList = append(grantKeyList.GrantPubKeyList, newGrantPubKey)
	return accManager.accountStorage.SetGrantKeyList(ctx, me, grantKeyList)
}

func (accManager AccountManager) CheckAuthenticatePubKeyOwner(
	ctx sdk.Context, me types.AccountKey, signKey crypto.PubKey) (types.AccountKey, sdk.Error) {
	pubKey, err := accManager.GetOwnerKey(ctx, me)
	if err != nil {
		return "", err
	}
	if reflect.DeepEqual(*pubKey, signKey) {
		return me, nil
	}
	pubKey, err = accManager.GetPostKey(ctx, me)
	if err != nil {
		return "", err
	}
	if reflect.DeepEqual(*pubKey, signKey) {
		return me, nil
	}

	grantKeyList, err := accManager.accountStorage.GetGrantKeyList(ctx, me)
	if err != nil {
		return "", err
	}
	idx := 0
	for idx < len(grantKeyList.GrantPubKeyList) {
		if grantKeyList.GrantPubKeyList[idx].Expire < ctx.BlockHeader().Time {
			grantKeyList.GrantPubKeyList = append(
				grantKeyList.GrantPubKeyList[:idx], grantKeyList.GrantPubKeyList[idx+1:]...)
			continue
		}

		if reflect.DeepEqual(grantKeyList.GrantPubKeyList[idx].PubKey, signKey) {
			return grantKeyList.GrantPubKeyList[idx].Username, nil
		}
		idx += 1
	}
	if err := accManager.accountStorage.SetGrantKeyList(ctx, me, grantKeyList); err != nil {
		return "", err
	}
	return "", ErrCheckAuthenticatePubKeyOwner(me)
}

func (accManager AccountManager) GetDonationRelationship(
	ctx sdk.Context, me, other types.AccountKey) (int64, sdk.Error) {
	relationship, err := accManager.accountStorage.GetRelationship(ctx, me, other)
	if err != nil {
		return 0, err
	}
	if relationship == nil {
		return 0, nil
	}
	return relationship.DonationTimes, nil
}

func (accManager AccountManager) addPendingStakeToQueue(
	ctx sdk.Context, address sdk.Address, bank *model.AccountBank,
	pendingStake model.PendingStake) sdk.Error {
	pendingStakeQueue, err := accManager.accountStorage.GetPendingStakeQueue(ctx, address)
	if err != nil {
		return err
	}
	accManager.updateTXFromPendingStakeQueue(ctx, bank, pendingStakeQueue)
	pendingStakeQueue.PendingStakeList = append(pendingStakeQueue.PendingStakeList, pendingStake)
	pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Plus(pendingStake.Coin)
	return accManager.accountStorage.SetPendingStakeQueue(ctx, address, pendingStakeQueue)
}

func (accManager AccountManager) updateTXFromPendingStakeQueue(
	ctx sdk.Context, bank *model.AccountBank, pendingStakeQueue *model.PendingStakeQueue) {
	// remove expired transaction
	for len(pendingStakeQueue.PendingStakeList) > 0 {
		pendingStake := pendingStakeQueue.PendingStakeList[0]
		if pendingStake.EndTime < ctx.BlockHeader().Time {
			// remove the transaction from queue, clean stake coin in queue and minus total coin
			pendingStakeQueue.StakeCoinInQueue =
				pendingStakeQueue.StakeCoinInQueue.Sub(sdk.NewRat(
					pendingStakeQueue.LastUpdateTime-pendingStake.StartTime,
					TotalCoinDaysSec).Mul(pendingStake.Coin.ToRat()))

			// update bank stake
			bank.Stake = bank.Stake.Plus(pendingStake.Coin)
			pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Minus(pendingStake.Coin)
			pendingStakeQueue.PendingStakeList = pendingStakeQueue.PendingStakeList[1:]
		} else {
			break
		}
	}
	if len(pendingStakeQueue.PendingStakeList) == 0 {
		pendingStakeQueue.TotalCoin = types.NewCoin(0)
		pendingStakeQueue.StakeCoinInQueue = sdk.ZeroRat
	} else {
		// update all pending stake at the same time
		pendingStakeQueue.StakeCoinInQueue =
			pendingStakeQueue.StakeCoinInQueue.Add(sdk.NewRat(
				ctx.BlockHeader().Time-pendingStakeQueue.LastUpdateTime,
				TotalCoinDaysSec).Mul(pendingStakeQueue.TotalCoin.ToRat()))
	}

	pendingStakeQueue.LastUpdateTime = ctx.BlockHeader().Time
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
