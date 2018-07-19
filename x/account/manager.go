package account

import (
	"reflect"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"

	"github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// linoaccount encapsulates all basic struct
type AccountManager struct {
	storage     model.AccountStorage `json:"account_manager"`
	paramHolder param.ParamHolder    `json:"param_holder"`
}

// NewLinoAccount return the account pointer
func NewAccountManager(key sdk.StoreKey, holder param.ParamHolder) AccountManager {
	return AccountManager{
		storage:     model.NewAccountStorage(key),
		paramHolder: holder,
	}
}

// check if account exist
func (accManager AccountManager) DoesAccountExist(ctx sdk.Context, username types.AccountKey) bool {
	return accManager.storage.DoesAccountExist(ctx, username)
}

// create account, caller should make sure the register fee is valid
func (accManager AccountManager) CreateAccount(
	ctx sdk.Context, referrer types.AccountKey, username types.AccountKey,
	masterKey, transactionKey, micropaymentKey, postKey crypto.PubKey,
	registerFee types.Coin) sdk.Error {
	if accManager.DoesAccountExist(ctx, username) {
		return ErrAccountAlreadyExists(username)
	}
	accParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}
	if accParams.RegisterFee.IsGT(registerFee) {
		return ErrRegisterFeeInsufficient()
	}
	if err := accManager.storage.SetPendingStakeQueue(
		ctx, username, &model.PendingStakeQueue{}); err != nil {
		return err
	}

	if err := accManager.storage.SetBankFromAccountKey(ctx, username, &model.AccountBank{}); err != nil {
		return err
	}

	accountInfo := &model.AccountInfo{
		Username:        username,
		CreatedAt:       ctx.BlockHeader().Time,
		MasterKey:       masterKey,
		TransactionKey:  transactionKey,
		MicropaymentKey: micropaymentKey,
		PostKey:         postKey,
	}
	if err := accManager.storage.SetInfo(ctx, username, accountInfo); err != nil {
		return err
	}

	accountMeta := &model.AccountMeta{
		LastActivityAt:       ctx.BlockHeader().Time,
		LastReportOrUpvoteAt: ctx.BlockHeader().Time,
		TransactionCapacity:  accParams.RegisterFee,
	}
	if err := accManager.storage.SetMeta(ctx, username, accountMeta); err != nil {
		return err
	}
	if err := accManager.storage.SetReward(ctx, username, &model.Reward{}); err != nil {
		return err
	}
	if err := accManager.AddSavingCoinWithFullStake(
		ctx, username, accParams.RegisterFee, referrer,
		types.InitAccountWithFullStakeMemo, types.TransferIn); err != nil {
		return ErrAddSavingCoinWithFullStake()
	}
	if err := accManager.AddSavingCoin(
		ctx, username, registerFee.Minus(accParams.RegisterFee), referrer, types.InitAccountRegisterDepositMemo, types.TransferIn); err != nil {
		return ErrAddSavingCoin()
	}
	return nil
}

// use coin to present stake to prevent overflow
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
	stakeInQueue, err := types.RatToCoin(pendingStakeQueue.StakeCoinInQueue)
	totalStake := stake.Plus(stakeInQueue)
	return totalStake, nil
}

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

	if err := accManager.AddBalanceHistory(ctx, username, bank.NumOfTx,
		model.Detail{
			Amount:     coin,
			DetailType: detailType,
			To:         username,
			From:       from,
			CreatedAt:  ctx.BlockHeader().Time,
			Memo:       memo,
		}); err != nil {
		return err
	}
	bank.Saving = bank.Saving.Plus(coin)
	bank.NumOfTx++

	coinDayParams, err := accManager.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		return err
	}

	pendingStake := model.PendingStake{
		StartTime: ctx.BlockHeader().Time,
		EndTime:   ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake,
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

	if err := accManager.AddBalanceHistory(ctx, username, bank.NumOfTx,
		model.Detail{
			Amount:     coin,
			DetailType: detailType,
			To:         username,
			From:       from,
			CreatedAt:  ctx.BlockHeader().Time,
			Memo:       memo,
		}); err != nil {
		return err
	}
	bank.Saving = bank.Saving.Plus(coin)
	bank.Stake = bank.Stake.Plus(coin)
	bank.NumOfTx++

	if err := accManager.storage.SetBankFromAccountKey(ctx, username, bank); err != nil {
		return err
	}
	return nil
}

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

	if err := accManager.AddBalanceHistory(
		ctx, username, accountBank.NumOfTx, model.Detail{
			Amount:     coin,
			DetailType: detailType,
			To:         to,
			From:       username,
			CreatedAt:  ctx.BlockHeader().Time,
			Memo:       memo,
		}); err != nil {
		return err
	}
	accountBank.NumOfTx++
	accountBank.Saving = accountBank.Saving.Minus(coin)

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

