package account

import (
	"reflect"
	"time"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"

	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountManager - account manager
type AccountManager struct {
	storage     model.AccountStorage
	paramHolder param.ParamHolder
}

// NewLinoAccount - new account manager
func NewAccountManager(key sdk.StoreKey, holder param.ParamHolder) AccountManager {
	return AccountManager{
		storage:     model.NewAccountStorage(key),
		paramHolder: holder,
	}
}

// DoesAccountExist - check if account exists in KVStore or not
func (accManager AccountManager) DoesAccountExist(ctx sdk.Context, username types.AccountKey) bool {
	return accManager.storage.DoesAccountExist(ctx, username)
}

// CreateAccount - create account, caller should make sure the register fee is valid
func (accManager AccountManager) CreateAccount(
	ctx sdk.Context, referrer types.AccountKey, username types.AccountKey,
	resetKey, transactionKey, appKey crypto.PubKey, registerDeposit types.Coin) sdk.Error {
	if accManager.DoesAccountExist(ctx, username) {
		return ErrAccountAlreadyExists(username)
	}
	accParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}

	depositWithFullStake := registerDeposit
	if registerDeposit.IsGT(accParams.FirstDepositFullStakeLimit) {
		depositWithFullStake = accParams.FirstDepositFullStakeLimit
	}
	if err := accManager.storage.SetPendingStakeQueue(
		ctx, username, &model.PendingStakeQueue{}); err != nil {
		return err
	}

	if err := accManager.storage.SetBankFromAccountKey(ctx, username, &model.AccountBank{}); err != nil {
		return err
	}

	accountInfo := &model.AccountInfo{
		Username:       username,
		CreatedAt:      ctx.BlockHeader().Time.Unix(),
		ResetKey:       resetKey,
		TransactionKey: transactionKey,
		AppKey:         appKey,
	}
	if err := accManager.storage.SetInfo(ctx, username, accountInfo); err != nil {
		return err
	}

	accountMeta := &model.AccountMeta{
		LastActivityAt:       ctx.BlockHeader().Time.Unix(),
		LastReportOrUpvoteAt: ctx.BlockHeader().Time.Unix(),
		TransactionCapacity:  depositWithFullStake,
	}
	if err := accManager.storage.SetMeta(ctx, username, accountMeta); err != nil {
		return err
	}
	if err := accManager.storage.SetReward(ctx, username, &model.Reward{}); err != nil {
		return err
	}
	// when open account, blockchain will give a certain amount lino with full stake.
	if err := accManager.AddSavingCoinWithFullStake(
		ctx, username, depositWithFullStake, referrer,
		types.InitAccountWithFullStakeMemo, types.TransferIn); err != nil {
		return ErrAddSavingCoinWithFullStake()
	}
	if err := accManager.AddSavingCoin(
		ctx, username, registerDeposit.Minus(depositWithFullStake), referrer,
		types.InitAccountRegisterDepositMemo, types.TransferIn); err != nil {
		return ErrAddSavingCoin()
	}
	return nil
}

// GetStake - recalculate and get user current stake
func (accManager AccountManager) GetStake(
	ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	bank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	pendingStakeQueue, err := accManager.storage.GetPendingStakeQueue(ctx, username)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}

	accManager.updateTXFromPendingStakeQueue(ctx, bank, pendingStakeQueue)

	if err := accManager.storage.SetPendingStakeQueue(
		ctx, username, pendingStakeQueue); err != nil {
		return types.NewCoinFromInt64(0), err
	}

	if err := accManager.storage.SetBankFromAccountKey(ctx, username, bank); err != nil {
		return types.NewCoinFromInt64(0), err
	}

	stake := bank.Stake
	stakeInQueue := types.RatToCoin(pendingStakeQueue.StakeCoinInQueue)
	totalStake := stake.Plus(stakeInQueue)
	return totalStake, nil
}

// AddSavingCoin - add coin to balance and pending stake
func (accManager AccountManager) AddSavingCoin(
	ctx sdk.Context, username types.AccountKey, coin types.Coin, from types.AccountKey, memo string,
	detailType types.TransferDetailType) (err sdk.Error) {
	if !accManager.DoesAccountExist(ctx, username) {
		return ErrAccountNotFound(username)
	}
	if coin.IsZero() {
		return nil
	}
	bank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return err
	}

	bank.Saving = bank.Saving.Plus(coin)
	if err := accManager.AddBalanceHistory(ctx, username, bank.NumOfTx,
		model.Detail{
			Amount:     coin,
			DetailType: detailType,
			To:         username,
			From:       from,
			Balance:    bank.Saving,
			CreatedAt:  ctx.BlockHeader().Time.Unix(),
			Memo:       memo,
		}); err != nil {
		return err
	}

	bank.NumOfTx++
	coinDayParams, err := accManager.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		return err
	}

	startTime := ctx.BlockHeader().Time.Unix() / types.CoinDayRecordIntervalSec * types.CoinDayRecordIntervalSec
	pendingStake := model.PendingStake{
		StartTime: startTime,
		EndTime:   startTime + coinDayParams.SecondsToRecoverCoinDayStake,
		Coin:      coin,
	}
	if err := accManager.addPendingStakeToQueue(ctx, username, bank, pendingStake); err != nil {
		return err
	}

	if err := accManager.storage.SetBankFromAccountKey(ctx, username, bank); err != nil {
		return err
	}
	return nil
}

