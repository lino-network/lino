package global

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func initGlobalManager(t *testing.T, ctx sdk.Context, gm GlobalManager) error {
	globalState := genesis.GlobalState{
		TotalLino:                types.LNO(sdk.NewRat(10000)),
		GrowthRate:               sdk.Rat{98, 1000},
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{20, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
		ConsumptionFrictionRate:  sdk.Rat{1, 100},
		FreezingPeriodHr:         24 * 7,
	}
	return gm.initGlobalState(ctx, globalState)
}

func checkGlobalStorage(t *testing.T, ctx sdk.Context, gm GlobalManager, expectGlobalStatistic GlobalStatistics,
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

func TestGlobalManagerGenesis(t *testing.T) {
	gm := NewGlobalManager(TestKVStoreKey)
	ctx := getContext()

	err := initGlobalManager(t, ctx, gm)
	assert.Nil(t, err)
	globalMeta := GlobalMeta{
		TotalLino:  types.LNO(sdk.NewRat(10000)),
		GrowthRate: sdk.Rat{98, 1000},
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
		FreezingPeriodHr:        24 * 7,
	}
	inflationPool := InflationPool{
		InfraInflationPool:          types.Coin{196 * types.Decimals},
		ContentCreatorInflationPool: types.Coin{539 * types.Decimals},
		DeveloperInflationPool:      types.Coin{196 * types.Decimals},
		ValidatorInflationPool:      types.Coin{49 * types.Decimals},
	}
	checkGlobalStorage(t, ctx, gm, globalStatistics, globalMeta, globalAllocation, consumptionMeta, inflationPool)
}
