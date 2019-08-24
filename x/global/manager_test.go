package global

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global/model"
)

const (
	eventTypeTestEvent = "1"
)

type testEvent struct{}

// Construct some global addrs and txs for tests.
var (
	TestGlobalKVStoreKey = sdk.NewKVStoreKey("global")
	TestParamKVStoreKey  = sdk.NewKVStoreKey("param")
	totalLino            = types.NewCoinFromInt64(10000 * types.Decimals)
)

func InitGlobalManager(ctx sdk.Context, gm GlobalManager) error {
	return gm.InitGlobalManager(ctx, totalLino)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, false, log.NewNopLogger())
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
	baseTime := time.Now()
	var initMaxTPS = sdk.NewDec(1000)

	testCases := []struct {
		testName            string
		baseTime            int64
		nextTime            time.Time
		numOfTx             int64
		expectCurrentTPS    sdk.Dec
		expectMaxTPS        sdk.Dec
		expectCapacityRatio sdk.Dec
	}{
		{
			testName:            "0 tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime,
			numOfTx:             0,
			expectCurrentTPS:    sdk.ZeroDec(),
			expectMaxTPS:        initMaxTPS,
			expectCapacityRatio: sdk.ZeroDec(),
		},
		{
			testName:            "2/2 got 1 current tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime.Add(time.Duration(2) * time.Second),
			numOfTx:             2,
			expectCurrentTPS:    sdk.OneDec(),
			expectMaxTPS:        initMaxTPS,
			expectCapacityRatio: types.NewDecFromRat(1, 1000),
		},
		{
			testName:            "1000/1 got max tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime.Add(time.Duration(1) * time.Second),
			numOfTx:             1000,
			expectCurrentTPS:    initMaxTPS,
			expectMaxTPS:        initMaxTPS,
			expectCapacityRatio: sdk.OneDec(),
		},
		{
			testName:            "2000/2 got max tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime.Add(time.Duration(2) * time.Second),
			numOfTx:             2000,
			expectCurrentTPS:    initMaxTPS,
			expectMaxTPS:        initMaxTPS,
			expectCapacityRatio: sdk.OneDec(),
		},
		{
			testName:            "3000/2 got 1500 current tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime.Add(time.Duration(2) * time.Second),
			numOfTx:             3000,
			expectCurrentTPS:    sdk.NewDec(1500),
			expectMaxTPS:        sdk.NewDec(1500),
			expectCapacityRatio: sdk.OneDec(),
		},
		{
			testName:            "2000/2 got 1000 current tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime.Add(time.Duration(2) * time.Second),
			numOfTx:             2000,
			expectCurrentTPS:    sdk.NewDec(1000),
			expectMaxTPS:        sdk.NewDec(1500),
			expectCapacityRatio: types.NewDecFromRat(2, 3),
		},
	}
	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: tc.nextTime, NumTxs: tc.numOfTx})
		err := gm.SetLastBlockTime(ctx, tc.baseTime)
		if err != nil {
			t.Errorf("%s: failed to set last block time, got err %v", tc.testName, err)
		}
		err = gm.UpdateTPS(ctx)
		if err != nil {
			t.Errorf("%s: failed to update TPS, got err %v", tc.testName, err)
		}

		storage := model.NewGlobalStorage(TestGlobalKVStoreKey)
		tps, _ := storage.GetTPS(ctx)
		if !tc.expectCurrentTPS.Equal(tps.CurrentTPS) {
			t.Errorf("%s: diff current tps, got %v, want %v", tc.testName, tps.CurrentTPS, tc.expectCurrentTPS)
		}
		if !tc.expectMaxTPS.Equal(tps.MaxTPS) {
			t.Errorf("%s: diff max tps, got %v, want %v", tc.testName, tps.MaxTPS, tc.expectMaxTPS)
		}

		ratio, err := gm.GetTPSCapacityRatio(ctx)
		if err != nil {
			t.Errorf("%s: failed to get TPS capacity ratio, got err %v", tc.testName, err)
		}
		if !tc.expectCapacityRatio.Equal(ratio) {
			t.Errorf("%s: diff ratio, got %v, want %v", tc.testName, ratio, tc.expectCapacityRatio)
		}
	}
}

func TestAddFrictionAndRegisterContentRewardEvent(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := ctx.BlockHeader().Time.Unix()
	testCases := []struct {
		testName              string
		frictionCoin          types.Coin
		evaluate              types.MiniDollar
		registerBaseTime      int64
		expectCoinInStatistic types.Coin
		expectInWindow        types.MiniDollar
	}{
		{
			testName:              "add 1 friction",
			frictionCoin:          types.NewCoinFromInt64(1),
			evaluate:              types.NewMiniDollar(1),
			registerBaseTime:      baseTime,
			expectCoinInStatistic: types.NewCoinFromInt64(1),
			expectInWindow:        types.NewMiniDollar(1),
		},
		{
			testName:              "add 100 more friction",
			frictionCoin:          types.NewCoinFromInt64(100),
			evaluate:              types.NewMiniDollar(1),
			registerBaseTime:      baseTime + 100,
			expectCoinInStatistic: types.NewCoinFromInt64(101),
			expectInWindow:        types.NewMiniDollar(2),
		},
		{
			testName:              "add 1 more friction",
			frictionCoin:          types.NewCoinFromInt64(1),
			evaluate:              types.NewMiniDollar(100),
			registerBaseTime:      baseTime + 1001,
			expectCoinInStatistic: types.NewCoinFromInt64(102),
			expectInWindow:        types.NewMiniDollar(102),
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: time.Unix(tc.registerBaseTime, 0)})
		err := gm.AddFrictionAndRegisterContentRewardEvent(
			ctx, testEvent{}, tc.frictionCoin, tc.evaluate)
		if err != nil {
			t.Errorf("%s: failed to add friction and register event, got err %v", tc.testName, err)
		}

		consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
		if err != nil {
			t.Errorf("%s: failed to get consumption meta, got err %v", tc.testName, err)
		}
		if !consumptionMeta.ConsumptionRewardPool.IsZero() {
			t.Errorf("%s: diff consumption reward pool, got %v, want zero",
				tc.testName, consumptionMeta.ConsumptionRewardPool)
		}
		if !consumptionMeta.ConsumptionWindow.Equal(tc.expectInWindow.Int) {
			t.Errorf("%s: diff consumption window, got %v, want %v", tc.testName,
				consumptionMeta.ConsumptionWindow, tc.expectInWindow)
		}
		pastDay, err := gm.GetPastDay(ctx, ctx.BlockHeader().Time.Unix())
		if err != nil {
			t.Errorf("%s: failed to get past day, got err %v", tc.testName, err)
		}
		linoStakeStatistic, err := gm.storage.GetLinoStakeStat(ctx, pastDay)
		if err != nil {
			t.Errorf("%s: failed to get lino power statistic, got err %v", tc.testName, err)
		}
		if !linoStakeStatistic.TotalConsumptionFriction.IsEqual(tc.expectCoinInStatistic) {
			t.Errorf("%s: diff total consumption friction, got %v, want %v", tc.testName,
				linoStakeStatistic.TotalConsumptionFriction, tc.expectCoinInStatistic)
		}

		if err = gm.CommitEventCache(ctx); err != nil {
			t.Errorf("%s: failed to commit event cache, %v", tc.testName, err)
		}

		timeEventList := gm.GetTimeEventListAtTime(ctx, tc.registerBaseTime+24*7*3600)
		if !assert.Equal(t, types.TimeEventList{Events: []types.Event{testEvent{}}}, *timeEventList) {
			t.Errorf("%s: diff event list, got %v, want %v", tc.testName,
				*timeEventList, types.TimeEventList{Events: []types.Event{testEvent{}}})
		}
	}
}

