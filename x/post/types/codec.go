package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(CreatePostMsg{}, "lino/createPost", nil)
	cdc.RegisterConcrete(UpdatePostMsg{}, "lino/updatePost", nil)
	cdc.RegisterConcrete(DeletePostMsg{}, "lino/deletePost", nil)
	cdc.RegisterConcrete(DonateMsg{}, "lino/donate", nil)
	cdc.RegisterConcrete(IDADonateMsg{}, "lino/idaDonate", nil)
}

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
