package types

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ErrReceiverNotFound - error when receiver user is not found
func ErrReceiverNotFound(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeReceiverNotFound, fmt.Sprintf("receiver %s not found", username))
}

// ErrSenderNotFound - error when sender user is not found
func ErrSenderNotFound(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeSenderNotFound, fmt.Sprintf("sender %s not found", username))
}

// ErrAccountNotFound - error when account is not found
func ErrAccountNotFound(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("account %s not found", username))
}

// ErrAccountBankNotFound - error if account bank is not found
func ErrAccountBankNotFound(addr sdk.Address) sdk.Error {
	return types.NewError(types.CodeAccountBankNotFound, fmt.Sprintf("account bank not found: %s", addr))
}

// ErrReferrerNotFound - error when referrer is not found
func ErrReferrerNotFound(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeReferrerNotFound, fmt.Sprintf("referrer %s not found", username))
}

// ErrAccountAlreadyExists - error when register user already exists
func ErrAccountAlreadyExists(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountAlreadyExists, fmt.Sprintf("account %v already exists", accKey))
}

// ErrAddressAlreadyTaken - error when register address is already took by other username
func ErrAddressAlreadyTaken(addr string) sdk.Error {
	return types.NewError(types.CodeAddressIsTaken, fmt.Sprintf("address %s is already taken", addr))
}

// ErrRegisterFeeInsufficient - error when register fee insufficient
func ErrRegisterFeeInsufficient() sdk.Error {
	return types.NewError(types.CodeRegisterFeeInsufficient, fmt.Sprintf("register fee insufficient"))
}

// ErrAddSavingCoinWithFullCoinDay - error when register deposit with full coin day failed
func ErrAddSavingCoinWithFullCoinDay() sdk.Error {
	return types.NewError(types.CodeAddSavingCoinWithFullCoinDay, fmt.Sprint("failed to add saving coin with full coin day"))
}

// ErrAddSavingCoin - error when register add deposit failed
func ErrAddSavingCoin() sdk.Error {
	return types.NewError(types.CodeAddSavingCoin, fmt.Sprint("failed to add saving coin"))
}

// ErrGetResetKey - error when get reset public key failed
func ErrGetResetKey(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGetResetKey, fmt.Sprintf("get %v reset key failed", accKey))
}

// ErrGetTransactionKey - error when get transaction public key failed
func ErrGetTransactionKey(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGetTransactionKey, fmt.Sprintf("get %v transaction key failed", accKey))
}

// ErrGetAppKey - error when get app public key failed
func ErrGetAppKey(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGetAppKey, fmt.Sprintf("get %v app key failed", accKey))
}

// ErrSigningKey - error when get signing public key failed
func ErrGetSigningKey(username types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGetSigningKeyFailed, fmt.Sprintf("get %v signing key failed", username))
}

// ErrGetSavingFromBank - error when get saving failed
func ErrGetSavingFromBank(err error) sdk.Error {
	return types.NewError(types.CodeGetSavingFromBank, fmt.Sprintf("failed to get saving from bank: %s", err.Error()))
}

// ErrGetSequence - error when get sequence number failed
func ErrGetSequence(err error) sdk.Error {
	return types.NewError(types.CodeGetSequence, fmt.Sprintf("failed to get sequence: %s", err.Error()))
}

// ErrGetAddress - error when get address failed
func ErrGetAddress(err error) sdk.Error {
	return types.NewError(types.CodeGetAddressFailed, fmt.Sprintf("failed to get address: %s", err.Error()))
}

// ErrGetLastReportOrUpvoteAt - error when get last report or upvote time failed
func ErrGetLastReportOrUpvoteAt(err error) sdk.Error {
	return types.NewError(types.CodeGetLastReportOrUpvoteAt, fmt.Sprintf("failed to get last report or upvote at: %s", err.Error()))
}

