package manager

import (
	"github.com/stretchr/testify/mock"

	linotypes "github.com/lino-network/lino/types"
	types "github.com/lino-network/lino/x/post/types"
)

func (suite *PostManagerTestSuite) TestRewardEvent() {
	user1 := suite.user1
	user2 := suite.user2
	app1 := suite.app1
	postID := "post1"
	err := suite.pm.CreatePost(suite.Ctx, user1, postID, app1, "content", "title")
	suite.Require().Nil(err)

	testCases := []struct {
		testName string
		event    types.RewardEvent
		reward   linotypes.Coin
		hasDev   bool
		hasPost  bool
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
			reward:  linotypes.NewCoinFromInt64(3333),
			hasDev:  true,
			hasPost: true,
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
			reward:  linotypes.NewCoinFromInt64(3333),
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
			reward:  linotypes.NewCoinFromInt64(3333),
			hasDev:  false,
			hasPost: true,
		},
	}

	for _, tc := range testCases {
		if tc.hasPost {
			suite.global.On("GetRewardAndPopFromWindow", mock.Anything, tc.event.Evaluate).Return(
				tc.reward, nil,
			).Once()
			suite.am.On("AddCoinToUsername", mock.Anything, tc.event.PostAuthor, tc.reward).Return(nil).Once()
			if tc.hasDev {
				suite.dev.On(
					"ReportConsumption", mock.Anything, tc.event.FromApp, tc.event.Evaluate).Return(nil).Once()
			}
		}
		err := suite.pm.ExecRewardEvent(suite.Ctx, tc.event)
		suite.Nil(err)
		suite.global.AssertExpectations(suite.T())
		suite.am.AssertExpectations(suite.T())
		suite.dev.AssertExpectations(suite.T())
	}
}
