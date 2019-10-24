package manager

import (
	"github.com/stretchr/testify/mock"

	linotypes "github.com/lino-network/lino/types"
	types "github.com/lino-network/lino/x/post/types"
)

func (suite *PostManagerTestSuite) TestExecRewardEvent() {
	user1 := suite.user1
	user2 := suite.user2
	app1 := suite.app1
	postID := "post1"
	err := suite.pm.CreatePost(suite.Ctx, user1, postID, app1, "content", "title")
	suite.Require().Nil(err)

	// fixed pool in this test.
	suite.am.On("GetPool", mock.Anything, linotypes.InflationConsumptionPool).Return(
		linotypes.NewCoinFromInt64(10000), nil).Maybe()

	testCases := []struct {
		testName        string
		event           types.RewardEvent
		consumptionPool linotypes.MiniDollar
		hasDev          bool
		hasPost         bool
		expectedReward  linotypes.Coin
	}{
		{
			testName: "OK",
			event: types.RewardEvent{
				PostAuthor: user1,
				PostID:     postID,
				Consumer:   user2,
				Evaluate:   linotypes.NewMiniDollar(100),
				FromApp:    app1,
			},
			consumptionPool: linotypes.NewMiniDollar(300),
			expectedReward:  linotypes.NewCoinFromInt64(3333),
			hasDev:          true,
			hasPost:         true,
		},
		{
			testName: "PostDeleted",
			event: types.RewardEvent{
				PostAuthor: user1,
				PostID:     "deletedpost",
				Consumer:   user2,
				Evaluate:   linotypes.NewMiniDollar(100),
				FromApp:    app1,
			},
			hasDev:  true,
			hasPost: false,
		},
		{
			testName: "NoDev",
			event: types.RewardEvent{
				PostAuthor: user1,
				PostID:     postID,
				Consumer:   user2,
				Evaluate:   linotypes.NewMiniDollar(100),
				FromApp:    user2,
			},
			consumptionPool: linotypes.NewMiniDollar(300),
			expectedReward:  linotypes.NewCoinFromInt64(3333),
			hasDev:          false,
			hasPost:         true,
		},
	}

	for _, tc := range testCases {
		suite.pm.postStorage.SetConsumptionWindow(suite.Ctx, tc.consumptionPool)
		if tc.hasPost {
			suite.am.On("MoveFromPool", mock.Anything, linotypes.InflationConsumptionPool,
				linotypes.NewAccOrAddrFromAcc(tc.event.PostAuthor),
				tc.expectedReward).Return(nil).Once()
			if tc.hasDev {
				suite.dev.On(
					"ReportConsumption", mock.Anything,
					tc.event.FromApp, tc.event.Evaluate).Return(nil).Once()
			}
		}
		err := suite.pm.ExecRewardEvent(suite.Ctx, tc.event)
		suite.Nil(err)
		suite.am.AssertExpectations(suite.T())
		suite.dev.AssertExpectations(suite.T())
	}
}

func (suite *PostManagerTestSuite) TestExecMultipleRewardEvents() {
	user1 := suite.user1
	user2 := suite.user2
	app1 := suite.app1
	postID := "post1"
	err := suite.pm.CreatePost(suite.Ctx, user1, postID, app1, "content", "title")
	suite.Require().Nil(err)
	err = suite.pm.CreatePost(suite.Ctx, user2, postID, app1, "content", "title")
	suite.Require().Nil(err)

	totalConsumption := linotypes.NewMiniDollar(100)
	totalInflation := linotypes.NewCoinFromInt64(10000)
	suite.pm.postStorage.SetConsumptionWindow(suite.Ctx, totalConsumption)
	events := []types.RewardEvent{
		{
			PostAuthor: user1,
			PostID:     postID,
			Evaluate:   linotypes.NewMiniDollar(20),
			FromApp:    app1,
		},
		{
			PostAuthor: user2,
			PostID:     postID,
			Evaluate:   linotypes.NewMiniDollar(30),
			FromApp:    app1,
		},
		{
			PostAuthor: user2,
			PostID:     postID,
			Evaluate:   linotypes.NewMiniDollar(50),
			FromApp:    app1,
		},
	}

	for _, e := range events {
		suite.am.On("GetPool", mock.Anything, linotypes.InflationConsumptionPool).Return(
			totalInflation, nil).Once()
		inflation := linotypes.DecToCoin(
			totalInflation.ToDec().Mul(e.Evaluate.ToDec().Quo(totalConsumption.ToDec())))
		totalInflation = totalInflation.Minus(inflation)
		suite.am.On("MoveFromPool", mock.Anything, linotypes.InflationConsumptionPool,
			linotypes.NewAccOrAddrFromAcc(e.PostAuthor), inflation).Return(nil).Once()
		suite.dev.On(
			"ReportConsumption", mock.Anything,
			e.FromApp, e.Evaluate).Return(nil).Once()
		err := suite.pm.ExecRewardEvent(suite.Ctx, e)
		suite.Nil(err)
		totalConsumption = totalConsumption.Minus(e.Evaluate)
	}
	suite.am.AssertExpectations(suite.T())
}