// AddSavingCoinWithFullStake - add coin to balance with full stake
func (accManager AccountManager) AddSavingCoinWithFullStake(
	ctx sdk.Context, username types.AccountKey, coin types.Coin, from types.AccountKey, memo string,
	detailType types.TransferDetailType) (err sdk.Error) {
	if !accManager.DoesAccountExist(ctx, username) {
		return ErrAccountNotFound(username)
	}
	if coin.IsZero() {
		return nil
	}
	bank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return err
	}

	bank.Saving = bank.Saving.Plus(coin)
	if err := accManager.AddBalanceHistory(ctx, username, bank.NumOfTx,
		model.Detail{
			Amount:     coin,
			DetailType: detailType,
			To:         username,
			From:       from,
			Balance:    bank.Saving,
			CreatedAt:  ctx.BlockHeader().Time.Unix(),
			Memo:       memo,
		}); err != nil {
		return err
	}
	bank.Stake = bank.Stake.Plus(coin)
	bank.NumOfTx++

	if err := accManager.storage.SetBankFromAccountKey(ctx, username, bank); err != nil {
		return err
	}
	return nil
}

// MinusSavingCoin - minus coin from balance, remove stake in the tail
func (accManager AccountManager) MinusSavingCoin(
	ctx sdk.Context, username types.AccountKey, coin types.Coin, to types.AccountKey,
	memo string, detailType types.TransferDetailType) (err sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return err
	}

	accountParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}
	remain := accountBank.Saving.Minus(coin)
	if !remain.IsGTE(accountParams.MinimumBalance) {
		return ErrAccountSavingCoinNotEnough()
	}

	if coin.IsZero() {
		return nil
	}
	accountBank.Saving = accountBank.Saving.Minus(coin)

	if err := accManager.AddBalanceHistory(
		ctx, username, accountBank.NumOfTx, model.Detail{
			Amount:     coin,
			DetailType: detailType,
			To:         to,
			From:       username,
			Balance:    accountBank.Saving,
			CreatedAt:  ctx.BlockHeader().Time.Unix(),
			Memo:       memo,
		}); err != nil {
		return err
	}
	accountBank.NumOfTx++

	pendingStakeQueue, err :=
		accManager.storage.GetPendingStakeQueue(ctx, username)
	if err != nil {
		return err
	}
	// update pending stake queue, remove expired transaction
	accManager.updateTXFromPendingStakeQueue(ctx, accountBank, pendingStakeQueue)

	coinDayParams, err := accManager.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		return err
	}

	for len(pendingStakeQueue.PendingStakeList) > 0 {
		lengthOfQueue := len(pendingStakeQueue.PendingStakeList)
		pendingStake := pendingStakeQueue.PendingStakeList[lengthOfQueue-1]
		recoverRatio := sdk.NewRat(
			pendingStakeQueue.LastUpdatedAt-pendingStake.StartTime,
			coinDayParams.SecondsToRecoverCoinDayStake)
		if coin.IsGTE(pendingStake.Coin) {
			// if withdraw money is much than last pending transaction, remove last transaction
			coin = coin.Minus(pendingStake.Coin)

			pendingStakeCoinWithoutLastTx :=
				pendingStakeQueue.StakeCoinInQueue.Sub((recoverRatio.Mul(pendingStake.Coin.ToRat())))
			pendingStakeQueue.StakeCoinInQueue = pendingStakeCoinWithoutLastTx

			pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Minus(pendingStake.Coin)
			pendingStakeQueue.PendingStakeList = pendingStakeQueue.PendingStakeList[:lengthOfQueue-1]
		} else {
			// otherwise try to cut last pending transaction
			pendingStakeCoinWithoutSpentCoin :=
				pendingStakeQueue.StakeCoinInQueue.Sub(
					recoverRatio.Mul(coin.ToRat()))
			pendingStakeQueue.StakeCoinInQueue = pendingStakeCoinWithoutSpentCoin

			pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Minus(coin)
			pendingStakeQueue.PendingStakeList[lengthOfQueue-1].Coin =
				pendingStakeQueue.PendingStakeList[lengthOfQueue-1].Coin.Minus(coin)
			coin = types.NewCoinFromInt64(0)
			break
		}
	}
	if coin.IsPositive() {
		accountBank.Stake = accountBank.Saving
	}
	if err := accManager.storage.SetPendingStakeQueue(
		ctx, username, pendingStakeQueue); err != nil {
		return err
	}

	if err := accManager.storage.SetBankFromAccountKey(
		ctx, username, accountBank); err != nil {
		return err
	}
	return nil
}

