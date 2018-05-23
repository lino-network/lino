package global

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
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
	TestParamKVStoreKey  = sdk.NewKVStoreKey("param")
)

func InitGlobalManager(ctx sdk.Context, gm GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoin(10000*types.Decimals))
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func setupTest(t *testing.T) (sdk.Context, GlobalManager) {
	ctx := getContext()
	holder := param.NewParamHolder(TestParamKVStoreKey)
	holder.InitParam(ctx)
	globalManager := NewGlobalManager(TestGlobalKVStoreKey, holder)
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
	paras, err := gm.paramHolder.GetEvaluateOfContentValueParam(ctx)
	assert.Nil(t, err)
	cases := []struct {
		createdTime                        int64
		evaluateTime                       int64
		expectedTimeAdjustment             float64
		totalConsumption                   types.Coin
		expectedTotalConsumptionAdjustment float64
		numOfConsumptionOnAuthor           int64
		expectedConumptionTimesAdjustment  float64
		Consumption                        types.Coin
		ExpectEvaluateResult               types.Coin
	}{
		{baseTime, baseTime + 3153600*5, 0.5, types.NewCoin(5000 * types.Decimals), 1.5, 7, 2.5,
			types.NewCoin(1000), types.NewCoin(470)},
		{baseTime, baseTime, 0.9933071490757153, types.NewCoin(0), 1.9933071490757153, 7, 2.5,
			types.NewCoin(1000), types.NewCoin(1243)},
		{baseTime, baseTime + 360*24*3600, 0.007667910249215412, types.NewCoin(0), 1.9933071490757153, 7, 2.5,
			types.NewCoin(1000), types.NewCoin(9)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoin(0), 1.9933071490757153, 7, 2.5,
			types.NewCoin(5 * types.Decimals), types.NewCoin(179346)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoin(0), 1.9933071490757153, 7, 2.5,
			types.NewCoin(1 * types.Decimals), types.NewCoin(49489)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoin(1000 * types.Decimals), 1.9820137900379085, 7, 2.5,
			types.NewCoin(1 * types.Decimals), types.NewCoin(49209)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoin(5000 * types.Decimals), 1.5, 7, 2.5,
			types.NewCoin(1 * types.Decimals), types.NewCoin(37242)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoin(5000 * types.Decimals), 1.5, 100, 2,
			types.NewCoin(1 * types.Decimals), types.NewCoin(29793)},
	}

	for _, cs := range cases {
		newCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: cs.evaluateTime})
		assert.Equal(t, cs.expectedTotalConsumptionAdjustment,
			PostTotalConsumptionAdjustment(cs.totalConsumption, paras))
		assert.Equal(t, cs.expectedTimeAdjustment,
			PostTimeAdjustment(cs.evaluateTime-cs.createdTime, paras))
		assert.Equal(t, cs.expectedConumptionTimesAdjustment,
			PostConsumptionTimesAdjustment(cs.numOfConsumptionOnAuthor, paras))
		evaluateResult, err := gm.EvaluateConsumption(
			newCtx, cs.Consumption, cs.numOfConsumptionOnAuthor,
			cs.createdTime, cs.totalConsumption)
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
		consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
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
		evaluate                    types.Coin
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
		// issue https://github.com/lino-network/lino/issues/150
		{types.NewCoin(77777777777777), sdk.ZeroRat, types.NewCoin(23300000),
			types.NewCoin(100000000), types.NewCoin(333333333333333),
			types.NewCoin(76700000), types.NewCoin(255555555555556)},
	}

	for _, cs := range cases {
		consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
		assert.Nil(t, err)
		consumptionMeta.ConsumptionRewardPool = cs.initConsumptionRewardPool
		consumptionMeta.ConsumptionWindow = cs.initConsumptionWindow
		err = gm.storage.SetConsumptionMeta(ctx, consumptionMeta)
		assert.Nil(t, err)
		reward, err := gm.GetRewardAndPopFromWindow(ctx, cs.evaluate, cs.penaltyScore)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectReward, reward)
		consumptionMeta, err = gm.storage.GetConsumptionMeta(ctx)
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
	err := gm.storage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)
	for i := 0; i < types.HoursPerYear; i++ {
		pool, err := gm.storage.GetInflationPool(ctx)
		assert.Nil(t, err)
		consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
		assert.Nil(t, err)
		err = gm.AddHourlyInflationToRewardPool(ctx, int64(i+1))
		assert.Nil(t, err)
		assert.Equal(t, totalConsumption,
			consumptionMeta.ConsumptionRewardPool.Plus(pool.ContentCreatorInflationPool))
	}
	pool, err := gm.storage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoin(0), pool.ContentCreatorInflationPool)

	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, globalMeta.TotalLinoCoin, types.NewCoin(10000*types.Decimals).Plus(totalConsumption))
}

func TestRecalculateAnnuallyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalLino := types.NewCoin(10000000000 * types.Decimals)
	ceiling := sdk.NewRat(98, 1000)
	floor := sdk.NewRat(30, 1000)

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
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(90000000 * types.Decimals),
			types.NewCoin(60000000 * types.Decimals), types.NewCoin(150000000 * types.Decimals),
			types.NewCoin(60000000 * types.Decimals), types.NewCoin(30000000 * types.Decimals)},
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
		err := gm.storage.SetGlobalMeta(ctx, globalMeta)
		assert.Nil(t, err)
		err = gm.RecalculateAnnuallyInflation(ctx)
		assert.Nil(t, err)
		pool, err := gm.storage.GetInflationPool(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectDeveloperInflation, pool.DeveloperInflationPool)
		assert.Equal(t, cs.expectContentCreatorInflation, pool.ContentCreatorInflationPool)
		assert.Equal(t, cs.expectInfraInflation, pool.InfraInflationPool)
		assert.Equal(t, cs.expectValidatorInflation, pool.ValidatorInflationPool)
		globalMeta, err = gm.storage.GetGlobalMeta(ctx)
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
		{types.NewCoin(100000000 * types.Decimals), types.NewCoin(0 * types.Decimals),
			sdk.NewRat(98, 1000), floor},
		{types.NewCoin(0 * types.Decimals), types.NewCoin(100000000 * types.Decimals),
			sdk.NewRat(98, 1000), sdk.NewRat(98, 1000)},
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
		// issue https://github.com/lino-network/lino/issues/150
		{types.NewCoin(77777777), types.NewCoin(77777777 + 5555555),
			sdk.NewRat(98, 1000), sdk.NewRat(71, 1000)},
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
		err := gm.storage.SetGlobalMeta(ctx, globalMeta)
		assert.Nil(t, err)
		growthRate, err := gm.getGrowthRate(ctx)
		assert.Nil(t, err)
		fmt.Println(cs.expectGrowthRate, growthRate)
		assert.True(t, cs.expectGrowthRate.Equal(growthRate))
		globalMeta, err = gm.storage.GetGlobalMeta(ctx)
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
	err := gm.storage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)
	for i := 0; i < types.HoursPerYear; i++ {
		pool, err := gm.storage.GetInflationPool(ctx)
		assert.Nil(t, err)
		coin, err := gm.GetValidatorHourlyInflation(ctx, int64(i+1))
		assert.Nil(t, err)
		assert.Equal(t, coin, types.RatToCoin(pool.ValidatorInflationPool.ToRat().
			Mul(sdk.NewRat(1, int64(types.HoursPerYear-i)))))
	}
	pool, err := gm.storage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoin(0), pool.ValidatorInflationPool)
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, globalMeta.TotalLinoCoin, types.NewCoin(10000*types.Decimals).Plus(totalValidatorInflation))
}

func TestGetInfraMonthlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalInfraInflation := types.NewCoin(10000 * 100)
	inflationPool := &model.InflationPool{
		InfraInflationPool: totalInfraInflation,
	}
	err := gm.storage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)
	for i := 1; i <= types.HoursPerYear*60; i++ {
		if i%types.MinutesPerMonth == 0 {
			pool, err := gm.storage.GetInflationPool(ctx)
			assert.Nil(t, err)
			coin, err := gm.GetInfraMonthlyInflation(ctx, int64(i/types.MinutesPerMonth-1)%12)
			assert.Nil(t, err)
			assert.Equal(t, coin, types.RatToCoin(pool.InfraInflationPool.ToRat().
				Mul(sdk.NewRat(1, int64(12-(i/types.MinutesPerMonth-1)%12)))))
		}
	}
	pool, err := gm.storage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoin(0), pool.InfraInflationPool)
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, globalMeta.TotalLinoCoin, types.NewCoin(10000*types.Decimals).Plus(totalInfraInflation))
}

func TestGetDeveloperMonthlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalDeveloperInflation := types.NewCoin(10000 * 100)
	inflationPool := &model.InflationPool{
		DeveloperInflationPool: totalDeveloperInflation,
	}
	err := gm.storage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)
	for i := 1; i <= types.HoursPerYear*60; i++ {
		if i%types.MinutesPerMonth == 0 {
			pool, err := gm.storage.GetInflationPool(ctx)
			assert.Nil(t, err)
			coin, err := gm.GetDeveloperMonthlyInflation(ctx, int64(i/types.MinutesPerMonth-1)%12)
			assert.Nil(t, err)
			assert.Equal(t, coin, types.RatToCoin(pool.DeveloperInflationPool.ToRat().
				Mul(sdk.NewRat(1, int64(12-(i/types.MinutesPerMonth-1)%12)))))
		}
	}
	pool, err := gm.storage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoin(0), pool.DeveloperInflationPool)
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, globalMeta.TotalLinoCoin, types.NewCoin(10000*types.Decimals).Plus(totalDeveloperInflation))
}

func TestAddToValidatorInflationPool(t *testing.T) {
	ctx, gm := setupTest(t)
	totalValidatorInflation := types.NewCoin(0)
	inflationPool := &model.InflationPool{
		ValidatorInflationPool: totalValidatorInflation,
	}
	err := gm.storage.SetInflationPool(ctx, inflationPool)
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
		pool, err := gm.storage.GetInflationPool(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expect, pool.ValidatorInflationPool)
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
		globalMeta, err := gm.storage.GetGlobalMeta(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expect, globalMeta.CumulativeConsumption)
	}
}
