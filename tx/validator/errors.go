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

func ErrIllegalWithdraw() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Illegal withdraw"))
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
