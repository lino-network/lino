package validator

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case types.CodeInvalidUsername:
		return "Invalid username format"
	case types.CodeAccRegisterFailed:
		return "Validator register failed"
	case types.CodeValidatorHandlerFailed:
		return "Validator handler failed"
	case types.CodeValidatorManagerFailed:
		return "Validator manager failed"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

// Error constructors
func ErrSetValidator() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Set validator failed"))
}

func ErrGetValidator() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Get validator failed"))
}

func ErrSetValidatorList() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Set validator list failed"))
}

func ErrGetValidatorList() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Get validator list failed"))
}

func ErrSetValidatorUpdateList() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Set validator update list failed"))
}

func ErrGetValidatorUpdateList() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Get validator update list failed"))
}

func ErrValidatorMarshalError(err error) sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Validator marshal error: %s", err.Error()))
}

func ErrValidatorUnmarshalError(err error) sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Validator unmarshal error: %s", err.Error()))
}

func ErrAbsentValidatorNotCorrect() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("absent validator index out of range"))
}

func ErrAlreayInTheList() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Account has alreay in the list"))
}

func ErrNotInTheList() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Account not in the list"))
}

func ErrUsernameNotFound() sdk.Error {
	return newError(types.CodeUsernameNotFound, fmt.Sprintf("Username not found"))
}

func ErrDepositNotAvailable() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Deposit not available"))
}

func ErrNoDeposit() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("No Deposit"))
}

func ErrRegisterFeeNotEnough() sdk.Error {
	return newError(types.CodeUsernameNotFound, fmt.Sprintf("Register fee not enough"))
}

func ErrInvalidUsername() sdk.Error {
	return newError(types.CodeInvalidUsername, fmt.Sprintf("Invalida Username"))
}

func ErrAccountCoinNotEnough() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Account bank's coins are not enough"))
}

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	} else {
		return codeToDefaultMsg(code)
	}
}

func newError(code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(code, msg)
}