// ErrGetLastReportOrUpvoteAt - error when update last report or upvote time failed
func ErrUpdateLastReportOrUpvoteAt(err error) sdk.Error {
	return types.NewError(types.CodeUpdateLastReportOrUpvoteAt, fmt.Sprintf("failed to update last report or upvote at: %s", err.Error()))
}

// ErrGetLastPostAt - error when get last post time failed
func ErrGetLastPostAt(err error) sdk.Error {
	return types.NewError(types.CodeGetLastPostAt, fmt.Sprintf("failed to get last post at: %s", err.Error()))
}

// ErrUpdateLastPostAt - error when update last post time failed
func ErrUpdateLastPostAt(err error) sdk.Error {
	return types.NewError(types.CodeUpdateLastPostAt, fmt.Sprintf("failed to update last post at: %s", err.Error()))
}

// ErrGetFrozenMoneyList - error when get frozen money list failed
func ErrGetFrozenMoneyList(err error) sdk.Error {
	return types.NewError(types.CodeGetFrozenMoneyList, fmt.Sprintf("failed to get frozen money list: %s", err.Error()))
}

// ErrFrozenMoneyListTooLong - error when the length of frozen money list exceeds the upper limit
func ErrFrozenMoneyListTooLong() sdk.Error {
	return types.NewError(types.CodeFrozenMoneyListTooLong, fmt.Sprintf("frozen money list too long"))
}

// ErrIncreaseSequenceByOne - error when increase sequence number failed
func ErrIncreaseSequenceByOne(err error) sdk.Error {
	return types.NewError(types.CodeIncreaseSequenceByOne, fmt.Sprintf("failed to increase sequence by one: %s", err.Error()))
}

// ErrCheckResetKey - error when transaction needs reset permission signed by other key
func ErrCheckResetKey() sdk.Error {
	return types.NewError(types.CodeCheckResetKey, fmt.Sprintf("transaction needs reset key"))
}

// ErrCheckTransactionKey - error when transaction needs transaction key permission signed by other key
func ErrCheckTransactionKey() sdk.Error {
	return types.NewError(types.CodeCheckTransactionKey, fmt.Sprintf("transaction needs transaction key"))
}

// ErrCheckGrantAppKey - error when transaction needs app key permission signed by other key
func ErrCheckGrantAppKey() sdk.Error {
	return types.NewError(types.CodeCheckGrantAppKey, fmt.Sprintf("only user's own app key or above can sign grant or revoke app permission msg"))
}

// ErrGrantPubKeyNotFound - error if grant public key is not found
func ErrGrantPubKeyNotFound() sdk.Error {
	return types.NewError(
		types.CodeGrantPubKeyNotFound, fmt.Sprintf("grant public key is not found"))
}

// ErrCheckAuthenticatePubKeyOwner - error when transaction signed by invalid public key
func ErrCheckAuthenticatePubKeyOwner(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeCheckAuthenticatePubKeyOwner, fmt.Sprintf("user %v authenticate public key match failed", accKey))
}

// ErrCheckAuthenticatePubKeyAddress - error when transaction signed by invalid public key
func ErrCheckAuthenticatePubKeyAddress(addr sdk.AccAddress) sdk.Error {
	return types.NewError(types.CodeCheckAuthenticatePubKeyOwner, fmt.Sprintf("address %v authenticate public key match failed", addr))
}

// ErrGrantKeyExpired - error when transaction signed by expired grant public key
func ErrGrantKeyExpired(owner types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGrantKeyExpired, fmt.Sprintf("grant user %v key expired", owner))
}

// ErrGrantKeyMismatch - error when transaction signed by mismatch permission grant public key
func ErrGrantKeyMismatch(owner types.AccountKey) sdk.Error {
	return types.NewError(types.CodeGrantKeyMismatch, fmt.Sprintf("grant user %v key can't match his own key", owner))
}

// ErrGrantKeyMismatch - error when transaction signed by mismatch app permission grant public key
func ErrAppGrantKeyMismatch(owner types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAppGrantKeyMismatch, fmt.Sprintf("grant user %v app key can't match his own key", owner))
}

