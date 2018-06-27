package vote

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(VoterDepositMsg{}, "lino/voteDeposit", nil)
	cdc.RegisterConcrete(VoterRevokeMsg{}, "lino/voteRevoke", nil)
	cdc.RegisterConcrete(VoterWithdrawMsg{}, "lino/voteWithdraw", nil)
	cdc.RegisterConcrete(DelegateMsg{}, "lino/delegate", nil)
	cdc.RegisterConcrete(DelegatorWithdrawMsg{}, "lino/delegateWithdraw", nil)
	cdc.RegisterConcrete(RevokeDelegationMsg{}, "lino/delegateRevoke", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
