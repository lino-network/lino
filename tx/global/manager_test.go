package global

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/global/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

const (
	eventTypeTestEvent = "1"
)

type testEvent struct{}

// Construct some global addrs and txs for tests.
var (
	TestGlobalKVStoreKey = sdk.NewKVStoreKey("global")
)

func InitGlobalManager(ctx sdk.Context, gm GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoin(10000*types.Decimals))
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func setupTest(t *testing.T) (sdk.Context, GlobalManager) {
	ctx := getContext()
	globalManager := NewGlobalManager(TestGlobalKVStoreKey)
	cdc := globalManager.WireCodec()
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(testEvent{}, "test", nil)
	err := InitGlobalManager(ctx, globalManager)
	assert.Nil(t, err)
	return ctx, globalManager
}

func TestTPS(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := time.Now().Unix()
	var initMaxTPS = sdk.NewRat(1000)

	cases := []struct {
		BaseTime            int64
		NextTime            int64
		NumOfTx             int32
		ExpectCurrentTPS    sdk.Rat
		ExpectMaxTPS        sdk.Rat
		ExpectCapacityRatio sdk.Rat
	}{
		{BaseTime: baseTime, NextTime: baseTime, NumOfTx: 0, ExpectCurrentTPS: sdk.ZeroRat,
			ExpectMaxTPS: initMaxTPS, ExpectCapacityRatio: sdk.ZeroRat},
		{BaseTime: baseTime, NextTime: baseTime + 2, NumOfTx: 2, ExpectCurrentTPS: sdk.OneRat,
			ExpectMaxTPS: initMaxTPS, ExpectCapacityRatio: sdk.NewRat(1, 1000)},
		{BaseTime: baseTime, NextTime: baseTime + 1, NumOfTx: 1000, ExpectCurrentTPS: initMaxTPS,
			ExpectMaxTPS: initMaxTPS, ExpectCapacityRatio: sdk.OneRat},
		{BaseTime: baseTime, NextTime: baseTime + 2, NumOfTx: 2000, ExpectCurrentTPS: initMaxTPS,
			ExpectMaxTPS: initMaxTPS, ExpectCapacityRatio: sdk.OneRat},
		{BaseTime: baseTime, NextTime: baseTime + 2, NumOfTx: 3000, ExpectCurrentTPS: sdk.NewRat(1500),
			ExpectMaxTPS: sdk.NewRat(1500), ExpectCapacityRatio: sdk.OneRat},
		{BaseTime: baseTime, NextTime: baseTime + 2, NumOfTx: 2000, ExpectCurrentTPS: sdk.NewRat(1000),
			ExpectMaxTPS: sdk.NewRat(1500), ExpectCapacityRatio: sdk.NewRat(2, 3)},
	}
	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: cs.NextTime, NumTxs: cs.NumOfTx})
		err := gm.UpdateTPS(ctx, cs.BaseTime)
		assert.Nil(t, err)
		storage := model.NewGlobalStorage(TestGlobalKVStoreKey)
		tps, err := storage.GetTPS(ctx)
		assert.Equal(t, true, cs.ExpectCurrentTPS.Equal(tps.CurrentTPS))
		assert.Equal(t, true, cs.ExpectMaxTPS.Equal(tps.MaxTPS))
		ratio, err := gm.GetTPSCapacityRatio(ctx)
		assert.Nil(t, err)
		assert.Equal(t, true, cs.ExpectCapacityRatio.Equal(ratio))
	}
}

