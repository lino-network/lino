package account

import (
	"fmt"
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
func (accManager AccountManager) IsAccountExist(
	ctx sdk.Context, accKey types.AccountKey) bool {
	accountInfo, _ := accManager.storage.GetInfo(ctx, accKey)
	return accountInfo != nil
}

// Implements types.AccountManager.
func (accManager AccountManager) CreateAccount(
	ctx sdk.Context, accKey types.AccountKey,
	masterKey crypto.PubKey, transactionKey crypto.PubKey, postKey crypto.PubKey) sdk.Error {
	if accManager.IsAccountExist(ctx, accKey) {
		return ErrAccountAlreadyExists(accKey)
	}
	bank, err := accManager.storage.GetBankFromAddress(ctx, masterKey.Address())
	if err != nil {
		return ErrAccountCreateFailed(accKey)
	}
	if bank.Username != "" {
		return ErrBankAlreadyRegistered()
	}
	accParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}
	if !bank.Saving.IsGTE(accParams.RegisterFee) {
		return ErrRegisterFeeInsufficient()
	}

	accountInfo := &model.AccountInfo{
		Username:       accKey,
		CreatedAt:      ctx.BlockHeader().Time,
		MasterKey:      masterKey,
		TransactionKey: transactionKey,
		PostKey:        postKey,
		Address:        masterKey.Address(),
	}
	if err := accManager.storage.SetInfo(ctx, accKey, accountInfo); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}

	bank.Username = accKey
	if err := accManager.storage.SetBankFromAddress(
		ctx, masterKey.Address(), bank); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}

	accountMeta := &model.AccountMeta{
		LastActivityAt:      ctx.BlockHeader().Time,
		TransactionCapacity: types.NewCoin(0),
	}
	if err := accManager.storage.SetMeta(ctx, accKey, accountMeta); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}
	reward :=
		&model.Reward{types.NewCoin(0), types.NewCoin(0), types.NewCoin(0), types.NewCoin(0)}
	if err := accManager.storage.SetReward(ctx, accKey, reward); err != nil {
		return ErrAccountCreateFailed(accKey).TraceCause(err, "")
	}

	if err := accManager.storage.SetGrantKeyList(
		ctx, accKey, &model.GrantKeyList{GrantPubKeyList: []model.GrantPubKey{}}); err != nil {
		return err
	}
	return nil
}

// use coin to present stake to prevent overflow
func (accManager AccountManager) GetStake(
	ctx sdk.Context, accKey types.AccountKey) (types.Coin, sdk.Error) {
	bank, err := accManager.storage.GetBankFromAccountKey(ctx, accKey)
	if err != nil {
		return types.NewCoin(0), ErrGetStake(accKey).TraceCause(err, "")
	}
	pendingStakeQueue, err := accManager.storage.GetPendingStakeQueue(ctx, bank.Address)
	if err != nil {
		return types.NewCoin(0), err
	}

	accManager.updateTXFromPendingStakeQueue(ctx, bank, pendingStakeQueue)

	stake := bank.Stake
	if err := accManager.storage.SetPendingStakeQueue(
		ctx, bank.Address, pendingStakeQueue); err != nil {
		return types.NewCoin(0), err
	}

	if err := accManager.storage.SetBankFromAddress(ctx, bank.Address, bank); err != nil {
		return types.NewCoin(0), err
	}
	fmt.Println(stake, pendingStakeQueue.StakeCoinInQueue, pendingStakeQueue, ctx.BlockHeader().Time)
	return stake.Plus(types.RatToCoin(pendingStakeQueue.StakeCoinInQueue)), nil
}

