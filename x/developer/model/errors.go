package model

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// not found error
func ErrDeveloperNotFound() sdk.Error {
	return types.NewError(types.CodeDeveloperNotFound, fmt.Sprintf("developer is not found"))
}

func ErrDeveloperListNotFound() sdk.Error {
	return types.NewError(types.CodeDeveloperListNotFound, fmt.Sprintf("developer list is not found"))
}

// marshal error
func ErrFailedToMarshalDeveloper(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalDeveloper, fmt.Sprintf("failed to marshal developer: %s", err.Error()))
}

func ErrFailedToMarshalDeveloperList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalDeveloperList, fmt.Sprintf("failed to marshal developer list: %s", err.Error()))
}

// unmarshal error
func ErrFailedToUnmarshalDeveloper(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalDeveloper, fmt.Sprintf("failed to unmarshal developer: %s", err.Error()))
}

func ErrFailedToUnmarshalDeveloperList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalDeveloperList, fmt.Sprintf("failed to unmarshal developer list: %s", err.Error()))
}