func TestEvaluateConsumption(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := ctx.BlockHeader().Time
	paras, err := gm.globalStorage.GetEvaluateOfContentValuePara(ctx)
	assert.Nil(t, err)
	cases := []struct {
		createdTime                        int64
		evaluateTime                       int64
		expectedTimeAdjustment             float64
		totalReward                        types.Coin
		expectedTotalConsumptionAdjustment float64
		numOfConsumptionOnAuthor           int64
		expectedConumptionTimesAdjustment  float64
		Consumption                        types.Coin
		ExpectEvaluateResult               types.Coin
	}{
		{baseTime, baseTime + 3153600*5, 0.5, types.NewCoin(5000 * types.Decimals), 1.5, 7, 1.5,
			types.NewCoin(1000), types.NewCoin(282)},
		{baseTime, baseTime, 0.9933071490757153, types.NewCoin(0), 1.9933071490757153, 7, 1.5,
			types.NewCoin(1000), types.NewCoin(746)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoin(0), 1.9933071490757153, 7, 1.5,
			types.NewCoin(1000), types.NewCoin(745)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoin(0), 1.9933071490757153, 7, 1.5,
			types.NewCoin(5 * types.Decimals), types.NewCoin(107607)},
	}

	for _, cs := range cases {
		newCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: cs.evaluateTime})
		assert.Equal(t, cs.expectedTotalConsumptionAdjustment,
			PostTotalConsumptionAdjustment(cs.totalReward, paras))
		assert.Equal(t, cs.expectedTimeAdjustment,
			PostTimeAdjustment(cs.evaluateTime-cs.createdTime, paras))
		assert.Equal(t, cs.expectedConumptionTimesAdjustment,
			PostConsumptionTimesAdjustment(cs.numOfConsumptionOnAuthor, paras))
		evaluateResult, err := gm.EvaluateConsumption(
			newCtx, cs.Consumption, cs.numOfConsumptionOnAuthor,
			cs.createdTime, cs.totalReward)
		assert.Nil(t, err)
		assert.Equal(t, cs.ExpectEvaluateResult, evaluateResult)
	}
}

func TestAddFrictionAndRegisterContentRewardEvent(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := ctx.BlockHeader().Time
	cases := []struct {
		frictionCoin           types.Coin
		evaluateCoin           types.Coin
		registerBaseTime       int64
		expectCoinInRewardPool types.Coin
		expectCoinInWindow     types.Coin
	}{
		{types.NewCoin(1), types.NewCoin(1), baseTime, types.NewCoin(1), types.NewCoin(1)},
		{types.NewCoin(100), types.NewCoin(1), baseTime + 100, types.NewCoin(101), types.NewCoin(2)},
		{types.NewCoin(1), types.NewCoin(100), baseTime + 1001, types.NewCoin(102), types.NewCoin(102)},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: cs.registerBaseTime})
		err := gm.AddFrictionAndRegisterContentRewardEvent(
			ctx, testEvent{}, cs.frictionCoin, cs.evaluateCoin)
		assert.Nil(t, err)
		consumptionMeta, err := gm.globalStorage.GetConsumptionMeta(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectCoinInRewardPool, consumptionMeta.ConsumptionRewardPool)
		assert.Equal(t, cs.expectCoinInWindow, consumptionMeta.ConsumptionWindow)
		timeEventList := gm.GetTimeEventListAtTime(ctx, cs.registerBaseTime+24*7*3600)
		assert.Equal(t, types.TimeEventList{[]types.Event{testEvent{}}}, *timeEventList)
	}
}

func TestGetRewardAndPopFromWindow(t *testing.T) {
	ctx, gm := setupTest(t)
	cases := []struct {
		coin                        types.Coin
		penaltyScore                sdk.Rat
		expectReward                types.Coin
		initConsumptionRewardPool   types.Coin
		initConsumptionWindow       types.Coin
		expectConsumptionRewardPool types.Coin
		expectConsumptionWindow     types.Coin
	}{
		{types.NewCoin(1), sdk.ZeroRat, types.NewCoin(100), types.NewCoin(1000),
			types.NewCoin(10), types.NewCoin(900), types.NewCoin(9)},
		{types.NewCoin(1), sdk.NewRat(1, 1000), types.NewCoin(100), types.NewCoin(1000),
			types.NewCoin(10), types.NewCoin(900), types.NewCoin(9)},
		{types.NewCoin(1), sdk.NewRat(6, 1000), types.NewCoin(99), types.NewCoin(1000),
			types.NewCoin(10), types.NewCoin(901), types.NewCoin(9)},
		{types.NewCoin(1), sdk.NewRat(1, 10), types.NewCoin(90), types.NewCoin(1000),
			types.NewCoin(10), types.NewCoin(910), types.NewCoin(9)},
		{types.NewCoin(1), sdk.NewRat(5, 10), types.NewCoin(50), types.NewCoin(1000),
			types.NewCoin(10), types.NewCoin(950), types.NewCoin(9)},
		{types.NewCoin(1), sdk.OneRat, types.NewCoin(0), types.NewCoin(1000),
			types.NewCoin(10), types.NewCoin(1000), types.NewCoin(9)},
		{types.NewCoin(0), sdk.ZeroRat, types.NewCoin(0), types.NewCoin(1000),
			types.NewCoin(10), types.NewCoin(1000), types.NewCoin(10)},
		{types.NewCoin(0), sdk.OneRat, types.NewCoin(0), types.NewCoin(1000),
			types.NewCoin(10), types.NewCoin(1000), types.NewCoin(10)},
	}

	for _, cs := range cases {
		consumptionMeta, err := gm.globalStorage.GetConsumptionMeta(ctx)
		assert.Nil(t, err)
		consumptionMeta.ConsumptionRewardPool = cs.initConsumptionRewardPool
		consumptionMeta.ConsumptionWindow = cs.initConsumptionWindow
		err = gm.globalStorage.SetConsumptionMeta(ctx, consumptionMeta)
		assert.Nil(t, err)
		reward, err := gm.GetRewardAndPopFromWindow(ctx, cs.coin, cs.penaltyScore)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectReward, reward)
		consumptionMeta, err = gm.globalStorage.GetConsumptionMeta(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectConsumptionRewardPool, consumptionMeta.ConsumptionRewardPool)
		assert.Equal(t, cs.expectConsumptionWindow, consumptionMeta.ConsumptionWindow)
	}
}

