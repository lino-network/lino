package event

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

type CodeType = sdk.CodeType

func codeToDefaultMsg(code CodeType) string {
	switch code {
	case types.CodeEventExecuteError:
		return "Event execute error"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

func ErrEventExecuteError(key []byte) sdk.Error {
	return newError(types.CodeEventExecuteError, fmt.Sprintf("Event execute failed"))
}

func ErrWrongEventType() sdk.Error {
	return newError(types.CodeEventExecuteError, fmt.Sprintf("Wrong event type"))
}

func ErrEventNotFound(key []byte) sdk.Error {
	return newError(types.CodeEventExecuteError, fmt.Sprintf("Event not found for key: %s", key))
}

func ErrEventMarshalError(err error) sdk.Error {
	return newError(types.CodeEventExecuteError, fmt.Sprintf("Event marshal error: %s", err.Error()))
}

func ErrEventUnmarshalError(err error) sdk.Error {
	return newError(types.CodeEventExecuteError, fmt.Sprintf("Event unmarshal error: %s", err.Error()))
}

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
