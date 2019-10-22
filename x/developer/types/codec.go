package types

import (
	wire "github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(DeveloperRegisterMsg{}, "lino/devRegister", nil)
	cdc.RegisterConcrete(DeveloperUpdateMsg{}, "lino/devUpdate", nil)
	cdc.RegisterConcrete(DeveloperRevokeMsg{}, "lino/devRevoke", nil)
	cdc.RegisterConcrete(IDAIssueMsg{}, "lino/IDAIssue", nil)
	cdc.RegisterConcrete(IDAMintMsg{}, "lino/IDAMint", nil)
	cdc.RegisterConcrete(IDATransferMsg{}, "lino/IDATransfer", nil)
	cdc.RegisterConcrete(IDAAuthorizeMsg{}, "lino/IDAAuthorize", nil)
	cdc.RegisterConcrete(UpdateAffiliatedMsg{}, "lino/UpdateAffiliated", nil)
}

var ModuleCdc = wire.New()

func init() {
	RegisterWire(ModuleCdc)
	wire.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
