package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
)

// ErrInvalidMsgQuota - error when message fee is not valid
func ErrInvalidMsgQuota() sdk.Error {
	return linotypes.NewError(linotypes.CodeInvalidMsgQuota, fmt.Sprintf("invalid message quota"))
}
