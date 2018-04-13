package infra

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
	case types.CodeInfraProviderHandlerFailed:
		return "Infra provider handler failed"
	case types.CodeInfraProviderManagerFailed:
		return "Infra manager failed"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

// Error constructors
func ErrProviderNotFound() sdk.Error {
	return newError(types.CodeUsernameNotFound, fmt.Sprintf("Provider not found"))
}

func ErrInvalidUsername() sdk.Error {
	return newError(types.CodeInvalidUsername, fmt.Sprintf("Invalida Username"))
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
