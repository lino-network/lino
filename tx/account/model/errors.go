package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrAccountInfoDoesntExist() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("account info doesn't exist"))
}

func ErrGetAccountInfo() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("get account info failed"))
}

func ErrAccountBankDoesntExist() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("account bank doesn't exist"))
}

func ErrGetBankFromAccountKey() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("get bank from account key failed"))
}

func ErrGetBankFromAddress() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("get bank from address failed"))
}

func ErrAccountStorageInternal() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("account storage internal err"))
}

func ErrGetInfo() sdk.Error {
	return sdk.NewError(types.CodeInvalidMsg, fmt.Sprintf("account storage operation failed"))
}

func ErrInvalidLinoAmount() sdk.Error {
	return sdk.NewError(types.CodeInvalidMsg, fmt.Sprintf("Invalid Lino amount"))
}

func ErrUsernameNotFound() sdk.Error {
	return sdk.NewError(types.CodeUsernameNotFound, fmt.Sprintf("Username not found"))
}

func ErrInvalidUsername() sdk.Error {
	return sdk.NewError(types.CodeInvalidUsername, fmt.Sprintf("Invalida Username"))
}

func ErrUsernameAddressMismatch() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("Username and address mismatch"))
}

func ErrSetInfoFailed() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set info failed"))
}

func ErrSetBankFailed() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set bank failed"))
}

func ErrGetMetaFailed() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage get meta failed"))
}

func ErrSetMetaFailed() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set meta failed"))
}

func ErrGetRewardFailed() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage get reward failed"))
}

func ErrSetRewardFailed() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set reward failed"))
}

func ErrGetRelationshipFailed() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage get relationship failed"))
}

func ErrSetRelationshipFailed() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set relationship failed"))
}

func ErrGetPendingStakeFailed() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage get pending stake failed"))
}

func ErrSetPendingStakeFailed() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set pending stake failed"))
}

func ErrAddMoneyFailed() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage add money to bank failed"))
}

func ErrSetFollowerMeta() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set follower meta failed"))
}

func ErrSetFollowingMeta() sdk.Error {
	return sdk.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set following meta failed"))
}
