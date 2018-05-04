package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("global")
)

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func InitGlobalStorage(t *testing.T, ctx sdk.Context, gm GlobalStorage) error {
	return gm.InitGlobalState(ctx, types.NewCoin(10000*types.Decimals))
}

func checkGlobalStorage(t *testing.T, ctx sdk.Context, gm GlobalStorage, expectGlobalStatistic GlobalStatistics,
	expectGlobalMeta GlobalMeta, expectGlobalAllocation GlobalAllocationParam, expectConsumptionMeta ConsumptionMeta,
	expectInflationPool InflationPool) {
	globalStatistic, err := gm.GetGlobalStatistics(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectGlobalStatistic, *globalStatistic)
	globalMeta, err := gm.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectGlobalMeta, *globalMeta)
	globalAllocation, err := gm.GetGlobalAllocationParam(ctx)
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
	gm := NewGlobalStorage(TestKVStoreKey)
	ctx := getContext()

	err := InitGlobalStorage(t, ctx, gm)
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

	globalAllocation := GlobalAllocationParam{
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
	checkGlobalStorage(t, ctx, gm, globalStatistics, globalMeta, globalAllocation, consumptionMeta, inflationPool)
}
