package developer

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(DeveloperRegisterMsg{}, "lino/devRegister", nil)
	cdc.RegisterConcrete(DeveloperUpdateMsg{}, "lino/devUpdate", nil)
	cdc.RegisterConcrete(DeveloperRevokeMsg{}, "lino/devRevoke", nil)
	cdc.RegisterConcrete(GrantPermissionMsg{}, "lino/grantPermission", nil)
	cdc.RegisterConcrete(RevokePermissionMsg{}, "lino/revokePermission", nil)
	cdc.RegisterConcrete(PreAuthorizationMsg{}, "lino/preAuthorizationPermission", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	wire.RegisterCrypto(msgCdc)
}
