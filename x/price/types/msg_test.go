package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/lino-network/lino/types"
)

type PriceMsgTestSuite struct {
	suite.Suite
}

func TestPriceMsgTestSuite(t *testing.T) {
	suite.Run(t, new(PriceMsgTestSuite))
}

func (suite *PriceMsgTestSuite) SetupTest() {
}

func (suite *PriceMsgTestSuite) TestFeedPriceMsgValidateBasic() {
	testCases := []struct {
		testName       string
		msg            FeedPriceMsg
		expectedResult sdk.Error
	}{
		{
			"invalid username",
			FeedPriceMsg{
				Username: "3v",
				Price:    types.NewMiniDollar(100),
			},
			types.ErrInvalidUsername("3v"),
		},
		{
			"invalid price zero",
			FeedPriceMsg{
				Username: "user1",
				Price:    types.NewMiniDollar(0),
			},
			ErrInvalidPriceFeed(types.NewMiniDollar(0)),
		},
		{
			"invalid price negative",
			FeedPriceMsg{
				Username: "user1",
				Price:    types.NewMiniDollar(-100),
			},
			ErrInvalidPriceFeed(types.NewMiniDollar(-100)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.testName, func() {
			result := tc.msg.ValidateBasic()
			suite.Equal(tc.expectedResult, result)
		})
	}
}
