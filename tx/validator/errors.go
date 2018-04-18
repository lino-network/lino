package validator

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrSetValidator() sdk.Error {
	return sdk.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Set validator failed"))
}

func ErrAbsentValidatorNotCorrect() sdk.Error {
	return sdk.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("absent validator index out of range"))
}

func ErrAlreayInTheList() sdk.Error {
	return sdk.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Account has alreay in the list"))
}

func ErrNotInTheList() sdk.Error {
	return sdk.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Account not in the list"))
}

func ErrUsernameNotFound() sdk.Error {
	return sdk.NewError(types.CodeUsernameNotFound, fmt.Sprintf("Username not found"))
}

func ErrIllegalWithdraw() sdk.Error {
	return sdk.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Illegal withdraw"))
}

func ErrCommitingDepositNotEnough() sdk.Error {
	return sdk.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Commiting deposit not enough"))
}

func ErrVotingDepositNotEnough() sdk.Error {
	return sdk.NewError(types.CodeValidatorHandlerFailed, fmt.Sprintf("Voting Deposit fee not enough"))
}

func ErrInvalidUsername() sdk.Error {
	return sdk.NewError(types.CodeInvalidUsername, fmt.Sprintf("Invalida Username"))
}

func ErrAccountCoinNotEnough() sdk.Error {
	return sdk.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Account bank's coins are not enough"))
}