func TestTimeEventList(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := ctx.BlockHeader().Time
	regcases := []struct {
		registerAtTime  int64
		expectResult    sdk.Error
		expectEventList *types.TimeEventList
	}{
		{baseTime, nil, &types.TimeEventList{[]types.Event{testEvent{}}}},
		{baseTime, nil, &types.TimeEventList{[]types.Event{testEvent{}, testEvent{}}}},
		{baseTime - 1, ErrGlobalManagerRegisterExpiredEvent(baseTime - 1), nil},
		{baseTime + 1, nil, &types.TimeEventList{[]types.Event{testEvent{}}}},
	}

	for _, cs := range regcases {
		err := gm.registerEventAtTime(ctx, cs.registerAtTime, testEvent{})
		assert.Equal(t, cs.expectResult, err)
		eventList := gm.GetTimeEventListAtTime(ctx, cs.registerAtTime)
		assert.Equal(t, cs.expectEventList, eventList)
	}

	rmcases := []struct {
		removeAtTime    int64
		expectEventList *types.TimeEventList
	}{
		{baseTime, nil},
		{baseTime - 1, nil},
		{baseTime + 1, nil},
	}

	for _, cs := range rmcases {
		err := gm.RemoveTimeEventList(ctx, cs.removeAtTime)
		assert.Nil(t, err)
		eventList := gm.GetTimeEventListAtTime(ctx, cs.removeAtTime)
		assert.Equal(t, cs.expectEventList, eventList)
	}
}

func TestRegisterCoinReturnEvent(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := ctx.BlockHeader().Time
	cases := []struct {
		registerAtTime         int64
		times                  int64
		interval               int64
		expectTimeWithTwoEvent []int64
		expectTimeWithOneEvent []int64
	}{
		{baseTime, 5, 10, []int64{},
			[]int64{baseTime + 10*3600, baseTime + 20*3600, baseTime + 30*3600,
				baseTime + 40*3600, baseTime + 50*3600}},
		{baseTime, 2, 10, []int64{baseTime + 10*3600, baseTime + 20*3600},
			[]int64{baseTime + 30*3600, baseTime + 40*3600, baseTime + 50*3600}},
		{baseTime + 20*3600, 4, 5,
			[]int64{baseTime + 10*3600, baseTime + 20*3600, baseTime + 30*3600,
				baseTime + 40*3600},
			[]int64{baseTime + 25*3600, baseTime + 35*3600, baseTime + 50*3600}},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: cs.registerAtTime})
		events := []types.Event{}
		for i := int64(0); i < cs.times; i++ {
			events = append(events, testEvent{})
		}
		err := gm.RegisterCoinReturnEvent(ctx, events, cs.times, cs.interval)
		assert.Nil(t, err)
		for _, time := range cs.expectTimeWithOneEvent {
			eventList := gm.GetTimeEventListAtTime(ctx, time)
			assert.Equal(t, 1, len(eventList.Events))
		}
		for _, time := range cs.expectTimeWithTwoEvent {
			eventList := gm.GetTimeEventListAtTime(ctx, time)
			assert.Equal(t, 2, len(eventList.Events))
		}
	}
}

