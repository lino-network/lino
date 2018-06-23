package model

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Error constructors
func ErrAccountInfoNotFound() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("account info is not found"))
}

func ErrGetAccountInfo() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("get account info failed"))
}

func ErrAccountBankNotFound() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("account bank is not found"))
}

func ErrGetBankFromAccountKey() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("get bank from account key failed"))
}

func ErrGetBankFromAddress() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("get bank from address failed"))
}

func ErrAccountStorageInternal() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("account storage internal err"))
}

func ErrGetInfo() sdk.Error {
	return types.NewError(types.CodeInvalidMsg, fmt.Sprintf("account storage operation failed"))
}

func ErrInvalidLinoAmount() sdk.Error {
	return types.NewError(types.CodeInvalidMsg, fmt.Sprintf("Invalid Lino amount"))
}

func ErrUsernameNotFound() sdk.Error {
	return types.NewError(types.CodeUsernameNotFound, fmt.Sprintf("Username not found"))
}

func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("Invalida Username"))
}

func ErrUsernameAddressMismatch() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("Username and address mismatch"))
}

func ErrSetInfoFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set info failed"))
}

func ErrSetBankFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set bank failed"))
}

func ErrGetMetaFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage get meta failed"))
}

func ErrSetMetaFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set meta failed"))
}

func ErrGetRewardFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage get reward failed"))
}

func ErrSetRewardFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set reward failed"))
}

func ErrGetRelationshipFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage get relationship failed"))
}

func ErrSetRelationshipFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set relationship failed"))
}

func ErrSetBalanceHistoryFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set balance history failed"))
}

func ErrGetBalanceHistoryFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage get balance history failed"))
}

func ErrGetPendingStakeFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage get pending stake failed"))
}

func ErrSetPendingStakeFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set pending stake failed"))
}

func ErrGetGrantUserFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage get grant user failed"))
}

func ErrSetGrantUserFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set grant user failed"))
}

func ErrGetGrantListFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage get grant key list failed"))
}

func ErrSetGrantListFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set grant key list failed"))
}

func ErrAddMoneyFailed() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage add money to bank failed"))
}

func ErrSetFollowerMeta() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set follower meta failed"))
}

func ErrSetFollowingMeta() sdk.Error {
	return types.NewError(types.CodeAccountStorageFail, fmt.Sprintf("AccountStorage set following meta failed"))
}