func TestGetRewardAndPopFromWindow(t *testing.T) {
	ctx, gm := setupTest(t)
	testCases := []struct {
		testName                    string
		evaluate                    types.MiniDollar
		expectReward                types.Coin
		initConsumptionRewardPool   types.Coin
		initConsumptionWindow       types.MiniDollar
		expectConsumptionRewardPool types.Coin
		expectConsumptionWindow     types.MiniDollar
	}{
		{
			testName:                    "1 evaluate, 0 penalty",
			evaluate:                    types.NewMiniDollar(1),
			expectReward:                types.NewCoinFromInt64(100),
			initConsumptionRewardPool:   types.NewCoinFromInt64(1000),
			initConsumptionWindow:       types.NewMiniDollar(10),
			expectConsumptionRewardPool: types.NewCoinFromInt64(900),
			expectConsumptionWindow:     types.NewMiniDollar(9),
		},
		{
			// issue https://github.com/lino-network/lino/issues/150
			testName:                    "test large number",
			evaluate:                    types.NewMiniDollar(77777777777777),
			expectReward:                types.NewCoinFromInt64(23333333),
			initConsumptionRewardPool:   types.NewCoinFromInt64(100000000),
			initConsumptionWindow:       types.NewMiniDollar(333333333333333),
			expectConsumptionRewardPool: types.NewCoinFromInt64(76666667),
			expectConsumptionWindow:     types.NewMiniDollar(255555555555556),
		},
	}

	for _, tc := range testCases {
		globalMeta, err := gm.storage.GetGlobalMeta(ctx)
		if err != nil {
			t.Errorf("%s: failed to get global meta, got err %v", tc.testName, err)
			return
		}
		totalLino := globalMeta.TotalLinoCoin
		consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
		if err != nil {
			t.Errorf("%s: failed to get consumption meta, got err %v", tc.testName, err)
			return
		}

		consumptionMeta.ConsumptionRewardPool = tc.initConsumptionRewardPool
		consumptionMeta.ConsumptionWindow = tc.initConsumptionWindow

		err = gm.storage.SetConsumptionMeta(ctx, consumptionMeta)
		if err != nil {
			t.Errorf("%s: failed to set consumption meta, got err %v", tc.testName, err)
			return
		}

		reward, err := gm.GetRewardAndPopFromWindow(ctx, tc.evaluate)
		if err != nil {
			t.Errorf("%s: failed to get reward and pop from window, got err %v", tc.testName, err)
			return
		}
		if !reward.IsEqual(tc.expectReward) {
			t.Errorf("%s: diff reward, got %v, want %v", tc.testName, reward, tc.expectReward)
			return
		}

		consumptionMeta, err = gm.storage.GetConsumptionMeta(ctx)
		if err != nil {
			t.Errorf("%s: failed to get consumption meta again, got err %v", tc.testName, err)
			return
		}
		if !consumptionMeta.ConsumptionRewardPool.IsEqual(tc.expectConsumptionRewardPool) {
			t.Errorf("%s: diff consumption reward pool, got %v, want %v", tc.testName,
				consumptionMeta.ConsumptionRewardPool, tc.expectConsumptionRewardPool)
			return
		}
		if !consumptionMeta.ConsumptionWindow.Equal(tc.expectConsumptionWindow.Int) {
			t.Errorf("%s: diff consumption window, got %v, want %v", tc.testName,
				consumptionMeta.ConsumptionWindow, tc.expectConsumptionWindow)
			return
		}
		globalMeta, err = gm.storage.GetGlobalMeta(ctx)
		if err != nil {
			t.Errorf("%s: failed to get global meta again, got err %v", tc.testName, err)
			return
		}
		if !globalMeta.TotalLinoCoin.IsEqual(totalLino.Plus(tc.expectReward)) {
			t.Errorf(
				"%s: total lino incorrect, expect %v, got %v",
				tc.testName, totalLino.Plus(tc.expectReward), globalMeta.TotalLinoCoin)
			return
		}
	}
}

func TestTimeEventList(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := ctx.BlockHeader().Time.Unix()
	regCases := []struct {
		testName        string
		registerAtTime  int64
		expectResult    sdk.Error
		expectEventList *types.TimeEventList
	}{
		{
			testName:        "register one event",
			registerAtTime:  baseTime,
			expectResult:    nil,
			expectEventList: &types.TimeEventList{Events: []types.Event{testEvent{}}},
		},
		{
			testName:        "register two events",
			registerAtTime:  baseTime,
			expectResult:    nil,
			expectEventList: &types.TimeEventList{Events: []types.Event{testEvent{}, testEvent{}}}},
		{
			testName:        "can't register expired event",
			registerAtTime:  baseTime - 1,
			expectResult:    ErrRegisterExpiredEvent(baseTime - 1),
			expectEventList: nil,
		},
		{
			testName:        "register one event again",
			registerAtTime:  baseTime + 1,
			expectResult:    nil,
			expectEventList: &types.TimeEventList{Events: []types.Event{testEvent{}}},
		},
	}

	for _, tc := range regCases {
		err := gm.registerEventAtTime(ctx, tc.registerAtTime, testEvent{})
		if !assert.Equal(t, tc.expectResult, err) {
			t.Errorf("%s: diff err result, got %v, want %v", tc.testName, err, tc.expectResult)
		}

		if err = gm.CommitEventCache(ctx); err != nil {
			t.Errorf("%s: failed to commit event cache, %v", tc.testName, err)
		}

		eventList := gm.GetTimeEventListAtTime(ctx, tc.registerAtTime)
		if !assert.Equal(t, tc.expectEventList, eventList) {
			t.Errorf("%s: diff event list, got %v, want %v", tc.testName, eventList, tc.expectEventList)
		}
	}

	rmCases := []struct {
		testName        string
		removeAtTime    int64
		expectEventList *types.TimeEventList
	}{
		{
			testName:        "remove event",
			removeAtTime:    baseTime,
			expectEventList: nil,
		},
		{
			testName:        "remove expired event",
			removeAtTime:    baseTime - 1,
			expectEventList: nil,
		},
		{
			testName:        "remove future event",
			removeAtTime:    baseTime + 1,
			expectEventList: nil,
		},
	}

	for _, tc := range rmCases {
		err := gm.RemoveTimeEventList(ctx, tc.removeAtTime)
		if err != nil {
			t.Errorf("%s: failed to remove time event list, got err %v", tc.testName, err)
		}
		eventList := gm.GetTimeEventListAtTime(ctx, tc.removeAtTime)
		if !assert.Equal(t, tc.expectEventList, eventList) {
			t.Errorf("%s: diff event list, got %v, want %v", tc.testName, eventList, tc.expectEventList)
		}
	}
}

