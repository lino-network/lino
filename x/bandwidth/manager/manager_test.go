package manager

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	parammodel "github.com/lino-network/lino/param"
	param "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/bandwidth/model"
	global "github.com/lino-network/lino/x/global/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BandwidthManagerTestSuite struct {
	testsuites.CtxTestSuite
	bm BandwidthManager
	ph *param.ParamKeeper
	// deps
	global *global.GlobalKeeper
}

func TestBandwidthManagerTestSuite(t *testing.T) {
	suite.Run(t, new(BandwidthManagerTestSuite))
}

func (suite *BandwidthManagerTestSuite) SetupTest() {
	baseTime := time.Now()
	testBandwidthKey := sdk.NewKVStoreKey("bandwidth")
	suite.SetupCtx(0, baseTime.Add(3*time.Second), testBandwidthKey)
	suite.global = &global.GlobalKeeper{}
	suite.ph = &param.ParamKeeper{}
	suite.bm = *NewBandwidthManager(testBandwidthKey, suite.ph, suite.global)
	suite.bm.InitGenesis(suite.Ctx)
	suite.ph.On("GetBandwidthParam", mock.Anything).Return(&parammodel.BandwidthParam{
		SecondsToRecoverBandwidth:   int64(7 * 24 * 3600),
		CapacityUsagePerTransaction: types.NewCoinFromInt64(1 * types.Decimals),
		VirtualCoin:                 types.NewCoinFromInt64(1 * types.Decimals),
		GeneralMsgQuotaRatio:        types.NewDecFromRat(20, 100),
		GeneralMsgEMAFactor:         types.NewDecFromRat(1, 10),
		AppMsgQuotaRatio:            types.NewDecFromRat(80, 100),
		AppMsgEMAFactor:             types.NewDecFromRat(1, 10),
		ExpectedMaxMPS:              types.NewDecFromRat(1000, 1),
		MsgFeeFactorA:               types.NewDecFromRat(6, 1),
		MsgFeeFactorB:               types.NewDecFromRat(10, 1),
	}, nil).Maybe()

	suite.global.On("GetLastBlockTime", mock.Anything).Return(baseTime.Unix(), nil).Maybe()

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
		expectedRes string
	}{
		{
			testName:    "test1",
			x:           sdk.NewDec(0),
			expectedRes: "1",
		},
		{
			testName:    "test2",
			x:           sdk.NewDec(-6),
			expectedRes: "0.002522536744331140", // truth: 0.00247875217
		},
		{
			testName:    "test3",
			x:           sdk.NewDec(10),
			expectedRes: "20983.411084513772091023", // truth: 22026.4657948
		},
		{
			testName:    "test4",
			x:           sdk.NewDec(30),
			expectedRes: "6944323407751.788994887415441546", // truth: 1.0686475e+13
		},
	}

	for _, tc := range testCases {
		res := suite.bm.approximateExp(tc.x)
		expectedFee, err := sdk.NewDecFromStr(tc.expectedRes)
		suite.Require().Nil(err)
		suite.Equal(expectedFee, res, "%s", tc.testName)
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
		expectedFee, err := sdk.NewDecFromStr(tc.expectMessageFee)
		suite.Require().Nil(err)
		suite.Equal(expectedFee, suite.bm.blockStatsCache.CurMsgFee, "%s", tc.testName)
	}
}

// only one step
func (suite *BandwidthManagerTestSuite) TestUpdateMaxMPSAndEMA() {
	testCases := []struct {
		testName              string
		blockStatsCache       model.BlockStatsCache
		bandwidthInfo         model.BandwidthInfo
		expectGeneralEMA      string
		expectAppEMA          string
		expectedLastBlockInfo model.LastBlockInfo
	}{
		{
			testName: "test general message ema",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(10),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(1000),
			},
			blockStatsCache: model.BlockStatsCache{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 60,
			},
			expectedLastBlockInfo: model.LastBlockInfo{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 60,
			},
			expectGeneralEMA: "11",
			expectAppEMA:     "0",
		},
		{
			testName: "test app message ema",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(0),
				AppMsgEMA:     sdk.NewDec(50),
				MaxMPS:        sdk.NewDec(1000),
			},
			blockStatsCache: model.BlockStatsCache{
				TotalMsgSignedByApp:  270,
				TotalMsgSignedByUser: 0,
			},
			expectedLastBlockInfo: model.LastBlockInfo{
				TotalMsgSignedByApp:  270,
				TotalMsgSignedByUser: 0,
			},
			expectGeneralEMA: "0",
			expectAppEMA:     "54",
		},
	}

	for _, tc := range testCases {
		suite.bm.storage.SetBandwidthInfo(suite.Ctx, &tc.bandwidthInfo)
		suite.bm.blockStatsCache = tc.blockStatsCache

		err := suite.bm.UpdateMaxMPSAndEMA(suite.Ctx)
		fmt.Println(err)
		suite.Nil(err, "%s", tc.testName)

		expectedGeneralEMA, err := sdk.NewDecFromStr(tc.expectGeneralEMA)
		suite.Nil(err, "%s", tc.testName)

		expectedAppEMA, err := sdk.NewDecFromStr(tc.expectAppEMA)
		suite.Nil(err, "%s", tc.testName)

		info, err := suite.bm.storage.GetBandwidthInfo(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(expectedGeneralEMA, info.GeneralMsgEMA, "%s", tc.testName)
		suite.Equal(expectedAppEMA, info.AppMsgEMA, "%s", tc.testName)

		// make sure store the cache into last block info
		lastBlockInfo, err := suite.bm.storage.GetLastBlockInfo(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(tc.expectedLastBlockInfo, *lastBlockInfo, "%s", tc.testName)
	}
}
