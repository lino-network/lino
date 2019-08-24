package account

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterWire - register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(RegisterMsg{}, "lino/register", nil)
	cdc.RegisterConcrete(TransferMsg{}, "lino/transfer", nil)
	cdc.RegisterConcrete(RecoverMsg{}, "lino/recover", nil)
	cdc.RegisterConcrete(UpdateAccountMsg{}, "lino/updateAcc", nil)
}

var msgCdc = wire.New()

func init() {
	RegisterWire(msgCdc)
	sdk.RegisterCodec(msgCdc)
	wire.RegisterCrypto(msgCdc)
}
