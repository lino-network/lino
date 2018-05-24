package global

import (
	"fmt"
	"math/big"
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
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
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
		{baseTime, baseTime + 3153600*5, 0.5, types.NewCoinFromInt64(5000 * types.Decimals), 1.5, 7, 2.5,
			types.NewCoinFromInt64(1000), types.NewCoinFromInt64(470)},
		{baseTime, baseTime, 0.9933071490757153, types.NewCoinFromInt64(0), 1.9933071490757153, 7, 2.5,
			types.NewCoinFromInt64(1000), types.NewCoinFromInt64(1243)},
		{baseTime, baseTime + 360*24*3600, 0.007667910249215412, types.NewCoinFromInt64(0), 1.9933071490757153, 7, 2.5,
			types.NewCoinFromInt64(1000), types.NewCoinFromInt64(9)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoinFromInt64(0), 1.9933071490757153, 7, 2.5,
			types.NewCoinFromInt64(5 * types.Decimals), types.NewCoinFromInt64(179346)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoinFromInt64(0), 1.9933071490757153, 7, 2.5,
			types.NewCoinFromInt64(1 * types.Decimals), types.NewCoinFromInt64(49489)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoinFromInt64(1000 * types.Decimals), 1.9820137900379085, 7, 2.5,
			types.NewCoinFromInt64(1 * types.Decimals), types.NewCoinFromInt64(49209)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoinFromInt64(5000 * types.Decimals), 1.5, 7, 2.5,
			types.NewCoinFromInt64(1 * types.Decimals), types.NewCoinFromInt64(37242)},
		{baseTime, baseTime + 24*3600, 0.9931225268669581, types.NewCoinFromInt64(5000 * types.Decimals), 1.5, 100, 2,
			types.NewCoinFromInt64(1 * types.Decimals), types.NewCoinFromInt64(29793)},
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
		{types.NewCoinFromInt64(1), types.NewCoinFromInt64(1), baseTime, types.NewCoinFromInt64(1), types.NewCoinFromInt64(1)},
		{types.NewCoinFromInt64(100), types.NewCoinFromInt64(1), baseTime + 100, types.NewCoinFromInt64(101), types.NewCoinFromInt64(2)},
		{types.NewCoinFromInt64(1), types.NewCoinFromInt64(100), baseTime + 1001, types.NewCoinFromInt64(102), types.NewCoinFromInt64(102)},
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
		penaltyScore                *big.Rat
		expectReward                types.Coin
		initConsumptionRewardPool   types.Coin
		initConsumptionWindow       types.Coin
		expectConsumptionRewardPool types.Coin
		expectConsumptionWindow     types.Coin
	}{
		{types.NewCoinFromInt64(1), big.NewRat(0, 1), types.NewCoinFromInt64(100), types.NewCoinFromInt64(1000),
			types.NewCoinFromInt64(10), types.NewCoinFromInt64(900), types.NewCoinFromInt64(9)},
		{types.NewCoinFromInt64(1), big.NewRat(1, 1000), types.NewCoinFromInt64(100), types.NewCoinFromInt64(1000),
			types.NewCoinFromInt64(10), types.NewCoinFromInt64(900), types.NewCoinFromInt64(9)},
		{types.NewCoinFromInt64(1), big.NewRat(6, 1000), types.NewCoinFromInt64(99), types.NewCoinFromInt64(1000),
			types.NewCoinFromInt64(10), types.NewCoinFromInt64(901), types.NewCoinFromInt64(9)},
		{types.NewCoinFromInt64(1), big.NewRat(1, 10), types.NewCoinFromInt64(90), types.NewCoinFromInt64(1000),
			types.NewCoinFromInt64(10), types.NewCoinFromInt64(910), types.NewCoinFromInt64(9)},
		{types.NewCoinFromInt64(1), big.NewRat(5, 10), types.NewCoinFromInt64(50), types.NewCoinFromInt64(1000),
			types.NewCoinFromInt64(10), types.NewCoinFromInt64(950), types.NewCoinFromInt64(9)},
		{types.NewCoinFromInt64(1), big.NewRat(1, 1), types.NewCoinFromInt64(0), types.NewCoinFromInt64(1000),
			types.NewCoinFromInt64(10), types.NewCoinFromInt64(1000), types.NewCoinFromInt64(9)},
		{types.NewCoinFromInt64(0), big.NewRat(0, 1), types.NewCoinFromInt64(0), types.NewCoinFromInt64(1000),
			types.NewCoinFromInt64(10), types.NewCoinFromInt64(1000), types.NewCoinFromInt64(10)},
		{types.NewCoinFromInt64(0), big.NewRat(1, 1), types.NewCoinFromInt64(0), types.NewCoinFromInt64(1000),
			types.NewCoinFromInt64(10), types.NewCoinFromInt64(1000), types.NewCoinFromInt64(10)},
		// issue https://github.com/lino-network/lino/issues/150
		{types.NewCoinFromInt64(77777777777777), big.NewRat(0, 1), types.NewCoinFromInt64(23333333),
			types.NewCoinFromInt64(100000000), types.NewCoinFromInt64(333333333333333),
			types.NewCoinFromInt64(76666667), types.NewCoinFromInt64(255555555555556)},
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
	totalConsumption := types.NewCoinFromInt64(10000 * 100)
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
	assert.Equal(t, types.NewCoinFromInt64(0), pool.ContentCreatorInflationPool)

	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, globalMeta.TotalLinoCoin, types.NewCoinFromInt64(10000*types.Decimals).Plus(totalConsumption))
}

func TestRecalculateAnnuallyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalLino := types.NewCoinFromInt64(10000000000 * types.Decimals)
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
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(100000000 * types.Decimals),
			types.NewCoinFromInt64(60000000 * types.Decimals), types.NewCoinFromInt64(150000000 * types.Decimals),
			types.NewCoinFromInt64(60000000 * types.Decimals), types.NewCoinFromInt64(30000000 * types.Decimals)},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(103000000 * types.Decimals),
			types.NewCoinFromInt64(60000000 * types.Decimals), types.NewCoinFromInt64(150000000 * types.Decimals),
			types.NewCoinFromInt64(60000000 * types.Decimals), types.NewCoinFromInt64(30000000 * types.Decimals)},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(1098000000 * types.Decimals),
			types.NewCoinFromInt64(196000000 * types.Decimals), types.NewCoinFromInt64(490000000 * types.Decimals),
			types.NewCoinFromInt64(196000000 * types.Decimals), types.NewCoinFromInt64(98000000 * types.Decimals)},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(1099000000 * types.Decimals),
			types.NewCoinFromInt64(196000000 * types.Decimals), types.NewCoinFromInt64(490000000 * types.Decimals),
			types.NewCoinFromInt64(196000000 * types.Decimals), types.NewCoinFromInt64(98000000 * types.Decimals)},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(90000000 * types.Decimals),
			types.NewCoinFromInt64(60000000 * types.Decimals), types.NewCoinFromInt64(150000000 * types.Decimals),
			types.NewCoinFromInt64(60000000 * types.Decimals), types.NewCoinFromInt64(30000000 * types.Decimals)},
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
		assert.Equal(t, types.NewCoinFromInt64(0), globalMeta.CumulativeConsumption)
	}
}

func TestGetGrowthRate(t *testing.T) {
	ctx, gm := setupTest(t)
	totalLino := types.NewCoinFromInt64(1000000)
	ceiling := sdk.NewRat(98, 1000)
	floor := sdk.NewRat(30, 1000)
	bigLastYearConsumption, _ := new(big.Int).SetString("77777777777777777777", 10)
	bigThisYearConsumption, _ := new(big.Int).SetString("83333333333333333332", 10)
	bigLastYearConsumptionCoin, _ := types.NewCoinFromBigInt(bigLastYearConsumption)
	bigThisYearConsumptionCoin, _ := types.NewCoinFromBigInt(bigThisYearConsumption)
	cases := []struct {
		lastYearConsumtion types.Coin
		thisYearConsumtion types.Coin
		lastYearGrowthRate sdk.Rat
		expectGrowthRate   sdk.Rat
	}{
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(0 * types.Decimals),
			sdk.NewRat(98, 1000), floor},
		{types.NewCoinFromInt64(0 * types.Decimals), types.NewCoinFromInt64(100000000 * types.Decimals),
			sdk.NewRat(98, 1000), sdk.NewRat(98, 1000)},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(100000000 * types.Decimals),
			sdk.NewRat(98, 1000), floor},
		{types.NewCoinFromInt64(0), types.NewCoinFromInt64(100000000 * types.Decimals),
			sdk.NewRat(98, 1000), sdk.NewRat(98, 1000)},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(100010000 * types.Decimals),
			sdk.NewRat(98, 1000), floor},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(102900000 * types.Decimals),
			sdk.NewRat(98, 1000), floor},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(103000000 * types.Decimals),
			sdk.NewRat(98, 1000), floor},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(103100000 * types.Decimals),
			sdk.NewRat(98, 1000), sdk.NewRat(31, 1000)},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(109800000 * types.Decimals),
			sdk.NewRat(98, 1000), ceiling},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(109900000 * types.Decimals),
			sdk.NewRat(98, 1000), ceiling},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(109700000 * types.Decimals),
			sdk.NewRat(98, 1000), sdk.NewRat(97, 1000)},
		{types.NewCoinFromInt64(100000000 * types.Decimals), types.NewCoinFromInt64(104700000 * types.Decimals),
			sdk.NewRat(98, 1000), sdk.NewRat(47, 1000)},
		// issue https://github.com/lino-network/lino/issues/150
		{bigLastYearConsumptionCoin, bigThisYearConsumptionCoin,
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
		fmt.Println(growthRate, cs.expectGrowthRate)
		assert.True(t, cs.expectGrowthRate.Equal(growthRate))
		globalMeta, err = gm.storage.GetGlobalMeta(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.thisYearConsumtion, globalMeta.LastYearCumulativeConsumption)
		assert.Equal(t, types.NewCoinFromInt64(0), globalMeta.CumulativeConsumption)
		assert.True(t, cs.expectGrowthRate.Equal(globalMeta.GrowthRate))
	}
}

func TestGetValidatorHourlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalValidatorInflation := types.NewCoinFromInt64(10000 * 100)
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
		hourlyCoinRat := new(big.Rat).Mul(pool.ValidatorInflationPool.ToRat(), big.NewRat(1, int64(types.HoursPerYear-i)))
		hourlyCoin, err := types.RatToCoin(hourlyCoinRat)
		assert.Nil(t, err)

		assert.Equal(t, coin, hourlyCoin)
	}
	pool, err := gm.storage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoinFromInt64(0), pool.ValidatorInflationPool)
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, globalMeta.TotalLinoCoin, types.NewCoinFromInt64(10000*types.Decimals).Plus(totalValidatorInflation))
}

func TestGetInfraMonthlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalInfraInflation := types.NewCoinFromInt64(10000 * 100)
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
			hourlyCoinRat := new(big.Rat).Mul(
				pool.InfraInflationPool.ToRat(),
				big.NewRat(1, int64(12-(i/types.MinutesPerMonth-1)%12)))
			hourlyCoin, err := types.RatToCoin(hourlyCoinRat)
			assert.Nil(t, err)
			assert.Equal(t, coin, hourlyCoin)
		}
	}
	pool, err := gm.storage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoinFromInt64(0), pool.InfraInflationPool)
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, globalMeta.TotalLinoCoin, types.NewCoinFromInt64(10000*types.Decimals).Plus(totalInfraInflation))
}

func TestGetDeveloperMonthlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalDeveloperInflation := types.NewCoinFromInt64(10000 * 100)
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
			hourlyCoinRat := new(big.Rat).Mul(
				pool.DeveloperInflationPool.ToRat(),
				big.NewRat(1, int64(12-(i/types.MinutesPerMonth-1)%12)))
			hourlyCoin, err := types.RatToCoin(hourlyCoinRat)
			assert.Equal(t, coin, hourlyCoin)
		}
	}
	pool, err := gm.storage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoinFromInt64(0), pool.DeveloperInflationPool)
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, globalMeta.TotalLinoCoin, types.NewCoinFromInt64(10000*types.Decimals).Plus(totalDeveloperInflation))
}

func TestAddToValidatorInflationPool(t *testing.T) {
	ctx, gm := setupTest(t)
	totalValidatorInflation := types.NewCoinFromInt64(0)
	inflationPool := &model.InflationPool{
		ValidatorInflationPool: totalValidatorInflation,
	}
	err := gm.storage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)

	cases := []struct {
		coin   types.Coin
		expect types.Coin
	}{
		{types.NewCoinFromInt64(100), types.NewCoinFromInt64(100)},
		{types.NewCoinFromInt64(1), types.NewCoinFromInt64(101)},
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
		{types.NewCoinFromInt64(100), types.NewCoinFromInt64(100)},
		{types.NewCoinFromInt64(1), types.NewCoinFromInt64(101)},
	}

	for _, cs := range cases {
		err := gm.AddConsumption(ctx, cs.coin)
		assert.Nil(t, err)
		globalMeta, err := gm.storage.GetGlobalMeta(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.expect, globalMeta.CumulativeConsumption)
	}
}
