package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker - called every begin blocker, udpate transaction per second
func BeginBlocker(
	ctx sdk.Context, req abci.RequestBeginBlock, gm *GlobalManager) (tags sdk.Tags) {
	if err := gm.ClearEventCache(ctx); err != nil {
		panic(err)
	}
	return
}

// EndBlocker - related to upgrade1update3.
func EndBlocker(
	ctx sdk.Context, req abci.RequestEndBlock, gm *GlobalManager) (tags sdk.Tags) {
	if err := gm.CommitEventCache(ctx); err != nil {
		panic(err)
	}
	return
}
