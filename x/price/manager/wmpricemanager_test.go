package manager

import (
	"fmt"
	"sort"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/lino-network/lino/param"
	mparam "github.com/lino-network/lino/param/mocks"
	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/testutils"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/price/model"
	"github.com/lino-network/lino/x/price/types"
	mval "github.com/lino-network/lino/x/validator/mocks"
	valmodel "github.com/lino-network/lino/x/validator/model"
)

type ValidatorAndVotes = valmodel.ReceivedVotesStatus

var (
	storeKeyStr = "priceStoreTestKey"
	storeKey    = sdk.NewKVStoreKey(storeKeyStr)

	genesisPrice = linotypes.NewMiniDollar(1200)
)

type PriceStoreDumper struct{}

func (dumper PriceStoreDumper) NewDumper() *testutils.Dumper {
	return model.NewPriceDumper(model.NewPriceStorage(storeKey))
}

type WMPriceManagerSuite struct {
	testsuites.GoldenTestSuite
	manager WeightedMedianPriceManager
	// read only, can be reseted at will.
	mParam *mparam.ParamKeeper
	mVal   *mval.ValidatorKeeper
}

func NewPriceManagerSuite() *WMPriceManagerSuite {
	return &WMPriceManagerSuite{
		GoldenTestSuite: testsuites.NewGoldenTestSuite(PriceStoreDumper{}, storeKey),
	}
}

func (suite *WMPriceManagerSuite) SetupTest() {
	suite.mParam = new(mparam.ParamKeeper)
	suite.mVal = new(mval.ValidatorKeeper)

	suite.manager = NewWeightedMedianPriceManager(
		storeKey, suite.mVal, suite.mParam)
	suite.SetupCtx(0, time.Unix(0, 0), storeKey)
}

func (suite *WMPriceManagerSuite) setValidators(vals []ValidatorAndVotes) {
	suite.mVal = new(mval.ValidatorKeeper)
	suite.manager.val = suite.mVal
	suite.mVal.On("GetCommittingValidatorVoteStatus", mock.Anything).Return(vals)
	names := toValNames(vals)
	suite.mVal.On("GetCommittingValidators", mock.Anything).Return(names)
}

// setup validators from "val1", "val2" ...."valx".
func (suite *WMPriceManagerSuite) setValidatorByDist(votes ...int64) {
	vals := []ValidatorAndVotes{}
	for i, v := range votes {
		vals = append(vals, ValidatorAndVotes{
			ValidatorName: linotypes.AccountKey(fmt.Sprintf("val%d", i+1)),
			ReceivedVotes: linotypes.NewCoinFromInt64(v),
		})
	}
	suite.setValidators(vals)
}

func (suite *WMPriceManagerSuite) setParam(priceParam *param.PriceParam) {
	suite.mParam.On("GetPriceParam", mock.Anything).Return(priceParam).Maybe()
}

var (
	basicParam = &param.PriceParam{
		TestnetMode:     false,
		UpdateEverySec:  int64(1 * time.Hour.Seconds()),
		FeedEverySec:    int64(10 * time.Minute.Seconds()),
		HistoryMaxLen:   71,
		PenaltyMissFeed: linotypes.NewCoinFromInt64(10000 * linotypes.Decimals),
	}
)

func (suite *WMPriceManagerSuite) setBasicParam(testnet bool) {
	suite.mParam = new(mparam.ParamKeeper)
	suite.manager.param = suite.mParam
	if testnet {
		param := *basicParam
		param.TestnetMode = true
		suite.setParam(&param)
	} else {
		suite.setParam(basicParam)
	}
}

func TestPriceManagerSuite(t *testing.T) {
	suite.Run(t, NewPriceManagerSuite())
}

func (suite *WMPriceManagerSuite) TestInitGenesis() {
	suite.setParam(&param.PriceParam{
		TestnetMode: false,
	})
	err := suite.manager.InitGenesis(suite.Ctx, linotypes.NewMiniDollar(-100))
	suite.NotNil(err)
	initPrice := linotypes.NewMiniDollar(1234)
	err = suite.manager.InitGenesis(suite.Ctx, initPrice)
	suite.Nil(err)
	price, err := suite.manager.CurrPrice(suite.Ctx)
	suite.NoError(err)
	suite.Equal(initPrice, price)
	suite.Golden()
}