// MinusSavingCoin - minus coin from balance, remove most charged stake coin
func (accManager AccountManager) MinusSavingCoinWithFullStake(
	ctx sdk.Context, username types.AccountKey, coin types.Coin, to types.AccountKey,
	memo string, detailType types.TransferDetailType) (err sdk.Error) {
	if coin.IsZero() {
		return nil
	}
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return err
	}

	accountParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}
	remain := accountBank.Saving.Minus(coin)
	if !remain.IsGTE(accountParams.MinimumBalance) {
		return ErrAccountSavingCoinNotEnough()
	}
	accountBank.Saving = remain

	if err := accManager.AddBalanceHistory(
		ctx, username, accountBank.NumOfTx, model.Detail{
			Amount:     coin,
			DetailType: detailType,
			To:         to,
			From:       username,
			Balance:    accountBank.Saving,
			CreatedAt:  ctx.BlockHeader().Time.Unix(),
			Memo:       memo,
		}); err != nil {
		return err
	}
	accountBank.NumOfTx++

	pendingStakeQueue, err :=
		accManager.storage.GetPendingStakeQueue(ctx, username)
	if err != nil {
		return err
	}
	// update pending stake queue, remove expired transaction
	accManager.updateTXFromPendingStakeQueue(ctx, accountBank, pendingStakeQueue)
	coinDayParams, err := accManager.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		return err
	}
	if accountBank.Stake.IsGTE(coin) {
		accountBank.Stake = accountBank.Stake.Minus(coin)
	} else {
		coin = coin.Minus(accountBank.Stake)
		accountBank.Stake = types.NewCoinFromInt64(0)

		for len(pendingStakeQueue.PendingStakeList) > 0 {
			pendingStake := pendingStakeQueue.PendingStakeList[0]
			recoverRatio := sdk.NewRat(
				pendingStakeQueue.LastUpdatedAt-pendingStake.StartTime,
				coinDayParams.SecondsToRecoverCoinDayStake)
			if coin.IsGTE(pendingStake.Coin) {
				// if withdraw money is much than first pending transaction, remove first transaction
				coin = coin.Minus(pendingStake.Coin)

				pendingStakeCoinWithoutLastTx :=
					pendingStakeQueue.StakeCoinInQueue.Sub((recoverRatio.Mul(pendingStake.Coin.ToRat())))
				pendingStakeQueue.StakeCoinInQueue = pendingStakeCoinWithoutLastTx

				pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Minus(pendingStake.Coin)
				pendingStakeQueue.PendingStakeList = pendingStakeQueue.PendingStakeList[1:]
			} else {
				// otherwise try to cut first pending transaction
				pendingStakeCoinWithoutSpentCoin :=
					pendingStakeQueue.StakeCoinInQueue.Sub(
						recoverRatio.Mul(coin.ToRat()))
				pendingStakeQueue.StakeCoinInQueue = pendingStakeCoinWithoutSpentCoin

				pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Minus(coin)
				pendingStakeQueue.PendingStakeList[0].Coin = pendingStakeQueue.PendingStakeList[0].Coin.Minus(coin)
				coin = types.NewCoinFromInt64(0)
				break
			}
		}
	}

	if err := accManager.storage.SetPendingStakeQueue(
		ctx, username, pendingStakeQueue); err != nil {
		return err
	}

	if err := accManager.storage.SetBankFromAccountKey(
		ctx, username, accountBank); err != nil {
		return err
	}
	return nil
}

// AddBalanceHistory - add each balance related tx to balance history
func (accManager AccountManager) AddBalanceHistory(
	ctx sdk.Context, username types.AccountKey, numOfTx int64,
	transactionDetail model.Detail) sdk.Error {
	// set balance history
	balanceHistory, err :=
		accManager.storage.GetBalanceHistory(
			ctx, username, numOfTx/types.BalanceHistoryBundleSize)
	if err != nil {
		return err
	}
	if balanceHistory == nil {
		balanceHistory = &model.BalanceHistory{Details: []model.Detail{}}
	}
	balanceHistory.Details = append(balanceHistory.Details, transactionDetail)
	if err := accManager.storage.SetBalanceHistory(
		ctx, username, numOfTx/types.BalanceHistoryBundleSize,
		balanceHistory); err != nil {
		return err
	}

	return nil
}

