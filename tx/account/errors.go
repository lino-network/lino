//nolint
package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type CodeType = sdk.CodeType

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code CodeType) string {
	switch code {
	case types.CodeInvalidUsername:
		return "Invalid username format"
	case types.CodeAccountManagerFail:
		return "Account manager internal error"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

//----------------------------------------
// Error constructors

func ErrInvalidUsername(msg string) sdk.Error {
	return newError(types.CodeInvalidUsername, msg)
}

func ErrCodeAccountManagerFail(msg string) sdk.Error {
	return newError(types.CodeAccountManagerFail, msg)
}

//----------------------------------------

func msgOrDefaultMsg(msg string, code CodeType) string {
	if msg != "" {
		return msg
	} else {
		return codeToDefaultMsg(code)
	}
}

func newError(code CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(code, msg)
}
