package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// not found error
func ErrInfraProviderNotFound() sdk.Error {
	return types.NewError(types.CodeInfraProviderNotFound, fmt.Sprintf("infra provider is not found"))
}

func ErrInfraProviderListNotFound() sdk.Error {
	return types.NewError(types.CodeInfraProviderListNotFound, fmt.Sprintf("infra provider list is not found"))
}

// marshal error
func ErrFailedToMarshalInfraProvider(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalInfraProvider, fmt.Sprintf("failed to marshal infra provider: %s", err.Error()))
}

func ErrFailedToMarshalInfraProviderList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToMarshalInfraProviderList, fmt.Sprintf("failed to marshal infra provider list: %s", err.Error()))
}

// unmarshal error
func ErrFailedToUnmarshalInfraProvider(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalInfraProvider, fmt.Sprintf("failed to unmarshal infra provider: %s", err.Error()))
}

func ErrFailedToUnmarshalInfraProviderList(err error) sdk.Error {
	return types.NewError(types.CodeFailedToUnmarshalInfraProviderList, fmt.Sprintf("failed to unmarshal infra provider list: %s", err.Error()))
}
