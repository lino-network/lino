package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker
func BeginBlocker(
	ctx sdk.Context, req abci.RequestBeginBlock, bm BandwidthKeeper) (tags sdk.Tags) {
	// calculate the new general msg fee for the current block
	if err := bm.CalculateCurMsgFee(ctx); err != nil {
		panic(err)
	}

	// clear stats for block info
	if err := bm.ClearBlockInfo(ctx); err != nil {
		panic(err)
	}
	return
}

// EndBlocker
func EndBlocker(
	ctx sdk.Context, req abci.RequestEndBlock, bm BandwidthKeeper) (tags sdk.Tags) {
	// update maxMPS and EMA for different msgs and store cur block info
	if err := bm.UpdateMaxMPSAndEMA(ctx); err != nil {
		panic(err)
	}
	return
}
