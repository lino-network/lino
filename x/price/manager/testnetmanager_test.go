package manager

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/lino-network/lino/types"
)

type testnetPriceManagerTestSuite struct {
	suite.Suite
	price TestnetPriceManager
}

func TestTestnetPriceManagerTestSuite(t *testing.T) {
	suite.Run(t, new(testnetPriceManagerTestSuite))
}

func (suite *testnetPriceManagerTestSuite) SetupTest() {
	suite.price = TestnetPriceManager{}
}

func (suite *testnetPriceManagerTestSuite) TestCoinToMiniDollar() {
	testCases := []struct {
		testName string
		coin     types.Coin
		expected types.MiniDollar
	}{
		{
			testName: "0 Coin",
			coin:     types.NewCoinFromInt64(0),
			expected: types.NewMiniDollar(0),
		},
		{
			testName: "1 Coin",
			coin:     types.NewCoinFromInt64(1),
			expected: types.NewMiniDollar(1200),
		},
		{
			testName: "3 Coin",
			coin:     types.NewCoinFromInt64(3),
			expected: types.NewMiniDollar(3600),
		},
		{
			testName: "1 LNO",
			coin:     types.NewCoinFromInt64(types.Decimals),
			expected: types.NewMiniDollar(12 * 10000000),
		},
		{
			testName: "2 LNO",
			coin:     types.NewCoinFromInt64(2 * types.Decimals),
			expected: types.NewMiniDollar(24 * 10000000),
		},
		{
			testName: "1000000 LNO",
			coin:     types.NewCoinFromInt64(1000000 * types.Decimals),
			expected: types.NewMiniDollar(1000000 * 12 * 10000000),
		},
	}

	for _, tc := range testCases {
		suite.Equal(tc.expected, suite.price.CoinToMiniDollar(tc.coin), "%s", tc.testName)
	}
}

func (suite *testnetPriceManagerTestSuite) TestMiniDollarToCoin() {
	testCases := []struct {
		testName       string
		mini           types.MiniDollar
		expectedBought types.Coin
		expectedUsed   types.MiniDollar
	}{
		{
			testName:       "0 minidollar",
			mini:           types.NewMiniDollar(0),
			expectedBought: types.NewCoinFromInt64(0),
			expectedUsed:   types.NewMiniDollar(0),
		},
		{
			testName:       "1 minidollar",
			mini:           types.NewMiniDollar(1),
			expectedBought: types.NewCoinFromInt64(0),
			expectedUsed:   types.NewMiniDollar(0),
		},
		{
			testName:       "1199 minidollar",
			mini:           types.NewMiniDollar(1199),
			expectedBought: types.NewCoinFromInt64(0),
			expectedUsed:   types.NewMiniDollar(0),
		},
		{
			testName:       "1200 minidollar",
			mini:           types.NewMiniDollar(1200),
			expectedBought: types.NewCoinFromInt64(1),
			expectedUsed:   types.NewMiniDollar(1200),
		},
		{
			testName:       "1201 minidollar",
			mini:           types.NewMiniDollar(1201),
			expectedBought: types.NewCoinFromInt64(1),
			expectedUsed:   types.NewMiniDollar(1200),
		},
		{
			testName:       "2399 minidollar",
			mini:           types.NewMiniDollar(2399),
			expectedBought: types.NewCoinFromInt64(1),
			expectedUsed:   types.NewMiniDollar(1200),
		},
		{
			testName:       "120000000 minidollar",
			mini:           types.NewMiniDollar(120000000),
			expectedBought: types.NewCoinFromInt64(100000),
			expectedUsed:   types.NewMiniDollar(120000000),
		},
		{
			testName:       "8755619048 minidollar",
			mini:           types.NewMiniDollar(8755619048),
			expectedBought: types.NewCoinFromInt64(7296349),
			expectedUsed:   types.NewMiniDollar(8755618800),
		},
	}

	for _, tc := range testCases {
		bought, used := suite.price.MiniDollarToCoin(tc.mini)
		suite.Equal(tc.expectedBought, bought, "%s", tc.testName)
		suite.Equal(tc.expectedUsed, used, "%s", tc.testName)
	}
}