func TestRegisterCoinReturnEvent(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := ctx.BlockHeader().Time.Unix()
	testCases := []struct {
		testName               string
		registerAtTime         int64
		times                  int64
		interval               int64
		expectTimeWithTwoEvent []int64
		expectTimeWithOneEvent []int64
	}{
		{
			testName:               "one event with 5 times",
			registerAtTime:         baseTime,
			times:                  5,
			interval:               10 * 3600,
			expectTimeWithTwoEvent: []int64{},
			expectTimeWithOneEvent: []int64{
				baseTime + 10*3600,
				baseTime + 20*3600,
				baseTime + 30*3600,
				baseTime + 40*3600,
				baseTime + 50*3600,
			},
		},
		{
			testName:       "two event with 2 times and one event with 3 times",
			registerAtTime: baseTime,
			times:          2,
			interval:       10 * 3600,
			expectTimeWithTwoEvent: []int64{
				baseTime + 10*3600,
				baseTime + 20*3600,
			},
			expectTimeWithOneEvent: []int64{
				baseTime + 30*3600,
				baseTime + 40*3600,
				baseTime + 50*3600,
			},
		},
		{
			testName:       "two event with 4 times and one event with 3 times",
			registerAtTime: baseTime + 20*3600,
			times:          4,
			interval:       5 * 3600,
			expectTimeWithTwoEvent: []int64{
				baseTime + 10*3600,
				baseTime + 20*3600,
				baseTime + 30*3600,
				baseTime + 40*3600,
			},
			expectTimeWithOneEvent: []int64{
				baseTime + 25*3600,
				baseTime + 35*3600,
				baseTime + 50*3600,
			},
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: time.Unix(tc.registerAtTime, 0)})
		events := []types.Event{}
		for i := int64(0); i < tc.times; i++ {
			events = append(events, testEvent{})
		}
		err := gm.RegisterCoinReturnEvent(ctx, events, tc.times, tc.interval)
		if err != nil {
			t.Errorf("%s: failed to register coin return event, got err %v", tc.testName, err)
		}

		if err = gm.CommitEventCache(ctx); err != nil {
			t.Errorf("%s: failed to commit event cache, %v", tc.testName, err)
		}

		for _, time := range tc.expectTimeWithOneEvent {
			eventList := gm.GetTimeEventListAtTime(ctx, time)
			if len(eventList.Events) != 1 {
				t.Errorf("%s: diff time one event, got %v, want 1", tc.testName, len(eventList.Events))
			}
		}
		for _, time := range tc.expectTimeWithTwoEvent {
			eventList := gm.GetTimeEventListAtTime(ctx, time)
			if len(eventList.Events) != 2 {
				t.Errorf("%s: diff time one event, got %v, want 1", tc.testName, len(eventList.Events))
			}
		}
	}
}

func TestDistributeHourlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, globalMeta.TotalLinoCoin, totalLino)
	assert.Equal(t, globalMeta.TotalLinoCoin, globalMeta.LastYearTotalLinoCoin)
	globalAllocationParam, _ := gm.paramHolder.GetGlobalAllocationParam(ctx)
	lastYearTotalLino := globalMeta.LastYearTotalLinoCoin
	expectContentCreatorInflation := types.NewCoinFromInt64(0)
	expectValidatorInflation := types.NewCoinFromInt64(0)
	expectDeveloperInflation := types.NewCoinFromInt64(0)
	expectInfraInflation := types.NewCoinFromInt64(0)

	globalAllocation, err := gm.paramHolder.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)
	for i := 0; i < types.HoursPerYear; i++ {
		err = gm.DistributeHourlyInflation(ctx)
		assert.Nil(t, err)

		consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
		assert.Nil(t, err)

		inflationPool, err := gm.storage.GetInflationPool(ctx)
		assert.Nil(t, err)

		hourlyInflation :=
			types.DecToCoin(lastYearTotalLino.ToDec().
				Mul(globalAllocationParam.GlobalGrowthRate).
				Mul(types.NewDecFromRat(1, int64(types.HoursPerYear))))
		expectContentCreatorInflation =
			expectContentCreatorInflation.Plus(
				types.DecToCoin(hourlyInflation.ToDec().Mul(globalAllocation.ContentCreatorAllocation)))
		expectValidatorInflation =
			expectValidatorInflation.Plus(
				types.DecToCoin(hourlyInflation.ToDec().Mul(globalAllocation.ValidatorAllocation)))
		expectDeveloperInflation =
			expectDeveloperInflation.Plus(
				types.DecToCoin(hourlyInflation.ToDec().Mul(globalAllocation.DeveloperAllocation)))
		expectInfraInflation =
			expectInfraInflation.Plus(
				types.DecToCoin(hourlyInflation.ToDec().Mul(globalAllocation.InfraAllocation)))
		assert.True(t, expectContentCreatorInflation.IsEqual(consumptionMeta.ConsumptionRewardPool))
		assert.True(t, expectInfraInflation.IsEqual(inflationPool.InfraInflationPool))
		assert.True(t, expectDeveloperInflation.IsEqual(inflationPool.DeveloperInflationPool))
		assert.True(t, expectValidatorInflation.IsEqual(inflationPool.ValidatorInflationPool))
	}
	globalMeta, err = gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, globalMeta.TotalLinoCoin, totalLino)
}

func AddToDeveloperInflationPool(t *testing.T) {
	ctx, gm := setupTest(t)
	testCases := []struct {
		testName                   string
		initDeveloperInflationPool types.Coin
		addAmount                  types.Coin
	}{
		{
			testName:                   "add to empty pool",
			initDeveloperInflationPool: types.NewCoinFromInt64(0 * types.Decimals),
			addAmount:                  types.NewCoinFromInt64(1 * types.Decimals),
		},
		{
			testName:                   "normal add operation",
			initDeveloperInflationPool: types.NewCoinFromInt64(10000 * types.Decimals),
			addAmount:                  types.NewCoinFromInt64(10 * types.Decimals),
		},
	}
	for _, tc := range testCases {
		inflationPool := &model.InflationPool{
			DeveloperInflationPool: tc.initDeveloperInflationPool,
		}
		err := gm.storage.SetInflationPool(ctx, inflationPool)
		if err != nil {
			t.Errorf("%s: failed to set inflation pool, got err %v", tc.testName, err)
			return
		}
		err = gm.AddToDeveloperInflationPool(ctx, tc.addAmount)
		if err != nil {
			t.Errorf("%s: failed to add to developer inflation pool, got err %v", tc.testName, err)
			return
		}
		inflationPool, err = gm.storage.GetInflationPool(ctx)
		if err != nil {
			t.Errorf("%s: failed to get inflation pool, got err %v", tc.testName, err)
			return
		}
		if !tc.initDeveloperInflationPool.Plus(tc.addAmount).IsEqual(inflationPool.DeveloperInflationPool) {
			t.Errorf(
				"%s: failed to add inflation to developer inflation pool, expect %v, got %v",
				tc.testName, tc.initDeveloperInflationPool.Plus(tc.addAmount), inflationPool.DeveloperInflationPool)
			return
		}
	}
}

