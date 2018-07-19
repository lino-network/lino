package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(
	ctx sdk.Context, req abci.RequestBeginBlock, gm GlobalManager,
	lastBlockTime int64) (tags sdk.Tags) {
	gm.UpdateTPS(ctx, lastBlockTime)
	return
}
