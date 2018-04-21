package model

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// // Error constructors
func ErrGetDeveloper() sdk.Error {
	return sdk.NewError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Get developer failed"))
}

func ErrSetDeveloperList() sdk.Error {
	return sdk.NewError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Set developer list failed"))
}

func ErrGetDeveloperList() sdk.Error {
	return sdk.NewError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Get developer list failed"))
}

func ErrDeveloperMarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Developer marshal error: %s", err.Error()))
}

func ErrDeveloperUnmarshalError(err error) sdk.Error {
	return sdk.NewError(types.CodeDeveloperManagerFailed, fmt.Sprintf("Developer unmarshal error: %s", err.Error()))
}
