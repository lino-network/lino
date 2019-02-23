package infra

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
)

// RegisterWire - register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(ProviderReportMsg{}, "lino/providerReport", nil)
}

var msgCdc = wire.New()

func init() {
	RegisterWire(msgCdc)
}
