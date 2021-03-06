package types

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

// ErrNotAVoterOrHasDuty
func ErrNotAVoterOrHasDuty() sdk.Error {
	return types.NewError(types.CodeNotAVoterOrHasDuty, fmt.Sprintf("not a voter or has duty"))
}

// ErrNoDuty
func ErrNoDuty() sdk.Error {
	return types.NewError(types.CodeNoDuty, fmt.Sprintf("voter doesn't have duty"))
}

// ErrFrozenAmountIsNotEmpty
func ErrFrozenAmountIsNotEmpty() sdk.Error {
	return types.NewError(types.CodeFrozenAmountIsNotEmpty, fmt.Sprintf("frozen money is not empty"))
}

// ErrNegativeFrozenAmount -
func ErrNegativeFrozenAmount() sdk.Error {
	return types.NewError(
		types.CodeNegativeFrozenAmount, fmt.Sprintf("fronzen amount is negative"))
}

// ErrInsufficientStake
func ErrInsufficientStake() sdk.Error {
	return types.NewError(types.CodeInsufficientStake, fmt.Sprintf("stake is not enough"))
}

// ErrStakeStatNotFound -
func ErrStakeStatNotFound(day int64) sdk.Error {
	return types.NewError(
		types.CodeStakeStatNotFound, fmt.Sprintf("stake stats not found: %d", day))
}
