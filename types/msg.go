package types

// nolint
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/codec"
)

// Transactions messages must fulfill the Msg
type Msg interface {
	sdk.Msg
	GetPermission() Permission
	GetConsumeAmount() Coin
}

// Register the lino message type
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterInterface((*Msg)(nil), nil)
}
