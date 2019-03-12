package global

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/lino-network/lino/types"
)

// BeginBlocker - called every begin blocker, udpate transaction per second
func BeginBlocker(
	ctx sdk.Context, req abci.RequestBeginBlock, gm *GlobalManager) (tags sdk.Tags) {
	if err := gm.UpdateTPS(ctx); err != nil {
		panic(err)
	}
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

	// one time execution for upgrade1update3, divided by 6 because it's the approximate right
	// reward pool amount.
	if ctx.BlockHeight() == types.BlockchainUpgrade1Update3Height {
		consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
		if err != nil {
			panic(err)
		}
		consumptionMeta.ConsumptionRewardPool = types.DecToCoin(
			consumptionMeta.ConsumptionRewardPool.ToDec().Quo(sdk.NewDec(6)))
		if err := gm.storage.SetConsumptionMeta(ctx, consumptionMeta); err != nil {
			panic(err)
		}
	}

	return
}
