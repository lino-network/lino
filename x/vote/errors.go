package vote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrAccountNotFound() sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("account is not found"))
}

func ErrIllegalWithdraw() sdk.Error {
	return types.NewError(types.CodeIllegalWithdraw, fmt.Sprintf("illegal withdraw"))
}

func ErrValidatorCannotRevoke() sdk.Error {
	return types.NewError(types.CodeValidatorCannotRevoke, fmt.Sprintf("invalid revoke"))
}

func ErrVoteAlreadyExist() sdk.Error {
	return types.NewError(types.CodeVoteAlreadyExist, fmt.Sprintf("Vote exist"))
}

func ErrInvalidCoin() sdk.Error {
	return types.NewError(types.CodeInvalidCoin, fmt.Sprintf("can't withdraw 0 coin"))
}

func ErrInsufficientDeposit() sdk.Error {
	return types.NewError(types.CodeInsufficientDeposit, fmt.Sprintf("deposit is not enough"))
}

func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("invalid username"))
}