// UpdateJSONMeta - update user JONS meta data
func (accManager AccountManager) UpdateJSONMeta(
	ctx sdk.Context, username types.AccountKey, JSONMeta string) sdk.Error {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return err
	}
	accountMeta.JSONMeta = JSONMeta

	return accManager.storage.SetMeta(ctx, username, accountMeta)
}

// GetResetKey - get reset public key
func (accManager AccountManager) GetResetKey(
	ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, ErrGetResetKey(username)
	}
	return accountInfo.ResetKey, nil
}

// GetTransactionKey - get transaction public key
func (accManager AccountManager) GetTransactionKey(
	ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, ErrGetTransactionKey(username)
	}
	return accountInfo.TransactionKey, nil
}

// GetAppKey - get app public key
func (accManager AccountManager) GetAppKey(
	ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, ErrGetAppKey(username)
	}
	return accountInfo.AppKey, nil
}

// GetSavingFromBank - get user balance
func (accManager AccountManager) GetSavingFromBank(
	ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return types.Coin{}, ErrGetSavingFromBank(err)
	}
	return accountBank.Saving, nil
}

// GetSequence - get user sequence number
func (accManager AccountManager) GetSequence(
	ctx sdk.Context, username types.AccountKey) (int64, sdk.Error) {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return 0, ErrGetSequence(err)
	}
	return accountMeta.Sequence, nil
}

// GetLastReportOrUpvoteAt - get user last report or upvote time
func (accManager AccountManager) GetLastReportOrUpvoteAt(
	ctx sdk.Context, username types.AccountKey) (int64, sdk.Error) {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return 0, ErrGetLastReportOrUpvoteAt(err)
	}
	return accountMeta.LastReportOrUpvoteAt, nil
}

// UpdateLastReportOrUpvoteAt - update user last report or upvote time to current block time
func (accManager AccountManager) UpdateLastReportOrUpvoteAt(
	ctx sdk.Context, username types.AccountKey) sdk.Error {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return ErrUpdateLastReportOrUpvoteAt(err)
	}
	accountMeta.LastReportOrUpvoteAt = ctx.BlockHeader().Time.Unix()
	return accManager.storage.SetMeta(ctx, username, accountMeta)
}

// GetLastPostAt - get user last post time
func (accManager AccountManager) GetLastPostAt(
	ctx sdk.Context, username types.AccountKey) (int64, sdk.Error) {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return 0, ErrGetLastPostAt(err)
	}
	return accountMeta.LastPostAt, nil
}

// UpdateLastPostAt - update user last post time to current block time
func (accManager AccountManager) UpdateLastPostAt(
	ctx sdk.Context, username types.AccountKey) sdk.Error {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return ErrUpdateLastPostAt(err)
	}
	accountMeta.LastPostAt = ctx.BlockHeader().Time.Unix()
	return accManager.storage.SetMeta(ctx, username, accountMeta)
}

// GetFrozenMoneyList - get user frozen money list
func (accManager AccountManager) GetFrozenMoneyList(
	ctx sdk.Context, username types.AccountKey) ([]model.FrozenMoney, sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return nil, ErrGetFrozenMoneyList(err)
	}
	return accountBank.FrozenMoneyList, nil
}

// IncreaseSequenceByOne - increase user sequence number by one
func (accManager AccountManager) IncreaseSequenceByOne(
	ctx sdk.Context, username types.AccountKey) sdk.Error {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return ErrIncreaseSequenceByOne(err)
	}
	accountMeta.Sequence++
	if err := accManager.storage.SetMeta(ctx, username, accountMeta); err != nil {
		return err
	}
	return nil
}

// AddDirectDeposit - when user received the donation, the donation except friction will be added to
// total income and original income
func (accManager AccountManager) AddDirectDeposit(
	ctx sdk.Context, username types.AccountKey, directDeposit types.Coin) sdk.Error {
	reward, err := accManager.storage.GetReward(ctx, username)
	if err != nil {
		return err
	}
	reward.TotalIncome = reward.TotalIncome.Plus(directDeposit)
	reward.OriginalIncome = reward.OriginalIncome.Plus(directDeposit)
	if err := accManager.storage.SetReward(ctx, username, reward); err != nil {
		return err
	}
	return nil
}

