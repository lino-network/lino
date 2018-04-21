package register

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrInvalidUsername(msg string) sdk.Error {
	return sdk.NewError(types.CodeInvalidUsername, msg)
}

func ErrAccRegisterFail(msg string) sdk.Error {
	return sdk.NewError(types.CodeAccRegisterFailed, msg)
}