// func TestRecalculateAnnuallyInflation(t *testing.T) {
// 	ctx, gm := setupTest(t)
// 	totalLino := types.NewCoinFromInt64(10000000000 * types.Decimals)
// 	ceiling := types.NewDecFromRat(98, 1000)
// 	floor := types.NewDecFromRat(30, 1000)

// 	testCases := []struct {
// 		testName            string
// 		lastYearConsumption types.Coin
// 		thisYearConsumption types.Coin
// 		expectGrowthRate    sdk.Dec
// 	}{
// 		{
// 			testName:            "same consumption of last year and this year, get floor growth rate",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			expectGrowthRate:    floor,
// 		},
// 		{
// 			testName:            "right equal to floor growth rate",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(103000000 * types.Decimals),
// 			expectGrowthRate:    floor,
// 		},
// 		{
// 			testName:            "right equal to ceiling growth rate",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(1098000000 * types.Decimals),
// 			expectGrowthRate:    ceiling,
// 		},
// 		{
// 			testName:            "bigger than ceiling will be rounded to ceiling growth rate",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(1099000000 * types.Decimals),
// 			expectGrowthRate:    ceiling,
// 		},
// 		{
// 			testName:            "less year consumption will use floor growth rate",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(90000000 * types.Decimals),
// 			expectGrowthRate:    floor,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		globalMeta := &model.GlobalMeta{
// 			TotalLinoCoin:                 totalLino,
// 			LastYearTotalLinoCoin:         totalLino,
// 		}
// 		err := gm.storage.SetGlobalMeta(ctx, globalMeta)
// 		if err != nil {
// 			t.Errorf("%s: failed to set global meta, got err %v", tc.testName, err)
// 		}

// 		err = gm.SetTotalLinoAndRecalculateGrowthRate(ctx)
// 		if err != nil {
// 			t.Errorf("%s: failed to recalculate annually inflation, got err %v", tc.testName, err)
// 		}

// 		globalMeta, err = gm.storage.GetGlobalMeta(ctx)
// 		if err != nil {
// 			t.Errorf("%s: failed to get global meta, got err %v", tc.testName, err)
// 		}
// 		if !globalMeta.LastYearTotalLinoCoin.IsEqual(globalMeta.TotalLinoCoin) {
// 			t.Errorf("%s: diff lino coin, got %v, want %v", tc.testName,
// 				globalMeta.LastYearTotalLinoCoin, globalMeta.TotalLinoCoin)
// 		}
// 		if !globalMeta.LastYearCumulativeConsumption.IsEqual(tc.thisYearConsumption) {
// 			t.Errorf("%s: diff last year cumulative consumption, got %v, want %v", tc.testName,
// 				globalMeta.LastYearCumulativeConsumption, tc.thisYearConsumption)
// 		}
// 		if !globalMeta.CumulativeConsumption.IsEqual(types.NewCoinFromInt64(0)) {
// 			t.Errorf("%s: diff cumulative consumption, got %v, want %v", tc.testName,
// 				globalMeta.CumulativeConsumption, types.NewCoinFromInt64(0))
// 		}
// 	}
// }

// func TestGetGrowthRate(t *testing.T) {
// 	ctx, gm := setupTest(t)
// 	totalLino := types.NewCoinFromInt64(1000000)
// 	ceiling := types.NewDecFromRat(98, 1000)
// 	floor := types.NewDecFromRat(30, 1000)
// 	bigLastYearConsumption, _ := new(big.Int).SetString("77777777777777777777", 10)
// 	bigThisYearConsumption, _ := new(big.Int).SetString("83333333333333333332", 10)
// 	bigLastYearConsumptionCoin := types.NewCoinFromBigInt(bigLastYearConsumption)
// 	bigThisYearConsumptionCoin := types.NewCoinFromBigInt(bigThisYearConsumption)
// 	bigGrowth := bigThisYearConsumptionCoin.Minus(bigLastYearConsumptionCoin)
// 	bigExpectedGrowth := sdk.NewDecFromBigInt(bigGrowth.Amount.BigInt()).Quo(sdk.NewDecFromBigInt(bigLastYearConsumptionCoin.Amount.BigInt()))

// 	testCases := []struct {
// 		testName            string
// 		lastYearConsumption types.Coin
// 		thisYearConsumption types.Coin
// 		lastYearGrowthRate  sdk.Dec
// 		expectGrowthRate    sdk.Dec
// 	}{
// 		{
// 			testName:            "floor growth rate, 0 thisYearConsumption",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(0 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    floor,
// 		},
// 		{
// 			testName:            "ceiling growth rate, 0 lastYearConsumption",
// 			lastYearConsumption: types.NewCoinFromInt64(0 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    types.NewDecFromRat(98, 1000),
// 		},
// 		{
// 			testName:            "floor growth rate, thisYearConsumption = lastYearConsumption",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    floor,
// 		},
// 		{
// 			testName:            "ceiling growth rate, 0 lastYearConsumption",
// 			lastYearConsumption: types.NewCoinFromInt64(0),
// 			thisYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    types.NewDecFromRat(98, 1000),
// 		},
// 		{
// 			testName:            "less than floor will be rounded to floor growth rate 1",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(100010000 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    floor,
// 		},
// 		{
// 			testName:            "less than floor will be rounded to floor growth rate 2",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(102900000 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    floor,
// 		},
// 		{
// 			testName:            "less than floor will be rounded to floor growth rate 3",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(103000000 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    floor,
// 		},
// 		{
// 			testName:            "growth rate between floor and ceiling 1",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(103100000 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    types.NewDecFromRat(31, 1000),
// 		},
// 		{
// 			testName:            "right equal to ceiling growth rate",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(109800000 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    ceiling,
// 		},
// 		{
// 			testName:            "higher than ceiling will be rouned to ceiling growth rate",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(109900000 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    ceiling,
// 		},
// 		{
// 			testName:            "growth rate between floor and ceiling 2",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(109700000 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    types.NewDecFromRat(97, 1000),
// 		},
// 		{
// 			testName:            "growth rate between floor and ceiling 3",
// 			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
// 			thisYearConsumption: types.NewCoinFromInt64(104700000 * types.Decimals),
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    types.NewDecFromRat(47, 1000),
// 		},
// 		// issue https://github.com/lino-network/lino/issues/150
// 		{
// 			testName:            "overflow testing",
// 			lastYearConsumption: bigLastYearConsumptionCoin,
// 			thisYearConsumption: bigThisYearConsumptionCoin,
// 			lastYearGrowthRate:  types.NewDecFromRat(98, 1000),
// 			expectGrowthRate:    bigExpectedGrowth,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		globalMeta := &model.GlobalMeta{
// 			TotalLinoCoin:                 totalLino,
// 			LastYearTotalLinoCoin:         totalLino,
// 			LastYearCumulativeConsumption: tc.lastYearConsumption,
// 			CumulativeConsumption:         tc.thisYearConsumption,
// 		}
// 		err := gm.storage.SetGlobalMeta(ctx, globalMeta)
// 		if err != nil {
// 			t.Errorf("%s: failed to set global meta, got err %v", tc.testName, err)
// 		}
// 		err = gm.paramHolder.UpdateGlobalGrowthRate(ctx, tc.lastYearGrowthRate)
// 		if err != nil {
// 			t.Errorf("%s: failed to set global growth rate, got err %v", tc.testName, err)
// 		}

