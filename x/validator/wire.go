package validator

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(ValidatorDepositMsg{}, "lino/valDeposit", nil)
	cdc.RegisterConcrete(ValidatorWithdrawMsg{}, "lino/valWithdraw", nil)
	cdc.RegisterConcrete(ValidatorRevokeMsg{}, "lino/valRevoke", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	wire.RegisterCrypto(msgCdc)
}
