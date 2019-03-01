package vote

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// ErrAccountNotFound - error if account is not found
func ErrAccountNotFound() sdk.Error {
	return types.NewError(types.CodeAccountNotFound, fmt.Sprintf("account is not found"))
}

// ErrIllegalWithdraw - error if withdraw is illegal
func ErrIllegalWithdraw() sdk.Error {
	return types.NewError(types.CodeIllegalWithdraw, fmt.Sprintf("illegal withdraw"))
}

// ErrValidatorCannotRevoke - error if voter is validator
func ErrValidatorCannotRevoke() sdk.Error {
	return types.NewError(types.CodeValidatorCannotRevoke, fmt.Sprintf("invalid revoke"))
}

// ErrVoteAlreadyExist - error if user already vote for a proposal
func ErrVoteAlreadyExist() sdk.Error {
	return types.NewError(types.CodeVoteAlreadyExist, fmt.Sprintf("Vote exist"))
}

// ErrVoteNotFound - error if voter is not found
func ErrVoterNotFound() sdk.Error {
	return types.NewError(types.CodeVoterNotFound, fmt.Sprintf("voter not found"))
}

// ErrInvalidCoin - error if coin is invalid
func ErrInvalidCoin() sdk.Error {
	return types.NewError(types.CodeInvalidCoin, fmt.Sprintf("can't withdraw 0 coin"))
}

// ErrInsufficientDeposit - error if voter deposit is insufficient
func ErrInsufficientDeposit() sdk.Error {
	return types.NewError(types.CodeInsufficientDeposit, fmt.Sprintf("deposit is not enough"))
}

// ErrInvalidUsername - error if username is invalid
func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("invalid username"))
}

// ErrQueryFailed - error when query vote store failed
func ErrQueryFailed() sdk.Error {
	return types.NewError(types.CodeVoteQueryFailed, fmt.Sprintf("query vote store failed"))
}
