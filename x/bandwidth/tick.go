package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker
func BeginBlocker(
	ctx sdk.Context, req abci.RequestBeginBlock, bm BandwidthKeeper) (tags sdk.Tags) {
	if err := bm.BeginBlocker(ctx); err != nil {
		panic(err)
	}
	return
}

// EndBlocker
func EndBlocker(
	ctx sdk.Context, req abci.RequestEndBlock, bm BandwidthKeeper) (tags sdk.Tags) {
	if err := bm.EndBlocker(ctx); err != nil {
		panic(err)
	}
	return
}
