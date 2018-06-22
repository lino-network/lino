package infra

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Error constructors
func ErrProviderNotFound() sdk.Error {
	return types.NewError(types.CodeUsernameNotFound, fmt.Sprintf("Provider not found"))
}

func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("Invalid Username"))
}

func ErrInvalidUsage() sdk.Error {
	return types.NewError(types.CodeInfraInvalidMsg, fmt.Sprintf("Invalid Usage"))
}