type feedAction struct {
	feeder linotypes.AccountKey
	t      time.Time
	price  linotypes.MiniDollar
	err    sdk.Error
}

func (suite *WMPriceManagerSuite) TestFeedPrice() {
	suite.setBasicParam(false)
	feedInterval := basicParam.FeedEverySec
	testCases := []struct {
		name    string
		valDist []int64
		actions []feedAction
		succ    bool
	}{
		{
			name:    "invalid price",
			valDist: []int64{100, 100, 100},
			actions: []feedAction{
				{
					feeder: "val1",
					t:      time.Unix(0, 0),
					price:  linotypes.NewMiniDollar(-100),
					err:    types.ErrInvalidPriceFeed(linotypes.NewMiniDollar(-100)),
				},
			},
			succ: false,
		},
		{
			name:    "feeder not exist",
			valDist: []int64{100, 100, 100},
			actions: []feedAction{
				{
					feeder: "val100",
					t:      time.Unix(0, 0),
					price:  linotypes.NewMiniDollar(100),
					err:    types.ErrNotAValidator("val100"),
				},
			},
			succ: false,
		},
		{
			name:    "3rd time rate limited",
			valDist: []int64{100, 100, 100},
			actions: []feedAction{
				{
					feeder: "val2",
					t:      time.Unix(0, 0),
					price:  linotypes.NewMiniDollar(100),
					err:    nil,
				},
				{
					feeder: "val2",
					t:      time.Unix(0+feedInterval, 0),
					price:  linotypes.NewMiniDollar(200),
					err:    nil,
				},
				{
					feeder: "val2",
					t:      time.Unix(0+feedInterval*2-1, 0),
					price:  linotypes.NewMiniDollar(200),
					err:    types.ErrPriceFeedRateLimited(),
				},
			},
			succ: true,
		},
		{
			name:    "succ two round",
			valDist: []int64{100, 100, 100},
			actions: []feedAction{
				{
					feeder: "val1",
					t:      time.Unix(0, 0),
					price:  linotypes.NewMiniDollar(100),
				},
				{
					feeder: "val2",
					t:      time.Unix(0, 0),
					price:  linotypes.NewMiniDollar(200),
				},
				{
					feeder: "val3",
					t:      time.Unix(0, 0),
					price:  linotypes.NewMiniDollar(300),
				},
				{
					feeder: "val1",
					t:      time.Unix(0+feedInterval, 0),
					price:  linotypes.NewMiniDollar(400),
				},
				{
					feeder: "val2",
					t:      time.Unix(0+feedInterval, 0),
					price:  linotypes.NewMiniDollar(500),
				},
				{
					feeder: "val3",
					t:      time.Unix(0+feedInterval, 0),
					price:  linotypes.NewMiniDollar(600),
				},
			},
			succ: true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.LoadState(false, "genesis")
			suite.setValidatorByDist(tc.valDist...)
			for _, act := range tc.actions {
				suite.NextBlock(act.t)
				err := suite.manager.FeedPrice(suite.Ctx, act.feeder, act.price)
				suite.Require().Equal(act.err, err)
			}
			if tc.succ {
				suite.Golden()
			}
		})
	}
}

type feedRound struct {
	updateTime time.Time
	valDist    []int64
	actions    []feedAction
	slashes    []linotypes.AccountKey
	priceAfter linotypes.MiniDollar
}

