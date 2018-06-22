package types

// nolint
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Transactions messages must fulfill the Msg
type Msg interface {
	sdk.Msg
	GetPermission() Permission
}
