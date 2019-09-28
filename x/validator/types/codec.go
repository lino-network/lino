package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(ValidatorRegisterMsg{}, "lino/valRegister", nil)
	cdc.RegisterConcrete(ValidatorRevokeMsg{}, "lino/valRevoke", nil)
	cdc.RegisterConcrete(VoteValidatorMsg{}, "lino/voteValidator", nil)
}

// ModuleCdc is the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