// AddIncomeAndReward - after the evaluate of content value, the original friction
// will be added to original income and friciton income. The actual inflation will
// be added to inflation income, total income and unclaim reward
func (accManager AccountManager) AddIncomeAndReward(
	ctx sdk.Context, username types.AccountKey,
	originalDonation, friction, actualReward types.Coin,
	consumer, postAuthor types.AccountKey, postID string) sdk.Error {
	reward, err := accManager.storage.GetReward(ctx, username)
	if err != nil {
		return err
	}
	reward.TotalIncome = reward.TotalIncome.Plus(actualReward)
	reward.OriginalIncome = reward.OriginalIncome.Plus(friction)
	reward.FrictionIncome = reward.FrictionIncome.Plus(friction)
	reward.InflationIncome = reward.InflationIncome.Plus(actualReward)
	reward.UnclaimReward = reward.UnclaimReward.Plus(actualReward)
	if err := accManager.storage.SetReward(ctx, username, reward); err != nil {
		return err
	}

	// add reward detail
	bank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return err
	}

	rewardDetail := model.RewardDetail{
		OriginalDonation: originalDonation,
		FrictionDonation: friction,
		ActualReward:     actualReward,
		Consumer:         consumer,
		PostAuthor:       postAuthor,
		PostID:           postID,
	}
	if err := accManager.AddRewardHistory(ctx, username, bank.NumOfReward,
		rewardDetail); err != nil {
		return err
	}

	bank.NumOfReward++
	if err := accManager.storage.SetBankFromAccountKey(ctx, username, bank); err != nil {
		return err
	}

	return nil
}

// AddRewardHistory - add reward detail to user reward history
func (accManager AccountManager) AddRewardHistory(
	ctx sdk.Context, username types.AccountKey, numOfReward int64,
	rewardDetail model.RewardDetail) sdk.Error {

	slotNum := numOfReward / types.RewardHistoryBundleSize

	rewardHistory, err := accManager.storage.GetRewardHistory(ctx, username, slotNum)
	if err != nil {
		return err
	}
	if rewardHistory == nil {
		rewardHistory = &model.RewardHistory{Details: []model.RewardDetail{}}
	}

	rewardHistory.Details = append(rewardHistory.Details, rewardDetail)

	if err := accManager.storage.SetRewardHistory(
		ctx, username, slotNum, rewardHistory); err != nil {
		return err
	}

	return nil
}

// ClaimReward - add content reward to user balance
func (accManager AccountManager) ClaimReward(
	ctx sdk.Context, username types.AccountKey) sdk.Error {
	reward, err := accManager.storage.GetReward(ctx, username)
	if err != nil {
		return err
	}
	if err := accManager.AddSavingCoin(
		ctx, username, reward.UnclaimReward, "", "", types.ClaimReward); err != nil {
		return err
	}
	reward.UnclaimReward = types.NewCoinFromInt64(0)
	if err := accManager.storage.SetReward(ctx, username, reward); err != nil {
		return err
	}

	// clear reward history
	if err := accManager.ClearRewardHistory(ctx, username); err != nil {
		return err
	}

	return nil
}

// ClearRewardHistory - clear user reward history
func (accManager AccountManager) ClearRewardHistory(
	ctx sdk.Context, username types.AccountKey) sdk.Error {
	bank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return err
	}

	slotNum := bank.NumOfReward / types.RewardHistoryBundleSize
	for i := int64(0); i <= slotNum; i++ {
		accManager.storage.DeleteRewardHistory(ctx, username, i)
	}

	bank.NumOfReward = 0
	if err := accManager.storage.SetBankFromAccountKey(ctx, username, bank); err != nil {
		return err
	}

	return nil
}

// IsMyFollower - check KV store to check if user in my follower list
func (accManager AccountManager) IsMyFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) bool {
	return accManager.storage.IsMyFollower(ctx, me, follower)
}

// IsMyFollowing - check KV store to check if user in my following list
func (accManager AccountManager) IsMyFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) bool {
	return accManager.storage.IsMyFollowing(ctx, me, following)
}

// SetFollower - update KV store to add follower if doesn't exist
func (accManager AccountManager) SetFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) sdk.Error {
	if accManager.storage.IsMyFollower(ctx, me, follower) {
		return nil
	}
	meta := model.FollowerMeta{
		CreatedAt:    ctx.BlockHeader().Time.Unix(),
		FollowerName: follower,
	}
	accManager.storage.SetFollowerMeta(ctx, me, meta)
	return nil
}

// SetFollowing - update KV store to add following if doesn't exist
func (accManager AccountManager) SetFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) sdk.Error {
	if accManager.storage.IsMyFollowing(ctx, me, following) {
		return nil
	}
	meta := model.FollowingMeta{
		CreatedAt:     ctx.BlockHeader().Time.Unix(),
		FollowingName: following,
	}
	accManager.storage.SetFollowingMeta(ctx, me, meta)
	return nil
}

// RemoveFollower - update KV store to remove follower if exist
func (accManager AccountManager) RemoveFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) sdk.Error {
	if !accManager.storage.IsMyFollower(ctx, me, follower) {
		return nil
	}
	accManager.storage.RemoveFollowerMeta(ctx, me, follower)
	return nil
}

// RemoveFollowing - update KV store to remove following if exist
func (accManager AccountManager) RemoveFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) sdk.Error {
	if !accManager.storage.IsMyFollowing(ctx, me, following) {
		return nil
	}
	accManager.storage.RemoveFollowingMeta(ctx, me, following)
	return nil
}