// ErrPreAuthGrantKeyMismatch - error when transaction signed by mismatch preauth permission grant public key
func ErrPreAuthGrantKeyMismatch(owner types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAppGrantKeyMismatch, fmt.Sprintf("grant user %v transaction key can't match his own key", owner))
}

// ErrPreAuthAmountInsufficient - error when transaction cost coin exceeds preauth amount
func ErrPreAuthAmountInsufficient(owner types.AccountKey, balance, consume types.Coin) sdk.Error {
	return types.NewError(
		types.CodeAppGrantKeyMismatch,
		fmt.Sprintf("grant user %v doesn't have enough preauthorization balance, have %v, wanna consume %v", owner, balance, consume))
}

// ErrUnsupportGrantLevel - error when grant permission not supported
func ErrUnsupportGrantLevel() sdk.Error {
	return types.NewError(types.CodeUnsupportGrantLevel, fmt.Sprintf("unsupport grant level"))
}

// ErrRevokePermissionLevelMismatch - error when revoke permission doesn't match target public key
func ErrRevokePermissionLevelMismatch(got, expect types.Permission) sdk.Error {
	return types.NewError(types.CodeRevokePermissionLevelMismatch, fmt.Sprintf("revoke permission level mismatch, got %v, expect %v", got, expect))
}

// ErrCheckUserTPSCapacity - error when check user capacity failed
func ErrCheckUserTPSCapacity(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeCheckUserTPSCapacity, fmt.Sprintf("update user %v transaction capacity failed", accKey))
}

// ErrAccountTPSCapacityNotEnough - error when user tps capacity not enough
func ErrAccountTPSCapacityNotEnough(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountTPSCapacityNotEnough, fmt.Sprintf("user %v transaction capacity not enough, please wait", accKey))
}

// ErrAccountSavingCoinNotEnough - error when user saving balance not enough
func ErrAccountSavingCoinNotEnough() sdk.Error {
	return types.NewError(types.CodeAccountSavingCoinNotEnough, fmt.Sprintf("account bank's saving coins not enough"))
}

// ErrPoolNotEnough - error if pool balance is not enough.
func ErrPoolNotEnough(name types.PoolName) sdk.Error {
	return types.NewError(types.CodePoolNotEnough, fmt.Sprintf("not enough balance: %s", name))
}

// ErrInvalidUsername - error when username is invalid
func ErrInvalidUsername(msg string) sdk.Error {
	return types.NewError(types.CodeInvalidUsername, msg)
}

// ErrInvalidMemo - error when memo is invalid (length too long)
func ErrInvalidMemo() sdk.Error {
	return types.NewError(types.CodeInvalidMemo, fmt.Sprintf("invalid memo"))
}

// ErrInvalidJSONMeta - error when JSON meta is invalid (length too long)
func ErrInvalidJSONMeta() sdk.Error {
	return types.NewError(types.CodeInvalidJSONMeta, fmt.Sprintf("invalid account JSON meta"))
}

// ErrQueryFailed - error when query account store failed
func ErrQueryFailed() sdk.Error {
	return types.NewError(types.CodeAccountQueryFailed, fmt.Sprintf("query account store failed"))
}

// ErrQueryTxFaild - error when query tx failed
func ErrQueryTxFailed(msg string) sdk.Error {
	return types.NewError(types.CodeAccountQueryFailed, fmt.Sprintf("query tx failed, err: %s", msg))
}

// ErrNegativeMoveAmount -
func ErrNegativeMoveAmount(amount types.Coin) sdk.Error {
	return types.NewError(
		types.CodeNegativeMoveAmount, fmt.Sprintf("negative amount: %+v", amount))
}

// ErrPoolNotFound - error if pool is not found
func ErrPoolNotFound(name types.PoolName) sdk.Error {
       return types.NewError(types.CodePoolNotFound, fmt.Sprintf("pool not found: %s", name))
}
