package manager

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/x/bandwidth/model"
	global "github.com/lino-network/lino/x/global/mocks"
	"github.com/stretchr/testify/suite"
)

type BandwidthManagerTestSuite struct {
	testsuites.CtxTestSuite
	bm          BandwidthManager
	paramHolder param.ParamHolder
	// deps
	global *global.GlobalKeeper
}

func TestBandwidthManagerTestSuite(t *testing.T) {
	suite.Run(t, new(BandwidthManagerTestSuite))
}

func (suite *BandwidthManagerTestSuite) SetupTest() {
	testBandwidthKey := sdk.NewKVStoreKey("bandwidth")
	suite.SetupCtx(0, time.Unix(0, 0), testBandwidthKey)
	suite.global = &global.GlobalKeeper{}
	suite.bm = NewBandwidthManager(testBandwidthKey, suite.paramHolder, suite.global)

}

func (suite *BandwidthManagerTestSuite) TestAddMsgSignedByUser() {
	testCases := []struct {
		testName              string
		amount                uint32
		expectBlockStatsCache model.BlockStatsCache
	}{
		{
			testName: "add user signed message",
			amount:   1,
			expectBlockStatsCache: model.BlockStatsCache{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 1,
			},
		},
	}

	for _, tc := range testCases {
		suite.bm.AddMsgSignedByUser(suite.Ctx, tc.amount)
		suite.Equal(tc.expectBlockStatsCache, suite.bm.blockStatsCache, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestAddMsgSignedByApp() {
	testCases := []struct {
		testName              string
		amount                uint32
		expectBlockStatsCache model.BlockStatsCache
	}{
		{
			testName: "add user signed message",
			amount:   1,
			expectBlockStatsCache: model.BlockStatsCache{
				TotalMsgSignedByApp:  1,
				TotalMsgSignedByUser: 0,
			},
		},
	}

	for _, tc := range testCases {
		suite.bm.AddMsgSignedByApp(suite.Ctx, tc.amount)
		suite.Equal(tc.expectBlockStatsCache, suite.bm.blockStatsCache, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestApproximateExp() {
	testCases := []struct {
		testName    string
		x           sdk.Dec
		expectedRes sdk.Dec
	}{
		{
			testName:    "test1",
			x:           sdk.NewDec(0),
			expectedRes: sdk.NewDec(1),
		},
		// TODO(zhimao): compare the result with math.Exp, and calculate the difference
	}

	for _, tc := range testCases {
		res := suite.bm.approximateExp(tc.x)
		suite.Equal(tc.expectedRes, res, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestCalculateCurMsgFee() {
	testCases := []struct {
		testName         string
		bandwidthInfo    model.BandwidthInfo
		expectMessageFee string
	}{
		{
			testName: "test1",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(0),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(1000),
			},
			expectMessageFee: "0.025225367443311400",
		},
		{
			testName: "test2",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(100),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(1000),
			},
			expectMessageFee: "0.500059123770510650",
		},
		{
			testName: "test3",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(200),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(1000),
			},
			expectMessageFee: "10.000000000000000000",
		},
		{
			testName: "test4",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(300),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(1000),
			},
			expectMessageFee: "199.976353287961290180",
		},
	}

	for _, tc := range testCases {
		suite.bm.storage.SetBandwidthInfo(suite.Ctx, &tc.bandwidthInfo)
		suite.bm.CalculateCurMsgFee(suite.Ctx)
		// expectedFee, err := sdk.NewDecFromStr(tc.expectMessageFee)
		// suite.Require().Nil(err)
		// suite.Equal(expectedFee, suite.bm.blockStatsCache.CurMsgFee, "%s", tc.testName)
	}
}

// // only one step
// func (suite *BandwidthManagerTestSuite) TestUpdateMaxMPSAndEMA(t *testing.T) {
// 	baseTime := time.Now()
// 	testCases := []struct {
// 		testName         string
// 		LastBlockInfo    model.LastBlockInfo
// 		bandwidthInfo    model.BandwidthInfo
// 		lastBlockTime    int64
// 		expectGeneralEMA string
// 		expectAppEMA     string
// 	}{
// 		{
// 			testName: "test general message ema",
// 			bandwidthInfo: model.BandwidthInfo{
// 				GeneralMsgEMA: sdk.NewDec(10),
// 				AppMsgEMA:     sdk.NewDec(0),
// 				MaxMPS:        sdk.NewDec(1000),
// 			},
// 			LastBlockInfo: model.LastBlockInfo{
// 				TotalMsgSignedByApp:  0,
// 				TotalMsgSignedByUser: 60,
// 			},
// 			lastBlockTime:    baseTime.Unix(),
// 			expectGeneralEMA: "11",
// 			expectAppEMA:     "0",
// 		},
// 		{
// 			testName: "test app message ema",
// 			bandwidthInfo: model.BandwidthInfo{
// 				GeneralMsgEMA: sdk.NewDec(0),
// 				AppMsgEMA:     sdk.NewDec(50),
// 				MaxMPS:        sdk.NewDec(1000),
// 			},
// 			LastBlockInfo: model.LastBlockInfo{
// 				TotalMsgSignedByApp:  270,
// 				TotalMsgSignedByUser: 0,
// 			},
// 			lastBlockTime:    baseTime.Unix(),
// 			expectGeneralEMA: "0",
// 			expectAppEMA:     "54",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: baseTime.Add(3 * time.Second)})
// 		bm.storage.SetBandwidthInfo(ctx, &tc.bandwidthInfo)
// 		bm.storage.SetLastBlockInfo(ctx, &tc.LastBlockInfo)

// 		bm.UpdateMaxMPSAndEMA(ctx, tc.lastBlockTime)
// 		expectedGeneralEMA, err := sdk.NewDecFromStr(tc.expectGeneralEMA)
// 		assert.Nil(t, err)

// 		expectedAppEMA, err := sdk.NewDecFromStr(tc.expectAppEMA)
// 		assert.Nil(t, err)

// 		info, err := bm.storage.GetBandwidthInfo(ctx)
// 		assert.Nil(t, err)
// 		assert.Equal(t, expectedGeneralEMA, info.GeneralMsgEMA, "%s: diff general EMA result, got %v, want %v", tc.testName, info.GeneralMsgEMA, expectedGeneralEMA)
// 		assert.Equal(t, expectedAppEMA, info.AppMsgEMA, "%s: diff app EMA result, got %v, want %v", tc.testName, info.AppMsgEMA, expectedAppEMA)
// 	}
// }
