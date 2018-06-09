package validator

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrAbsentValidatorNotCorrect() sdk.Error {
	return sdk.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Absent validator index out of range"))
}

func ErrGetPubKeyFailed() sdk.Error {
	return sdk.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Get ABCI public key failed"))
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

func ErrNoCoinToWithdraw() sdk.Error {
	return sdk.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("No coin to withdraw"))
}

func ErrCommitingDepositNotEnough() sdk.Error {
	return sdk.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Commiting deposit not enough"))
}

func ErrVotingDepositNotEnough() sdk.Error {
	return sdk.NewError(types.CodeValidatorHandlerFailed, fmt.Sprintf("Voting Deposit fee not enough"))
}

func ErrCommitingDepositExceedVotingDeposit() sdk.Error {
	return sdk.NewError(types.CodeValidatorHandlerFailed, fmt.Sprintf("Commiting deposit exceed voting deposit"))
}

func ErrInvalidUsername() sdk.Error {
	return sdk.NewError(types.CodeInvalidUsername, fmt.Sprintf("Invalida Username"))
}

func ErrPubKeyHasBeenRegistered() sdk.Error {
	return sdk.NewError(types.CodeValidatorHandlerFailed, fmt.Sprintf("Public key has been registered"))
}