// 		err = gm.SetTotalLinoAndRecalculateGrowthRate(ctx)
// 		if err != nil {
// 			t.Errorf("%s: failed to get growth rate, got err %v", tc.testName, err)
// 		}

// 		globalMeta, err = gm.storage.GetGlobalMeta(ctx)
// 		if err != nil {
// 			t.Errorf("%s: failed to get global meta, got err %v", tc.testName, err)
// 		}
// 		if !globalMeta.LastYearCumulativeConsumption.IsEqual(tc.thisYearConsumption) {
// 			t.Errorf("%s: diff last year cumulative consumption, got %v, want %v", tc.testName,
// 				globalMeta.LastYearCumulativeConsumption, tc.thisYearConsumption)
// 		}
// 		if !globalMeta.CumulativeConsumption.IsEqual(types.NewCoinFromInt64(0)) {
// 			t.Errorf("%s: diff cumulative consumption, got %v, want %v", tc.testName,
// 				globalMeta.CumulativeConsumption, types.NewCoinFromInt64(0))
// 		}
// 		globalParam, _ := gm.paramHolder.GetGlobalAllocationParam(ctx)
// 		if !tc.expectGrowthRate.Equal(globalParam.GlobalGrowthRate) {
// 			t.Errorf("%s: diff growth rate, got %v, want %v", tc.testName,
// 				globalParam.GlobalGrowthRate, tc.expectGrowthRate)
// 		}
// 	}
// }

func TestGetValidatorHourlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalValidatorInflation := types.NewCoinFromInt64(10000 * 100)
	inflationPool := &model.InflationPool{
		ValidatorInflationPool: totalValidatorInflation,
	}
	err := gm.storage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)
	coin, err := gm.GetValidatorHourlyInflation(ctx)
	assert.Nil(t, err)
	assert.Equal(t, totalValidatorInflation, coin)

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
	coin, err := gm.GetInfraMonthlyInflation(ctx)
	assert.Nil(t, err)
	assert.Equal(t, totalInfraInflation, coin)
	pool, err := gm.storage.GetInflationPool(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoinFromInt64(0), pool.InfraInflationPool)
	globalMeta, err := gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, globalMeta.TotalLinoCoin, types.NewCoinFromInt64(10000*types.Decimals).Plus(totalInfraInflation))
}

func TestPopDeveloperMonthlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalDeveloperInflation := types.NewCoinFromInt64(10000 * 100)
	inflationPool := &model.InflationPool{
		DeveloperInflationPool: totalDeveloperInflation,
	}
	err := gm.storage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)
	coin, err := gm.PopDeveloperMonthlyInflation(ctx)
	assert.Nil(t, err)
	assert.Equal(t, totalDeveloperInflation, coin)
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

	testCases := []struct {
		testName string
		coin     types.Coin
		expect   types.Coin
	}{
		{
			testName: "add 100 inflation",
			coin:     types.NewCoinFromInt64(100),
			expect:   types.NewCoinFromInt64(100),
		},
		{
			testName: "add 1 more inflation",
			coin:     types.NewCoinFromInt64(1),
			expect:   types.NewCoinFromInt64(101),
		},
	}

	for _, tc := range testCases {
		err := gm.AddToValidatorInflationPool(ctx, tc.coin)
		if err != nil {
			t.Errorf("%s: failed to add validator inflation pool, got err %v", tc.testName, err)
		}
		pool, err := gm.storage.GetInflationPool(ctx)
		if err != nil {
			t.Errorf("%s: failed to get inflation pool, got err %v", tc.testName, err)
		}
		if !pool.ValidatorInflationPool.IsEqual(tc.expect) {
			t.Errorf("%s: diff validator inflation pool, got %v, want %v", tc.testName,
				pool.ValidatorInflationPool, tc.expect)
		}
	}
}

func TestChainStartTime(t *testing.T) {
	ctx, gm := setupTest(t)

	testCases := []struct {
		testName  string
		startTime int64
	}{
		{
			testName:  "set start time to zero",
			startTime: 0,
		},
		{
			testName:  "normal case",
			startTime: time.Now().Unix(),
		},
	}

	for _, tc := range testCases {
		err := gm.SetChainStartTime(ctx, tc.startTime)
		if err != nil {
			t.Errorf("%s: failed to set chain start time, got err %v", tc.testName, err)
		}
		chainStartTime, err := gm.GetChainStartTime(ctx)
		if err != nil {
			t.Errorf("%s: failed to get chain start time, got err %v", tc.testName, err)
		}
		if chainStartTime != tc.startTime {
			t.Errorf("%s: diff chain start time, got %v, want %v", tc.testName, chainStartTime, tc.startTime)
			return
		}
	}
}

func TestPastMinutes(t *testing.T) {
	ctx, gm := setupTest(t)

	testCases := []struct {
		testName    string
		pastMinutes int64
	}{
		{
			testName:    "set past minutes to zero",
			pastMinutes: 0,
		},
		{
			testName:    "normal case",
			pastMinutes: time.Now().Unix() / 60,
		},
	}

	for _, tc := range testCases {
		err := gm.SetPastMinutes(ctx, tc.pastMinutes)
		if err != nil {
			t.Errorf("%s: failed to set past minutes, got err %v", tc.testName, err)
		}
		pastMinutes, err := gm.GetPastMinutes(ctx)
		if err != nil {
			t.Errorf("%s: failed to get past minutes, got err %v", tc.testName, err)
		}
		if pastMinutes != tc.pastMinutes {
			t.Errorf("%s: diff past minutes, got %v, want %v", tc.testName, pastMinutes, tc.pastMinutes)
			return
		}
	}
}

