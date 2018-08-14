package global

import (
	"math/big"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global/model"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
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

	return sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
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
	var initMaxTPS = sdk.NewRat(1000)

	testCases := []struct {
		testName            string
		baseTime            int64
		nextTime            time.Time
		numOfTx             int32
		expectCurrentTPS    sdk.Rat
		expectMaxTPS        sdk.Rat
		expectCapacityRatio sdk.Rat
	}{
		{
			testName:            "0 tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime,
			numOfTx:             0,
			expectCurrentTPS:    sdk.ZeroRat(),
			expectMaxTPS:        initMaxTPS,
			expectCapacityRatio: sdk.ZeroRat(),
		},
		{
			testName:            "2/2 got 1 current tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime.Add(time.Duration(2) * time.Second),
			numOfTx:             2,
			expectCurrentTPS:    sdk.OneRat(),
			expectMaxTPS:        initMaxTPS,
			expectCapacityRatio: sdk.NewRat(1, 1000),
		},
		{
			testName:            "1000/1 got max tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime.Add(time.Duration(1) * time.Second),
			numOfTx:             1000,
			expectCurrentTPS:    initMaxTPS,
			expectMaxTPS:        initMaxTPS,
			expectCapacityRatio: sdk.OneRat(),
		},
		{
			testName:            "2000/2 got max tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime.Add(time.Duration(2) * time.Second),
			numOfTx:             2000,
			expectCurrentTPS:    initMaxTPS,
			expectMaxTPS:        initMaxTPS,
			expectCapacityRatio: sdk.OneRat(),
		},
		{
			testName:            "3000/2 got 1500 current tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime.Add(time.Duration(2) * time.Second),
			numOfTx:             3000,
			expectCurrentTPS:    sdk.NewRat(1500),
			expectMaxTPS:        sdk.NewRat(1500),
			expectCapacityRatio: sdk.OneRat(),
		},
		{
			testName:            "2000/2 got 1000 current tps",
			baseTime:            baseTime.Unix(),
			nextTime:            baseTime.Add(time.Duration(2) * time.Second),
			numOfTx:             2000,
			expectCurrentTPS:    sdk.NewRat(1000),
			expectMaxTPS:        sdk.NewRat(1500),
			expectCapacityRatio: sdk.NewRat(6666667, 10000000),
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
		tps, err := storage.GetTPS(ctx)
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

func TestEvaluateConsumption(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := ctx.BlockHeader().Time.Unix()
	paras, err := gm.paramHolder.GetEvaluateOfContentValueParam(ctx)
	if err != nil {
		t.Errorf("TestEvaluateConsumption: failed to get evaluate of content value param, got err %v", err)
	}

	testCases := []struct {
		testName                           string
		createdTime                        int64
		evaluateTime                       int64
		expectedTimeAdjustment             float64
		totalConsumption                   types.Coin
		expectedTotalConsumptionAdjustment float64
		numOfConsumptionOnAuthor           int64
		expectedConsumptionTimesAdjustment float64
		consumption                        types.Coin
		expectEvaluateResult               types.Coin
	}{
		{
			testName:                           "evaluate in 182 days",
			createdTime:                        baseTime,
			evaluateTime:                       baseTime + 3153600*5,
			expectedTimeAdjustment:             0.5,
			totalConsumption:                   types.NewCoinFromInt64(5000 * types.Decimals),
			expectedTotalConsumptionAdjustment: 1.5,
			numOfConsumptionOnAuthor:           7,
			expectedConsumptionTimesAdjustment: 2.5,
			consumption:                        types.NewCoinFromInt64(1000),
			expectEvaluateResult:               types.NewCoinFromInt64(470),
		},
		{
			testName:                           "evaluate immediately",
			createdTime:                        baseTime,
			evaluateTime:                       baseTime,
			expectedTimeAdjustment:             0.9933071490757153,
			totalConsumption:                   types.NewCoinFromInt64(0),
			expectedTotalConsumptionAdjustment: 1.9933071490757153,
			numOfConsumptionOnAuthor:           7,
			expectedConsumptionTimesAdjustment: 2.5,
			consumption:                        types.NewCoinFromInt64(1000),
			expectEvaluateResult:               types.NewCoinFromInt64(1243),
		},
		{
			testName:                           "evaluate in 360 days",
			createdTime:                        baseTime,
			evaluateTime:                       baseTime + 360*24*3600,
			expectedTimeAdjustment:             0.007667910249215412,
			totalConsumption:                   types.NewCoinFromInt64(0),
			expectedTotalConsumptionAdjustment: 1.9933071490757153,
			numOfConsumptionOnAuthor:           7,
			expectedConsumptionTimesAdjustment: 2.5,
			consumption:                        types.NewCoinFromInt64(1000),
			expectEvaluateResult:               types.NewCoinFromInt64(9),
		},
		{
			testName:                           "evaluate in 1 day with 5 consumption",
			createdTime:                        baseTime,
			evaluateTime:                       baseTime + 24*3600,
			expectedTimeAdjustment:             0.9931225268669581,
			totalConsumption:                   types.NewCoinFromInt64(0),
			expectedTotalConsumptionAdjustment: 1.9933071490757153,
			numOfConsumptionOnAuthor:           7,
			expectedConsumptionTimesAdjustment: 2.5,
			consumption:                        types.NewCoinFromInt64(5 * types.Decimals),
			expectEvaluateResult:               types.NewCoinFromInt64(179346),
		},
		{
			testName:                           "evaluate in 1 day with 1 consumption",
			createdTime:                        baseTime,
			evaluateTime:                       baseTime + 24*3600,
			expectedTimeAdjustment:             0.9931225268669581,
			totalConsumption:                   types.NewCoinFromInt64(0),
			expectedTotalConsumptionAdjustment: 1.9933071490757153,
			numOfConsumptionOnAuthor:           7,
			expectedConsumptionTimesAdjustment: 2.5,
			consumption:                        types.NewCoinFromInt64(1 * types.Decimals),
			expectEvaluateResult:               types.NewCoinFromInt64(49489),
		},
		{
			testName:                           "evaluate in 1 day with 1000 total consumption",
			createdTime:                        baseTime,
			evaluateTime:                       baseTime + 24*3600,
			expectedTimeAdjustment:             0.9931225268669581,
			totalConsumption:                   types.NewCoinFromInt64(1000 * types.Decimals),
			expectedTotalConsumptionAdjustment: 1.9820137900379085,
			numOfConsumptionOnAuthor:           7,
			expectedConsumptionTimesAdjustment: 2.5,
			consumption:                        types.NewCoinFromInt64(1 * types.Decimals),
			expectEvaluateResult:               types.NewCoinFromInt64(49209)},
		{
			testName:                           "evaluate in 1 day with 5000 total consumption",
			createdTime:                        baseTime,
			evaluateTime:                       baseTime + 24*3600,
			expectedTimeAdjustment:             0.9931225268669581,
			totalConsumption:                   types.NewCoinFromInt64(5000 * types.Decimals),
			expectedTotalConsumptionAdjustment: 1.5,
			numOfConsumptionOnAuthor:           7,
			expectedConsumptionTimesAdjustment: 2.5,
			consumption:                        types.NewCoinFromInt64(1 * types.Decimals),
			expectEvaluateResult:               types.NewCoinFromInt64(37242),
		},
		{
			testName:                           "evaluate in 1 day with 5000 total consumption and 100 consumption on author",
			createdTime:                        baseTime,
			evaluateTime:                       baseTime + 24*3600,
			expectedTimeAdjustment:             0.9931225268669581,
			totalConsumption:                   types.NewCoinFromInt64(5000 * types.Decimals),
			expectedTotalConsumptionAdjustment: 1.5,
			numOfConsumptionOnAuthor:           100,
			expectedConsumptionTimesAdjustment: 2,
			consumption:                        types.NewCoinFromInt64(1 * types.Decimals),
			expectEvaluateResult:               types.NewCoinFromInt64(29793),
		},
	}

	for _, tc := range testCases {
		newCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: time.Unix(tc.evaluateTime, 0)})
		if PostTotalConsumptionAdjustment(tc.totalConsumption, paras) != tc.expectedTotalConsumptionAdjustment {
			t.Errorf("%s: diff total consumption adjustment, got %v, want %v", tc.testName,
				PostTotalConsumptionAdjustment(tc.totalConsumption, paras), tc.expectedTotalConsumptionAdjustment)
		}
		if PostTimeAdjustment(tc.evaluateTime-tc.createdTime, paras) != tc.expectedTimeAdjustment {
			t.Errorf("%s: diff time adjustment, got %v, want %v", tc.testName,
				PostTimeAdjustment(tc.evaluateTime-tc.createdTime, paras), tc.expectedTimeAdjustment)
		}
		if PostConsumptionTimesAdjustment(tc.numOfConsumptionOnAuthor, paras) != tc.expectedConsumptionTimesAdjustment {
			t.Errorf("%s: diff consumption times adjustment, got %v, want %v", tc.testName,
				PostConsumptionTimesAdjustment(tc.numOfConsumptionOnAuthor, paras), tc.expectedConsumptionTimesAdjustment)
		}

		evaluateResult, err := gm.EvaluateConsumption(
			newCtx, tc.consumption, tc.numOfConsumptionOnAuthor,
			tc.createdTime, tc.totalConsumption)
		if err != nil {
			t.Errorf("%s: failed to evaluate consumption, got err %v", tc.testName, err)
		}
		if !evaluateResult.IsEqual(tc.expectEvaluateResult) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, evaluateResult, tc.expectEvaluateResult)
		}
	}
}

