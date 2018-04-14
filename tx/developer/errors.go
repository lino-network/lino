package developer

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
	case types.CodeDeveloperHandlerFailed:
		return "Developer handler failed"
	case types.CodeDeveloperManagerFailed:
		return "Developer manager failed"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

// Error constructors
func ErrDeveloperNotFound() sdk.Error {
	return newError(types.CodeUsernameNotFound, fmt.Sprintf("Developer not found"))
}

func ErrUsernameNotFound() sdk.Error {
	return newError(types.CodeUsernameNotFound, fmt.Sprintf("Username not found"))
}

func ErrDeveloperDepositNotEnough() sdk.Error {
	return newError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Developer deposit not enough"))
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
