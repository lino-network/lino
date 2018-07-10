package infra

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrProviderNotFound() sdk.Error {
	return types.NewError(types.CodeInfraProviderNotFound, fmt.Sprintf("provider is not found"))
}

func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("invalid Username"))
}

func ErrInvalidUsage() sdk.Error {
	return types.NewError(types.CodeInvalidUsage, fmt.Sprintf("invalid Usage"))
}
