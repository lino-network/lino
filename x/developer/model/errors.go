package model

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// // Error constructors
func ErrGetDeveloper() sdk.Error {
	return types.NewError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Get developer failed"))
}

func ErrSetDeveloperList() sdk.Error {
	return types.NewError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Set developer list failed"))
}

func ErrGetDeveloperList() sdk.Error {
	return types.NewError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Get developer list failed"))
}

func ErrDeveloperMarshalError(err error) sdk.Error {
	return types.NewError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Developer marshal error: %s", err.Error()))
}

func ErrDeveloperUnmarshalError(err error) sdk.Error {
	return types.NewError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Developer unmarshal error: %s", err.Error()))
}
