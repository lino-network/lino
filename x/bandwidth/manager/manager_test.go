package manager

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	parammodel "github.com/lino-network/lino/param"
	param "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	linotypes "github.com/lino-network/lino/types"
	account "github.com/lino-network/lino/x/account/mocks"
	"github.com/lino-network/lino/x/bandwidth/model"
	"github.com/lino-network/lino/x/bandwidth/types"
	developer "github.com/lino-network/lino/x/developer/mocks"
	devModel "github.com/lino-network/lino/x/developer/model"
	global "github.com/lino-network/lino/x/global/mocks"
	vote "github.com/lino-network/lino/x/vote/mocks"
)

type BandwidthManagerTestSuite struct {
	testsuites.CtxTestSuite
	bm       BandwidthManager
	ph       *param.ParamKeeper
	baseTime time.Time
	// deps
	global *global.GlobalKeeper
	vm     *vote.VoteKeeper
	dm     *developer.DeveloperKeeper
	am     *account.AccountKeeper
}

func TestBandwidthManagerTestSuite(t *testing.T) {
	suite.Run(t, new(BandwidthManagerTestSuite))
}

func (suite *BandwidthManagerTestSuite) SetupTest() {
	suite.baseTime = time.Now()
	testBandwidthKey := sdk.NewKVStoreKey("bandwidth")
	suite.SetupCtx(0, suite.baseTime.Add(3*time.Second), testBandwidthKey)
	suite.global = &global.GlobalKeeper{}
	suite.ph = &param.ParamKeeper{}
	suite.vm = &vote.VoteKeeper{}
	suite.dm = &developer.DeveloperKeeper{}
	suite.am = &account.AccountKeeper{}
	suite.bm = *NewBandwidthManager(testBandwidthKey, suite.ph, suite.global, suite.vm, suite.dm, suite.am)
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
		AppBandwidthPoolSize:        linotypes.NewDecFromRat(10, 1),
		AppVacancyFactor:            linotypes.NewDecFromRat(69, 100),
		AppPunishmentFactor:         linotypes.NewDecFromRat(14, 5),
	}, nil).Maybe()

	suite.vm.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("AppX")).Return(linotypes.NewCoinFromInt64(10), nil).Maybe()
	suite.vm.On("GetLinoStake", suite.Ctx, linotypes.AccountKey("AppY")).Return(linotypes.NewCoinFromInt64(90), nil).Maybe()
	suite.dm.On("GetLiveDevelopers", mock.Anything).Return([]devModel.Developer{
		{
			Username: "AppX",
		},
		{
			Username: "AppY",
		},
	}, nil).Maybe()
	suite.dm.On("GetAffiliatingApp", suite.Ctx, linotypes.AccountKey("AppY")).Return(linotypes.AccountKey("AppY"), nil).Maybe()
	suite.dm.On("GetAffiliatingApp", suite.Ctx, linotypes.AccountKey("UserX")).Return(linotypes.AccountKey("dummy"), types.ErrUserMsgFeeNotEnough()).Maybe()
	suite.global.On("GetLastBlockTime", mock.Anything).Return(suite.baseTime.Unix(), nil).Maybe()
	suite.global.On("AddToValidatorInflationPool", mock.Anything, mock.Anything).Return(nil).Maybe()
	suite.am.On("MinusCoinFromUsername", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

}