func (accManager AccountManager) AddSavingCoin(
	ctx sdk.Context, accKey types.AccountKey, coin types.Coin) (err sdk.Error) {
	address, err := accManager.GetBankAddress(ctx, accKey)
	if err != nil {
		return ErrAddCoinToAccountSaving(accKey).TraceCause(err, "")
	}
	if err := accManager.AddSavingCoinToAddress(ctx, address, coin); err != nil {
		return ErrAddCoinToAccountSaving(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) AddSavingCoinToAddress(
	ctx sdk.Context, address sdk.Address, coin types.Coin) (err sdk.Error) {
	if coin.IsZero() {
		return nil
	}
	bank, _ := accManager.storage.GetBankFromAddress(ctx, address)
	if bank == nil {
		bank = &model.AccountBank{
			Address: address,
			Saving:  coin,
		}
		if err := accManager.storage.SetPendingStakeQueue(
			ctx, address, &model.PendingStakeQueue{}); err != nil {
			return err
		}
	} else {
		bank.Saving = bank.Saving.Plus(coin)
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
	if err := accManager.addPendingStakeToQueue(ctx, address, bank, pendingStake); err != nil {
		return ErrAddCoinToAddress(address).TraceCause(err, "")
	}
	if err := accManager.storage.SetBankFromAddress(ctx, bank.Address, bank); err != nil {
		return ErrAddCoinToAddress(address).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) AddCheckingCoin(
	ctx sdk.Context, accKey types.AccountKey, coin types.Coin) (err sdk.Error) {
	address, err := accManager.GetBankAddress(ctx, accKey)
	if err != nil {
		return ErrAddCoinToAccountChecking(accKey).TraceCause(err, "")
	}
	if err := accManager.AddCheckingCoinToAddress(ctx, address, coin); err != nil {
		return ErrAddCoinToAccountChecking(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) AddCheckingCoinToAddress(
	ctx sdk.Context, address sdk.Address, coin types.Coin) (err sdk.Error) {
	if coin.IsZero() {
		return nil
	}
	bank, _ := accManager.storage.GetBankFromAddress(ctx, address)
	if bank == nil {
		bank = &model.AccountBank{
			Address:  address,
			Checking: coin,
		}
		if err := accManager.storage.SetPendingStakeQueue(
			ctx, address, &model.PendingStakeQueue{}); err != nil {
			return err
		}
	} else {
		bank.Checking = bank.Checking.Plus(coin)
	}
	if err := accManager.storage.SetBankFromAddress(ctx, bank.Address, bank); err != nil {
		return ErrAddCoinToAddress(address).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) MinusSavingCoin(
	ctx sdk.Context, accKey types.AccountKey, coin types.Coin) (err sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, accKey)
	if err != nil {
		return ErrMinusCoinToAccount(accKey).TraceCause(err, "")
	}

	accountParams, err := accManager.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err
	}
	fmt.Println(accountBank, coin, accountParams.MinimumBalance)
	if !accountBank.Saving.Minus(coin).IsGTE(accountParams.MinimumBalance) {
		return ErrAccountSavingCoinNotEnough()
	}
	pendingStakeQueue, err :=
		accManager.storage.GetPendingStakeQueue(ctx, accountBank.Address)
	if err != nil {
		return err
	}
	accountBank.Saving = accountBank.Saving.Minus(coin)

	// update pending stake queue, remove expired transaction
	accManager.updateTXFromPendingStakeQueue(ctx, accountBank, pendingStakeQueue)

	coinDayParams, err := accManager.paramHolder.GetCoinDayParam(ctx)
	if err != nil {
		return err
	}

	for len(pendingStakeQueue.PendingStakeList) > 0 {
		lengthOfQueue := len(pendingStakeQueue.PendingStakeList)
		pendingStake := pendingStakeQueue.PendingStakeList[lengthOfQueue-1]
		if coin.IsGTE(pendingStake.Coin) {
			// if withdraw money more than last pending transaction, remove last transaction
			coin = coin.Minus(pendingStake.Coin)
			pendingStakeQueue.StakeCoinInQueue =
				pendingStakeQueue.StakeCoinInQueue.Sub(sdk.NewRat(
					pendingStakeQueue.LastUpdatedAt-pendingStake.StartTime,
					coinDayParams.SecondsToRecoverCoinDayStake).Mul(pendingStake.Coin.ToRat()))
			pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Minus(pendingStake.Coin)
			pendingStakeQueue.PendingStakeList =
				pendingStakeQueue.PendingStakeList[:lengthOfQueue-1]
		} else {
			// otherwise try to cut last pending transaction
			pendingStakeQueue.StakeCoinInQueue =
				pendingStakeQueue.StakeCoinInQueue.Sub(sdk.NewRat(
					pendingStakeQueue.LastUpdatedAt-pendingStake.StartTime,
					coinDayParams.SecondsToRecoverCoinDayStake).Mul(coin.ToRat()))
			pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Minus(coin)
			pendingStakeQueue.PendingStakeList[lengthOfQueue-1].Coin =
				pendingStakeQueue.PendingStakeList[lengthOfQueue-1].Coin.Minus(coin)
			coin = types.NewCoin(0)
			break
		}
	}
	if coin.IsPositive() {
		accountBank.Stake = accountBank.Saving
	}
	if err := accManager.storage.SetPendingStakeQueue(
		ctx, accountBank.Address, pendingStakeQueue); err != nil {
		return err
	}

	if err := accManager.storage.SetBankFromAddress(
		ctx, accountBank.Address, accountBank); err != nil {
		return ErrMinusCoinToAccount(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) MinusCheckingCoin(
	ctx sdk.Context, accKey types.AccountKey, coin types.Coin) (err sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, accKey)
	if err != nil {
		return ErrMinusCoinToAccount(accKey).TraceCause(err, "")
	}
	if !accountBank.Checking.IsGTE(coin) {
		return ErrAccountCheckingCoinNotEnough()
	}
	accountBank.Checking = accountBank.Checking.Minus(coin)
	if err := accManager.storage.SetBankFromAddress(
		ctx, accountBank.Address, accountBank); err != nil {
		return ErrMinusCoinToAccount(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) GetBankAddress(
	ctx sdk.Context, accKey types.AccountKey) (sdk.Address, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, accKey)
	if err != nil {
		return nil, ErrGetBankAddress(accKey).TraceCause(err, "")
	}
	return accountInfo.Address, nil
}

func (accManager AccountManager) GetTransactionKey(
	ctx sdk.Context, accKey types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, accKey)
	if err != nil {
		return nil, ErrGetTransactionKey(accKey).TraceCause(err, "")
	}
	return accountInfo.TransactionKey, nil
}

func (accManager AccountManager) GetMasterKey(
	ctx sdk.Context, accKey types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, accKey)
	if err != nil {
		return nil, ErrGetMasterKey(accKey).TraceCause(err, "")
	}
	return accountInfo.MasterKey, nil
}

func (accManager AccountManager) GetPostKey(
	ctx sdk.Context, accKey types.AccountKey) (crypto.PubKey, sdk.Error) {
	accountInfo, err := accManager.storage.GetInfo(ctx, accKey)
	if err != nil {
		return nil, ErrGetPostKey(accKey).TraceCause(err, "")
	}
	return accountInfo.PostKey, nil
}

func (accManager AccountManager) GetBankSaving(
	ctx sdk.Context, accKey types.AccountKey) (types.Coin, sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, accKey)
	if err != nil {
		return types.Coin{}, ErrGetBankSaving(accKey).TraceCause(err, "")
	}
	return accountBank.Saving, nil
}

func (accManager AccountManager) GetBankChecking(
	ctx sdk.Context, accKey types.AccountKey) (types.Coin, sdk.Error) {
	accountBank, err := accManager.storage.GetBankFromAccountKey(ctx, accKey)
	if err != nil {
		return types.Coin{}, ErrGetBankSaving(accKey).TraceCause(err, "")
	}
	return accountBank.Checking, nil
}

func (accManager AccountManager) GetSequence(
	ctx sdk.Context, accKey types.AccountKey) (int64, sdk.Error) {
	accountMeta, err := accManager.storage.GetMeta(ctx, accKey)
	if err != nil {
		return 0, ErrGetSequence(accKey).TraceCause(err, "")
	}
	return accountMeta.Sequence, nil
}

func (accManager AccountManager) IncreaseSequenceByOne(
	ctx sdk.Context, accKey types.AccountKey) sdk.Error {
	accountMeta, err := accManager.storage.GetMeta(ctx, accKey)
	if err != nil {
		return ErrGetSequence(accKey).TraceCause(err, "")
	}
	accountMeta.Sequence += 1
	if err := accManager.storage.SetMeta(ctx, accKey, accountMeta); err != nil {
		return ErrIncreaseSequenceByOne(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) AddIncomeAndReward(
	ctx sdk.Context, accKey types.AccountKey,
	originIncome, friction, actualReward types.Coin) sdk.Error {
	reward, err := accManager.storage.GetReward(ctx, accKey)
	if err != nil {
		return ErrAddIncomeAndReward(accKey).TraceCause(err, "")
	}
	reward.OriginalIncome = reward.OriginalIncome.Plus(originIncome)
	reward.FrictionIncome = reward.FrictionIncome.Plus(friction)
	reward.ActualReward = reward.ActualReward.Plus(actualReward)
	reward.UnclaimReward = reward.UnclaimReward.Plus(actualReward)
	if err := accManager.storage.SetReward(ctx, accKey, reward); err != nil {
		return ErrAddIncomeAndReward(accKey).TraceCause(err, "")
	}
	return nil
}

func (accManager AccountManager) ClaimReward(ctx sdk.Context, accKey types.AccountKey) sdk.Error {
	reward, err := accManager.storage.GetReward(ctx, accKey)
	if err != nil {
		return ErrClaimReward(accKey).TraceCause(err, "")
	}
	if err := accManager.AddSavingCoin(ctx, accKey, reward.UnclaimReward); err != nil {
		return ErrClaimReward(accKey).TraceCause(err, "")
	}
	reward.UnclaimReward = types.NewCoin(0)
	if err := accManager.storage.SetReward(ctx, accKey, reward); err != nil {
		return ErrClaimReward(accKey).TraceCause(err, "")
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
		incrementRatio := sdk.NewRat(
			ctx.BlockHeader().Time-accountMeta.LastActivityAt,
			bandwidthParams.SecondsToRecoverBandwidth)
		if incrementRatio.GT(sdk.OneRat) {
			incrementRatio = sdk.OneRat
		}
		accountMeta.TransactionCapacity =
			accountMeta.TransactionCapacity.Plus(types.RatToCoin(
				stake.Minus(accountMeta.TransactionCapacity).ToRat().Mul(incrementRatio)))
	}
	currentTxCost := types.RatToCoin(
		bandwidthParams.CapacityUsagePerTransaction.ToRat().Mul(tpsCapacityRatio))
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
	validityPeriod int64, grantLevel int64) sdk.Error {
	pubKey, err := accManager.GetPostKey(ctx, authorizedUser)
	if err != nil {
		return err
	}

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
	ctx sdk.Context, address sdk.Address, bank *model.AccountBank,
	pendingStake model.PendingStake) sdk.Error {
	pendingStakeQueue, err := accManager.storage.GetPendingStakeQueue(ctx, address)
	if err != nil {
		return err
	}
	accManager.updateTXFromPendingStakeQueue(ctx, bank, pendingStakeQueue)
	pendingStakeQueue.PendingStakeList = append(pendingStakeQueue.PendingStakeList, pendingStake)
	pendingStakeQueue.TotalCoin = pendingStakeQueue.TotalCoin.Plus(pendingStake.Coin)
	return accManager.storage.SetPendingStakeQueue(ctx, address, pendingStakeQueue)
}

func (accManager AccountManager) RecoverAccount(
	ctx sdk.Context, username types.AccountKey, newPostKey, newTransactionKey crypto.PubKey) sdk.Error {
	accInfo, err := accManager.storage.GetInfo(ctx, username)
	if err != nil {
		return err
	}
	accInfo.PostKey = newPostKey
	accInfo.TransactionKey = newTransactionKey
	return accManager.storage.SetInfo(ctx, username, accInfo)
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
				pendingStakeQueue.StakeCoinInQueue.Sub(stakeRatioOfThisTransaction.Mul(pendingStake.Coin.ToRat()))
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
		// recoverRatio = (currentTime - lastUpdateTime)/totalRecoverSeconds
		recoverRatio := sdk.NewRat(
			ctx.BlockHeader().Time-pendingStakeQueue.LastUpdatedAt,
			coinDayParams.SecondsToRecoverCoinDayStake)
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