// CheckUserTPSCapacity - to prevent user spam the chain, every user has a TPS capacity
func (accManager AccountManager) CheckUserTPSCapacity(
	ctx sdk.Context, me types.AccountKey, tpsCapacityRatio sdk.Rat) sdk.Error {
	accountMeta, err := accManager.storage.GetMeta(ctx, me)
	if err != nil {
		return err
	}
	// get update to date user stake
	stake, err := accManager.GetStake(ctx, me)
	if err != nil {
		return err
	}

	// get bandwidth parameters
	bandwidthParams, err := accManager.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return err
	}

	// add virtual coin as the upper limit for capacity
	stake = stake.Plus(bandwidthParams.VirtualCoin)
	// if stake less than last update transaction capacity, set to stake
	if accountMeta.TransactionCapacity.IsGTE(stake) {
		accountMeta.TransactionCapacity = stake
	} else {
		// otherwise try to increase user capacity
		incrementRatio := sdk.NewRat(
			ctx.BlockHeader().Time.Unix()-accountMeta.LastActivityAt,
			bandwidthParams.SecondsToRecoverBandwidth)
		if incrementRatio.GT(sdk.OneRat()) {
			incrementRatio = sdk.OneRat()
		}
		capacityTillStake := stake.Minus(accountMeta.TransactionCapacity)
		increaseCapacity := types.RatToCoin(capacityTillStake.ToRat().Mul(incrementRatio))
		accountMeta.TransactionCapacity =
			accountMeta.TransactionCapacity.Plus(increaseCapacity)
	}
	// based on current tps, calculate current transaction cost
	currentTxCost := types.RatToCoin(
		bandwidthParams.CapacityUsagePerTransaction.ToRat().Mul(tpsCapacityRatio))
	// check if user current capacity is enough or not
	if currentTxCost.IsGT(accountMeta.TransactionCapacity) {
		return ErrAccountTPSCapacityNotEnough(me)
	}
	accountMeta.TransactionCapacity = accountMeta.TransactionCapacity.Minus(currentTxCost)
	accountMeta.LastActivityAt = ctx.BlockHeader().Time.Unix()
	if err := accManager.storage.SetMeta(ctx, me, accountMeta); err != nil {
		return err
	}
	return nil
}

// UpdateDonationRelationship - increase donation relationship times by 1
func (accManager AccountManager) UpdateDonationRelationship(
	ctx sdk.Context, me, other types.AccountKey) sdk.Error {
	relationship, err := accManager.storage.GetRelationship(ctx, me, other)
	if err != nil {
		return err
	}
	if relationship == nil {
		relationship = &model.Relationship{
			DonationTimes: 0,
		}
	}
	relationship.DonationTimes++
	if err := accManager.storage.SetRelationship(ctx, me, other, relationship); err != nil {
		return err
	}
	return nil
}

// AuthorizePermission - userA authorize permission to userB (currently only support auth to a developer)
func (accManager AccountManager) AuthorizePermission(
	ctx sdk.Context, me types.AccountKey, authorizedUser types.AccountKey,
	validityPeriod int64, grantLevel types.Permission, amount types.Coin) sdk.Error {
	d := time.Duration(validityPeriod) * time.Second
	newGrantPubKey := model.GrantPubKey{
		Username:   authorizedUser,
		Permission: grantLevel,
		CreatedAt:  ctx.BlockHeader().Time.Unix(),
		ExpiresAt:  ctx.BlockHeader().Time.Add(d).Unix(),
		Amount:     amount,
	}

	// If grant preauth permission, grant to developer's tx key
	if grantLevel == types.PreAuthorizationPermission {
		txKey, err := accManager.GetTransactionKey(ctx, authorizedUser)
		if err != nil {
			return err
		}
		return accManager.storage.SetGrantPubKey(ctx, me, txKey, &newGrantPubKey)
	}

	// If grant app permission, grant to developer's app key
	if grantLevel == types.AppPermission {
		appKey, err := accManager.GetAppKey(ctx, authorizedUser)
		if err != nil {
			return err
		}
		return accManager.storage.SetGrantPubKey(ctx, me, appKey, &newGrantPubKey)
	}
	return ErrUnsupportGrantLevel()
}

// RevokePermission - revoke permission from a developer
func (accManager AccountManager) RevokePermission(
	ctx sdk.Context, me types.AccountKey, pubKey crypto.PubKey) sdk.Error {
	_, err := accManager.storage.GetGrantPubKey(ctx, me, pubKey)
	if err != nil {
		return err
	}
	accManager.storage.DeleteGrantPubKey(ctx, me, pubKey)
	return nil
}

