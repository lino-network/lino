package model

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ErrDeveloperNotFound - error if developer is not found in KVStore
func ErrDeveloperNotFound() sdk.Error {
	return types.NewError(types.CodeDeveloperNotFound, fmt.Sprintf("developer is not found"))
}

// ErrDeveloperListNotFound - error if developer list is not found in KVStore
func ErrDeveloperListNotFound() sdk.Error {
	return types.NewError(types.CodeDeveloperListNotFound, fmt.Sprintf("developer list is not found"))
}

// ErrFailedToMarshalDeveloper - error if marshal developer failed
func ErrFailedToMarshalDeveloper(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalDeveloper, fmt.Sprintf("failed to marshal developer: %s", err.Error()))
}

// ErrFailedToMarshalDeveloperList - error if marshal developer list failed
func ErrFailedToMarshalDeveloperList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalDeveloperList, fmt.Sprintf("failed to marshal developer list: %s", err.Error()))
}

// ErrFailedToUnmarshalDeveloper - error if unmarshal developer failed
func ErrFailedToUnmarshalDeveloper(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalDeveloper, fmt.Sprintf("failed to unmarshal developer: %s", err.Error()))
}

// ErrFailedToUnmarshalDeveloperList - error if unmarshal developer list failed
func ErrFailedToUnmarshalDeveloperList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalDeveloperList, fmt.Sprintf("failed to unmarshal developer list: %s", err.Error()))
}
