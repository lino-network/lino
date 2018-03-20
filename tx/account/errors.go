package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case types.CodeInvalidUsername:
		return "Invalid username format"
	case types.CodeAccountManagerFail:
		return "Account manager internal error"
	case types.CodeUsernameNotFound:
		return "Username not found"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

// Error constructors
func ErrUsernameNotFound(msg string) sdk.Error {
	return newError(types.CodeUsernameNotFound, msg)
}

func ErrInvalidUsername(msg string) sdk.Error {
	return newError(types.CodeInvalidUsername, msg)
}

func ErrAccountManagerFail(msg string) sdk.Error {
	return newError(types.CodeAccountManagerFail, msg)
}

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(code, msg)
}