// CheckSigningPubKeyOwner - given a public key, check if it is valid for given permission
func (accManager AccountManager) CheckSigningPubKeyOwner(
	ctx sdk.Context, me types.AccountKey, signKey crypto.PubKey,
	permission types.Permission, amount types.Coin) (types.AccountKey, sdk.Error) {
	if !accManager.DoesAccountExist(ctx, me) {
		return "", ErrAccountNotFound(me)
	}
	// if permission is reset, only reset key can sign for the msg
	if permission == types.ResetPermission {
		pubKey, err := accManager.GetResetKey(ctx, me)
		if err != nil {
			return "", err
		}
		if reflect.DeepEqual(pubKey, signKey) {
			return me, nil
		}
		return "", ErrCheckResetKey()
	}

	// otherwise transaction key has the highest permission
	pubKey, err := accManager.GetTransactionKey(ctx, me)
	if err != nil {
		return "", err
	}
	if reflect.DeepEqual(pubKey, signKey) {
		return me, nil
	}
	if permission == types.TransactionPermission {
		return "", ErrCheckTransactionKey()
	}

	// if all above keys not matched, check last one, app key
	if permission == types.AppPermission || permission == types.GrantAppPermission {
		pubKey, err = accManager.GetAppKey(ctx, me)
		if err != nil {
			return "", err
		}
		if reflect.DeepEqual(pubKey, signKey) {
			return me, nil
		}
	}

	if permission == types.GrantAppPermission {
		return "", ErrCheckGrantAppKey()
	}

	// if user doesn't use his own key, check his grant user pubkey
	grantPubKey, err := accManager.storage.GetGrantPubKey(ctx, me, signKey)
	if err != nil {
		return "", err
	}
	if grantPubKey.ExpiresAt < ctx.BlockHeader().Time.Unix() {
		accManager.storage.DeleteGrantPubKey(ctx, me, signKey)
		return "", ErrGrantKeyExpired(me)
	}
	if permission != grantPubKey.Permission {
		ErrGrantKeyMismatch(grantPubKey.Username)
	}
	if permission == types.PreAuthorizationPermission {
		txKey, err := accManager.GetTransactionKey(ctx, grantPubKey.Username)
		if err != nil {
			return "", err
		}
		if !reflect.DeepEqual(signKey, txKey) {
			accManager.storage.DeleteGrantPubKey(ctx, me, signKey)
			return "", ErrPreAuthGrantKeyMismatch(grantPubKey.Username)
		}
		if amount.IsGT(grantPubKey.Amount) {
			return "", ErrPreAuthAmountInsufficient(grantPubKey.Username, grantPubKey.Amount, amount)
		}
		grantPubKey.Amount = grantPubKey.Amount.Minus(amount)
		if grantPubKey.Amount.IsEqual(types.NewCoinFromInt64(0)) {
			accManager.storage.DeleteGrantPubKey(ctx, me, signKey)
		} else {
			if err := accManager.storage.SetGrantPubKey(ctx, me, signKey, grantPubKey); err != nil {
				return "", nil
			}
		}
		return grantPubKey.Username, nil
	}
	if permission == types.AppPermission {
		appKey, err := accManager.GetAppKey(ctx, grantPubKey.Username)
		if err != nil {
			return "", err
		}
		if !reflect.DeepEqual(signKey, appKey) {
			accManager.storage.DeleteGrantPubKey(ctx, me, signKey)
			return "", ErrAppGrantKeyMismatch(grantPubKey.Username)
		}
		return grantPubKey.Username, nil
	}
	return "", ErrCheckAuthenticatePubKeyOwner(me)
}

// GetDonationRelationship - get donation relationship between two user
func (accManager AccountManager) GetDonationRelationship(
	ctx sdk.Context, me, other types.AccountKey) (int64, sdk.Error) {
	relationship, err := accManager.storage.GetRelationship(ctx, me, other)
	if err != nil {
		return 0, err
	}
	if relationship == nil {
		return 0, nil
	}
	return relationship.DonationTimes, nil
}

func (accManager AccountManager) addPendingStakeToQueue(
	ctx sdk.Context, username types.AccountKey, bank *model.AccountBank,
	pendingStake model.PendingStake) sdk.Error {
	pendingStakeQueue, err := accManager.storage.GetPendingStakeQueue(ctx, username)
	if err != nil {
		return err
	}
	accManager.updateTXFromPendingStakeQueue(ctx, bank, pendingStakeQueue)
	if len(pendingStakeQueue.PendingStakeList) > 0 && pendingStakeQueue.PendingStakeList[len(pendingStakeQueue.PendingStakeList)-1].StartTime == pendingStake.StartTime {
		pendingStakeQueue.PendingStakeList[len(pendingStakeQueue.PendingStakeList)-1].Coin.Plus(pendingStake.Coin)
	} else {
		pendingStakeQueue.PendingStakeList = append(pendingStakeQueue.PendingStakeList, pendingStake)
		pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Plus(pendingStake.Coin)
	}
	return accManager.storage.SetPendingStakeQueue(ctx, username, pendingStakeQueue)
}

