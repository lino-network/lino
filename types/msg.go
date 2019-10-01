package types

// nolint
import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Transactions messages must fulfill the Msg
type Msg interface {
	sdk.Msg
	GetPermission() Permission
	GetConsumeAmount() Coin
}

type AddrMsg interface {
	sdk.Msg
	GetAccOrAddrSigners() []AccOrAddr
}

// Register the lino message type
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterInterface((*AddrMsg)(nil), nil)
}
