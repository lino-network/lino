package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// ErrInfraProviderNotFound - error if infra provider is not found
func ErrInfraProviderNotFound() sdk.Error {
	return types.NewError(types.CodeInfraProviderNotFound, fmt.Sprintf("infra provider is not found"))
}

// ErrInfraProviderListNotFound - error if infra provider list is not found
func ErrInfraProviderListNotFound() sdk.Error {
	return types.NewError(types.CodeInfraProviderListNotFound, fmt.Sprintf("infra provider list is not found"))
}

// ErrFailedToMarshalInfraProvider - error if marshal infra provider failed
func ErrFailedToMarshalInfraProvider(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalInfraProvider, fmt.Sprintf("failed to marshal infra provider: %s", err.Error()))
}

// ErrFailedToMarshalInfraProviderList - error if marshal infra provider list failed
func ErrFailedToMarshalInfraProviderList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalInfraProviderList, fmt.Sprintf("failed to marshal infra provider list: %s", err.Error()))
}

// ErrFailedToUnmarshalInfraProvider - error if unmarshal infra provider failed
func ErrFailedToUnmarshalInfraProvider(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalInfraProvider, fmt.Sprintf("failed to unmarshal infra provider: %s", err.Error()))
}

// ErrFailedToUnmarshalInfraProviderList - error if unmarshal infra provider list failed
func ErrFailedToUnmarshalInfraProviderList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalInfraProviderList, fmt.Sprintf("failed to unmarshal infra provider list: %s", err.Error()))
}
