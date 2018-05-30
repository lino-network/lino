package account

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/tx/account/model"
	"github.com/lino-network/lino/types"

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
func (accManager AccountManager) IsAccountExist(ctx sdk.Context, username types.AccountKey) bool {
	return accManager.storage.AccountExist(ctx, username)
}

// create account, caller should make sure the register fee is valid
func (accManager AccountManager) CreateAccount(
	ctx sdk.Context, username types.AccountKey,
	masterKey crypto.PubKey, transactionKey crypto.PubKey, postKey crypto.PubKey,
	registerFee types.Coin) sdk.Error {
	if accManager.IsAccountExist(ctx, username) {
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
		return ErrAccountCreateFailed(username).TraceCause(err, "")
	}

	accountInfo := &model.AccountInfo{
		Username:       username,
		CreatedAt:      ctx.BlockHeader().Time,
		MasterKey:      masterKey,
		TransactionKey: transactionKey,
		PostKey:        postKey,
	}
	if err := accManager.storage.SetInfo(ctx, username, accountInfo); err != nil {
		return ErrAccountCreateFailed(username).TraceCause(err, "")
	}

	accountMeta := &model.AccountMeta{
		LastActivityAt:      ctx.BlockHeader().Time,
		TransactionCapacity: types.NewCoinFromInt64(0),
	}
	if err := accManager.storage.SetMeta(ctx, username, accountMeta); err != nil {
		return ErrAccountCreateFailed(username).TraceCause(err, "")
	}
	if err := accManager.storage.SetReward(ctx, username, &model.Reward{}); err != nil {
		return ErrAccountCreateFailed(username).TraceCause(err, "")
	}
	if err := accManager.storage.SetGrantKeyList(
		ctx, username, &model.GrantKeyList{GrantPubKeyList: []model.GrantPubKey{}}); err != nil {
		return err
	}
	if err := accManager.AddSavingCoin(ctx, username, registerFee, types.TransferIn); err != nil {
		return err
	}
	return nil
}

// use coin to present stake to prevent overflow
func (accManager AccountManager) GetStake(
	ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	bank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return types.NewCoinFromInt64(0), ErrGetStake(username).TraceCause(err, "")
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
	stakeInQueue, err := types.RatToCoin(pendingStakeQueue.StakeCoinInQueue.GetRat())
	totalStake := stake.Plus(stakeInQueue)
	return totalStake, nil
}

