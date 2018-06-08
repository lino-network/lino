package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// // Error constructors
func ErrGetInfraProvider() sdk.Error {
	return sdk.NewError(types.CodeInfraProviderManagerFailed, fmt.Sprintf("Get infra provider failed"))
}

func ErrSetInfraProviderList() sdk.Error {
	return sdk.NewError(types.CodeInfraProviderManagerFailed, fmt.Sprintf("Set infra provider list failed"))
}

func ErrGetInfraProviderList() sdk.Error {
	return sdk.NewError(types.CodeInfraProviderManagerFailed, fmt.Sprintf("Get infra provider list failed"))
}

func ErrInfraProviderMarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeInfraProviderManagerFailed, fmt.Sprintf("Infra provider marshal error: %s", err.Error()))
}

func ErrInfraProviderUnmarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeInfraProviderManagerFailed, fmt.Sprintf("Infra provider unmarshal error: %s", err.Error()))
}