func (suite *BandwidthManagerTestSuite) TestAddMsgSignedByUser() {
	testCases := []struct {
		testName        string
		amount          int64
		expectBlockInfo model.BlockInfo
	}{
		{
			testName: "add user signed message",
			amount:   1,
			expectBlockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 1,
				CurMsgFee:            linotypes.NewCoinFromInt64(int64(0)),
				CurU:                 sdk.NewDec(1),
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
		amount          int64
		expectBlockInfo model.BlockInfo
		expectAppInfo   model.AppBandwidthInfo
	}{
		{
			testName: "add user signed message",
			amount:   1,
			expectBlockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  1,
				TotalMsgSignedByUser: 0,
				CurMsgFee:            linotypes.NewCoinFromInt64(int64(0)),
				CurU:                 sdk.NewDec(1),
			},
			expectAppInfo: model.AppBandwidthInfo{
				MessagesInCurBlock: 1,
				MaxBandwidthCredit: sdk.NewDec(0),
				CurBandwidthCredit: sdk.NewDec(0),
				ExpectedMPS:        sdk.NewDec(0),
				LastRefilledAt:     0,
			},
		},
	}

	for _, tc := range testCases {
		appName := linotypes.AccountKey("test")
		appBandwidthInfo := model.AppBandwidthInfo{}

		err := suite.bm.storage.SetAppBandwidthInfo(suite.Ctx, appName, &appBandwidthInfo)
		suite.Require().Nil(err)

		suite.bm.AddMsgSignedByApp(suite.Ctx, appName, tc.amount)
		info, err := suite.bm.storage.GetBlockInfo(suite.Ctx)
		suite.Require().Nil(err)
		suite.Equal(tc.expectBlockInfo, *info, "%s", tc.testName)

		appInfo, err := suite.bm.storage.GetAppBandwidthInfo(suite.Ctx, appName)
		suite.Equal(tc.expectAppInfo, *appInfo, "%s", tc.testName)
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
		{
			testName: "test5",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(300),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(500),
			},
			expectMessageFeeCoin: int64(19997635),
		},
	}

	for _, tc := range testCases {
		suite.bm.storage.SetBandwidthInfo(suite.Ctx, &tc.bandwidthInfo)
		err := suite.bm.CalculateCurMsgFee(suite.Ctx)
		suite.Require().Nil(err)

		info, getErr := suite.bm.storage.GetBlockInfo(suite.Ctx)
		suite.Require().Nil(getErr)
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
		expectMaxMPS     sdk.Dec
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
			expectMaxMPS:     sdk.NewDec(1000),
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
			expectMaxMPS:     sdk.NewDec(1000),
		},
		{
			testName: "test update max MPS",
			bandwidthInfo: model.BandwidthInfo{
				GeneralMsgEMA: sdk.NewDec(0),
				AppMsgEMA:     sdk.NewDec(0),
				MaxMPS:        sdk.NewDec(1000),
			},
			blockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  15000,
				TotalMsgSignedByUser: 15000,
			},
			expectGeneralEMA: "500",
			expectAppEMA:     "500",
			expectMaxMPS:     sdk.NewDec(10000),
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
		suite.Equal(tc.expectMaxMPS, info.MaxMPS, "%s", tc.testName)

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