func TestGetConsumptionFrictionRate(t *testing.T) {
	ctx, gm := setupTest(t)

	testCases := []struct {
		testName        string
		consumptionMeta model.ConsumptionMeta
	}{
		{
			testName: "normal friction rate",
			consumptionMeta: model.ConsumptionMeta{
				ConsumptionFrictionRate: types.NewDecFromRat(5, 100),
			},
		},
		{
			testName: "10% friction rate",
			consumptionMeta: model.ConsumptionMeta{
				ConsumptionFrictionRate: types.NewDecFromRat(1, 10),
			},
		},
	}

	for _, tc := range testCases {
		err := gm.storage.SetConsumptionMeta(ctx, &tc.consumptionMeta)
		assert.Nil(t, err)
		frictionRate, err := gm.GetConsumptionFrictionRate(ctx)
		assert.Nil(t, err)

		if !frictionRate.Equal(tc.consumptionMeta.ConsumptionFrictionRate) {
			t.Errorf("%s: diff friction rate, got %v, want %v",
				tc.testName, frictionRate, tc.consumptionMeta.ConsumptionFrictionRate)
			return
		}
	}
}

func TestAddLinoStakeToStat(t *testing.T) {
	ctx, gm := setupTest(t)

	testCases := []struct {
		testName string
		day      int64
		stat     model.LinoStakeStat
		addStake types.Coin
	}{
		{
			testName: "add 10 stake to empty stat",
			day:      0,
			stat: model.LinoStakeStat{
				TotalConsumptionFriction: types.NewCoinFromInt64(0),
				TotalLinoStake:           types.NewCoinFromInt64(0),
				UnclaimedFriction:        types.NewCoinFromInt64(0),
				UnclaimedLinoStake:       types.NewCoinFromInt64(0),
			},
			addStake: types.NewCoinFromInt64(10 * types.Decimals),
		},
		{
			testName: "add 10 stake to stat with 20 stake",
			day:      0,
			stat: model.LinoStakeStat{
				TotalConsumptionFriction: types.NewCoinFromInt64(0),
				TotalLinoStake:           types.NewCoinFromInt64(20 * types.Decimals),
				UnclaimedFriction:        types.NewCoinFromInt64(0),
				UnclaimedLinoStake:       types.NewCoinFromInt64(20 * types.Decimals),
			},
			addStake: types.NewCoinFromInt64(10 * types.Decimals),
		},
	}

	for _, tc := range testCases {
		err := gm.storage.SetLinoStakeStat(ctx, tc.day, &tc.stat)
		assert.Nil(t, err)
		err = gm.AddLinoStakeToStat(ctx, tc.addStake)
		assert.Nil(t, err)

		stats, err := gm.storage.GetLinoStakeStat(ctx, tc.day)
		assert.Nil(t, err)
		if !stats.TotalLinoStake.IsEqual(
			tc.stat.TotalLinoStake.Plus(tc.addStake)) {
			t.Errorf("%s: diff friction rate, got %v, want %v",
				tc.testName, stats.TotalLinoStake,
				tc.stat.TotalLinoStake.Plus(tc.addStake))
			return
		}
		if !stats.UnclaimedLinoStake.IsEqual(
			tc.stat.UnclaimedLinoStake.Plus(tc.addStake)) {
			t.Errorf("%s: diff friction rate, got %v, want %v",
				tc.testName, stats.UnclaimedLinoStake,
				tc.stat.UnclaimedLinoStake.Plus(tc.addStake))
			return
		}
	}
}

func TestMinusLinoStakeFromStat(t *testing.T) {
	ctx, gm := setupTest(t)

	testCases := []struct {
		testName   string
		day        int64
		stat       model.LinoStakeStat
		minusStake types.Coin
	}{
		{
			testName: "minus 10 stake to from stat with 10 stake",
			day:      0,
			stat: model.LinoStakeStat{
				TotalConsumptionFriction: types.NewCoinFromInt64(0),
				TotalLinoStake:           types.NewCoinFromInt64(0),
				UnclaimedFriction:        types.NewCoinFromInt64(0),
				UnclaimedLinoStake:       types.NewCoinFromInt64(0),
			},
			minusStake: types.NewCoinFromInt64(10 * types.Decimals),
		},
		{
			testName: "minus 10 stake to stat with more than 10 stake",
			day:      0,
			stat: model.LinoStakeStat{
				TotalConsumptionFriction: types.NewCoinFromInt64(0),
				TotalLinoStake:           types.NewCoinFromInt64(10000 * types.Decimals),
				UnclaimedFriction:        types.NewCoinFromInt64(0),
				UnclaimedLinoStake:       types.NewCoinFromInt64(10000 * types.Decimals),
			},
			minusStake: types.NewCoinFromInt64(10 * types.Decimals),
		},
	}

	for _, tc := range testCases {
		err := gm.storage.SetLinoStakeStat(ctx, tc.day, &tc.stat)
		assert.Nil(t, err)
		err = gm.MinusLinoStakeFromStat(ctx, tc.minusStake)
		assert.Nil(t, err)

		stats, err := gm.storage.GetLinoStakeStat(ctx, tc.day)
		assert.Nil(t, err)
		if !stats.TotalLinoStake.IsEqual(
			tc.stat.TotalLinoStake.Minus(tc.minusStake)) {
			t.Errorf("%s: diff friction rate, got %v, want %v",
				tc.testName, stats.TotalLinoStake,
				tc.stat.TotalLinoStake.Minus(tc.minusStake))
			return
		}
		if !stats.UnclaimedLinoStake.IsEqual(
			tc.stat.UnclaimedLinoStake.Minus(tc.minusStake)) {
			t.Errorf("%s: diff friction rate, got %v, want %v",
				tc.testName, stats.UnclaimedLinoStake,
				tc.stat.UnclaimedLinoStake.Minus(tc.minusStake))
			return
		}
	}
}

func TestAddToDeveloperInflationPool(t *testing.T) {
	ctx, gm := setupTest(t)

	testCases := []struct {
		testName      string
		inflationPool model.InflationPool
		addCoin       types.Coin
	}{
		{
			testName: "add 10 LNO to empty inflation pool",
			inflationPool: model.InflationPool{
				DeveloperInflationPool: types.NewCoinFromInt64(0),
			},
			addCoin: types.NewCoinFromInt64(10 * types.Decimals),
		},
		{
			testName: "add 10 LNO to a pool with 10000 LNO",
			inflationPool: model.InflationPool{
				DeveloperInflationPool: types.NewCoinFromInt64(10000 * types.Decimals),
			},
			addCoin: types.NewCoinFromInt64(10 * types.Decimals),
		},
	}

	for _, tc := range testCases {
		err := gm.storage.SetInflationPool(ctx, &tc.inflationPool)
		assert.Nil(t, err)
		err = gm.AddToDeveloperInflationPool(ctx, tc.addCoin)
		assert.Nil(t, err)

		inflationPool, err := gm.storage.GetInflationPool(ctx)
		assert.Nil(t, err)
		if !inflationPool.DeveloperInflationPool.IsEqual(
			tc.inflationPool.DeveloperInflationPool.Plus(tc.addCoin)) {
			t.Errorf("%s: diff friction rate, got %v, want %v",
				tc.testName, inflationPool.DeveloperInflationPool,
				tc.inflationPool.DeveloperInflationPool.Plus(tc.addCoin))
			return
		}
	}
}

