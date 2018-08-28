package infra

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// RegisterWire - register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(ProviderReportMsg{}, "lino/providerReport", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