func (accManager AccountManager) AddBalanceHistory(
	ctx sdk.Context, username types.AccountKey, numOfTx int64,
	transactionDetail model.Detail) sdk.Error {
	// set balance history
	accParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}
	balanceHistory, err :=
		accManager.storage.GetBalanceHistory(
			ctx, username, numOfTx/accParams.BalanceHistoryBundleSize)
	if err != nil {
		return err
	}
	if balanceHistory == nil {
		balanceHistory = &model.BalanceHistory{Details: []model.Detail{}}
	}
	balanceHistory.Details = append(balanceHistory.Details, transactionDetail)
	if err := accManager.storage.SetBalanceHistory(
		ctx, username, numOfTx/accParams.BalanceHistoryBundleSize,
		balanceHistory); err != nil {
		return err
	}

	return nil
}

func (accManager AccountManager) UpdateJSONMeta(
	ctx sdk.Context, username types.AccountKey, JSONMeta string) sdk.Error {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return err
	}
	accountMeta.JSONMeta = JSONMeta

	return accManager.storage.SetMeta(ctx, username, accountMeta)
}

func (accManager AccountManager) GetMasterKey(
	ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, ErrGetMasterKey(username)
	}
	return accountInfo.MasterKey, nil
}

func (accManager AccountManager) GetTransactionKey(
	ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, ErrGetTransactionKey(username)
	}
	return accountInfo.TransactionKey, nil
}

func (accManager AccountManager) GetMicropaymentKey(
	ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, ErrGetMicropaymentKey(username)
	}
	return accountInfo.MicropaymentKey, nil
}

func (accManager AccountManager) GetPostKey(
	ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, ErrGetPostKey(username)
	}
	return accountInfo.PostKey, nil
}

func (accManager AccountManager) GetSavingFromBank(
	ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return types.Coin{}, ErrGetSavingFromBank(err)
	}
	return accountBank.Saving, nil
}

func (accManager AccountManager) GetSequence(
	ctx sdk.Context, username types.AccountKey) (int64, sdk.Error) {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return 0, ErrGetSequence(err)
	}
	return accountMeta.Sequence, nil
}

// check if account exist
func (accManager AccountManager) GetLastReportOrUpvoteAt(
	ctx sdk.Context, username types.AccountKey) (int64, sdk.Error) {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return 0, ErrGetLastReportOrUpvoteAt(err)
	}
	return accountMeta.LastReportOrUpvoteAt, nil
}

// check if account exist
func (accManager AccountManager) UpdateLastReportOrUpvoteAt(
	ctx sdk.Context, username types.AccountKey) sdk.Error {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return ErrUpdateLastReportOrUpvoteAt(err)
	}
	accountMeta.LastReportOrUpvoteAt = ctx.BlockHeader().Time
	return accManager.storage.SetMeta(ctx, username, accountMeta)
}

func (accManager AccountManager) GetFrozenMoneyList(
	ctx sdk.Context, username types.AccountKey) ([]model.FrozenMoney, sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return nil, ErrGetFrozenMoneyList(err)
	}
	return accountBank.FrozenMoneyList, nil
}

func (accManager AccountManager) IncreaseSequenceByOne(
	ctx sdk.Context, username types.AccountKey) sdk.Error {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return ErrIncreaseSequenceByOne(err)
	}
	accountMeta.Sequence += 1
	if err := accManager.storage.SetMeta(ctx, username, accountMeta); err != nil {
		return err
	}
	return nil
}

