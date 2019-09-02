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

	// calculate vacancy coefficient
	u, err := bm.GetVacancyCoeff(ctx)
	if err != nil {
		panic(err)
	}

	// get all app bandwidth info
	allInfo, err := bm.GetAllAppInfo(ctx)
	if err != nil {
		panic(err)
	}

	for _, info := range allInfo {
		if info.MessagesInCurBlock == 0 {
			continue
		}
		// refill bandwidth for apps with messages in current block
		if err := bm.RefillAppBandwidthCredit(ctx, info.Username); err != nil {
			panic(err)
		}
		// calculate cost and consume bandwidth credit
		p, err := bm.GetPunishmentCoeff(ctx, info.Username)
		if err != nil {
			panic(err)
		}
		costPerMsg := bm.GetBandwidthCostPerMsg(ctx, u, p)
		if err := bm.ConsumeBandwidthCredit(ctx, costPerMsg, info.Username); err != nil {
			panic(err)
		}
	}
	return
}
