package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker - called every end blocker, udpate new round
func EndBlocker(ctx sdk.Context, req abci.RequestEndBlock, am AccountKeeper) {
	err := am.Mint(ctx)
	if err != nil {
		panic(err)
	}
	return
}
