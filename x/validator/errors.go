package validator

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// ErrAccountNotFound - error when account is not found
func ErrAccountNotFound() sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("username is not found"))
}

// ErrInsufficientDeposit - error when deposit is insufficient
func ErrInsufficientDeposit() sdk.Error {
	return types.NewError(types.CodeInsufficientDeposit, fmt.Sprintf("voting deposit fee not enough"))
}

// ErrInsufficientDeposit - error if required voting deposit less than committing deposit
func ErrUnbalancedAccount() sdk.Error {
	return types.NewError(types.CodeUnbalancedAccount, fmt.Sprintf("committing deposit not enough"))
}

// ErrIllegalWithdraw - error if withdraw less than minimum withdraw requirement
func ErrIllegalWithdraw() sdk.Error {
	return types.NewError(types.CodeIllegalWithdraw, fmt.Sprintf("illegal withdraw"))
}

// ErrInvalidCoin - error if coin in msg is invalid
func ErrInvalidCoin() sdk.Error {
	return types.NewError(types.CodeInvalidCoin, fmt.Sprintf("no coin to withdraw"))
}

// ErrInvalidUsername - error if username is invalid
func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("Invalida Username"))
}

// ErrInvalidWebsite - error if website is invalid
func ErrInvalidWebsite() sdk.Error {
	return types.NewError(types.CodeInvalidWebsite, fmt.Sprintf("Invalida website"))
}

// ErrValidatorPubKeyAlreadyExist - error if validator public key is already exist
func ErrValidatorPubKeyAlreadyExist() sdk.Error {
	return types.NewError(types.CodeValidatorPubKeyAlreadyExist, fmt.Sprintf("validator public key has been registered"))
}