func TestAddHourlyInflationToRewardPool(t *testing.T) {
	ctx, gm := setupTest(t)
	totalConsumption := types.NewCoin(10000 * 100)
	inflationPool := &model.InflationPool{
		ContentCreatorInflationPool: totalConsumption,
	}
	err := gm.globalStorage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)
	for i := 0; i < types.HoursPerYear; i++ {
		pool, err := gm.globalStorage.GetInflationPool(ctx)
		assert.Nil(t, err)
		consumptionMeta, err := gm.globalStorage.GetConsumptionMeta(ctx)
		assert.Nil(t, err)
		err = gm.AddHourlyInflationToRewardPool(ctx, int64(i+1))
		assert.Nil(t, err)
		assert.Equal(t, totalConsumption,
			consumptionMeta.ConsumptionRewardPool.Plus(pool.ContentCreatorInflationPool))
	}
	pool, err := gm.globalStorage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoin(0), pool.ContentCreatorInflationPool)
}

func TestRecalculateAnnuallyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalLino := types.NewCoin(10000000000 * types.Decimals)
	ceiling := sdk.NewRat(98, 1000)
	floor := sdk.NewRat(30, 1000)

	globalAllocation := &model.GlobalAllocation{
		InfraAllocation:          sdk.NewRat(20, 100),
		ContentCreatorAllocation: sdk.NewRat(50, 100),
		DeveloperAllocation:      sdk.NewRat(20, 100),
		ValidatorAllocation:      sdk.NewRat(10, 100),
	}
	err := gm.globalStorage.SetGlobalAllocation(ctx, globalAllocation)
	assert.Nil(t, err)

	cases := []struct {
		lastYearConsumtion            types.Coin
		thisYearConsumtion            types.Coin
		expectInfraInflation          types.Coin
		expectContentCreatorInflation types.Coin
		expectDeveloperInflation      types.Coin
		expectValidatorInflation      types.Coin
	}{
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(100000000 * types.Decimals),
			types.NewCoin(60000000 * types.Decimals), types.NewCoin(150000000 * types.Decimals),
			types.NewCoin(60000000 * types.Decimals), types.NewCoin(30000000 * types.Decimals)},
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(103000000 * types.Decimals),
			types.NewCoin(60000000 * types.Decimals), types.NewCoin(150000000 * types.Decimals),
			types.NewCoin(60000000 * types.Decimals), types.NewCoin(30000000 * types.Decimals)},
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(1098000000 * types.Decimals),
			types.NewCoin(196000000 * types.Decimals), types.NewCoin(490000000 * types.Decimals),
			types.NewCoin(196000000 * types.Decimals), types.NewCoin(98000000 * types.Decimals)},
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(1099000000 * types.Decimals),
			types.NewCoin(196000000 * types.Decimals), types.NewCoin(490000000 * types.Decimals),
			types.NewCoin(196000000 * types.Decimals), types.NewCoin(98000000 * types.Decimals)},
	}

	for _, cs := range cases {
		globalMeta := &model.GlobalMeta{
			TotalLinoCoin:                 totalLino,
			LastYearCumulativeConsumption: cs.lastYearConsumtion,
			CumulativeConsumption:         cs.thisYearConsumtion,
			GrowthRate:                    sdk.NewRat(98, 1000),
			Ceiling:                       ceiling,
			Floor:                         floor,
		}
		err := gm.globalStorage.SetGlobalMeta(ctx, globalMeta)
		assert.Nil(t, err)
		err = gm.RecalculateAnnuallyInflation(ctx)
		assert.Nil(t, err)
		pool, err := gm.globalStorage.GetInflationPool(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectDeveloperInflation, pool.DeveloperInflationPool)
		assert.Equal(t, cs.expectContentCreatorInflation, pool.ContentCreatorInflationPool)
		assert.Equal(t, cs.expectInfraInflation, pool.InfraInflationPool)
		assert.Equal(t, cs.expectValidatorInflation, pool.ValidatorInflationPool)
		globalMeta, err = gm.globalStorage.GetGlobalMeta(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.thisYearConsumtion, globalMeta.LastYearCumulativeConsumption)
		assert.Equal(t, types.NewCoin(0), globalMeta.CumulativeConsumption)
	}
}

