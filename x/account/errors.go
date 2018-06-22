package account

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Error constructors
func ErrInvalidLinoAmount() sdk.Error {
	return types.NewError(types.CodeInvalidMsg, fmt.Sprintf("invalid Lino amount"))
}

func ErrUsernameNotFound() sdk.Error {
	return types.NewError(types.CodeUsernameNotFound, fmt.Sprintf("username not found"))
}

func ErrUsernameExist() sdk.Error {
	return types.NewError(types.CodeAccountHandlerFail, fmt.Sprintf("username exist"))
}

func ErrReferrerNotFound() sdk.Error {
	return types.NewError(types.CodeAccountHandlerFail, fmt.Sprintf("referrer not found"))
}

func ErrInvalidUsername(msg string) sdk.Error {
	return types.NewError(types.CodeInvalidUsername, msg)
}

func ErrInvalidMemo() sdk.Error {
	return types.NewError(types.CodeInvalidMemo, fmt.Sprintf("invalid memo"))
}

func ErrInvalidJSONMeta() sdk.Error {
	return types.NewError(types.CodeInvalidMsg, fmt.Sprintf("invalid account JSON meta"))
}

func ErrTransferHandler(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("transfer from account %v failed", accKey))
}

func ErrOpenBankFeeInsufficient(provide types.Coin, expect types.Coin) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail,
		fmt.Sprintf("open bank failed, fee insufficient, need %v, but only %v provided", expect, provide))
}

func ErrAddCoinAccountNotFound(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("add saving coin to account %v failed, account not found", accKey))
}

func ErrAddCoinToAccountSaving(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("add saving coin to account %v failed", accKey))
}

func ErrAddBalanceHistory(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("add balance history to account %v failed", accKey))
}

func ErrMinusCoinToAccount(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("minus coin to account %v failed", accKey))
}

func ErrUpdateFrozenMoney(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("update %v frozen money failed", accKey))
}

func ErrGetBankAddress(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v bank address failed", accKey))
}

func ErrCheckMasterKey() sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("transaction needs master key"))
}

func ErrCheckTransactionKey() sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("transaction needs transaction key"))
}

func ErrCheckAuthenticatePubKeyOwner(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("user %v authenticate public key match failed", accKey))
}

func ErrGetTransactionKey(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v transaction key failed", accKey))
}

func ErrUpdateJSONMeta(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("update %v JSON meta failed", accKey))
}

func ErrGetMasterKey(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v master key failed", accKey))
}

func ErrGetPostKey(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v post key failed", accKey))
}

func ErrGetBankSaving(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v bank saving failed", accKey))
}

func ErrGetBankChecking(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v bank checking failed", accKey))
}

func ErrGetSequence(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v sequence failed", accKey))
}

func ErrGetFrozenMoneyList(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v frozen money list failed", accKey))
}

func ErrIncreaseSequenceByOne(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("increase account %v sequence failed", accKey))
}

func ErrAddIncomeAndReward(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("add income and reward for user %v failed", accKey))
}

func ErrClaimReward(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("claim user %v reward failed", accKey))
}

func ErrGetStake(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get user %v stake failed", accKey))
}

func ErrCheckUserTPSCapacity(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("update user %v transaction capacity failed", accKey))
}

func ErrAccountTPSCapacityNotEnough(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("user %v transaction capacity not enough, please wait", accKey))
}

func ErrAccountAlreadyExists(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("account %v exists", accKey))
}

func ErrRecoverMasterKeyAlreadyOccupied() sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, "recover account failed: master key occupied")
}

func ErrBankAlreadyRegistered() sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("bank connection exists"))
}

func ErrRegisterFeeInsufficient() sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("register fee insufficient"))
}

func ErrAccountCreateFailed(accKey types.AccountKey) sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("create account %v failed", accKey))
}

func ErrUsernameAddressMismatch() sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("username and address mismatch"))
}

func ErrAccountSavingCoinNotEnough() sdk.Error {
	return types.NewError(types.CodeAccountManagerFail, fmt.Sprintf("account bank's saving coins not enough"))
}
