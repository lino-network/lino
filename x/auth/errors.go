package auth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrIncorrectStdTxType() sdk.Error {
	return types.NewError(types.CodeIncorrectStdTxType, fmt.Sprint("incorrect stdTx type"))
}

func ErrNoSignatures() sdk.Error {
	return types.NewError(types.CodeNoSignatures, fmt.Sprint("no signatures"))
}

func ErrUnknownMsgType() sdk.Error {
	return types.NewError(types.CodeUnknownMsgType, fmt.Sprint("unknown msg type"))
}

func ErrWrongNumberOfSigners() sdk.Error {
	return types.NewError(types.CodeWrongNumberOfSigners, fmt.Sprint("the number of siners is wrong"))
}

func ErrInvalidSequence(msg string) sdk.Error {
	return types.NewError(types.CodeInvalidSequence, fmt.Sprintf("msg: %v", msg))
}

func ErrUnverifiedBytes(msg string) sdk.Error {
	return types.NewError(types.CodeUnverifiedBytes, fmt.Sprintf("msg: %v", msg))
}
