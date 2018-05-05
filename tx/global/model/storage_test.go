package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestGlobalKVStoreKey = sdk.NewKVStoreKey("global")
	TestParamKVStoreKey  = sdk.NewKVStoreKey("param")
)

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func InitGlobalStorage(
	t *testing.T, ctx sdk.Context, gm GlobalStorage, param *param.GlobalAllocationParam) sdk.Error {
	return gm.InitGlobalState(ctx, types.NewCoin(10000*types.Decimals), param)
}

func checkGlobalStorage(t *testing.T, ctx sdk.Context, gm GlobalStorage, ph param.ParamHolder, expectGlobalStatistic GlobalStatistics,
	expectGlobalMeta GlobalMeta, expectGlobalAllocation param.GlobalAllocationParam, expectConsumptionMeta ConsumptionMeta,
	expectInflationPool InflationPool) {
	globalStatistic, err := gm.GetGlobalStatistics(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectGlobalStatistic, *globalStatistic)
	globalMeta, err := gm.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectGlobalMeta, *globalMeta)
	globalAllocation, err := ph.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectGlobalAllocation, *globalAllocation)
	consumptionMeta, err := gm.GetConsumptionMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectConsumptionMeta, *consumptionMeta)
	inflationPool, err := gm.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectInflationPool, *inflationPool)
}

func TestGlobalStorageGenesis(t *testing.T) {
	gm := NewGlobalStorage(TestGlobalKVStoreKey)
	ph := param.NewParamHolder(TestParamKVStoreKey)
	ctx := getContext()

	ph.InitParam(ctx)

	allocationParam, err := ph.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)
	err = InitGlobalStorage(t, ctx, gm, allocationParam)
	assert.Nil(t, err)
	globalMeta := GlobalMeta{
		TotalLinoCoin:                 types.NewCoin(10000 * types.Decimals),
		GrowthRate:                    sdk.NewRat(98, 1000),
		CumulativeConsumption:         types.NewCoin(0),
		LastYearCumulativeConsumption: types.NewCoin(0),
		Ceiling: sdk.NewRat(98, 1000),
		Floor:   sdk.NewRat(30, 1000),
	}

	globalStatistics := GlobalStatistics{}

	globalAllocation := param.GlobalAllocationParam{
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{50, 100},
		DeveloperAllocation:      sdk.Rat{20, 100},
		ValidatorAllocation:      sdk.Rat{10, 100},
	}
	consumptionMeta := ConsumptionMeta{
		ConsumptionFrictionRate:     sdk.Rat{5, 100},
		ReportStakeWindow:           sdk.ZeroRat,
		DislikeStakeWindow:          sdk.ZeroRat,
		ConsumptionWindow:           types.NewCoin(0),
		ConsumptionRewardPool:       types.NewCoin(0),
		ConsumptionFreezingPeriodHr: 24 * 7,
	}
	inflationPool := InflationPool{
		InfraInflationPool:          types.NewCoin(196 * types.Decimals),
		ContentCreatorInflationPool: types.NewCoin(490 * types.Decimals),
		DeveloperInflationPool:      types.NewCoin(196 * types.Decimals),
		ValidatorInflationPool:      types.NewCoin(98 * types.Decimals),
	}
	checkGlobalStorage(t, ctx, gm, ph, globalStatistics, globalMeta, globalAllocation, consumptionMeta, inflationPool)
}