func TestGetInterestSince(t *testing.T) {
	ctx, gm := setupTest(t)

	testCases := []struct {
		testName       string
		pastRecord     []model.LinoStakeStat
		since          int64
		current        int64
		stake          types.Coin
		expectInterest types.Coin
	}{
		{
			testName: "get past 3 days interest",
			pastRecord: []model.LinoStakeStat{
				model.LinoStakeStat{
					TotalConsumptionFriction: types.NewCoinFromInt64(10000 * types.Decimals),
					TotalLinoStake:           types.NewCoinFromInt64(1000 * types.Decimals),
					UnclaimedFriction:        types.NewCoinFromInt64(10000 * types.Decimals),
					UnclaimedLinoStake:       types.NewCoinFromInt64(1000 * types.Decimals),
				},
				model.LinoStakeStat{
					TotalConsumptionFriction: types.NewCoinFromInt64(0 * types.Decimals),
					TotalLinoStake:           types.NewCoinFromInt64(1000 * types.Decimals),
					UnclaimedFriction:        types.NewCoinFromInt64(0 * types.Decimals),
					UnclaimedLinoStake:       types.NewCoinFromInt64(1000 * types.Decimals),
				},
				model.LinoStakeStat{
					TotalConsumptionFriction: types.NewCoinFromInt64(7777 * types.Decimals),
					TotalLinoStake:           types.NewCoinFromInt64(1000 * types.Decimals),
					UnclaimedFriction:        types.NewCoinFromInt64(1111 * types.Decimals),
					UnclaimedLinoStake:       types.NewCoinFromInt64(143 * types.Decimals),
				},
			},
			since:          0,
			current:        3600 * 24 * 3,
			stake:          types.NewCoinFromInt64(77 * types.Decimals),
			expectInterest: types.NewCoinFromInt64(77000000 + 0 + 59823077),
		},
		{
			testName: "get past 2 days interest",
			pastRecord: []model.LinoStakeStat{
				model.LinoStakeStat{
					TotalConsumptionFriction: types.NewCoinFromInt64(10000 * types.Decimals),
					TotalLinoStake:           types.NewCoinFromInt64(1000 * types.Decimals),
					UnclaimedFriction:        types.NewCoinFromInt64(10000 * types.Decimals),
					UnclaimedLinoStake:       types.NewCoinFromInt64(1000 * types.Decimals),
				},
				model.LinoStakeStat{
					TotalConsumptionFriction: types.NewCoinFromInt64(8000 * types.Decimals),
					TotalLinoStake:           types.NewCoinFromInt64(1000 * types.Decimals),
					UnclaimedFriction:        types.NewCoinFromInt64(8000 * types.Decimals),
					UnclaimedLinoStake:       types.NewCoinFromInt64(1000 * types.Decimals),
				},
				model.LinoStakeStat{
					TotalConsumptionFriction: types.NewCoinFromInt64(7777 * types.Decimals),
					TotalLinoStake:           types.NewCoinFromInt64(1000 * types.Decimals),
					UnclaimedFriction:        types.NewCoinFromInt64(1111 * types.Decimals),
					UnclaimedLinoStake:       types.NewCoinFromInt64(143 * types.Decimals),
				},
			},
			since:          3600 * 24 * 1,
			current:        3600 * 24 * 3,
			stake:          types.NewCoinFromInt64(77 * types.Decimals),
			expectInterest: types.NewCoinFromInt64(61600000 + 59823077),
		},
		{
			testName: "get one of days doesn't have stake",
			pastRecord: []model.LinoStakeStat{
				model.LinoStakeStat{
					TotalConsumptionFriction: types.NewCoinFromInt64(10000 * types.Decimals),
					TotalLinoStake:           types.NewCoinFromInt64(1000 * types.Decimals),
					UnclaimedFriction:        types.NewCoinFromInt64(10000 * types.Decimals),
					UnclaimedLinoStake:       types.NewCoinFromInt64(1000 * types.Decimals),
				},
				model.LinoStakeStat{
					TotalConsumptionFriction: types.NewCoinFromInt64(1000 * types.Decimals),
					TotalLinoStake:           types.NewCoinFromInt64(0 * types.Decimals),
					UnclaimedFriction:        types.NewCoinFromInt64(1000 * types.Decimals),
					UnclaimedLinoStake:       types.NewCoinFromInt64(0 * types.Decimals),
				},
				model.LinoStakeStat{
					TotalConsumptionFriction: types.NewCoinFromInt64(7777 * types.Decimals),
					TotalLinoStake:           types.NewCoinFromInt64(1000 * types.Decimals),
					UnclaimedFriction:        types.NewCoinFromInt64(1111 * types.Decimals),
					UnclaimedLinoStake:       types.NewCoinFromInt64(143 * types.Decimals),
				},
			},
			since:          3600 * 24 * 1,
			current:        3600 * 24 * 3,
			stake:          types.NewCoinFromInt64(0),
			expectInterest: types.NewCoinFromInt64(0),
		},
	}

	for _, tc := range testCases {
		startDay, _ := gm.GetPastDay(ctx, tc.since)
		for i := 0; i < len(tc.pastRecord); i++ {
			err := gm.storage.SetLinoStakeStat(ctx, int64(i), &tc.pastRecord[i])
			if err != nil {
				t.Errorf("%s: failed to set lino stake stat, got err %v", tc.testName, err)
			}
		}
		ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(tc.current, 0)})
		interest, err := gm.GetInterestSince(ctx, tc.since, tc.stake)
		if err != nil {
			t.Errorf("%s: failed to get interest, got err %v", tc.testName, err)
		}
		if !interest.IsEqual(tc.expectInterest) {
			t.Errorf("%s: diff interest, got %v, want %v",
				tc.testName, interest, tc.expectInterest)
			return
		}

		for i := int(startDay); i < len(tc.pastRecord); i++ {
			stat, err := gm.storage.GetLinoStakeStat(ctx, int64(i))
			if err != nil {
				t.Errorf("%s: failed to set lino stake stat, got err %v", tc.testName, err)
			}
			if !stat.TotalLinoStake.IsEqual(tc.pastRecord[i].TotalLinoStake) {
				t.Errorf("%s: diff total lino stake, got %v, want %v",
					tc.testName, stat.TotalLinoStake, tc.pastRecord[i].TotalLinoStake)
				return
			}
			if !stat.TotalConsumptionFriction.IsEqual(tc.pastRecord[i].TotalConsumptionFriction) {
				t.Errorf("%s: diff total consumption friction, got %v, want %v",
					tc.testName, stat.TotalConsumptionFriction, tc.pastRecord[i].TotalConsumptionFriction)
				return
			}
			if !stat.TotalConsumptionFriction.IsEqual(tc.pastRecord[i].TotalConsumptionFriction) {
				t.Errorf("%s: diff total consumption friction, got %v, want %v",
					tc.testName, stat.TotalConsumptionFriction, tc.pastRecord[i].TotalConsumptionFriction)
				return
			}
			if !stat.UnclaimedLinoStake.IsEqual(tc.pastRecord[i].UnclaimedLinoStake.Minus(tc.stake)) {
				t.Errorf("%s: diff total consumption friction, got %v, want %v",
					tc.testName, stat.UnclaimedLinoStake, tc.pastRecord[i].UnclaimedLinoStake.Minus(tc.stake))
				return
			}
			interest := types.NewCoinFromInt64(0)
			if !tc.pastRecord[i].UnclaimedLinoStake.IsZero() {
				interest =
					types.DecToCoin(tc.pastRecord[i].UnclaimedFriction.ToDec().Mul(
						tc.stake.ToDec().Quo(tc.pastRecord[i].UnclaimedLinoStake.ToDec())))
			}
			if !stat.UnclaimedFriction.IsEqual(tc.pastRecord[i].UnclaimedFriction.Minus(interest)) {
				t.Errorf("%s: diff total consumption friction, got %v, want %v",
					tc.testName, stat.UnclaimedFriction, tc.pastRecord[i].UnclaimedFriction.Minus(interest))
				return
			}
		}
	}
}

