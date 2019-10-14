package types

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(StakeInMsg{}, "lino/stakeIn", nil)
	cdc.RegisterConcrete(StakeOutMsg{}, "lino/stakeOut", nil)
	cdc.RegisterConcrete(ClaimInterestMsg{}, "lino/claimInterest", nil)
	cdc.RegisterConcrete(StakeInForMsg{}, "lino/stakeInFor", nil)
}

var msgCdc = wire.New()

func init() {
	RegisterWire(msgCdc)
}