func TestAddFrictionAndRegisterContentRewardEvent(t *testing.T) {
	ctx, gm := setupTest(t)
	baseTime := ctx.BlockHeader().Time.Unix()
	testCases := []struct {
		testName               string
		frictionCoin           types.Coin
		evaluateCoin           types.Coin
		registerBaseTime       int64
		expectCoinInRewardPool types.Coin
		expectCoinInWindow     types.Coin
	}{
		{
			testName:               "add 1 friction",
			frictionCoin:           types.NewCoinFromInt64(1),
			evaluateCoin:           types.NewCoinFromInt64(1),
			registerBaseTime:       baseTime,
			expectCoinInRewardPool: types.NewCoinFromInt64(1),
			expectCoinInWindow:     types.NewCoinFromInt64(1),
		},
		{
			testName:               "add 100 more friction",
			frictionCoin:           types.NewCoinFromInt64(100),
			evaluateCoin:           types.NewCoinFromInt64(1),
			registerBaseTime:       baseTime + 100,
			expectCoinInRewardPool: types.NewCoinFromInt64(101),
			expectCoinInWindow:     types.NewCoinFromInt64(2),
		},
		{
			testName:               "add 1 more friction",
			frictionCoin:           types.NewCoinFromInt64(1),
			evaluateCoin:           types.NewCoinFromInt64(100),
			registerBaseTime:       baseTime + 1001,
			expectCoinInRewardPool: types.NewCoinFromInt64(102),
			expectCoinInWindow:     types.NewCoinFromInt64(102),
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: time.Unix(tc.registerBaseTime, 0)})
		err := gm.AddFrictionAndRegisterContentRewardEvent(
			ctx, testEvent{}, tc.frictionCoin, tc.evaluateCoin)
		if err != nil {
			t.Errorf("%s: failed to add friction and register event, got err %v", tc.testName, err)
		}

		consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
		if err != nil {
			t.Errorf("%s: failed to get consumption meta, got err %v", tc.testName, err)
		}
		if !consumptionMeta.ConsumptionRewardPool.IsEqual(tc.expectCoinInRewardPool) {
			t.Errorf("%s: diff consumption reward pool, got %v, want %v", tc.testName,
				consumptionMeta.ConsumptionRewardPool, tc.expectCoinInRewardPool)
		}
		if !consumptionMeta.ConsumptionWindow.IsEqual(tc.expectCoinInWindow) {
			t.Errorf("%s: diff consumption window, got %v, want %v", tc.testName,
				consumptionMeta.ConsumptionWindow, tc.expectCoinInWindow)
		}

		timeEventList := gm.GetTimeEventListAtTime(ctx, tc.registerBaseTime+24*7*3600)
		if !assert.Equal(t, types.TimeEventList{[]types.Event{testEvent{}}}, *timeEventList) {
			t.Errorf("%s: diff event list, got %v, want %v", tc.testName,
				*timeEventList, types.TimeEventList{[]types.Event{testEvent{}}})
		}
	}
}

