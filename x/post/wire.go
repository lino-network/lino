package post

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(CreatePostMsg{}, "lino/createPost", nil)
	cdc.RegisterConcrete(UpdatePostMsg{}, "lino/updatePost", nil)
	cdc.RegisterConcrete(DeletePostMsg{}, "lino/deletePost", nil)
	cdc.RegisterConcrete(DonateMsg{}, "lino/donate", nil)
	cdc.RegisterConcrete(ViewMsg{}, "lino/view", nil)
	cdc.RegisterConcrete(ReportOrUpvoteMsg{}, "lino/reportOrUpvote", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
}
