package auth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// ErrIncorrectStdTxType - error if parse std tx failed
func ErrIncorrectStdTxType() sdk.Error {
	return types.NewError(types.CodeIncorrectStdTxType, fmt.Sprint("incorrect stdTx type"))
}

// ErrNoSignatures - error if transaction without signatures
func ErrNoSignatures() sdk.Error {
	return types.NewError(types.CodeNoSignatures, fmt.Sprint("no signatures"))
}

// ErrUnknownMsgType - error if msg can't be recognized
func ErrUnknownMsgType() sdk.Error {
	return types.NewError(types.CodeUnknownMsgType, fmt.Sprint("unknown msg type"))
}

// ErrWrongNumberOfSigners - error if number of signers and signatures mismatch
func ErrWrongNumberOfSigners() sdk.Error {
	return types.NewError(types.CodeWrongNumberOfSigners, fmt.Sprint("the number of siners is wrong"))
}

// XXX(yumin): no longer available, due to upgrade-1.
// ErrInvalidSequence - error if sequence number mismatch
// func ErrInvalidSequence(msg string) sdk.Error {
// 	return types.NewError(types.CodeInvalidSequence, fmt.Sprintf("msg: %v", msg))
// }

// ErrUnverifiedBytes - error if signbyte verification failed
func ErrUnverifiedBytes(msg string) sdk.Error {
	return types.NewError(types.CodeUnverifiedBytes, fmt.Sprintf("msg: %v", msg))
}

// ErrMsgFeeNotEnough - error if the provided message fee is not enough
func ErrMsgFeeNotEnough() sdk.Error {
	return types.NewError(types.CodeMsgFeeNotEnough, fmt.Sprint("message fee is not enough"))
}
