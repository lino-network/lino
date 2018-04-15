package global

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/global/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	oldwire "github.com/tendermint/go-wire"
	dbm "github.com/tendermint/tmlibs/db"
)

const (
	eventTypeTestEvent = 0x1
)

type testEvent struct{}

// Construct some global addrs and txs for tests.
var (
	TestGlobalKVStoreKey = sdk.NewKVStoreKey("global")

	_ = oldwire.RegisterInterface(
		struct{ types.Event }{},
		oldwire.ConcreteType{testEvent{}, eventTypeTestEvent},
	)
)

func InitGlobalManager(ctx sdk.Context, gm *GlobalManager) error {
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
	return gm.InitGlobalManager(ctx, globalState)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func setupTest(t *testing.T) (sdk.Context, *GlobalManager) {
	ctx := getContext()
	globalManager := NewGlobalManager(TestGlobalKVStoreKey)
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
		{BaseTime: baseTime, NextTime: baseTime, NumOfTx: 0, ExpectCurrentTPS: initMaxTPS,
			ExpectMaxTPS: initMaxTPS, ExpectCapacityRatio: sdk.OneRat},
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
