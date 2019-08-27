package bandwidth

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/x/bandwidth/model"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAddMsgSignedByUser(t *testing.T) {
	ctx, bm := setupTest(t, 1)

	testCases := []struct {
		testName           string
		amount             uint32
		expectCurBlockInfo model.CurBlockInfo
	}{
		{
			testName: "add user signed message",
			amount:   1,
			expectCurBlockInfo: model.CurBlockInfo{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 1,
				CurMsgFee:            sdk.NewDec(0),
			},
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: time.Now()})
		bm.AddMsgSignedByUser(ctx, tc.amount)
		checkCurBlockInfo(t, ctx, tc.testName, bm, tc.expectCurBlockInfo)
	}
}

func TestAddMsgSignedByApp(t *testing.T) {
	ctx, bm := setupTest(t, 1)

	testCases := []struct {
		testName           string
		amount             uint32
		expectCurBlockInfo model.CurBlockInfo
	}{
		{
			testName: "add user signed message",
			amount:   1,
			expectCurBlockInfo: model.CurBlockInfo{
				TotalMsgSignedByApp:  1,
				TotalMsgSignedByUser: 0,
				CurMsgFee:            sdk.NewDec(0),
			},
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: time.Now()})
		bm.AddMsgSignedByApp(ctx, tc.amount)
		checkCurBlockInfo(t, ctx, tc.testName, bm, tc.expectCurBlockInfo)
	}
}

func TestApproximateExp(t *testing.T) {
	_, bm := setupTest(t, 1)
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
		res := bm.approximateExp(tc.x)
		assert.Equal(t, tc.expectedRes, res, "%s: diff exp result, got %v, want %v", tc.testName, res, tc.expectedRes)
	}
}

func TestCalculateCurMsgFee(t *testing.T) {
	ctx, bm := setupTest(t, 1)

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
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: time.Now()})
		bm.storage.SetBandwidthInfo(ctx, &tc.bandwidthInfo)
		bm.CalculateCurMsgFee(ctx)
		expectedFee, err := sdk.NewDecFromStr(tc.expectMessageFee)
		assert.Nil(t, err)

		info, err := bm.storage.GetCurBlockInfo(ctx)
		assert.Nil(t, err)
		assert.Equal(t, expectedFee, info.CurMsgFee, "%s: diff exp result, got %v, want %v", tc.testName, info.CurMsgFee, expectedFee)
	}
}

// only one step
func TestUpdateMaxMPSAndEMA(t *testing.T) {
	ctx, bm := setupTest(t, 1)
	baseTime := time.Now()
	testCases := []struct {
		testName         string
		curBlockInfo     model.CurBlockInfo
		bandwidthInfo    model.BandwidthInfo
		lastBlockTime    int64
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
			curBlockInfo: model.CurBlockInfo{
				TotalMsgSignedByApp:  0,
				TotalMsgSignedByUser: 60,
			},
			lastBlockTime:    baseTime.Unix(),
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
			curBlockInfo: model.CurBlockInfo{
				TotalMsgSignedByApp:  270,
				TotalMsgSignedByUser: 0,
			},
			lastBlockTime:    baseTime.Unix(),
			expectGeneralEMA: "0",
			expectAppEMA:     "54",
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: baseTime.Add(3 * time.Second)})
		bm.storage.SetBandwidthInfo(ctx, &tc.bandwidthInfo)
		bm.storage.SetCurBlockInfo(ctx, &tc.curBlockInfo)

		bm.UpdateMaxMPSAndEMA(ctx, tc.lastBlockTime)
		expectedGeneralEMA, err := sdk.NewDecFromStr(tc.expectGeneralEMA)
		assert.Nil(t, err)

		expectedAppEMA, err := sdk.NewDecFromStr(tc.expectAppEMA)
		assert.Nil(t, err)

		info, err := bm.storage.GetBandwidthInfo(ctx)
		assert.Nil(t, err)
		assert.Equal(t, expectedGeneralEMA, info.GeneralMsgEMA, "%s: diff general EMA result, got %v, want %v", tc.testName, info.GeneralMsgEMA, expectedGeneralEMA)
		assert.Equal(t, expectedAppEMA, info.AppMsgEMA, "%s: diff app EMA result, got %v, want %v", tc.testName, info.AppMsgEMA, expectedAppEMA)
	}
}
