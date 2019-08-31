package manager

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	parammodel "github.com/lino-network/lino/param"
	param "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	linotypes "github.com/lino-network/lino/types"
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
		CapacityUsagePerTransaction: linotypes.NewCoinFromInt64(1 * linotypes.Decimals),
		VirtualCoin:                 linotypes.NewCoinFromInt64(1 * linotypes.Decimals),
		GeneralMsgQuotaRatio:        linotypes.NewDecFromRat(20, 100),
		GeneralMsgEMAFactor:         linotypes.NewDecFromRat(1, 10),
		AppMsgQuotaRatio:            linotypes.NewDecFromRat(80, 100),
		AppMsgEMAFactor:             linotypes.NewDecFromRat(1, 10),
		ExpectedMaxMPS:              linotypes.NewDecFromRat(1000, 1),
		MsgFeeFactorA:               linotypes.NewDecFromRat(6, 1),
		MsgFeeFactorB:               linotypes.NewDecFromRat(10, 1),
		MaxMPSDecayRate:             linotypes.NewDecFromRat(99, 100),
	}, nil).Maybe()

	suite.global.On("GetLastBlockTime", mock.Anything).Return(baseTime.Unix(), nil).Maybe()

}

func (suite *BandwidthManagerTestSuite) TestAddMsgSignedByUser() {
	testCases := []struct {
		testName        string
		amount          uint32
		expectBlockInfo model.BlockInfo
	}{
		{
			testName: "add user signed message",
			amount:   1,
			expectBlockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 1,
				CurMsgFee:            linotypes.NewCoinFromInt64(int64(0)),
			},
		},
	}

	for _, tc := range testCases {
		suite.bm.AddMsgSignedByUser(suite.Ctx, tc.amount)
		info, err := suite.bm.storage.GetBlockInfo(suite.Ctx)
		suite.Require().Nil(err)
		suite.Equal(tc.expectBlockInfo, *info, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestAddMsgSignedByApp() {
	testCases := []struct {
		testName        string
		amount          uint32
		expectBlockInfo model.BlockInfo
	}{
		{
			testName: "add user signed message",
			amount:   1,
			expectBlockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  1,
				TotalMsgSignedByUser: 0,
				CurMsgFee:            linotypes.NewCoinFromInt64(int64(0)),
			},
		},
	}

	for _, tc := range testCases {
		suite.bm.AddMsgSignedByApp(suite.Ctx, tc.amount)
		info, err := suite.bm.storage.GetBlockInfo(suite.Ctx)
		suite.Require().Nil(err)
		suite.Equal(tc.expectBlockInfo, *info, "%s", tc.testName)
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
		testName             string
		bandwidthInfo        model.BandwidthInfo
		expectMessageFeeCoin int64
	}{
		{
			testName: "test1",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(0),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(1000),
			},
			expectMessageFeeCoin: int64(2523),
		},
		{
			testName: "test2",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(100),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(1000),
			},
			expectMessageFeeCoin: int64(50006),
		},
		{
			testName: "test3",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(200),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(1000),
			},
			expectMessageFeeCoin: int64(1000000),
		},
		{
			testName: "test4",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(300),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(1000),
			},
			expectMessageFeeCoin: int64(19997635),
		},
	}

	for _, tc := range testCases {
		suite.bm.storage.SetBandwidthInfo(suite.Ctx, &tc.bandwidthInfo)
		suite.bm.CalculateCurMsgFee(suite.Ctx)
		info, getErr := suite.bm.storage.GetBlockInfo(suite.Ctx)
		suite.Require().Nil(getErr)
		fmt.Println(info.CurMsgFee)
		suite.Equal(linotypes.NewCoinFromInt64(tc.expectMessageFeeCoin), info.CurMsgFee, "%s", tc.testName)
	}
}

// only one step
func (suite *BandwidthManagerTestSuite) TestUpdateMaxMPSAndEMA() {
	testCases := []struct {
		testName         string
		blockInfo        model.BlockInfo
		bandwidthInfo    model.BandwidthInfo
		expectGeneralEMA string
		expectAppEMA     string
	}{
		{
			testName: "test general message ema",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(10),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(1000),
			},
			blockInfo: model.BlockInfo{
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
			blockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  270,
				TotalMsgSignedByUser: 0,
			},
			expectGeneralEMA: "0",
			expectAppEMA:     "54",
		},
	}

	for _, tc := range testCases {
		suite.bm.storage.SetBandwidthInfo(suite.Ctx, &tc.bandwidthInfo)
		suite.bm.storage.SetBlockInfo(suite.Ctx, &tc.blockInfo)

		err := suite.bm.UpdateMaxMPSAndEMA(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)

		expectedGeneralEMA, err := sdk.NewDecFromStr(tc.expectGeneralEMA)
		suite.Nil(err, "%s", tc.testName)

		expectedAppEMA, err := sdk.NewDecFromStr(tc.expectAppEMA)
		suite.Nil(err, "%s", tc.testName)

		info, err := suite.bm.storage.GetBandwidthInfo(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(expectedGeneralEMA, info.GeneralMsgEMA, "%s", tc.testName)
		suite.Equal(expectedAppEMA, info.AppMsgEMA, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestDecayMaxMPS() {
	testCases := []struct {
		testName     string
		info         model.BandwidthInfo
		expectedInfo model.BandwidthInfo
	}{
		{
			testName: "test1",
			info: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(0),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(100),
			},
			expectedInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(0),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(99),
			},
		},
	}

	for _, tc := range testCases {
		err := suite.bm.storage.SetBandwidthInfo(suite.Ctx, &tc.info)
		suite.Require().Nil(err)
		err = suite.bm.DecayMaxMPS(suite.Ctx)
		suite.Require().Nil(err)
		info, err := suite.bm.storage.GetBandwidthInfo(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(tc.expectedInfo, *info, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestCalculateEMA() {
	testCases := []struct {
		testName    string
		prevEMA     sdk.Dec
		k           sdk.Dec
		curMPS      sdk.Dec
		expectedEMA sdk.Dec
	}{
		{
			testName:    "test1",
			prevEMA:     linotypes.NewDecFromRat(100, 1),
			k:           linotypes.NewDecFromRat(1, 10),
			curMPS:      linotypes.NewDecFromRat(200, 1),
			expectedEMA: linotypes.NewDecFromRat(110, 1),
		},
	}

	for _, tc := range testCases {
		res := suite.bm.calculateEMA(tc.prevEMA, tc.k, tc.curMPS)
		suite.Equal(tc.expectedEMA, res, "%s", tc.testName)
	}
}