// When user received the donation, the donation except friction will be added to
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

// After the evaluate of content value, the original friction will be added to
// original income and friciton income. The actual inflation will be added to
// inflation income, total income and unclaim reward
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

func (accManager AccountManager) AddRewardHistory(
	ctx sdk.Context, username types.AccountKey, numOfReward int64,
	rewardDetail model.RewardDetail) sdk.Error {
	// set reward history
	accParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}

	slotNum := numOfReward / accParams.RewardHistoryBundleSize

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

func (accManager AccountManager) ClearRewardHistory(
	ctx sdk.Context, username types.AccountKey) sdk.Error {
	accParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}

	bank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return err
	}

	slotNum := bank.NumOfReward / accParams.RewardHistoryBundleSize
	for i := int64(0); i <= slotNum; i++ {
		accManager.storage.DeleteRewardHistory(ctx, username, i)
	}

	bank.NumOfReward = 0
	if err := accManager.storage.SetBankFromAccountKey(ctx, username, bank); err != nil {
		return err
	}

	return nil
}

func (accManager AccountManager) IsMyFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) bool {
	return accManager.storage.IsMyFollower(ctx, me, follower)
}

func (accManager AccountManager) IsMyFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) bool {
	return accManager.storage.IsMyFollowing(ctx, me, following)
}

func (accManager AccountManager) SetFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) sdk.Error {
	if accManager.storage.IsMyFollower(ctx, me, follower) {
		return nil
	}
	meta := model.FollowerMeta{
		CreatedAt:    ctx.BlockHeader().Time,
		FollowerName: follower,
	}
	accManager.storage.SetFollowerMeta(ctx, me, meta)
	return nil
}

func (accManager AccountManager) SetFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) sdk.Error {
	if accManager.storage.IsMyFollowing(ctx, me, following) {
		return nil
	}
	meta := model.FollowingMeta{
		CreatedAt:     ctx.BlockHeader().Time,
		FollowingName: following,
	}
	accManager.storage.SetFollowingMeta(ctx, me, meta)
	return nil
}

func (accManager AccountManager) RemoveFollower(
	ctx sdk.Context, me types.AccountKey, follower types.AccountKey) sdk.Error {
	if !accManager.storage.IsMyFollower(ctx, me, follower) {
		return nil
	}
	accManager.storage.RemoveFollowerMeta(ctx, me, follower)
	return nil
}

func (accManager AccountManager) RemoveFollowing(
	ctx sdk.Context, me types.AccountKey, following types.AccountKey) sdk.Error {
	if !accManager.storage.IsMyFollowing(ctx, me, following) {
		return nil
	}
	accManager.storage.RemoveFollowingMeta(ctx, me, following)
	return nil
}

func (accManager AccountManager) CheckUserTPSCapacity(
	ctx sdk.Context, me types.AccountKey, tpsCapacityRatio sdk.Rat) sdk.Error {
	accountMeta, err := accManager.storage.GetMeta(ctx, me)
	if err != nil {
		return err
	}
	stake, err := accManager.GetStake(ctx, me)
	if err != nil {
		return err
	}

	bandwidthParams, err := accManager.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return err
	}

	if accountMeta.TransactionCapacity.IsGTE(stake) {
		accountMeta.TransactionCapacity = stake
	} else {
		incrementRatio := sdk.NewRat(
			ctx.BlockHeader().Time-accountMeta.LastActivityAt,
			bandwidthParams.SecondsToRecoverBandwidth)
		if incrementRatio.GT(sdk.OneRat()) {
			incrementRatio = sdk.OneRat()
		}
		capacityTillStake := stake.Minus(accountMeta.TransactionCapacity)
		increaseCapacity, err := types.RatToCoin(capacityTillStake.ToRat().Mul(incrementRatio))
		if err != nil {
			return err
		}
		accountMeta.TransactionCapacity =
			accountMeta.TransactionCapacity.Plus(increaseCapacity)
	}
	currentTxCost, err := types.RatToCoin(
		bandwidthParams.CapacityUsagePerTransaction.ToRat().Mul(tpsCapacityRatio))
	if err != nil {
		return err
	}
	if currentTxCost.IsGT(accountMeta.TransactionCapacity) {
		return ErrAccountTPSCapacityNotEnough(me)
	}
	accountMeta.TransactionCapacity = accountMeta.TransactionCapacity.Minus(currentTxCost)
	accountMeta.LastActivityAt = ctx.BlockHeader().Time
	if err := accManager.storage.SetMeta(ctx, me, accountMeta); err != nil {
		return err
	}
	return nil
}

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
	relationship.DonationTimes += 1
	if err := accManager.storage.SetRelationship(ctx, me, other, relationship); err != nil {
		return err
	}
	return nil
}

