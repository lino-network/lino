package types

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(StakeInMsg{}, "lino/stakeIn", nil)
	cdc.RegisterConcrete(StakeOutMsg{}, "lino/stakeOut", nil)
	cdc.RegisterConcrete(DelegateMsg{}, "lino/delegate", nil)
	cdc.RegisterConcrete(DelegatorWithdrawMsg{}, "lino/delegateWithdraw", nil)
	cdc.RegisterConcrete(ClaimInterestMsg{}, "lino/claimInterest", nil)
}

var msgCdc = wire.New()

func init() {
	RegisterWire(msgCdc)
}
