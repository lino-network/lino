package app

import (
	"fmt"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Error constructors
func ErrGenesisFailed(msg string) sdk.Error {
	return types.NewError(types.CodeGenesisFailed, fmt.Sprintf("genesis failed: %s", msg))
}