func (accManager AccountManager) AddSavingCoin(
	ctx sdk.Context, username types.AccountKey, coin types.Coin,
	coinFrom types.BalanceHistoryDetailType) (err sdk.Error) {
	if !accManager.IsAccountExist(ctx, username) {
		return ErrAddCoinAccountNotFound(username)
	}
	bank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return ErrAddCoinToAccountSaving(username).TraceCause(err, "")
	}
	bank.Saving = bank.Saving.Plus(coin)

	if err := accManager.AddBalanceHistory(ctx, username, coin, coinFrom); err != nil {
		return err
	}

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
		return ErrAddCoinToAccountSaving(username).TraceCause(err, "")
	}

	if err := accManager.storage.SetBankFromAccountKey(ctx, username, bank); err != nil {
		return ErrAddCoinToAccountSaving(username).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) AddBalanceHistory(
	ctx sdk.Context, username types.AccountKey, coin types.Coin,
	detailType types.BalanceHistoryDetailType) sdk.Error {
	// set balance history
	accParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}
	balanceHistory, err :=
		accManager.storage.GetBalanceHistory(
			ctx, username, ctx.BlockHeader().Time/accParams.BalanceHistoryIntervalTime)
	if err != nil {
		return err
	}
	if balanceHistory == nil {
		balanceHistory = &model.BalanceHistory{Details: []model.Detail{}}
	}
	balanceHistory.Details = append(balanceHistory.Details,
		model.Detail{
			Amount:     coin,
			CreatedAt:  ctx.BlockHeader().Time,
			DetailType: detailType,
		})
	if err := accManager.storage.SetBalanceHistory(
		ctx, username, ctx.BlockHeader().Time/accParams.BalanceHistoryIntervalTime, balanceHistory); err != nil {
		return ErrAddCoinToAccountSaving(username).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) MinusSavingCoin(
	ctx sdk.Context, username types.AccountKey, coin types.Coin,
	coinFor types.BalanceHistoryDetailType) (err sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return ErrMinusCoinToAccount(username).TraceCause(err, "")
	}

	accountParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}
	remain := accountBank.Saving.Minus(coin)
	if !remain.IsGTE(accountParams.MinimumBalance) {
		return ErrAccountSavingCoinNotEnough()
	}

	if err := accManager.AddBalanceHistory(ctx, username, coin, coinFor); err != nil {
		return err
	}

	pendingStakeQueue, err :=
		accManager.storage.GetPendingStakeQueue(ctx, username)
	if err != nil {
		return err
	}
	accountBank.Saving = accountBank.Saving.Minus(coin)

	// update pending stake queue, remove expired transaction
	fmt.Println("minus before update", pendingStakeQueue)
	accManager.updateTXFromPendingStakeQueue(ctx, accountBank, pendingStakeQueue)
	fmt.Println("minus after update", pendingStakeQueue)

	coinDayParams, err := accManager.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		return err
	}

	for len(pendingStakeQueue.PendingStakeList) > 0 {
		lengthOfQueue := len(pendingStakeQueue.PendingStakeList)
		pendingStake := pendingStakeQueue.PendingStakeList[lengthOfQueue-1]
		recoverRatio := big.NewRat(
			pendingStakeQueue.LastUpdatedAt-pendingStake.StartTime,
			coinDayParams.SecondsToRecoverCoinDayStake)
		if coin.IsGTE(pendingStake.Coin) {
			// if withdraw money more than last pending transaction, remove last transaction
			coin = coin.Minus(pendingStake.Coin)

			pendingStakeCoinWithoutLastTx :=
				new(big.Rat).Sub(
					pendingStakeQueue.StakeCoinInQueue.GetRat(),
					(new(big.Rat).Mul(recoverRatio, pendingStake.Coin.ToRat())))
			pendingStakeQueue.StakeCoinInQueue = sdk.ToRat(pendingStakeCoinWithoutLastTx)

			pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Minus(pendingStake.Coin)
			pendingStakeQueue.PendingStakeList = pendingStakeQueue.PendingStakeList[:lengthOfQueue-1]
		} else {
			// otherwise try to cut last pending transaction
			pendingStakeCoinWithoutSpentCoin :=
				new(big.Rat).Sub(
					pendingStakeQueue.StakeCoinInQueue.GetRat(),
					(new(big.Rat).Mul(recoverRatio, coin.ToRat())))
			pendingStakeQueue.StakeCoinInQueue = sdk.ToRat(pendingStakeCoinWithoutSpentCoin)

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
		return ErrMinusCoinToAccount(username).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) GetTransactionKey(
	ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, ErrGetTransactionKey(username).TraceCause(err, "")
	}
	return accountInfo.TransactionKey, nil
}

func (accManager AccountManager) GetMasterKey(
	ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, ErrGetMasterKey(username).TraceCause(err, "")
	}
	return accountInfo.MasterKey, nil
}

func (accManager AccountManager) GetPostKey(
	ctx sdk.Context, username types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return nil, ErrGetPostKey(username).TraceCause(err, "")
	}
	return accountInfo.PostKey, nil
}

func (accManager AccountManager) GetSavingFromBank(
	ctx sdk.Context, username types.AccountKey) (types.Coin, sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return types.Coin{}, ErrGetBankSaving(username).TraceCause(err, "")
	}
	return accountBank.Saving, nil
}

func (accManager AccountManager) GetSequence(
	ctx sdk.Context, username types.AccountKey) (int64, sdk.Error) {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return 0, ErrGetSequence(username).TraceCause(err, "")
	}
	return accountMeta.Sequence, nil
}

func (accManager AccountManager) GetFrozenMoneyList(
	ctx sdk.Context, username types.AccountKey) ([]model.FrozenMoney, sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return nil, ErrGetFrozenMoneyList(username).TraceCause(err, "")
	}
	return accountBank.FrozenMoneyList, nil
}

