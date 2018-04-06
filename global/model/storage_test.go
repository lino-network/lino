package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/genesis"
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

func InitGlobalStorage(t *testing.T, ctx sdk.Context, gm *GlobalStorage) error {
	globalState := genesis.GlobalState{
		TotalLino:                10000,
		GrowthRate:               sdk.Rat{98, 1000},
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{20, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
		ConsumptionFrictionRate:  sdk.Rat{1, 100},
		FreezingPeriodHr:         24 * 7,
	}
	return gm.InitGlobalState(ctx, globalState)
}

func checkGlobalStorage(t *testing.T, ctx sdk.Context, gm GlobalStorage, expectGlobalStatistic GlobalStatistics,
	expectGlobalMeta GlobalMeta, expectGlobalAllocation GlobalAllocation, expectConsumptionMeta ConsumptionMeta,
	expectInflationPool InflationPool) {
	globalStatistic, err := gm.GetGlobalStatistics(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectGlobalStatistic, *globalStatistic)
	globalMeta, err := gm.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectGlobalMeta, *globalMeta)
	globalAllocation, err := gm.GetGlobalAllocation(ctx)
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
		TotalLino:             types.LNO(sdk.NewRat(10000)),
		GrowthRate:            sdk.Rat{98, 1000},
		CumulativeConsumption: types.NewCoin(0),
	}

	globalStatistics := GlobalStatistics{}

	globalAllocation := GlobalAllocation{
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{20, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
	}
	consumptionMeta := ConsumptionMeta{
		ConsumptionFrictionRate: sdk.Rat{1, 100},
		ReportStakeWindow:       sdk.ZeroRat,
		DislikeStakeWindow:      sdk.ZeroRat,
		ConsumptionWindow:       types.NewCoin(0),
		ConsumptionRewardPool:   types.NewCoin(0),
		FreezingPeriodHr:        24 * 7,
	}
	inflationPool := InflationPool{
		InfraInflationPool:          types.NewCoin(196 * types.Decimals),
		ContentCreatorInflationPool: types.NewCoin(539 * types.Decimals),
		DeveloperInflationPool:      types.NewCoin(196 * types.Decimals),
		ValidatorInflationPool:      types.NewCoin(49 * types.Decimals),
	}
	checkGlobalStorage(t, ctx, gm, globalStatistics, globalMeta, globalAllocation, consumptionMeta, inflationPool)
}
