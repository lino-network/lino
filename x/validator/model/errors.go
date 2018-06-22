package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrGetValidator() sdk.Error {
	return types.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Get validator failed"))
}

func ErrSetValidatorList() sdk.Error {
	return types.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Set validator list failed"))
}

func ErrGetValidatorList() sdk.Error {
	return types.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Get validator list failed"))
}

func ErrValidatorMarshalError(err error) sdk.Error {
	return types.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Validator marshal error: %s", err.Error()))
}

func ErrValidatorUnmarshalError(err error) sdk.Error {
	return types.NewError(types.CodeValidatorManagerFailed, fmt.Sprintf("Validator unmarshal error: %s", err.Error()))
}
