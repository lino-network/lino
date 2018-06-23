package developer

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(DeveloperRegisterMsg{}, "lino/devRegister", nil)
	cdc.RegisterConcrete(DeveloperRevokeMsg{}, "lino/devRevoke", nil)
	cdc.RegisterConcrete(GrantDeveloperMsg{}, "lino/grantDeveloper", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