// RecoverAccount - reset three public key pairs
func (accManager AccountManager) RecoverAccount(
	ctx sdk.Context, username types.AccountKey,
	newResetPubKey, newTransactionPubKey, newAppPubKey crypto.PubKey) sdk.Error {
	accInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return err
	}

	accInfo.ResetKey = newResetPubKey
	accInfo.TransactionKey = newTransactionPubKey
	accInfo.AppKey = newAppPubKey
	if err := accManager.storage.SetInfo(ctx, username, accInfo); err != nil {
		return err
	}
	return nil
}

func (accManager AccountManager) updateTXFromPendingStakeQueue(
	ctx sdk.Context, bank *model.AccountBank, pendingStakeQueue *model.PendingStakeQueue) sdk.Error {
	// remove expired transaction
	coinDayParams, err := accManager.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		return err
	}

	currentTimeSlot := ctx.BlockHeader().Time.Unix() / types.CoinDayRecordIntervalSec * types.CoinDayRecordIntervalSec
	for len(pendingStakeQueue.PendingStakeList) > 0 {
		pendingStake := pendingStakeQueue.PendingStakeList[0]
		if pendingStake.EndTime <= currentTimeSlot {
			// remove the transaction from queue, clean stake coin in queue and minus total coin
			//stakeRatioOfThisTransaction means the ratio of stake of this transaction was added last time
			stakeRatioOfThisTransaction := sdk.NewRat(
				pendingStakeQueue.LastUpdatedAt-pendingStake.StartTime,
				coinDayParams.SecondsToRecoverCoinDayStake)
			// remote the stake in the queue of this transaction
			pendingStakeQueue.StakeCoinInQueue =
				pendingStakeQueue.StakeCoinInQueue.Sub(
					stakeRatioOfThisTransaction.Mul(pendingStake.Coin.ToRat()))
			// update bank stake
			bank.Stake = bank.Stake.Plus(pendingStake.Coin)

			pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Minus(pendingStake.Coin)

			pendingStakeQueue.PendingStakeList = pendingStakeQueue.PendingStakeList[1:]
		} else {
			break
		}
	}
	if len(pendingStakeQueue.PendingStakeList) == 0 {
		pendingStakeQueue.TotalCoin = types.NewCoinFromInt64(0)
		pendingStakeQueue.StakeCoinInQueue = sdk.ZeroRat()
	} else {
		// update all pending stake at the same time
		// recoverRatio = (currentTime - lastUpdateTime)/totalRecoverSeconds
		recoverRatio := sdk.NewRat(
			currentTimeSlot-pendingStakeQueue.LastUpdatedAt,
			coinDayParams.SecondsToRecoverCoinDayStake)

		if err != nil {
			return err
		}
		pendingStakeQueue.StakeCoinInQueue =
			pendingStakeQueue.StakeCoinInQueue.Add(
				recoverRatio.Mul(pendingStakeQueue.TotalCoin.ToRat()))
	}

	pendingStakeQueue.LastUpdatedAt = currentTimeSlot
	return nil
}

// AddFrozenMoney - add frozen money to user's frozen money list
func (accManager AccountManager) AddFrozenMoney(
	ctx sdk.Context, username types.AccountKey,
	amount types.Coin, start, interval, times int64) sdk.Error {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return err
	}
	accManager.cleanExpiredFrozenMoney(ctx, accountBank)
	frozenMoney := model.FrozenMoney{
		Amount:   amount,
		StartAt:  start,
		Interval: interval,
		Times:    times,
	}
	accountBank.FrozenMoneyList = append(accountBank.FrozenMoneyList, frozenMoney)

	if err := accManager.storage.SetBankFromAccountKey(ctx, username, accountBank); err != nil {
		return err
	}
	return nil
}

func (accManager AccountManager) cleanExpiredFrozenMoney(ctx sdk.Context, bank *model.AccountBank) {
	idx := 0
	for idx < len(bank.FrozenMoneyList) {
		frozenMoney := bank.FrozenMoneyList[idx]
		if ctx.BlockHeader().Time.Unix() > frozenMoney.StartAt+3600*frozenMoney.Interval*frozenMoney.Times {
			bank.FrozenMoneyList = append(bank.FrozenMoneyList[:idx], bank.FrozenMoneyList[idx+1:]...)
			continue
		}

		idx++
	}
}

// IterateAccounts - iterate accounts in KVStore
func (accManager AccountManager) IterateAccounts(ctx sdk.Context, process func(model.AccountInfo, model.AccountBank) (stop bool)) {
	accManager.storage.IterateAccounts(ctx, process)
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
