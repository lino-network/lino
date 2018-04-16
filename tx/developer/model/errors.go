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
	case types.CodeDeveloperHandlerFailed:
		return "Developer handler failed"
	case types.CodeDeveloperManagerFailed:
		return "Developer manager failed"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

// // Error constructors
func ErrGetDeveloper() sdk.Error {
	return newError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Get developer failed"))
}

func ErrSetDeveloperList() sdk.Error {
	return newError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Set developer list failed"))
}

func ErrGetDeveloperList() sdk.Error {
	return newError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Get developer list failed"))
}

func ErrDeveloperMarshalError(err error) sdk.Error {
	return newError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Developer marshal error: %s", err.Error()))
}

func ErrDeveloperUnmarshalError(err error) sdk.Error {
	return newError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Developer unmarshal error: %s", err.Error()))
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