func (accManager AccountManager) IncreaseSequenceByOne(
	ctx sdk.Context, username types.AccountKey) sdk.Error {
	accountMeta, err := accManager.storage.GetMeta(ctx, username)
	if err != nil {
		return ErrGetSequence(username).TraceCause(err, "")
	}
	accountMeta.Sequence += 1
	if err := accManager.storage.SetMeta(ctx, username, accountMeta); err != nil {
		return ErrIncreaseSequenceByOne(username).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) AddIncomeAndReward(
	ctx sdk.Context, username types.AccountKey,
	originIncome, friction, actualReward types.Coin) sdk.Error {
	reward, err := accManager.storage.GetReward(ctx, username)
	if err != nil {
		return ErrAddIncomeAndReward(username).TraceCause(err, "")
	}
	reward.OriginalIncome = reward.OriginalIncome.Plus(originIncome)
	reward.FrictionIncome = reward.FrictionIncome.Plus(friction)
	reward.ActualReward = reward.ActualReward.Plus(actualReward)
	reward.UnclaimReward = reward.UnclaimReward.Plus(actualReward)
	if err := accManager.storage.SetReward(ctx, username, reward); err != nil {
		return ErrAddIncomeAndReward(username).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) ClaimReward(ctx sdk.Context, username types.AccountKey) sdk.Error {
	reward, err := accManager.storage.GetReward(ctx, username)
	if err != nil {
		return ErrClaimReward(username).TraceCause(err, "")
	}
	if err := accManager.AddSavingCoin(ctx, username, reward.UnclaimReward, types.ClaimReward); err != nil {
		return ErrClaimReward(username).TraceCause(err, "")
	}
	reward.UnclaimReward = types.NewCoinFromInt64(0)
	if err := accManager.storage.SetReward(ctx, username, reward); err != nil {
		return ErrClaimReward(username).TraceCause(err, "")
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
		return ErrCheckUserTPSCapacity(me).TraceCause(err, "")
	}
	stake, err := accManager.GetStake(ctx, me)
	if err != nil {
		return ErrCheckUserTPSCapacity(me).TraceCause(err, "")
	}

	bandwidthParams, err := accManager.paramHolder.GetBandwidthParam(ctx)
	if err != nil {
		return err
	}

	if accountMeta.TransactionCapacity.IsGTE(stake) {
		accountMeta.TransactionCapacity = stake
	} else {
		incrementRatio := big.NewRat(
			ctx.BlockHeader().Time-accountMeta.LastActivityAt,
			bandwidthParams.SecondsToRecoverBandwidth)
		if incrementRatio.Cmp(types.OneRat) > 0 {
			incrementRatio = types.OneRat
		}
		capacityTillStake := stake.Minus(accountMeta.TransactionCapacity)
		increateCapacity, err := types.RatToCoin(
			new(big.Rat).Mul(capacityTillStake.ToRat(), incrementRatio))
		if err != nil {
			return err
		}
		accountMeta.TransactionCapacity =
			accountMeta.TransactionCapacity.Plus(increateCapacity)
	}
	currentTxCost, err := types.RatToCoin(
		new(big.Rat).Mul(bandwidthParams.CapacityUsagePerTransaction.ToRat(), tpsCapacityRatio.GetRat()))
	if err != nil {
		return err
	}
	if currentTxCost.IsGT(accountMeta.TransactionCapacity) {
		return ErrAccountTPSCapacityNotEnough(me)
	}
	accountMeta.TransactionCapacity = accountMeta.TransactionCapacity.Minus(currentTxCost)
	accountMeta.LastActivityAt = ctx.BlockHeader().Time
	if err := accManager.storage.SetMeta(ctx, me, accountMeta); err != nil {
		return ErrIncreaseSequenceByOne(me).TraceCause(err, "")
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
		relationship = &model.Relationship{0}
	}
	relationship.DonationTimes += 1
	if err := accManager.storage.SetRelationship(ctx, me, other, relationship); err != nil {
		return err
	}
	return nil
}

func (accManager AccountManager) AuthorizePermission(
	ctx sdk.Context, me types.AccountKey, authorizedUser types.AccountKey,
	validityPeriod int64, grantLevel types.Permission) sdk.Error {
	grantKeyList, err := accManager.storage.GetGrantKeyList(ctx, me)
	if err != nil {
		return err
	}

	idx := 0
	for idx < len(grantKeyList.GrantPubKeyList) {
		if grantKeyList.GrantPubKeyList[idx].ExpiresAt < ctx.BlockHeader().Time ||
			grantKeyList.GrantPubKeyList[idx].Username == authorizedUser {
			grantKeyList.GrantPubKeyList = append(
				grantKeyList.GrantPubKeyList[:idx], grantKeyList.GrantPubKeyList[idx+1:]...)
			continue
		}
		idx += 1
	}

	pubKey, err := accManager.GetPostKey(ctx, authorizedUser)
	if err != nil {
		return err
	}

	newGrantPubKey := model.GrantPubKey{
		Username:  authorizedUser,
		PubKey:    pubKey,
		ExpiresAt: ctx.BlockHeader().Time + validityPeriod,
	}
	grantKeyList.GrantPubKeyList = append(grantKeyList.GrantPubKeyList, newGrantPubKey)
	return accManager.storage.SetGrantKeyList(ctx, me, grantKeyList)
}

func (accManager AccountManager) CheckAuthenticatePubKeyOwner(
	ctx sdk.Context, me types.AccountKey, signKey crypto.PubKey,
	permission types.Permission) (types.AccountKey, sdk.Error) {
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
	pubKey, err = accManager.GetPostKey(ctx, me)
	if err != nil {
		return "", err
	}
	if reflect.DeepEqual(pubKey, signKey) {
		return me, nil
	}

	grantKeyList, err := accManager.storage.GetGrantKeyList(ctx, me)
	if err != nil {
		return "", err
	}
	idx := 0
	for idx < len(grantKeyList.GrantPubKeyList) {
		if grantKeyList.GrantPubKeyList[idx].ExpiresAt < ctx.BlockHeader().Time {
			grantKeyList.GrantPubKeyList = append(
				grantKeyList.GrantPubKeyList[:idx], grantKeyList.GrantPubKeyList[idx+1:]...)
			continue
		}

		if reflect.DeepEqual(grantKeyList.GrantPubKeyList[idx].PubKey, signKey) {
			return grantKeyList.GrantPubKeyList[idx].Username, nil
		}
		idx += 1
	}
	if err := accManager.storage.SetGrantKeyList(ctx, me, grantKeyList); err != nil {
		return "", err
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
	newMasterPubKey, newTransactionPubKey, newPostPubKey crypto.PubKey) sdk.Error {
	accInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return err
	}

	accInfo.MasterKey = newMasterPubKey
	accInfo.PostKey = newPostPubKey
	accInfo.TransactionKey = newTransactionPubKey
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
			stakeRatioOfThisTransaction := big.NewRat(
				pendingStakeQueue.LastUpdatedAt-pendingStake.StartTime,
				coinDayParams.SecondsToRecoverCoinDayStake)
			// remote the stake in the queue of this transaction
			pendingStakeQueue.StakeCoinInQueue =
				pendingStakeQueue.StakeCoinInQueue.Sub(
					sdk.ToRat(new(big.Rat).Mul(stakeRatioOfThisTransaction, pendingStake.Coin.ToRat())))
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
		pendingStakeQueue.StakeCoinInQueue = sdk.ZeroRat
	} else {
		// update all pending stake at the same time
		// recoverRatio = (currentTime - lastUpdateTime)/totalRecoverSeconds
		recoverRatio := big.NewRat(
			ctx.BlockHeader().Time-pendingStakeQueue.LastUpdatedAt,
			coinDayParams.SecondsToRecoverCoinDayStake)

		if err != nil {
			return err
		}
		pendingStakeQueue.StakeCoinInQueue =
			pendingStakeQueue.StakeCoinInQueue.Add(
				sdk.ToRat(new(big.Rat).Mul(recoverRatio, pendingStakeQueue.TotalCoin.ToRat())))
	}

	pendingStakeQueue.LastUpdatedAt = ctx.BlockHeader().Time
	return nil
}

func (accManager AccountManager) AddFrozenMoney(
	ctx sdk.Context, username types.AccountKey,
	amount types.Coin, start, interval, times int64) sdk.Error {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		return ErrUpdateFrozenMoney(username).TraceCause(err, "")
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
	for i, frozenMoney := range bank.FrozenMoneyList {
		if ctx.BlockHeader().Time > frozenMoney.StartAt+frozenMoney.Interval*frozenMoney.Times {
			bank.FrozenMoneyList = append(bank.FrozenMoneyList[:i], bank.FrozenMoneyList[i+1:]...)
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