// validators get slashed if not feeding price on time.
func (suite *WMPriceManagerSuite) TestUpdatePriceSlash() {
	suite.setBasicParam(false)
	feedInterval := basicParam.FeedEverySec
	updateInterval := basicParam.UpdateEverySec
	testCases := []struct {
		name      string
		rounds    []feedRound
		isTestnet bool
	}{
		{
			name: "no validator feed at init and no slash",
			rounds: []feedRound{
				{
					updateTime: time.Unix(updateInterval*1, 0),
					valDist:    []int64{100, 100, 100},
					actions:    []feedAction{},
					slashes:    nil,
				},
			},
		},
		{
			name: "newly joined no feed no slash",
			rounds: []feedRound{
				{
					updateTime: time.Unix(updateInterval*1, 0),
					valDist:    []int64{100, 100},
					actions:    []feedAction{},
					slashes:    nil,
				},
				{
					updateTime: time.Unix(updateInterval*2, 0),
					valDist:    []int64{100, 100, 100},
					actions: []feedAction{
						{
							feeder: "val1",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(10),
						},
						{
							feeder: "val2",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(10),
						},
					},
					slashes: nil,
				},
			},
		},
		{
			name: "no feed get slashed",
			rounds: []feedRound{
				{
					updateTime: time.Unix(updateInterval, 0),
					valDist:    []int64{100, 100, 100},
					actions:    []feedAction{},
					slashes:    nil,
				},
				{
					updateTime: time.Unix(updateInterval*2, 0),
					valDist:    []int64{100, 100, 100},
					actions: []feedAction{
						{
							feeder: "val1",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(10),
						},
						{
							feeder: "val2",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(10),
						},
					},
					slashes: []linotypes.AccountKey{"val3"},
				},
				{
					updateTime: time.Unix(updateInterval*3, 0),
					valDist:    []int64{100, 100, 100},
					actions: []feedAction{
						{
							// too frequent, does not count.
							feeder: "val1",
							t:      time.Unix(updateInterval+feedInterval+1, 0),
							price:  linotypes.NewMiniDollar(10),
							err:    types.ErrPriceFeedRateLimited(),
						},
						{
							feeder: "val2",
							t:      time.Unix(updateInterval*2+feedInterval, 0),
							price:  linotypes.NewMiniDollar(10),
						},
						{
							feeder: "val3",
							t:      time.Unix(updateInterval*3+feedInterval, 0),
							price:  linotypes.NewMiniDollar(10),
						},
					},
					slashes: []linotypes.AccountKey{"val1"},
				},
			},
		},
		{
			name: "all slashed",
			rounds: []feedRound{
				{
					updateTime: time.Unix(updateInterval, 0),
					valDist:    []int64{100, 100, 100},
					actions:    []feedAction{},
					slashes:    nil,
				},
				{
					updateTime: time.Unix(updateInterval*2, 0),
					valDist:    []int64{100, 100, 100},
					actions: []feedAction{
						{
							feeder: "val1",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(10),
						},
						{
							feeder: "val2",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(10),
						},
					},
					slashes: []linotypes.AccountKey{"val3"},
				},
				{
					updateTime: time.Unix(updateInterval*3, 0),
					valDist:    []int64{100, 100, 100},
					actions:    []feedAction{},
					slashes:    []linotypes.AccountKey{"val1", "val2", "val3"},
				},
			},
		},
		{
			name:      "testnet no slash",
			isTestnet: true,
			rounds: []feedRound{
				{
					updateTime: time.Unix(updateInterval, 0),
					valDist:    []int64{100, 100, 100},
					actions:    []feedAction{},
					slashes:    nil,
				},
				{
					updateTime: time.Unix(updateInterval*2, 0),
					valDist:    []int64{100, 100, 100},
					actions:    []feedAction{},
					slashes:    []linotypes.AccountKey{},
				},
				{
					updateTime: time.Unix(updateInterval*3, 0),
					valDist:    []int64{100, 100, 100},
					actions:    []feedAction{},
					slashes:    []linotypes.AccountKey{},
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.LoadState(false, "genesis")
			if tc.isTestnet {
				suite.setBasicParam(true)
			}
			for _, round := range tc.rounds {
				suite.setValidatorByDist(round.valDist...)
				for _, slash := range round.slashes {
					suite.mVal.On("PunishCommittingValidator",
						mock.Anything,
						slash,
						basicParam.PenaltyMissFeed,
						linotypes.PunishNoPriceFed).Return(nil).Once()
				}
				for _, act := range round.actions {
					suite.NextBlock(act.t)
					err := suite.manager.FeedPrice(suite.Ctx, act.feeder, act.price)
					suite.Equal(act.err, err)
				}
				suite.NextBlock(round.updateTime)
				err := suite.manager.UpdatePrice(suite.Ctx)
				suite.Nil(err)
				suite.mVal.AssertExpectations(suite.T())
			}
			suite.Golden()
		})
	}
}

// current price is correct.
func (suite *WMPriceManagerSuite) TestUpdatePriceCurrPrice() {
	suite.setBasicParam(false)
	feedInterval := basicParam.FeedEverySec
	updateInterval := basicParam.UpdateEverySec
	testCases := []struct {
		name      string
		rounds    []feedRound
		isTestnet bool
	}{
		{
			name: "initial price is kept no feed",
			rounds: []feedRound{
				{
					updateTime: time.Unix(updateInterval*1, 0),
					valDist:    []int64{100, 100, 100},
					actions:    []feedAction{},
					priceAfter: genesisPrice,
				},
				{
					updateTime: time.Unix(updateInterval*2, 0),
					valDist:    []int64{100, 100, 100, 100, 100},
					actions:    []feedAction{},
					priceAfter: genesisPrice,
				},
				{
					updateTime: time.Unix(updateInterval*3, 0),
					valDist:    []int64{100, 100, 100, 100, 100, 100, 100},
					actions:    []feedAction{},
					priceAfter: genesisPrice,
				},
			},
		},
		{
			name: "weighted median of fed price case 1",
			rounds: []feedRound{
				{
					updateTime: time.Unix(updateInterval*1, 0),
					valDist:    []int64{100, 100, 100},
					actions: []feedAction{
						{
							feeder: "val1",
							t:      time.Unix(feedInterval, 0),
							price:  linotypes.NewMiniDollar(10),
						},
						{
							feeder: "val2",
							t:      time.Unix(feedInterval, 0),
							price:  linotypes.NewMiniDollar(20),
						},
						{
							feeder: "val3",
							t:      time.Unix(feedInterval, 0),
							price:  linotypes.NewMiniDollar(30),
						},
					},
					priceAfter: genesisPrice,
				},
				{
					updateTime: time.Unix(updateInterval*2, 0),
					valDist:    []int64{100, 100, 100},
					actions: []feedAction{
						{
							feeder: "val1",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(10),
						},
						{
							feeder: "val2",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(20),
						},
						{
							feeder: "val3",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(30),
						},
					},
					priceAfter: linotypes.NewMiniDollar(20),
				},
			},
		},
		{
			name: "weighted median of fed price case 2",
			rounds: []feedRound{
				{
					updateTime: time.Unix(updateInterval*1, 0),
					valDist:    []int64{10, 10, 100},
					actions: []feedAction{
						{
							feeder: "val1",
							t:      time.Unix(feedInterval, 0),
							price:  linotypes.NewMiniDollar(588),
						},
						{
							feeder: "val2",
							t:      time.Unix(feedInterval, 0),
							price:  linotypes.NewMiniDollar(3451),
						},
						{
							feeder: "val3",
							t:      time.Unix(feedInterval, 0),
							price:  linotypes.NewMiniDollar(1234567),
						},
					},
					priceAfter: linotypes.NewMiniDollar(1234567),
				},
				{
					updateTime: time.Unix(updateInterval*2, 0),
					valDist:    []int64{1000000, 10, 1},
					actions: []feedAction{
						{
							feeder: "val1",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(1201),
						},
						{
							feeder: "val2",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(20),
						},
						{
							feeder: "val3",
							t:      time.Unix(updateInterval+feedInterval, 0),
							price:  linotypes.NewMiniDollar(1234567),
						},
					},
					priceAfter: linotypes.NewMiniDollar(1201),
				},
			},
		},
		{
			name:      "testnet price fixed",
			isTestnet: true,
			rounds: []feedRound{
				{
					updateTime: time.Unix(updateInterval, 0),
					valDist:    []int64{100, 100, 100},
					actions: []feedAction{
						{
							feeder: "val1",
							t:      time.Unix(feedInterval, 0),
							price:  linotypes.NewMiniDollar(1000),
						},
						{
							feeder: "val2",
							t:      time.Unix(feedInterval, 0),
							price:  linotypes.NewMiniDollar(1000),
						},
						{
							feeder: "val3",
							t:      time.Unix(feedInterval, 0),
							price:  linotypes.NewMiniDollar(1000),
						},
					},
					priceAfter: linotypes.TestnetPrice,
				},
				{
					updateTime: time.Unix(updateInterval*2, 0),
					valDist:    []int64{100, 100, 100},
					actions: []feedAction{
						{
							feeder: "val1",
							t:      time.Unix(updateInterval*1+feedInterval, 0),
							price:  linotypes.NewMiniDollar(2000),
						},
						{
							feeder: "val2",
							t:      time.Unix(updateInterval*1+feedInterval, 0),
							price:  linotypes.NewMiniDollar(2000),
						},
						{
							feeder: "val3",
							t:      time.Unix(updateInterval*1+feedInterval, 0),
							price:  linotypes.NewMiniDollar(2000),
						},
					},
					priceAfter: linotypes.TestnetPrice,
				},
				{
					updateTime: time.Unix(updateInterval*3, 0),
					valDist:    []int64{100, 100, 100},
					actions: []feedAction{
						{
							feeder: "val1",
							t:      time.Unix(updateInterval*2+feedInterval, 0),
							price:  linotypes.NewMiniDollar(3000),
						},
						{
							feeder: "val2",
							t:      time.Unix(updateInterval*2+feedInterval, 0),
							price:  linotypes.NewMiniDollar(3000),
						},
						{
							feeder: "val3",
							t:      time.Unix(updateInterval*2+feedInterval, 0),
							price:  linotypes.NewMiniDollar(3000),
						},
					},
					priceAfter: linotypes.TestnetPrice,
				},
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.LoadState(false, "genesis")
			if tc.isTestnet {
				suite.setBasicParam(true)
			}
			for _, round := range tc.rounds {
				suite.setValidatorByDist(round.valDist...)
				// slashing is not tested here.
				suite.mVal.On("PunishCommittingValidator",
					mock.Anything,
					mock.Anything,
					basicParam.PenaltyMissFeed,
					linotypes.PunishNoPriceFed).Return(nil).Maybe()

				for _, act := range round.actions {
					suite.NextBlock(act.t)
					err := suite.manager.FeedPrice(suite.Ctx, act.feeder, act.price)
					suite.Equal(act.err, err, "%s", err)
				}
				suite.NextBlock(round.updateTime)
				err := suite.manager.UpdatePrice(suite.Ctx)
				suite.Nil(err)
				price, err := suite.manager.CurrPrice(suite.Ctx)
				suite.Nil(err)
				suite.Equal(round.priceAfter, price)
				suite.mVal.AssertExpectations(suite.T())
			}
			suite.Golden()
		})
	}
}

// current price is correct.
func (suite *WMPriceManagerSuite) TestUpdatePriceHistoryRolling() {
	suite.setBasicParam(false)
	feedInterval := basicParam.FeedEverySec
	updateInterval := basicParam.UpdateEverySec

	history := []int64{genesisPrice.Int64()}
	// 100 validators, weighted from 1..100
	var validators []int64
	for i := 1; i <= 100; i++ {
		validators = append(validators, int64(i))
	}
	suite.setValidatorByDist(validators...)

	rounds := []feedRound{}
	for i := 0; i < 300; i++ {
		actions := []feedAction{}
		for j := 0; j < 100; j++ {
			actions = append(actions, feedAction{
				feeder: linotypes.AccountKey(fmt.Sprintf("val%d", j+1)),
				t:      time.Unix(updateInterval*int64(i)+feedInterval, 0),
				price:  linotypes.NewMiniDollar(int64(i + j + 1)),
			})
		}
		mid := actions[70].price.Int64() // sum 1..70 = 2484, mid weight = 2525, the 71th is mid.
		history = append(history, mid)
		min := basicParam.HistoryMaxLen
		if len(history) < min {
			min = len(history)
		}
		historyTail := append([]int64(nil), history[len(history)-min:]...)
		sort.SliceStable(historyTail,
			func(i, j int) bool { return historyTail[i] < historyTail[j] })
		rounds = append(rounds, feedRound{
			updateTime: time.Unix(updateInterval*int64(i+1), 0),
			actions:    actions,
			priceAfter: linotypes.NewMiniDollar(historyTail[len(historyTail)/2]),
		})
	}

	suite.LoadState(false, "genesis")
	for _, round := range rounds {
		for _, act := range round.actions {
			suite.NextBlock(act.t)
			err := suite.manager.FeedPrice(suite.Ctx, act.feeder, act.price)
			suite.Equal(act.err, err, "%s", err)
		}
		suite.NextBlock(round.updateTime)
		err := suite.manager.UpdatePrice(suite.Ctx)
		suite.Nil(err)
		price, err := suite.manager.CurrPrice(suite.Ctx)
		suite.Nil(err)
		suite.Require().Equal(round.priceAfter, price, "%s != %s", round.priceAfter, price)
	}
	suite.Golden()
}

func (suite *WMPriceManagerSuite) TestCoinToMiniDollar() {
	suite.LoadState(false, "genesis")
	suite.setBasicParam(true)
	testCases := []struct {
		testName string
		coin     linotypes.Coin
		expected linotypes.MiniDollar
	}{
		{
			testName: "0 Coin",
			coin:     linotypes.NewCoinFromInt64(0),
			expected: linotypes.NewMiniDollar(0),
		},
		{
			testName: "1 Coin",
			coin:     linotypes.NewCoinFromInt64(1),
			expected: linotypes.NewMiniDollar(1200),
		},
		{
			testName: "3 Coin",
			coin:     linotypes.NewCoinFromInt64(3),
			expected: linotypes.NewMiniDollar(3600),
		},
		{
			testName: "1 LNO",
			coin:     linotypes.NewCoinFromInt64(linotypes.Decimals),
			expected: linotypes.NewMiniDollar(12 * 10000000),
		},
		{
			testName: "2 LNO",
			coin:     linotypes.NewCoinFromInt64(2 * linotypes.Decimals),
			expected: linotypes.NewMiniDollar(24 * 10000000),
		},
		{
			testName: "1000000 LNO",
			coin:     linotypes.NewCoinFromInt64(1000000 * linotypes.Decimals),
			expected: linotypes.NewMiniDollar(1000000 * 12 * 10000000),
		},
	}

	for _, tc := range testCases {
		rst, err := suite.manager.CoinToMiniDollar(suite.Ctx, tc.coin)
		suite.Nil(err)
		suite.Equal(tc.expected, rst, "%s", tc.testName)
	}
}

func (suite *WMPriceManagerSuite) TestMiniDollarToCoin() {
	suite.LoadState(false, "genesis")
	suite.setBasicParam(true)
	testCases := []struct {
		testName       string
		mini           linotypes.MiniDollar
		expectedBought linotypes.Coin
		expectedUsed   linotypes.MiniDollar
	}{
		{
			testName:       "0 minidollar",
			mini:           linotypes.NewMiniDollar(0),
			expectedBought: linotypes.NewCoinFromInt64(0),
			expectedUsed:   linotypes.NewMiniDollar(0),
		},
		{
			testName:       "1 minidollar",
			mini:           linotypes.NewMiniDollar(1),
			expectedBought: linotypes.NewCoinFromInt64(0),
			expectedUsed:   linotypes.NewMiniDollar(0),
		},
		{
			testName:       "1199 minidollar",
			mini:           linotypes.NewMiniDollar(1199),
			expectedBought: linotypes.NewCoinFromInt64(0),
			expectedUsed:   linotypes.NewMiniDollar(0),
		},
		{
			testName:       "1200 minidollar",
			mini:           linotypes.NewMiniDollar(1200),
			expectedBought: linotypes.NewCoinFromInt64(1),
			expectedUsed:   linotypes.NewMiniDollar(1200),
		},
		{
			testName:       "1201 minidollar",
			mini:           linotypes.NewMiniDollar(1201),
			expectedBought: linotypes.NewCoinFromInt64(1),
			expectedUsed:   linotypes.NewMiniDollar(1200),
		},
		{
			testName:       "2399 minidollar",
			mini:           linotypes.NewMiniDollar(2399),
			expectedBought: linotypes.NewCoinFromInt64(1),
			expectedUsed:   linotypes.NewMiniDollar(1200),
		},
		{
			testName:       "120000000 minidollar",
			mini:           linotypes.NewMiniDollar(120000000),
			expectedBought: linotypes.NewCoinFromInt64(100000),
			expectedUsed:   linotypes.NewMiniDollar(120000000),
		},
		{
			testName:       "8755619048 minidollar",
			mini:           linotypes.NewMiniDollar(8755619048),
			expectedBought: linotypes.NewCoinFromInt64(7296349),
			expectedUsed:   linotypes.NewMiniDollar(8755618800),
		},
	}

	for _, tc := range testCases {
		bought, used, err := suite.manager.MiniDollarToCoin(suite.Ctx, tc.mini)
		suite.Nil(err)
		suite.Equal(tc.expectedBought, bought, "%s", tc.testName)
		suite.Equal(tc.expectedUsed, used, "%s", tc.testName)
	}
}

func toValNames(vals []ValidatorAndVotes) (rst []linotypes.AccountKey) {
	for _, val := range vals {
		rst = append(rst, val.ValidatorName)
	}
	return
}
