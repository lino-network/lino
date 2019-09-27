package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(FeedPriceMsg{}, "lino/feedprice", nil)
}

// ModuleCdc is the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