func (accManager AccountManager) AuthorizePermission(
	ctx sdk.Context, me types.AccountKey, authorizedUser types.AccountKey,
	validityPeriod int64, times int64, grantLevel types.Permission) sdk.Error {
	if grantLevel == types.MicropaymentPermission {
		accParams, err := accManager.paramHolder.GetAccountParam(ctx)
		if err != nil {
			return err
		}
		if times > accParams.MaximumMicropaymentGrantTimes {
			return ErrGrantTimesExceedsLimitation(accParams.MaximumMicropaymentGrantTimes)
		}
	}
	newGrantPubKey := model.GrantPubKey{
		Username:   authorizedUser,
		Permission: grantLevel,
		LeftTimes:  times,
		CreatedAt:  ctx.BlockHeader().Time,
		ExpiresAt:  ctx.BlockHeader().Time + validityPeriod,
	}
	if grantLevel == types.MicropaymentPermission {
		micropaymentKey, err := accManager.GetMicropaymentKey(ctx, authorizedUser)
		if err != nil {
			return err
		}
		return accManager.storage.SetGrantPubKey(ctx, me, micropaymentKey, &newGrantPubKey)
	}
	if grantLevel == types.PostPermission {
		postKey, err := accManager.GetPostKey(ctx, authorizedUser)
		if err != nil {
			return err
		}
		return accManager.storage.SetGrantPubKey(ctx, me, postKey, &newGrantPubKey)
	}
	return ErrUnsupportGrantLevel()
}

func (accManager AccountManager) RevokePermission(
	ctx sdk.Context, me types.AccountKey, pubKey crypto.PubKey, grantLevel types.Permission) sdk.Error {
	grantPubKey, err := accManager.storage.GetGrantPubKey(ctx, me, pubKey)
	if err != nil {
		return err
	}
	if grantPubKey.ExpiresAt < ctx.BlockHeader().Time {
		accManager.storage.DeleteGrantPubKey(ctx, me, pubKey)
		return nil
	}
	if grantLevel != grantPubKey.Permission {
		return ErrRevokePermissionLevelMismatch(grantLevel, grantPubKey.Permission)
	}
	accManager.storage.DeleteGrantPubKey(ctx, me, pubKey)
	return nil
}

