package vote

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(StakeInMsg{}, "lino/stakeIn", nil)
	cdc.RegisterConcrete(RevokeStakeMsg{}, "lino/revokeStake", nil)
	cdc.RegisterConcrete(StakeOutMsg{}, "lino/stakeOut", nil)
	cdc.RegisterConcrete(DelegateMsg{}, "lino/delegate", nil)
	cdc.RegisterConcrete(DelegatorWithdrawMsg{}, "lino/delegateWithdraw", nil)
	cdc.RegisterConcrete(RevokeDelegationMsg{}, "lino/delegateRevoke", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
