package model

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

func ErrGetValidator() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Get validator failed"))
}

func ErrSetValidatorList() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Set validator list failed"))
}

func ErrGetValidatorList() sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Get validator list failed"))
}

func ErrValidatorMarshalError(err error) sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Validator marshal error: %s", err.Error()))
}

func ErrValidatorUnmarshalError(err error) sdk.Error {
	return newError(types.CodeValidatorManagerFailed, fmt.Sprintf("Validator unmarshal error: %s", err.Error()))
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
