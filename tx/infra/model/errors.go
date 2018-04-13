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
		return "InfraProvider register failed"
	case types.CodeInfraProviderHandlerFailed:
		return "InfraProvider handler failed"
	case types.CodeInfraProviderManagerFailed:
		return "InfraProvider manager failed"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

// // Error constructors
func ErrGetInfraProvider() sdk.Error {
	return newError(types.CodeInfraProviderManagerFailed, fmt.Sprintf("Get infra provider failed"))
}

func ErrSetInfraProviderList() sdk.Error {
	return newError(types.CodeInfraProviderManagerFailed, fmt.Sprintf("Set infra provider list failed"))
}

func ErrGetInfraProviderList() sdk.Error {
	return newError(types.CodeInfraProviderManagerFailed, fmt.Sprintf("Get infra provider list failed"))
}

func ErrInfraProviderMarshalError(err error) sdk.Error {
	return newError(types.CodeInfraProviderManagerFailed, fmt.Sprintf("Infra provider marshal error: %s", err.Error()))
}

func ErrInfraProviderUnmarshalError(err error) sdk.Error {
	return newError(types.CodeInfraProviderManagerFailed, fmt.Sprintf("Infra provider unmarshal error: %s", err.Error()))
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
