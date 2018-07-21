package account

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ErrFolloweeNotFound(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeFolloweeNotFound, fmt.Sprintf("followee %s not found", username))
}

func ErrFollowerNotFound(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeFollowerNotFound, fmt.Sprintf("follower %s not found", username))
}

func ErrReceiverNotFound(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeReceiverNotFound, fmt.Sprintf("receiver %s not found", username))
}

func ErrSenderNotFound(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeSenderNotFound, fmt.Sprintf("sender %s not found", username))
}

func ErrAccountNotFound(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("account %s not found", username))
}

func ErrReferrerNotFound(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeReferrerNotFound, fmt.Sprintf("referrer %s not found", username))
}

func ErrAccountAlreadyExists(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountAlreadyExists, fmt.Sprintf("account %v already exists", accKey))
}

func ErrRegisterFeeInsufficient() sdk.Error {
	return types.NewError(types.CodeRegisterFeeInsufficient, fmt.Sprintf("register fee insufficient"))
}

func ErrAddSavingCoinWithFullStake() sdk.Error {
	return types.NewError(types.CodeAddSavingCoinWithFullStake, fmt.Sprint("failed to add saving coin with full stake"))
}

func ErrAddSavingCoin() sdk.Error {
	return types.NewError(types.CodeAddSavingCoin, fmt.Sprint("failed to add saving coin"))
}

func ErrGetResetKey(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGetResetKey, fmt.Sprintf("get %v reset key failed", accKey))
}

func ErrGetTransactionKey(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGetTransactionKey, fmt.Sprintf("get %v transaction key failed", accKey))
}

func ErrGetPostKey(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGetPostKey, fmt.Sprintf("get %v post key failed", accKey))
}

func ErrGetSavingFromBank(err error) sdk.Error {
	return types.NewError(types.CodeGetSavingFromBank, fmt.Sprintf("failed to get saving from bank: %s", err.Error()))
}

func ErrGetSequence(err error) sdk.Error {
	return types.NewError(types.CodeGetSequence, fmt.Sprintf("failed to get sequence: %s", err.Error()))
}

func ErrGetLastReportOrUpvoteAt(err error) sdk.Error {
	return types.NewError(types.CodeGetLastReportOrUpvoteAt, fmt.Sprintf("failed to get last report or upvote at: %s", err.Error()))
}

func ErrUpdateLastReportOrUpvoteAt(err error) sdk.Error {
	return types.NewError(types.CodeUpdateLastReportOrUpvoteAt, fmt.Sprintf("failed to update last report or upvote at: %s", err.Error()))
}

func ErrGetFrozenMoneyList(err error) sdk.Error {
	return types.NewError(types.CodeGetFrozenMoneyList, fmt.Sprintf("failed to get frozen money list: %s", err.Error()))
}

func ErrIncreaseSequenceByOne(err error) sdk.Error {
	return types.NewError(types.CodeIncreaseSequenceByOne, fmt.Sprintf("failed to increase sequence by one: %s", err.Error()))
}

func ErrCheckResetKey() sdk.Error {
	return types.NewError(types.CodeCheckResetKey, fmt.Sprintf("transaction needs reset key"))
}

func ErrCheckTransactionKey() sdk.Error {
	return types.NewError(types.CodeCheckTransactionKey, fmt.Sprintf("transaction needs transaction key"))
}

func ErrCheckGrantPostKey() sdk.Error {
	return types.NewError(types.CodeCheckGrantPostKey, fmt.Sprintf("only user's own post key or above can sign grant or revoke post permission msg"))
}

func ErrCheckAuthenticatePubKeyOwner(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeCheckAuthenticatePubKeyOwner, fmt.Sprintf("user %v authenticate public key match failed", accKey))
}

func ErrGrantKeyExpired(owner types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGrantKeyExpired, fmt.Sprintf("grant user %v key expired", owner))
}

func ErrGrantKeyNoLeftTimes(owner types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGrantKeyNoLeftTimes, fmt.Sprintf("grant user %v key no left times", owner))
}

func ErrGrantKeyMismatch(owner types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGrantKeyMismatch, fmt.Sprintf("grant user %v key can't match his own key", owner))
}

func ErrPostGrantKeyMismatch(owner types.AccountKey) sdk.Error {
	return types.NewError(types.CodePostGrantKeyMismatch, fmt.Sprintf("grant user %v post key can't match his own key", owner))
}

func ErrGrantTimesExceedsLimitation(limitation int64) sdk.Error {
	return types.NewError(types.CodeGrantTimesExceedsLimitation, fmt.Sprintf("grant times exceeds %v limitation", limitation))
}

func ErrUnsupportGrantLevel() sdk.Error {
	return types.NewError(types.CodeUnsupportGrantLevel, fmt.Sprintf("unsupport grant level"))
}

func ErrRevokePermissionLevelMismatch(got, expect types.Permission) sdk.Error {
	return types.NewError(types.CodeRevokePermissionLevelMismatch, fmt.Sprintf("revoke permission level mismatch, got %v, expect %v", got, expect))
}

func ErrCheckUserTPSCapacity(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeCheckUserTPSCapacity, fmt.Sprintf("update user %v transaction capacity failed", accKey))
}

func ErrAccountTPSCapacityNotEnough(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountTPSCapacityNotEnough, fmt.Sprintf("user %v transaction capacity not enough, please wait", accKey))
}

func ErrAccountSavingCoinNotEnough() sdk.Error {
	return types.NewError(types.CodeAccountSavingCoinNotEnough, fmt.Sprintf("account bank's saving coins not enough"))
}

func ErrInvalidUsername(msg string) sdk.Error {
	return types.NewError(types.CodeInvalidUsername, msg)
}

func ErrInvalidMemo() sdk.Error {
	return types.NewError(types.CodeInvalidMemo, fmt.Sprintf("invalid memo"))
}

func ErrInvalidJSONMeta() sdk.Error {
	return types.NewError(types.CodeInvalidJSONMeta, fmt.Sprintf("invalid account JSON meta"))
}
