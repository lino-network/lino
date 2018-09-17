package vote

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(StakeInMsg{}, "lino/stakeIn", nil)
	cdc.RegisterConcrete(StakeOutMsg{}, "lino/stakeOut", nil)
	cdc.RegisterConcrete(DelegateMsg{}, "lino/delegate", nil)
	cdc.RegisterConcrete(DelegatorWithdrawMsg{}, "lino/delegateWithdraw", nil)
	cdc.RegisterConcrete(ClaimInterestMsg{}, "lino/claimInterest", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