func TestGetRewardAndPopFromWindow(t *testing.T) {
	ctx, gm := setupTest(t)
	testCases := []struct {
		testName                    string
		evaluate                    types.Coin
		penaltyScore                sdk.Rat
		expectReward                types.Coin
		initConsumptionRewardPool   types.Coin
		initConsumptionWindow       types.Coin
		expectConsumptionRewardPool types.Coin
		expectConsumptionWindow     types.Coin
	}{
		{
			testName:                    "1 evaluate, 0 penalty",
			evaluate:                    types.NewCoinFromInt64(1),
			penaltyScore:                sdk.ZeroRat(),
			expectReward:                types.NewCoinFromInt64(100),
			initConsumptionRewardPool:   types.NewCoinFromInt64(1000),
			initConsumptionWindow:       types.NewCoinFromInt64(10),
			expectConsumptionRewardPool: types.NewCoinFromInt64(900),
			expectConsumptionWindow:     types.NewCoinFromInt64(9),
		},
		{
			testName:                    "1/1000 penalty",
			evaluate:                    types.NewCoinFromInt64(1),
			penaltyScore:                sdk.NewRat(1, 1000),
			expectReward:                types.NewCoinFromInt64(100),
			initConsumptionRewardPool:   types.NewCoinFromInt64(1000),
			initConsumptionWindow:       types.NewCoinFromInt64(10),
			expectConsumptionRewardPool: types.NewCoinFromInt64(900),
			expectConsumptionWindow:     types.NewCoinFromInt64(9),
		},
		{
			testName:                    "6/1000 penalty",
			evaluate:                    types.NewCoinFromInt64(1),
			penaltyScore:                sdk.NewRat(6, 1000),
			expectReward:                types.NewCoinFromInt64(99),
			initConsumptionRewardPool:   types.NewCoinFromInt64(1000),
			initConsumptionWindow:       types.NewCoinFromInt64(10),
			expectConsumptionRewardPool: types.NewCoinFromInt64(901),
			expectConsumptionWindow:     types.NewCoinFromInt64(9)},
		{
			testName:                    "1/10 penalty",
			evaluate:                    types.NewCoinFromInt64(1),
			penaltyScore:                sdk.NewRat(1, 10),
			expectReward:                types.NewCoinFromInt64(90),
			initConsumptionRewardPool:   types.NewCoinFromInt64(1000),
			initConsumptionWindow:       types.NewCoinFromInt64(10),
			expectConsumptionRewardPool: types.NewCoinFromInt64(910),
			expectConsumptionWindow:     types.NewCoinFromInt64(9),
		},
		{
			testName:                    "5/10 penalty",
			evaluate:                    types.NewCoinFromInt64(1),
			penaltyScore:                sdk.NewRat(5, 10),
			expectReward:                types.NewCoinFromInt64(50),
			initConsumptionRewardPool:   types.NewCoinFromInt64(1000),
			initConsumptionWindow:       types.NewCoinFromInt64(10),
			expectConsumptionRewardPool: types.NewCoinFromInt64(950),
			expectConsumptionWindow:     types.NewCoinFromInt64(9),
		},
		{
			testName:                    "1/1 penalty",
			evaluate:                    types.NewCoinFromInt64(1),
			penaltyScore:                sdk.NewRat(1, 1),
			expectReward:                types.NewCoinFromInt64(0),
			initConsumptionRewardPool:   types.NewCoinFromInt64(1000),
			initConsumptionWindow:       types.NewCoinFromInt64(10),
			expectConsumptionRewardPool: types.NewCoinFromInt64(1000),
			expectConsumptionWindow:     types.NewCoinFromInt64(9),
		},
		{
			testName:                    "6/1000 penalty",
			evaluate:                    types.NewCoinFromInt64(0),
			penaltyScore:                sdk.ZeroRat(),
			expectReward:                types.NewCoinFromInt64(0),
			initConsumptionRewardPool:   types.NewCoinFromInt64(1000),
			initConsumptionWindow:       types.NewCoinFromInt64(10),
			expectConsumptionRewardPool: types.NewCoinFromInt64(1000),
			expectConsumptionWindow:     types.NewCoinFromInt64(10),
		},
		{
			testName:                    "0 evaluate, 0 penalty",
			evaluate:                    types.NewCoinFromInt64(0),
			penaltyScore:                sdk.OneRat(),
			expectReward:                types.NewCoinFromInt64(0),
			initConsumptionRewardPool:   types.NewCoinFromInt64(1000),
			initConsumptionWindow:       types.NewCoinFromInt64(10),
			expectConsumptionRewardPool: types.NewCoinFromInt64(1000),
			expectConsumptionWindow:     types.NewCoinFromInt64(10),
		},
		// issue https://github.com/lino-network/lino/issues/150
		{
			testName:                    "test large number",
			evaluate:                    types.NewCoinFromInt64(77777777777777),
			penaltyScore:                sdk.ZeroRat(),
			expectReward:                types.NewCoinFromInt64(23333330),
			initConsumptionRewardPool:   types.NewCoinFromInt64(100000000),
			initConsumptionWindow:       types.NewCoinFromInt64(333333333333333),
			expectConsumptionRewardPool: types.NewCoinFromInt64(76666670),
			expectConsumptionWindow:     types.NewCoinFromInt64(255555555555556),
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

		reward, err := gm.GetRewardAndPopFromWindow(ctx, tc.evaluate, tc.penaltyScore)
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
		if !consumptionMeta.ConsumptionWindow.IsEqual(tc.expectConsumptionWindow) {
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
			expectEventList: &types.TimeEventList{[]types.Event{testEvent{}}},
		},
		{
			testName:        "register two events",
			registerAtTime:  baseTime,
			expectResult:    nil,
			expectEventList: &types.TimeEventList{[]types.Event{testEvent{}, testEvent{}}}},
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
			expectEventList: &types.TimeEventList{[]types.Event{testEvent{}}},
		},
	}

	for _, tc := range regCases {
		err := gm.registerEventAtTime(ctx, tc.registerAtTime, testEvent{})
		if !assert.Equal(t, tc.expectResult, err) {
			t.Errorf("%s: diff err result, got %v, want %v", tc.testName, err, tc.expectResult)
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
			interval:               10,
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
			interval:       10,
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
			interval:       5,
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
	assert.Equal(
		t, globalMeta.AnnualInflation, types.RatToCoin(totalLino.ToRat().Mul(globalMeta.GrowthRate)))
	expectInflation := globalMeta.AnnualInflation
	expectContentCreatorInflation := types.NewCoinFromInt64(0)
	expectValidatorInflation := types.NewCoinFromInt64(0)
	expectDeveloperInflation := types.NewCoinFromInt64(0)
	expectInfraInflation := types.NewCoinFromInt64(0)

	globalAllocation, err := gm.paramHolder.GetGlobalAllocationParam(ctx)
	assert.Nil(t, err)
	for i := 0; i < types.HoursPerYear; i++ {
		err = gm.DistributeHourlyInflation(ctx, int64(i))
		assert.Nil(t, err)

		consumptionMeta, err := gm.storage.GetConsumptionMeta(ctx)
		assert.Nil(t, err)

		inflationPool, err := gm.storage.GetInflationPool(ctx)
		assert.Nil(t, err)

		hourlyInflation :=
			types.RatToCoin(expectInflation.ToRat().Mul(sdk.NewRat(1, int64(types.HoursPerYear-i))))
		expectInflation = expectInflation.Minus(hourlyInflation)
		expectContentCreatorInflation =
			expectContentCreatorInflation.Plus(
				types.RatToCoin(hourlyInflation.ToRat().Mul(globalAllocation.ContentCreatorAllocation)))
		expectValidatorInflation =
			expectValidatorInflation.Plus(
				types.RatToCoin(hourlyInflation.ToRat().Mul(globalAllocation.ValidatorAllocation)))
		expectDeveloperInflation =
			expectDeveloperInflation.Plus(
				types.RatToCoin(hourlyInflation.ToRat().Mul(globalAllocation.DeveloperAllocation)))
		expectInfraInflation =
			expectInfraInflation.Plus(
				types.RatToCoin(hourlyInflation.ToRat().Mul(globalAllocation.InfraAllocation)))
		assert.True(t, expectContentCreatorInflation.IsEqual(consumptionMeta.ConsumptionRewardPool))
		assert.True(t, expectInfraInflation.IsEqual(inflationPool.InfraInflationPool))
		assert.True(t, expectDeveloperInflation.IsEqual(inflationPool.DeveloperInflationPool))
		assert.True(t, expectValidatorInflation.IsEqual(inflationPool.ValidatorInflationPool))
	}
	globalMeta, err = gm.storage.GetGlobalMeta(ctx)
	assert.Nil(t, err)
	assert.Equal(t, types.NewCoinFromInt64(0), globalMeta.AnnualInflation)
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

func TestRecalculateAnnuallyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalLino := types.NewCoinFromInt64(10000000000 * types.Decimals)
	ceiling := sdk.NewRat(98, 1000)
	floor := sdk.NewRat(30, 1000)

	testCases := []struct {
		testName            string
		lastYearConsumption types.Coin
		thisYearConsumption types.Coin
		expectGrowthRate    sdk.Rat
	}{
		{
			testName:            "same consumption of last year and this year, get floor growth rate",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			expectGrowthRate:    floor,
		},
		{
			testName:            "right equal to floor growth rate",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(103000000 * types.Decimals),
			expectGrowthRate:    floor,
		},
		{
			testName:            "right equal to ceiling growth rate",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(1098000000 * types.Decimals),
			expectGrowthRate:    ceiling,
		},
		{
			testName:            "bigger than ceiling will be rounded to ceiling growth rate",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(1099000000 * types.Decimals),
			expectGrowthRate:    ceiling,
		},
		{
			testName:            "less year consumption will use floor growth rate",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(90000000 * types.Decimals),
			expectGrowthRate:    floor,
		},
	}

	for _, tc := range testCases {
		globalMeta := &model.GlobalMeta{
			TotalLinoCoin:                 totalLino,
			LastYearCumulativeConsumption: tc.lastYearConsumption,
			CumulativeConsumption:         tc.thisYearConsumption,
			GrowthRate:                    ceiling,
			Ceiling:                       ceiling,
			Floor:                         floor,
		}
		err := gm.storage.SetGlobalMeta(ctx, globalMeta)
		if err != nil {
			t.Errorf("%s: failed to set global meta, got err %v", tc.testName, err)
		}

		err = gm.RecalculateAnnuallyInflation(ctx)
		if err != nil {
			t.Errorf("%s: failed to recalculate annually inflation, got err %v", tc.testName, err)
		}

		globalMeta, err = gm.storage.GetGlobalMeta(ctx)
		if err != nil {
			t.Errorf("%s: failed to get global meta, got err %v", tc.testName, err)
		}

		expectAnnualInflation := types.RatToCoin(totalLino.ToRat().Mul(tc.expectGrowthRate))
		if !expectAnnualInflation.IsEqual(globalMeta.AnnualInflation) {
			t.Errorf("%s: expect annual inflation err, got %v, want %v", tc.testName,
				expectAnnualInflation, globalMeta.AnnualInflation)
		}

		if !globalMeta.LastYearCumulativeConsumption.IsEqual(tc.thisYearConsumption) {
			t.Errorf("%s: diff last year cumulative consumption, got %v, want %v", tc.testName,
				globalMeta.LastYearCumulativeConsumption, tc.thisYearConsumption)
		}
		if !globalMeta.CumulativeConsumption.IsEqual(types.NewCoinFromInt64(0)) {
			t.Errorf("%s: diff cumulative consumption, got %v, want %v", tc.testName,
				globalMeta.CumulativeConsumption, types.NewCoinFromInt64(0))
		}
	}
}

func TestGetGrowthRate(t *testing.T) {
	ctx, gm := setupTest(t)
	totalLino := types.NewCoinFromInt64(1000000)
	ceiling := sdk.NewRat(98, 1000)
	floor := sdk.NewRat(30, 1000)
	bigLastYearConsumption, _ := new(big.Int).SetString("77777777777777777777", 10)
	bigThisYearConsumption, _ := new(big.Int).SetString("83333333333333333332", 10)
	bigLastYearConsumptionCoin := types.NewCoinFromBigInt(bigLastYearConsumption)
	bigThisYearConsumptionCoin := types.NewCoinFromBigInt(bigThisYearConsumption)

	testCases := []struct {
		testName            string
		lastYearConsumption types.Coin
		thisYearConsumption types.Coin
		lastYearGrowthRate  sdk.Rat
		expectGrowthRate    sdk.Rat
	}{
		{
			testName:            "floor growth rate, 0 thisYearConsumption",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(0 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    floor,
		},
		{
			testName:            "ceiling growth rate, 0 lastYearConsumption",
			lastYearConsumption: types.NewCoinFromInt64(0 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    sdk.NewRat(98, 1000),
		},
		{
			testName:            "floor growth rate, thisYearConsumption = lastYearConsumption",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    floor,
		},
		{
			testName:            "ceiling growth rate, 0 lastYearConsumption",
			lastYearConsumption: types.NewCoinFromInt64(0),
			thisYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    sdk.NewRat(98, 1000),
		},
		{
			testName:            "less than floor will be rounded to floor growth rate 1",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(100010000 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    floor,
		},
		{
			testName:            "less than floor will be rounded to floor growth rate 2",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(102900000 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    floor,
		},
		{
			testName:            "less than floor will be rounded to floor growth rate 3",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(103000000 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    floor,
		},
		{
			testName:            "growth rate between floor and ceiling 1",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(103100000 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    sdk.NewRat(31, 1000),
		},
		{
			testName:            "right equal to ceiling growth rate",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(109800000 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    ceiling,
		},
		{
			testName:            "higher than ceiling will be rouned to ceiling growth rate",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(109900000 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    ceiling,
		},
		{
			testName:            "growth rate between floor and ceiling 2",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(109700000 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    sdk.NewRat(97, 1000),
		},
		{
			testName:            "growth rate between floor and ceiling 3",
			lastYearConsumption: types.NewCoinFromInt64(100000000 * types.Decimals),
			thisYearConsumption: types.NewCoinFromInt64(104700000 * types.Decimals),
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    sdk.NewRat(47, 1000),
		},
		// issue https://github.com/lino-network/lino/issues/150
		{
			testName:            "overflow testing",
			lastYearConsumption: bigLastYearConsumptionCoin,
			thisYearConsumption: bigThisYearConsumptionCoin,
			lastYearGrowthRate:  sdk.NewRat(98, 1000),
			expectGrowthRate:    sdk.NewRat(357143, 5000000),
		},
	}

	for _, tc := range testCases {
		globalMeta := &model.GlobalMeta{
			TotalLinoCoin:                 totalLino,
			LastYearCumulativeConsumption: tc.lastYearConsumption,
			CumulativeConsumption:         tc.thisYearConsumption,
			GrowthRate:                    tc.lastYearGrowthRate,
			Ceiling:                       ceiling,
			Floor:                         floor,
		}
		err := gm.storage.SetGlobalMeta(ctx, globalMeta)
		if err != nil {
			t.Errorf("%s: failed to set global meta, got err %v", tc.testName, err)
		}

		growthRate, err := gm.getGrowthRate(ctx)
		if err != nil {
			t.Errorf("%s: failed to get growth rate, got err %v", tc.testName, err)
		}
		if !tc.expectGrowthRate.Equal(growthRate) {
			t.Errorf("%s: diff growth rate, got %v, want %v", tc.testName, growthRate, tc.expectGrowthRate)
		}

		globalMeta, err = gm.storage.GetGlobalMeta(ctx)
		if err != nil {
			t.Errorf("%s: failed to get global meta, got err %v", tc.testName, err)
		}
		if !globalMeta.LastYearCumulativeConsumption.IsEqual(tc.thisYearConsumption) {
			t.Errorf("%s: diff last year cumulative consumption, got %v, want %v", tc.testName,
				globalMeta.LastYearCumulativeConsumption, tc.thisYearConsumption)
		}
		if !globalMeta.CumulativeConsumption.IsEqual(types.NewCoinFromInt64(0)) {
			t.Errorf("%s: diff cumulative consumption, got %v, want %v", tc.testName,
				globalMeta.CumulativeConsumption, types.NewCoinFromInt64(0))
		}
		if !tc.expectGrowthRate.Equal(globalMeta.GrowthRate) {
			t.Errorf("%s: diff growth rate, got %v, want %v", tc.testName,
				globalMeta.GrowthRate, tc.expectGrowthRate)
		}
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

func TestGetDeveloperMonthlyInflation(t *testing.T) {
	ctx, gm := setupTest(t)
	totalDeveloperInflation := types.NewCoinFromInt64(10000 * 100)
	inflationPool := &model.InflationPool{
		DeveloperInflationPool: totalDeveloperInflation,
	}
	err := gm.storage.SetInflationPool(ctx, inflationPool)
	assert.Nil(t, err)
	coin, err := gm.GetDeveloperMonthlyInflation(ctx)
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

func TestAddConsumption(t *testing.T) {
	ctx, gm := setupTest(t)

	testCases := []struct {
		testName string
		coin     types.Coin
		expect   types.Coin
	}{
		{
			testName: "add 100 consumption",
			coin:     types.NewCoinFromInt64(100),
			expect:   types.NewCoinFromInt64(100),
		},
		{
			testName: "add 1 more consumption",
			coin:     types.NewCoinFromInt64(1),
			expect:   types.NewCoinFromInt64(101),
		},
	}

	for _, tc := range testCases {
		err := gm.AddConsumption(ctx, tc.coin)
		if err != nil {
			t.Errorf("%s: failed to add consumption, got err %v", tc.testName, err)
		}

		globalMeta, err := gm.storage.GetGlobalMeta(ctx)
		if err != nil {
			t.Errorf("%s: failed to get global meta, got err %v", tc.testName, err)
		}
		if !globalMeta.CumulativeConsumption.IsEqual(tc.expect) {
			t.Errorf("%s: diff cumulative consumption, got %v, want %v", tc.testName,
				globalMeta.CumulativeConsumption, tc.expect)
		}
	}
}
