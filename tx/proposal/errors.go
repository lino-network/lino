package proposal

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func ErrUsernameNotFound() sdk.Error {
	return sdk.NewError(types.CodeProposalManagerError, fmt.Sprintf("Username not found"))
}

func ErrInvalidUsername() sdk.Error {
	return sdk.NewError(types.CodeProposalManagerError, fmt.Sprintf("Invalid Username"))
}

func ErrIllegalParameter() sdk.Error {
	return sdk.NewError(types.CodeProposalManagerError, fmt.Sprintf("Invalid parameter"))
}
