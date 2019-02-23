package account

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
)

// RegisterWire - register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(RegisterMsg{}, "lino/register", nil)
	cdc.RegisterConcrete(FollowMsg{}, "lino/follow", nil)
	cdc.RegisterConcrete(UnfollowMsg{}, "lino/unfollow", nil)
	cdc.RegisterConcrete(TransferMsg{}, "lino/transfer", nil)
	cdc.RegisterConcrete(ClaimMsg{}, "lino/claim", nil)
	cdc.RegisterConcrete(RecoverMsg{}, "lino/recover", nil)
	cdc.RegisterConcrete(UpdateAccountMsg{}, "lino/updateAcc", nil)
}

var msgCdc = wire.New()

func init() {
	RegisterWire(msgCdc)
	wire.RegisterCrypto(msgCdc)
}