func (suite *BandwidthManagerTestSuite) TestIsUserMsgFeeEnough() {
	testCases := []struct {
		testName    string
		providedFee auth.StdFee
		curMsgFee   linotypes.Coin
		expectedRes bool
	}{
		{
			testName: "test1",
			providedFee: auth.StdFee{
				Amount: sdk.NewCoins(sdk.NewCoin(linotypes.LinoCoinDenom, sdk.NewInt(12))),
			},
			curMsgFee:   linotypes.NewCoinFromInt64(int64(10)),
			expectedRes: true,
		},
		{
			testName: "test2",
			providedFee: auth.StdFee{
				Amount: sdk.NewCoins(sdk.NewCoin(linotypes.LinoCoinDenom, sdk.NewInt(2))),
			},
			curMsgFee:   linotypes.NewCoinFromInt64(int64(10)),
			expectedRes: false,
		},
		{
			testName: "test3",
			providedFee: auth.StdFee{
				Amount: sdk.NewCoins(sdk.NewCoin("dummy", sdk.NewInt(12))),
			},
			curMsgFee:   linotypes.NewCoinFromInt64(int64(10)),
			expectedRes: false,
		},
	}

	for _, tc := range testCases {
		info := model.BlockInfo{
			CurMsgFee: tc.curMsgFee,
		}
		suite.bm.storage.SetBlockInfo(suite.Ctx, &info)
		res := suite.bm.IsUserMsgFeeEnough(suite.Ctx, tc.providedFee)
		suite.Equal(tc.expectedRes, res, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestClearBlockInfo() {
	testCases := []struct {
		testName          string
		curBlockInfo      model.BlockInfo
		expectedBlockInfo model.BlockInfo
	}{
		{
			testName: "test1",
			curBlockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  1,
				TotalMsgSignedByUser: 2,
				CurMsgFee:            linotypes.NewCoinFromInt64(int64(32)),
				CurU:                 sdk.NewDec(2),
			},
			expectedBlockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 0,
				CurMsgFee:            linotypes.NewCoinFromInt64(int64(32)),
				CurU:                 sdk.NewDec(2),
			},
		},
	}

	for _, tc := range testCases {
		err := suite.bm.storage.SetBlockInfo(suite.Ctx, &tc.curBlockInfo)
		suite.Nil(err, "%s", tc.testName)

		err = suite.bm.ClearBlockInfo(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)

		res, err := suite.bm.storage.GetBlockInfo(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(tc.expectedBlockInfo, *res, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestGetBandwidthCostPerMsg() {
	testCases := []struct {
		testName     string
		u            sdk.Dec
		p            sdk.Dec
		expectedCost sdk.Dec
	}{
		{
			testName:     "test1",
			u:            sdk.NewDec(5),
			p:            sdk.NewDec(2),
			expectedCost: sdk.NewDec(10),
		},
	}

	for _, tc := range testCases {
		res := suite.bm.GetBandwidthCostPerMsg(suite.Ctx, tc.u, tc.p)
		suite.Equal(tc.expectedCost, res, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestGetPunishmentCoeff() {
	testCases := []struct {
		testName         string
		appBandwidthInfo model.AppBandwidthInfo
		expectedP        string
	}{
		{
			testName: "test1",
			appBandwidthInfo: model.AppBandwidthInfo{
				ExpectedMPS:        sdk.NewDec(100),
				MessagesInCurBlock: 300,
			},
			expectedP: "1",
		},
		{
			testName: "test2",
			appBandwidthInfo: model.AppBandwidthInfo{
				ExpectedMPS:        sdk.NewDec(100),
				MessagesInCurBlock: 375,
			},
			expectedP: "2.013271178444183139",
		},
	}

	for _, tc := range testCases {
		appName := linotypes.AccountKey("test")
		err := suite.bm.storage.SetAppBandwidthInfo(suite.Ctx, appName, &tc.appBandwidthInfo)
		suite.Nil(err, "%s", tc.testName)

		res, getErr := suite.bm.GetPunishmentCoeff(suite.Ctx, appName)
		suite.Nil(getErr, "%s", tc.testName)

		expectedRes, err := sdk.NewDecFromStr(tc.expectedP)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(expectedRes, res, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestGetVacancyCoeff() {
	testCases := []struct {
		testName      string
		bandwidthInfo model.BandwidthInfo
		expectedU     string
	}{
		{
			testName: "test1",
			bandwidthInfo: model.BandwidthInfo{
				MaxMPS:    sdk.NewDec(1000),
				AppMsgEMA: sdk.NewDec(0),
			},
			expectedU: "0.501692631996395802",
		},
		{
			testName: "test2",
			bandwidthInfo: model.BandwidthInfo{
				MaxMPS:    sdk.NewDec(800),
				AppMsgEMA: sdk.NewDec(800),
			},
			expectedU: "1",
		},
	}

	for _, tc := range testCases {
		err := suite.bm.storage.SetBandwidthInfo(suite.Ctx, &tc.bandwidthInfo)
		suite.Nil(err, "%s", tc.testName)

		err = suite.bm.CalculateCurU(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)

		expectedRes, err := sdk.NewDecFromStr(tc.expectedU)
		suite.Nil(err, "%s", tc.testName)

		blockInfo, err := suite.bm.storage.GetBlockInfo(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(expectedRes, blockInfo.CurU, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestRefillAppBandwidthCredit() {
	testCases := []struct {
		testName         string
		appBandwidthInfo model.AppBandwidthInfo
		curTime          time.Time
		expectedInfo     model.AppBandwidthInfo
	}{
		{
			testName: "test1",
			appBandwidthInfo: model.AppBandwidthInfo{
				MaxBandwidthCredit: sdk.NewDec(3000),
				ExpectedMPS:        sdk.NewDec(200),
				CurBandwidthCredit: sdk.NewDec(0),
				LastRefilledAt:     suite.baseTime.Unix(),
			},
			curTime: suite.baseTime.Add(3 * time.Second),
			expectedInfo: model.AppBandwidthInfo{
				MaxBandwidthCredit: sdk.NewDec(3000),
				ExpectedMPS:        sdk.NewDec(200),
				CurBandwidthCredit: sdk.NewDec(600),
				LastRefilledAt:     suite.baseTime.Add(3 * time.Second).Unix(),
			},
		},
		{
			testName: "test2",
			appBandwidthInfo: model.AppBandwidthInfo{
				MaxBandwidthCredit: sdk.NewDec(3000),
				ExpectedMPS:        sdk.NewDec(200),
				CurBandwidthCredit: sdk.NewDec(0),
				LastRefilledAt:     suite.baseTime.Unix(),
			},
			curTime: suite.baseTime.Add(1000 * time.Second),
			expectedInfo: model.AppBandwidthInfo{
				MaxBandwidthCredit: sdk.NewDec(3000),
				ExpectedMPS:        sdk.NewDec(200),
				CurBandwidthCredit: sdk.NewDec(3000),
				LastRefilledAt:     suite.baseTime.Add(1000 * time.Second).Unix(),
			},
		},
		{
			testName: "test3",
			appBandwidthInfo: model.AppBandwidthInfo{
				MaxBandwidthCredit: sdk.NewDec(3000),
				ExpectedMPS:        sdk.NewDec(200),
				CurBandwidthCredit: sdk.NewDec(0),
				LastRefilledAt:     suite.baseTime.Add(3 * time.Second).Unix(),
			},
			curTime: suite.baseTime,
			expectedInfo: model.AppBandwidthInfo{
				MaxBandwidthCredit: sdk.NewDec(3000),
				ExpectedMPS:        sdk.NewDec(200),
				CurBandwidthCredit: sdk.NewDec(0),
				LastRefilledAt:     suite.baseTime.Add(3 * time.Second).Unix(),
			},
		},
	}

	for _, tc := range testCases {
		suite.Ctx = suite.Ctx.WithBlockHeader(abci.Header{Time: tc.curTime})
		appName := linotypes.AccountKey("test")
		err := suite.bm.storage.SetAppBandwidthInfo(suite.Ctx, appName, &tc.appBandwidthInfo)
		suite.Nil(err, "%s", tc.testName)

		err = suite.bm.RefillAppBandwidthCredit(suite.Ctx, appName)
		suite.Nil(err, "%s", tc.testName)

		res, err := suite.bm.storage.GetAppBandwidthInfo(suite.Ctx, appName)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(tc.expectedInfo, *res, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestConsumeBandwidthCredit() {
	testCases := []struct {
		testName     string
		appInfo      model.AppBandwidthInfo
		u            sdk.Dec
		p            sdk.Dec
		expectedInfo model.AppBandwidthInfo
	}{
		{
			testName: "test1",
			u:        linotypes.NewDecFromRat(5, 10),
			p:        linotypes.NewDecFromRat(5, 10),
			appInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("test"),
				MaxBandwidthCredit: sdk.NewDec(1000),
				CurBandwidthCredit: sdk.NewDec(500),
				MessagesInCurBlock: 100,
				ExpectedMPS:        sdk.NewDec(50),
				LastRefilledAt:     0,
			},
			expectedInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("test"),
				MaxBandwidthCredit: sdk.NewDec(1000),
				CurBandwidthCredit: sdk.NewDec(525),
				MessagesInCurBlock: 0,
				ExpectedMPS:        sdk.NewDec(50),
				LastRefilledAt:     0,
			},
		},
		{
			testName: "test2",
			u:        linotypes.NewDecFromRat(1, 1),
			p:        linotypes.NewDecFromRat(200, 1),
			appInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("test"),
				MaxBandwidthCredit: sdk.NewDec(1000),
				CurBandwidthCredit: sdk.NewDec(500),
				MessagesInCurBlock: 20,
				ExpectedMPS:        sdk.NewDec(50),
				LastRefilledAt:     0,
			},
			expectedInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("test"),
				MaxBandwidthCredit: sdk.NewDec(1000),
				CurBandwidthCredit: sdk.NewDec(-3480),
				MessagesInCurBlock: 0,
				ExpectedMPS:        sdk.NewDec(50),
				LastRefilledAt:     0,
			},
		},
	}

	for _, tc := range testCases {
		appName := linotypes.AccountKey("test")
		err := suite.bm.storage.SetAppBandwidthInfo(suite.Ctx, appName, &tc.appInfo)
		suite.Nil(err, "%s", tc.testName)

		err = suite.bm.ConsumeBandwidthCredit(suite.Ctx, tc.u, tc.p, appName)
		suite.Nil(err, "%s", tc.testName)

		res, err := suite.bm.storage.GetAppBandwidthInfo(suite.Ctx, appName)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(tc.expectedInfo, *res, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestPrecheckAndConsumeBandwidthCredit() {
	testCases := []struct {
		testName        string
		appInfo         model.AppBandwidthInfo
		blockInfo       model.BlockInfo
		expectedAppInfo model.AppBandwidthInfo
		expectedErr     sdk.Error
	}{
		{
			testName: "test1",
			blockInfo: model.BlockInfo{
				CurU: sdk.NewDec(1),
			},
			appInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("test"),
				MaxBandwidthCredit: sdk.NewDec(1000),
				CurBandwidthCredit: sdk.NewDec(500),
				MessagesInCurBlock: 100,
				ExpectedMPS:        sdk.NewDec(50),
				LastRefilledAt:     0,
			},
			expectedAppInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("test"),
				MaxBandwidthCredit: sdk.NewDec(1000),
				CurBandwidthCredit: sdk.NewDec(499),
				MessagesInCurBlock: 100,
				ExpectedMPS:        sdk.NewDec(50),
				LastRefilledAt:     0,
			},
			expectedErr: nil,
		},
		{
			testName: "test2",
			blockInfo: model.BlockInfo{
				CurU: sdk.NewDec(2),
			},
			appInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("test"),
				MaxBandwidthCredit: sdk.NewDec(1000),
				CurBandwidthCredit: sdk.NewDec(1),
				MessagesInCurBlock: 100,
				ExpectedMPS:        sdk.NewDec(50),
				LastRefilledAt:     0,
			},
			expectedAppInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("test"),
				MaxBandwidthCredit: sdk.NewDec(1000),
				CurBandwidthCredit: sdk.NewDec(1),
				MessagesInCurBlock: 100,
				ExpectedMPS:        sdk.NewDec(50),
				LastRefilledAt:     0,
			},
			expectedErr: types.ErrAppBandwidthNotEnough(),
		},
	}

	for _, tc := range testCases {
		appName := linotypes.AccountKey("test")
		err := suite.bm.storage.SetAppBandwidthInfo(suite.Ctx, appName, &tc.appInfo)
		suite.Nil(err, "%s", tc.testName)

		err = suite.bm.storage.SetBlockInfo(suite.Ctx, &tc.blockInfo)
		suite.Nil(err, "%s", tc.testName)

		err = suite.bm.PrecheckAndConsumeBandwidthCredit(suite.Ctx, appName)
		suite.Equal(tc.expectedErr, err, "%s", tc.testName)

		appInfo, err := suite.bm.storage.GetAppBandwidthInfo(suite.Ctx, appName)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(tc.expectedAppInfo, *appInfo, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestReCalculateAppBandwidthInfo() {
	testCases := []struct {
		testName         string
		prevAppYInfo     model.AppBandwidthInfo
		expectedAppXInfo model.AppBandwidthInfo
		expectedAppYInfo model.AppBandwidthInfo
	}{
		{
			testName: "test1",
			prevAppYInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("AppY"),
				MaxBandwidthCredit: sdk.NewDec(800 * 10),
				CurBandwidthCredit: sdk.NewDec(0),
				MessagesInCurBlock: 0,
				ExpectedMPS:        sdk.NewDec(800),
				LastRefilledAt:     suite.baseTime.Unix(),
			},
			expectedAppXInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("AppX"),
				MaxBandwidthCredit: sdk.NewDec(80 * 10),
				CurBandwidthCredit: sdk.NewDec(80 * 10),
				MessagesInCurBlock: 0,
				ExpectedMPS:        sdk.NewDec(80),
				LastRefilledAt:     suite.Ctx.BlockHeader().Time.Unix(),
			},
			expectedAppYInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("AppY"),
				MaxBandwidthCredit: sdk.NewDec(720 * 10),
				CurBandwidthCredit: sdk.NewDec(2400),
				MessagesInCurBlock: 0,
				ExpectedMPS:        sdk.NewDec(720),
				LastRefilledAt:     suite.Ctx.BlockHeader().Time.Unix(),
			},
		},
	}

	for _, tc := range testCases {
		appX := linotypes.AccountKey("AppX")
		appY := linotypes.AccountKey("AppY")

		err := suite.bm.storage.SetAppBandwidthInfo(suite.Ctx, appY, &tc.prevAppYInfo)
		suite.Nil(err, "%s", tc.testName)

		err = suite.bm.ReCalculateAppBandwidthInfo(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)

		appXInfo, err := suite.bm.storage.GetAppBandwidthInfo(suite.Ctx, appX)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(tc.expectedAppXInfo, *appXInfo, "%s", tc.testName)

		appYInfo, err := suite.bm.storage.GetAppBandwidthInfo(suite.Ctx, appY)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(tc.expectedAppYInfo, *appYInfo, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestCheckAppBandwidth() {
	testCases := []struct {
		testName          string
		fee               auth.StdFee
		prevAppYInfo      model.AppBandwidthInfo
		expectedAppYInfo  model.AppBandwidthInfo
		expectedBlockInfo model.BlockInfo
		expectedErr       sdk.Error
	}{
		{
			testName: "test1",
			fee:      auth.StdFee{},
			prevAppYInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("AppY"),
				MaxBandwidthCredit: sdk.NewDec(800 * 10),
				CurBandwidthCredit: sdk.NewDec(4000),
				MessagesInCurBlock: 0,
				ExpectedMPS:        sdk.NewDec(800),
				LastRefilledAt:     0,
			},
			expectedAppYInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("AppY"),
				MaxBandwidthCredit: sdk.NewDec(800 * 10),
				CurBandwidthCredit: sdk.NewDec(7999),
				MessagesInCurBlock: 1,
				ExpectedMPS:        sdk.NewDec(800),
				LastRefilledAt:     suite.Ctx.BlockHeader().Time.Unix(),
			},
			expectedBlockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  1,
				TotalMsgSignedByUser: 0,
				CurMsgFee:            linotypes.NewCoinFromInt64(0),
				CurU:                 sdk.NewDec(1),
			},
			expectedErr: nil,
		},
		{
			testName: "test2",
			fee:      auth.StdFee{},
			prevAppYInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("AppY"),
				MaxBandwidthCredit: sdk.NewDec(800 * 10),
				CurBandwidthCredit: sdk.NewDec(0),
				MessagesInCurBlock: 0,
				ExpectedMPS:        sdk.NewDec(800),
				LastRefilledAt:     suite.Ctx.BlockHeader().Time.Unix(),
			},
			expectedAppYInfo: model.AppBandwidthInfo{
				Username:           linotypes.AccountKey("AppY"),
				MaxBandwidthCredit: sdk.NewDec(800 * 10),
				CurBandwidthCredit: sdk.NewDec(0),
				MessagesInCurBlock: 0,
				ExpectedMPS:        sdk.NewDec(800),
				LastRefilledAt:     suite.Ctx.BlockHeader().Time.Unix(),
			},
			expectedBlockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  1,
				TotalMsgSignedByUser: 0,
				CurMsgFee:            linotypes.NewCoinFromInt64(0),
				CurU:                 sdk.NewDec(1),
			},
			expectedErr: types.ErrAppBandwidthNotEnough(),
		},
	}

	for _, tc := range testCases {
		appY := linotypes.AccountKey("AppY")
		err := suite.bm.storage.SetAppBandwidthInfo(suite.Ctx, appY, &tc.prevAppYInfo)
		suite.Nil(err, "%s", tc.testName)

		err = suite.bm.CheckBandwidth(suite.Ctx, appY, tc.fee)
		suite.Equal(tc.expectedErr, err, "%s", tc.testName)

		appYInfo, err := suite.bm.storage.GetAppBandwidthInfo(suite.Ctx, appY)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(tc.expectedAppYInfo, *appYInfo, "%s", tc.testName)

		blockInfo, err := suite.bm.storage.GetBlockInfo(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(tc.expectedBlockInfo, *blockInfo, "%s", tc.testName)
	}
}

func (suite *BandwidthManagerTestSuite) TestCheckMsgFee() {
	testCases := []struct {
		testName          string
		username          linotypes.AccountKey
		fee               auth.StdFee
		blockInfo         model.BlockInfo
		expectedBlockInfo model.BlockInfo
		expectedErr       sdk.Error
	}{
		{
			testName: "test1",
			fee:      auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin(linotypes.LinoCoinDenom, sdk.NewInt(1000)))},
			username: linotypes.AccountKey("UserX"),
			blockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 0,
				CurMsgFee:            linotypes.NewCoinFromInt64(100),
				CurU:                 sdk.NewDec(1),
			},
			expectedBlockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 1,
				CurMsgFee:            linotypes.NewCoinFromInt64(100),
				CurU:                 sdk.NewDec(1),
			},
			expectedErr: nil,
		},
		{
			testName: "test2",
			fee:      auth.StdFee{Amount: sdk.NewCoins(sdk.NewCoin(linotypes.LinoCoinDenom, sdk.NewInt(10)))},
			username: linotypes.AccountKey("UserX"),
			blockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 0,
				CurMsgFee:            linotypes.NewCoinFromInt64(100),
				CurU:                 sdk.NewDec(1),
			},
			expectedBlockInfo: model.BlockInfo{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 0,
				CurMsgFee:            linotypes.NewCoinFromInt64(100),
				CurU:                 sdk.NewDec(1),
			},
			expectedErr: types.ErrUserMsgFeeNotEnough(),
		},
	}

	for _, tc := range testCases {
		err := suite.bm.storage.SetBlockInfo(suite.Ctx, &tc.blockInfo)
		suite.Nil(err, "%s", tc.testName)

		err = suite.bm.CheckBandwidth(suite.Ctx, tc.username, tc.fee)
		suite.Equal(tc.expectedErr, err, "%s", tc.testName)

		blockInfo, err := suite.bm.storage.GetBlockInfo(suite.Ctx)
		suite.Nil(err, "%s", tc.testName)
		suite.Equal(tc.expectedBlockInfo, *blockInfo, "%s", tc.testName)

	}
}
