package infra

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// ErrProviderNotFound - error if infra provider is not found
func ErrProviderNotFound() sdk.Error {
	return types.NewError(types.CodeInfraProviderNotFound, fmt.Sprintf("provider is not found"))
}

// ErrInvalidUsername - error if username is invalid
func ErrInvalidUsername() sdk.Error {
	return types.NewError(types.CodeInvalidUsername, fmt.Sprintf("invalid Username"))
}

// ErrInvalidUsage - error if report usgae is invalid
func ErrInvalidUsage() sdk.Error {
	return types.NewError(types.CodeInvalidUsage, fmt.Sprintf("invalid Usage"))
}

// ErrQueryFailed - error when query infra store failed
func ErrQueryFailed() sdk.Error {
	return types.NewError(types.CodeInfraQueryFailed, fmt.Sprintf("query infra store failed"))
}
