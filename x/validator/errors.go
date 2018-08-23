package validator

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrAccountNotFound() sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("username is not found"))
}

func ErrInsufficientDeposit() sdk.Error {
	return types.NewError(types.CodeInsufficientDeposit, fmt.Sprintf("voting deposit fee not enough"))
}

func ErrUnbalancedAccount() sdk.Error {
	return types.NewError(types.CodeUnbalancedAccount, fmt.Sprintf("committing deposit not enough"))
}

func ErrIllegalWithdraw() sdk.Error {
	return types.NewError(types.CodeIllegalWithdraw, fmt.Sprintf("illegal withdraw"))
}

func ErrInvalidCoin() sdk.Error {
	return types.NewError(types.CodeInvalidCoin, fmt.Sprintf("no coin to withdraw"))
}

func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("Invalida Username"))
}

func ErrInvalidWebsite() sdk.Error {
	return types.NewError(types.CodeInvalidWebsite, fmt.Sprintf("Invalida website"))
}

func ErrValidatorPubKeyAlreadyExist() sdk.Error {
	return types.NewError(types.CodeValidatorPubKeyAlreadyExist, fmt.Sprintf("validator public key has been registered"))
}
