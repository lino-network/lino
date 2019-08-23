package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewError - create an error
func NewError(code sdk.CodeType, msg string) sdk.Error {
	return sdk.NewError(LinoErrorCodeSpace, code, msg)
}

// ErrInvalidCoins - error if convert LNO to Coin failed
func ErrInvalidCoins(msg string) sdk.Error {
	return NewError(CodeInvalidCoins, msg)
}

// ErrAmountOverflow - error if coin amount int64 overflow
func ErrAmountOverflow() sdk.Error {
	return NewError(CodeInvalidInt64Number, "coin amount can't be represented as an int64")
}

// ErrInvalidQueryPath - error if query path length is incorrect or content is invalid
func ErrInvalidQueryPath() sdk.Error {
	return NewError(CodeInvalidQueryPath, "query path is invalid")
}

// ErrInvalidIDAAmount - error if the IDA amount is invalid.
func ErrInvalidIDAAmount(msg string) sdk.Error {
	return NewError(CodeInvalidIDAAmount, msg)
}
