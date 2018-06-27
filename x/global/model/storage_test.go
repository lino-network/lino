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
	"github.com/tendermint/tmlibs/log"
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

	return sdk.NewContext(ms, abci.Header{}, false, nil, log.NewNopLogger())
}

func InitGlobalStorage(
	t *testing.T, ctx sdk.Context, gm GlobalStorage, param *param.GlobalAllocationParam) sdk.Error {
	return gm.InitGlobalState(ctx, types.NewCoinFromInt64(10000*types.Decimals), param)
}

func checkGlobalStorage(t *testing.T, ctx sdk.Context, gm GlobalStorage,
	expectGlobalMeta GlobalMeta, expectConsumptionMeta ConsumptionMeta,
	expectInflationPool InflationPool) {
	globalMeta, err := gm.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expectGlobalMeta, *globalMeta)
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
		TotalLinoCoin:                 types.NewCoinFromInt64(10000 * types.Decimals),
		GrowthRate:                    sdk.NewRat(98, 1000),
		CumulativeConsumption:         types.NewCoinFromInt64(0),
		LastYearCumulativeConsumption: types.NewCoinFromInt64(0),
		Ceiling: sdk.NewRat(98, 1000),
		Floor:   sdk.NewRat(30, 1000),
	}

	consumptionMeta := ConsumptionMeta{
		ConsumptionFrictionRate:     sdk.NewRat(5, 100),
		ConsumptionWindow:           types.NewCoinFromInt64(0),
		ConsumptionRewardPool:       types.NewCoinFromInt64(0),
		ConsumptionFreezingPeriodHr: 24 * 7,
	}
	infraInflationPool, _ := types.RatToCoin(globalMeta.GrowthRate.Mul(
		globalMeta.TotalLinoCoin.ToRat().Mul(allocationParam.InfraAllocation)))
	contentCreatorInflationPool, _ := types.RatToCoin(
		globalMeta.GrowthRate.Mul(
			globalMeta.TotalLinoCoin.ToRat().Mul(allocationParam.ContentCreatorAllocation)))
	developerInflaionPool, _ := types.RatToCoin(
		globalMeta.GrowthRate.Mul(
			globalMeta.TotalLinoCoin.ToRat().Mul(allocationParam.DeveloperAllocation)))
	validatorInflaionPool, _ := types.RatToCoin(
		globalMeta.GrowthRate.Mul(
			globalMeta.TotalLinoCoin.ToRat().Mul(allocationParam.ValidatorAllocation)))
	inflationPool := InflationPool{
		InfraInflationPool:          infraInflationPool,
		ContentCreatorInflationPool: contentCreatorInflationPool,
		DeveloperInflationPool:      developerInflaionPool,
		ValidatorInflationPool:      validatorInflaionPool,
	}
	checkGlobalStorage(t, ctx, gm, globalMeta, consumptionMeta, inflationPool)
}