func TestGetGrowthRate(t *testing.T) {
	ctx, gm := setupTest(t)
	totalLino := types.NewCoin(1000000)
	ceiling := sdk.NewRat(98, 1000)
	floor := sdk.NewRat(30, 1000)

	cases := []struct {
		lastYearConsumtion types.Coin
		thisYearConsumtion types.Coin
		lastYearGrowthRate sdk.Rat
		expectGrowthRate   sdk.Rat
	}{
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(100000000 * types.Decimals),
			sdk.NewRat(98, 1000), floor},
		{types.NewCoin(0), types.NewCoin(100000000 * types.Decimals),
			sdk.NewRat(98, 1000), sdk.NewRat(98, 1000)},
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(100010000 * types.Decimals),
			sdk.NewRat(98, 1000), floor},
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(102900000 * types.Decimals),
			sdk.NewRat(98, 1000), floor},
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(103000000 * types.Decimals),
			sdk.NewRat(98, 1000), floor},
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(103100000 * types.Decimals),
			sdk.NewRat(98, 1000), sdk.NewRat(31, 1000)},
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(109800000 * types.Decimals),
			sdk.NewRat(98, 1000), ceiling},
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(109900000 * types.Decimals),
			sdk.NewRat(98, 1000), ceiling},
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(109700000 * types.Decimals),
			sdk.NewRat(98, 1000), sdk.NewRat(97, 1000)},
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(104700000 * types.Decimals),
			sdk.NewRat(98, 1000), sdk.NewRat(47, 1000)},
	}

	for _, cs := range cases {
		globalMeta := &model.GlobalMeta{
			TotalLinoCoin:                 totalLino,
			LastYearCumulativeConsumption: cs.lastYearConsumtion,
			CumulativeConsumption:         cs.thisYearConsumtion,
			GrowthRate:                    cs.lastYearGrowthRate,
			Ceiling:                       ceiling,
			Floor:                         floor,
		}
		err := gm.globalStorage.SetGlobalMeta(ctx, globalMeta)
		assert.Nil(t, err)
		growthRate, err := gm.getGrowthRate(ctx)
		assert.Nil(t, err)
		assert.True(t, cs.expectGrowthRate.Equal(growthRate))
		globalMeta, err = gm.globalStorage.GetGlobalMeta(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.thisYearConsumtion, globalMeta.LastYearCumulativeConsumption)
		assert.Equal(t, types.NewCoin(0), globalMeta.CumulativeConsumption)
		assert.True(t, cs.expectGrowthRate.Equal(globalMeta.GrowthRate))
	}
}

func TestGetValidatorHourlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalValidatorInflation := types.NewCoin(10000 * 100)
	inflationPool := &model.InflationPool{
		ValidatorInflationPool: totalValidatorInflation,
	}
	err := gm.globalStorage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)
	for i := 0; i < types.HoursPerYear; i++ {
		pool, err := gm.globalStorage.GetInflationPool(ctx)
		assert.Nil(t, err)
		coin, err := gm.GetValidatorHourlyInflation(ctx, int64(i+1))
		assert.Nil(t, err)
		assert.Equal(t, coin, types.RatToCoin(pool.ValidatorInflationPool.ToRat().
			Mul(sdk.NewRat(1, int64(types.HoursPerYear-i)))))
	}
	pool, err := gm.globalStorage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoin(0), pool.ValidatorInflationPool)
}

func TestGetInfraMonthlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalInfraInflation := types.NewCoin(10000 * 100)
	inflationPool := &model.InflationPool{
		InfraInflationPool: totalInfraInflation,
	}
	err := gm.globalStorage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)
	for i := 1; i <= types.HoursPerYear*60; i++ {
		if i%types.MinutesPerMonth == 0 {
			pool, err := gm.globalStorage.GetInflationPool(ctx)
			assert.Nil(t, err)
			coin, err := gm.GetInfraMonthlyInflation(ctx, int64(i/types.MinutesPerMonth-1)%12)
			assert.Nil(t, err)
			assert.Equal(t, coin, types.RatToCoin(pool.InfraInflationPool.ToRat().
				Mul(sdk.NewRat(1, int64(12-(i/types.MinutesPerMonth-1)%12)))))
		}
	}
	pool, err := gm.globalStorage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoin(0), pool.InfraInflationPool)
}

func TestGetDeveloperMonthlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalDeveloperInflation := types.NewCoin(10000 * 100)
	inflationPool := &model.InflationPool{
		DeveloperInflationPool: totalDeveloperInflation,
	}
	err := gm.globalStorage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)
	for i := 1; i <= types.HoursPerYear*60; i++ {
		if i%types.MinutesPerMonth == 0 {
			pool, err := gm.globalStorage.GetInflationPool(ctx)
			assert.Nil(t, err)
			coin, err := gm.GetDeveloperMonthlyInflation(ctx, int64(i/types.MinutesPerMonth-1)%12)
			assert.Nil(t, err)
			assert.Equal(t, coin, types.RatToCoin(pool.DeveloperInflationPool.ToRat().
				Mul(sdk.NewRat(1, int64(12-(i/types.MinutesPerMonth-1)%12)))))
		}
	}
	pool, err := gm.globalStorage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoin(0), pool.DeveloperInflationPool)
}

func TestAddToValidatorInflationPool(t *testing.T) {
	ctx, gm := setupTest(t)
	totalValidatorInflation := types.NewCoin(0)
	inflationPool := &model.InflationPool{
		ValidatorInflationPool: totalValidatorInflation,
	}
	err := gm.globalStorage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)

	cases := []struct {
		coin   types.Coin
		expect types.Coin
	}{
		{types.NewCoin(100), types.NewCoin(100)},
		{types.NewCoin(1), types.NewCoin(101)},
	}

	for _, cs := range cases {
		err := gm.AddToValidatorInflationPool(ctx, cs.coin)
		assert.Nil(t, err)
		pool, err := gm.globalStorage.GetInflationPool(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expect, pool.ValidatorInflationPool)
	}
}

func TestChangeInfraInternalInflation(t *testing.T) {
	ctx, gm := setupTest(t)

	cases := []struct {
		storageAllocation sdk.Rat
		CDNAllocation     sdk.Rat
	}{
		{sdk.NewRat(1, 100), sdk.NewRat(99, 100)},
		{sdk.ZeroRat, sdk.OneRat},
	}

	for _, cs := range cases {
		err := gm.ChangeInfraInternalInflation(ctx, cs.storageAllocation, cs.CDNAllocation)
		assert.Nil(t, err)
		allocation, err := gm.globalStorage.GetInfraInternalAllocation(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.storageAllocation, allocation.StorageAllocation)
		assert.Equal(t, cs.CDNAllocation, allocation.CDNAllocation)
	}
}

func TestAddConsumption(t *testing.T) {
	ctx, gm := setupTest(t)

	cases := []struct {
		coin   types.Coin
		expect types.Coin
	}{
		{types.NewCoin(100), types.NewCoin(100)},
		{types.NewCoin(1), types.NewCoin(101)},
	}

	for _, cs := range cases {
		err := gm.AddConsumption(ctx, cs.coin)
		assert.Nil(t, err)
		globalMeta, err := gm.globalStorage.GetGlobalMeta(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expect, globalMeta.CumulativeConsumption)
	}
}

func TestChangeGlobalInflation(t *testing.T) {
	ctx, gm := setupTest(t)

	cases := []struct {
		contentCreatorAllocation sdk.Rat
		developerAllocation      sdk.Rat
		infraAllocation          sdk.Rat
		validatorAllocation      sdk.Rat
	}{
		{sdk.NewRat(1, 100), sdk.NewRat(50, 100), sdk.NewRat(20, 100), sdk.NewRat(29, 100)},
	}

	for _, cs := range cases {
		err := gm.ChangeGlobalInflation(
			ctx, cs.infraAllocation, cs.contentCreatorAllocation,
			cs.developerAllocation, cs.validatorAllocation)
		assert.Nil(t, err)
		allocation, err := gm.globalStorage.GetGlobalAllocation(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.contentCreatorAllocation, allocation.ContentCreatorAllocation)
		assert.Equal(t, cs.developerAllocation, allocation.DeveloperAllocation)
		assert.Equal(t, cs.validatorAllocation, allocation.ValidatorAllocation)
		assert.Equal(t, cs.infraAllocation, allocation.InfraAllocation)
	}
}