func (accManager AccountManager) CheckSigningPubKeyOwner(
	ctx sdk.Context, me types.AccountKey, signKey crypto.PubKey,
	permission types.Permission) (types.AccountKey, sdk.Error) {
	if !accManager.DoesAccountExist(ctx, me) {
		return "", ErrAccountNotFound(me)
	}
	// if permission is master, only master key can sign for the msg
	if permission == types.MasterPermission {
		pubKey, err := accManager.GetMasterKey(ctx, me)
		if err != nil {
			return "", err
		}
		if reflect.DeepEqual(pubKey, signKey) {
			return me, nil
		}
		return "", ErrCheckMasterKey()
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

	// then check user's micropayment key
	pubKey, err = accManager.GetMicropaymentKey(ctx, me)
	if err != nil {
		return "", err
	}
	if reflect.DeepEqual(pubKey, signKey) {
		return me, nil
	}

	if permission == types.GrantMicropaymentPermission {
		return "", ErrCheckGrantMicropaymentKey()
	}

	// if all above keys not matched, check last one, post key
	if permission == types.PostPermission || permission == types.GrantPostPermission {
		pubKey, err = accManager.GetPostKey(ctx, me)
		if err != nil {
			return "", err
		}
		if reflect.DeepEqual(pubKey, signKey) {
			return me, nil
		}
	}

	if permission == types.GrantPostPermission {
		return "", ErrCheckGrantPostKey()
	}

	// if user doesn't use his own key, check his grant user pubkey
	grantPubKey, err := accManager.storage.GetGrantPubKey(ctx, me, signKey)
	if err != nil {
		return "", err
	}
	if grantPubKey.ExpiresAt < ctx.BlockHeader().Time {
		accManager.storage.DeleteGrantPubKey(ctx, me, signKey)
		return "", ErrGrantKeyExpired(me)
	}
	if permission != grantPubKey.Permission {
		ErrGrantKeyMismatch(grantPubKey.Username)
	}

	// check again if public key matched
	if permission == types.MicropaymentPermission {
		if grantPubKey.LeftTimes <= 0 {
			accManager.storage.DeleteGrantPubKey(ctx, me, signKey)
			return "", ErrGrantKeyNoLeftTimes(me)
		}
		micropaymentKey, err := accManager.GetMicropaymentKey(ctx, grantPubKey.Username)
		if err != nil {
			return "", err
		}
		if !reflect.DeepEqual(signKey, micropaymentKey) {
			accManager.storage.DeleteGrantPubKey(ctx, me, signKey)
			return "", ErrMicropaymentGrantKeyMismatch(grantPubKey.Username)
		}
		grantPubKey.LeftTimes--
		if err := accManager.storage.SetGrantPubKey(ctx, me, signKey, grantPubKey); err != nil {
			return "", nil
		}
		return grantPubKey.Username, nil
	}
	if permission == types.PostPermission {
		postKey, err := accManager.GetPostKey(ctx, grantPubKey.Username)
		if err != nil {
			return "", err
		}
		if !reflect.DeepEqual(signKey, postKey) {
			accManager.storage.DeleteGrantPubKey(ctx, me, signKey)
			return "", ErrPostGrantKeyMismatch(grantPubKey.Username)
		}
		return grantPubKey.Username, nil
	}
	return "", ErrCheckAuthenticatePubKeyOwner(me)
}

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
	pendingStakeQueue.PendingStakeList = append(pendingStakeQueue.PendingStakeList, pendingStake)
	pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Plus(pendingStake.Coin)
	return accManager.storage.SetPendingStakeQueue(ctx, username, pendingStakeQueue)
}

func (accManager AccountManager) RecoverAccount(
	ctx sdk.Context, username types.AccountKey,
	newMasterPubKey, newTransactionPubKey, newMicropaymentPubKey, newPostPubKey crypto.PubKey) sdk.Error {
	accInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return err
	}

	accInfo.MasterKey = newMasterPubKey
	accInfo.TransactionKey = newTransactionPubKey
	accInfo.MicropaymentKey = newMicropaymentPubKey
	accInfo.PostKey = newPostPubKey
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

	for len(pendingStakeQueue.PendingStakeList) > 0 {
		pendingStake := pendingStakeQueue.PendingStakeList[0]
		if pendingStake.EndTime < ctx.BlockHeader().Time {
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
			ctx.BlockHeader().Time-pendingStakeQueue.LastUpdatedAt,
			coinDayParams.SecondsToRecoverCoinDayStake)

		if err != nil {
			return err
		}
		pendingStakeQueue.StakeCoinInQueue =
			pendingStakeQueue.StakeCoinInQueue.Add(
				recoverRatio.Mul(pendingStakeQueue.TotalCoin.ToRat()))
	}

	pendingStakeQueue.LastUpdatedAt = ctx.BlockHeader().Time
	return nil
}

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
		if ctx.BlockHeader().Time > frozenMoney.StartAt+3600*frozenMoney.Interval*frozenMoney.Times {
			bank.FrozenMoneyList = append(bank.FrozenMoneyList[:idx], bank.FrozenMoneyList[idx+1:]...)
			continue
		}

		idx += 1
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