func TestRecordConsumptionAndLinoStake(t *testing.T) {
	ctx, gm := setupTest(t)

	testCases := []struct {
		testName     string
		atDay        int64
		previousStat model.LinoStakeStat
	}{
		{
			testName: "record consumption and empty lino stake at day 0",
			atDay:    1,
			previousStat: model.LinoStakeStat{
				TotalConsumptionFriction: types.NewCoinFromInt64(0),
				TotalLinoStake:           types.NewCoinFromInt64(0),
				UnclaimedFriction:        types.NewCoinFromInt64(0),
				UnclaimedLinoStake:       types.NewCoinFromInt64(0),
			},
		},
		{
			testName: "record consumption and normal lino stake at day 0",
			atDay:    1,
			previousStat: model.LinoStakeStat{
				TotalConsumptionFriction: types.NewCoinFromInt64(10000 * types.Decimals),
				TotalLinoStake:           types.NewCoinFromInt64(1000 * types.Decimals),
				UnclaimedFriction:        types.NewCoinFromInt64(10000 * types.Decimals),
				UnclaimedLinoStake:       types.NewCoinFromInt64(1000 * types.Decimals),
			},
		},
		{
			testName: "record consumption and normal lino stake at day 1000",
			atDay:    1000,
			previousStat: model.LinoStakeStat{
				TotalConsumptionFriction: types.NewCoinFromInt64(9 * types.Decimals),
				TotalLinoStake:           types.NewCoinFromInt64(1000 * types.Decimals),
				UnclaimedFriction:        types.NewCoinFromInt64(9 * types.Decimals),
				UnclaimedLinoStake:       types.NewCoinFromInt64(1000 * types.Decimals),
			},
		},
		{
			testName: "previous day no stake in stat",
			atDay:    1000,
			previousStat: model.LinoStakeStat{
				TotalConsumptionFriction: types.NewCoinFromInt64(9 * types.Decimals),
				TotalLinoStake:           types.NewCoinFromInt64(0),
				UnclaimedFriction:        types.NewCoinFromInt64(9 * types.Decimals),
				UnclaimedLinoStake:       types.NewCoinFromInt64(0),
			},
		},
	}

	for _, tc := range testCases {
		err := gm.storage.SetLinoStakeStat(ctx, tc.atDay-1, &tc.previousStat)
		assert.Nil(t, err)
		err = gm.SetPastMinutes(ctx, tc.atDay*24*60)
		assert.Nil(t, err)
		ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(tc.atDay*3600*24, 0)})
		err = gm.RecordConsumptionAndLinoStake(ctx)
		assert.Nil(t, err)

		linoStat, err := gm.storage.GetLinoStakeStat(ctx, tc.atDay)
		assert.Nil(t, err)
		if tc.previousStat.TotalLinoStake.IsZero() {
			if !linoStat.TotalConsumptionFriction.IsEqual(tc.previousStat.TotalConsumptionFriction) {
				t.Errorf("%s: diff total consumption friction rate, got %v, want %v",
					tc.testName, linoStat.TotalConsumptionFriction, tc.previousStat.TotalConsumptionFriction)
				return
			}
		} else {
			if !linoStat.TotalConsumptionFriction.IsZero() {
				t.Errorf("%s: diff total consumption friction rate, got %v, want zero",
					tc.testName, linoStat.TotalConsumptionFriction)
				return
			}
		}
		if tc.previousStat.TotalLinoStake.IsZero() {
			if !linoStat.UnclaimedFriction.IsEqual(tc.previousStat.UnclaimedFriction) {
				t.Errorf("%s: diff unclaimed friction, got %v, want %v",
					tc.testName, linoStat.UnclaimedFriction, tc.previousStat.UnclaimedFriction)
				return
			}
		} else {
			if !linoStat.UnclaimedFriction.IsZero() {
				t.Errorf("%s: diff unclaimed friction, got %v, want zero",
					tc.testName, linoStat.UnclaimedFriction)
				return
			}
		}
		if !linoStat.UnclaimedLinoStake.IsEqual(tc.previousStat.UnclaimedLinoStake) {
			t.Errorf("%s: diff unclaim lino stake, got %v, want %v",
				tc.testName, linoStat.UnclaimedLinoStake, tc.previousStat.UnclaimedLinoStake)
			return
		}
		if !linoStat.TotalLinoStake.IsEqual(tc.previousStat.TotalLinoStake) {
			t.Errorf("%s: diff total lino stake, got %v, want %v",
				tc.testName, linoStat.TotalLinoStake, tc.previousStat.TotalLinoStake)
			return
		}
	}
}

func TestRegisterParamChangeEvent(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := time.Now().Unix()

	testCases := []struct {
		testName        string
		atTime          int64
		expectEventList []types.Event
	}{
		{
			testName: "register parameter change event at empty time slot",
			atTime:   baseTime,
			expectEventList: []types.Event{
				testEvent{},
			},
		},
		{
			testName: "register second parameter change event",
			atTime:   baseTime,
			expectEventList: []types.Event{
				testEvent{},
				testEvent{},
			},
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(tc.atTime, 0)})
		err := gm.RegisterParamChangeEvent(ctx, testEvent{})
		if err != nil {
			t.Errorf("%s: failed to register parameter change event, got err %v", tc.testName, err)
		}

		if err = gm.CommitEventCache(ctx); err != nil {
			t.Errorf("%s: failed to commit event cache, %v", tc.testName, err)
		}
		timeEventList := gm.GetTimeEventListAtTime(ctx, tc.atTime+types.ParamChangeTimeout)
		assert.NotNil(t, timeEventList)
		assert.Equal(t, timeEventList.Events, tc.expectEventList)
	}
}
