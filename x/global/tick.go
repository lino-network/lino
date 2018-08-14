package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(
	ctx sdk.Context, req abci.RequestBeginBlock, gm GlobalManager) (tags sdk.Tags) {
	gm.UpdateTPS(ctx)
	return
}
