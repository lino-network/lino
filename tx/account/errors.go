package account

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrInvalidLinoAmount() sdk.Error {
	return sdk.NewError(types.CodeInvalidMsg, fmt.Sprintf("invalid Lino amount"))
}

func ErrUsernameNotFound() sdk.Error {
	return sdk.NewError(types.CodeUsernameNotFound, fmt.Sprintf("username not found"))
}

func ErrInvalidUsername() sdk.Error {
	return sdk.NewError(types.CodeInvalidUsername, fmt.Sprintf("invalida Username"))
}

func ErrTransferHandler(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("transfer from account %v failed", accKey))
}

func ErrOpenBankFeeInsufficient(provide types.Coin, expect types.Coin) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail,
		fmt.Sprintf("open bank failed, fee insufficient, need %v, but only %v provided", expect, provide))
}

func ErrAddCoinToAddress(addr sdk.Address) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("add coin to address %v failed", addr))
}

func ErrAddCoinToAccount(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("add coin to account %v failed", accKey))
}

func ErrMinusCoinToAccount(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("minus coin to account %v failed", accKey))
}

func ErrGetBankAddress(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v bank address failed", accKey))
}

func ErrGetOwnerKey(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v owner key failed", accKey))
}

func ErrGetPostKey(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v post key failed", accKey))
}

func ErrGetBankBalance(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v bank balance failed", accKey))
}

func ErrGetSequence(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get %v sequence failed", accKey))
}

func ErrIncreaseSequenceByOne(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("increase account %v sequence failed", accKey))
}

func ErrAddIncomeAndReward(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("add income and reward for user %v failed", accKey))
}

func ErrClaimReward(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("claim user %v reward failed", accKey))
}

func ErrGetStake(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("get user %v stake failed", accKey))
}

func ErrCheckUserTPSCapacity(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("update user %v transaction capacity failed", accKey))
}

func ErrAccountTPSCapacityNotEnough(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("user %v transaction capacity not enough, please wait", accKey))
}

func ErrAccountAlreadyExists(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("account %v exists", accKey))
}

func ErrBankAlreadyRegistered() sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("bank connection exists"))
}

func ErrRegisterFeeInsufficient() sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("register fee insufficient"))
}

func ErrAccountCreateFailed(accKey types.AccountKey) sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("create account %v failed", accKey))
}

func ErrUsernameAddressMismatch() sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("username and address mismatch"))
}

func ErrAccountCoinNotEnough() sdk.Error {
	return sdk.NewError(types.CodeAccountManagerFail, fmt.Sprintf("Account bank's coins are not enough"))
}
