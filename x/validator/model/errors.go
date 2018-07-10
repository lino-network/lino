package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// not found
func ErrValidatorNotFound() sdk.Error {
	return types.NewError(types.CodeValidatorNotFound, fmt.Sprintf("validator is not found"))
}

func ErrValidatorListNotFound() sdk.Error {
	return types.NewError(types.CodeValidatorListNotFound, fmt.Sprintf("validator list is not found"))
}

// marshal error
func ErrFailedToMarshalValidator(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalValidator, fmt.Sprintf("failed to marshal validator: %s", err.Error()))
}

func ErrFailedToMarshalValidatorList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalValidatorList, fmt.Sprintf("failed to marshal validator list: %s", err.Error()))
}

// unmarshal error
func ErrFailedToUnmarshalValidator(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalValidator, fmt.Sprintf("failed to unmarshal validator: %s", err.Error()))
}

func ErrFailedToUnmarshalValidatorList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalValidatorList, fmt.Sprintf("failed to unmarshal validator list: %s", err.Error()))
}
