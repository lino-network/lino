package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker - execute before every block.
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, vm ValidatorKeeper) {
	vm.OnBeginBlock(ctx, req)
}
