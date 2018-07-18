package post

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	globalModel "github.com/lino-network/lino/x/global/model"
	postModel "github.com/lino-network/lino/x/post/model"
	"github.com/stretchr/testify/assert"
)

func TestRewardEvent(t *testing.T) {
	ctx, am, _, pm, gm, dm := setupTest(t, 1)
	gs := globalModel.NewGlobalStorage(TestGlobalKVStoreKey)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user1 := createTestAccount(t, ctx, am, "user1")
	err := dm.RegisterDeveloper(ctx, "LinoApp1", types.NewCoinFromInt64(1000000*types.Decimals), "", "", "")
	assert.Nil(t, err)
	err = dm.RegisterDeveloper(ctx, "LinoApp2", types.NewCoinFromInt64(1000000*types.Decimals), "", "", "")
	assert.Nil(t, err)

	testCases := []struct {
		testName             string
		rewardEvent          RewardEvent
		totalReportOfThePost types.Coin
		totalUpvoteOfThePost types.Coin
		initRewardPool       types.Coin
		initRewardWindow     types.Coin
		expectPostMeta       postModel.PostMeta
		expectAppWeight      sdk.Rat
	}{
		{
			testName: "normal event",
			rewardEvent: RewardEvent{
				PostAuthor: user,
				PostID:     postID,
				Consumer:   user1,
				Evaluate:   types.NewCoinFromInt64(100),
				Original:   types.NewCoinFromInt64(100),
				Friction:   types.NewCoinFromInt64(15),
				FromApp:    types.AccountKey("LinoApp1"),
			},
			totalReportOfThePost: types.NewCoinFromInt64(0),
			totalUpvoteOfThePost: types.NewCoinFromInt64(100),
			initRewardPool:       types.NewCoinFromInt64(100),
			initRewardWindow:     types.NewCoinFromInt64(100),
			expectPostMeta: postModel.PostMeta{
				TotalUpvoteStake:        types.NewCoinFromInt64(100),
				TotalReportStake:        types.NewCoinFromInt64(0),
				TotalDonateCount:        1,
				TotalReward:             types.NewCoinFromInt64(100),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			expectAppWeight: sdk.OneRat(),
		},
		{
			testName: "100% panelty reward post",
			rewardEvent: RewardEvent{
				PostAuthor: user,
				PostID:     postID,
				Consumer:   user1,
				Evaluate:   types.NewCoinFromInt64(100),
				Original:   types.NewCoinFromInt64(100),
				Friction:   types.NewCoinFromInt64(15),
				FromApp:    types.AccountKey("LinoApp1"),
			},
			totalReportOfThePost: types.NewCoinFromInt64(100),
			totalUpvoteOfThePost: types.NewCoinFromInt64(100),
			initRewardPool:       types.NewCoinFromInt64(100),
			initRewardWindow:     types.NewCoinFromInt64(100),
			expectPostMeta: postModel.PostMeta{
				TotalUpvoteStake:        types.NewCoinFromInt64(100),
				TotalReportStake:        types.NewCoinFromInt64(100),
				TotalDonateCount:        1,
				TotalReward:             types.NewCoinFromInt64(0),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			expectAppWeight: sdk.OneRat(),
		},
		{
			testName: "50% panelty reward post",
			rewardEvent: RewardEvent{
				PostAuthor: user,
				PostID:     postID,
				Consumer:   user1,
				Evaluate:   types.NewCoinFromInt64(100),
				Original:   types.NewCoinFromInt64(100),
				Friction:   types.NewCoinFromInt64(15),
				FromApp:    types.AccountKey("LinoApp1"),
			},
			totalReportOfThePost: types.NewCoinFromInt64(50),
			totalUpvoteOfThePost: types.NewCoinFromInt64(100),
			initRewardPool:       types.NewCoinFromInt64(100),
			initRewardWindow:     types.NewCoinFromInt64(100),
			expectPostMeta: postModel.PostMeta{
				TotalUpvoteStake:        types.NewCoinFromInt64(100),
				TotalReportStake:        types.NewCoinFromInt64(50),
				TotalDonateCount:        1,
				TotalReward:             types.NewCoinFromInt64(50),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			expectAppWeight: sdk.OneRat(),
		},
		{
			testName: "evaluate as 1% of total window",
			rewardEvent: RewardEvent{
				PostAuthor: user,
				PostID:     postID,
				Consumer:   user1,
				Evaluate:   types.NewCoinFromInt64(1),
				Original:   types.NewCoinFromInt64(100),
				Friction:   types.NewCoinFromInt64(15),
				FromApp:    types.AccountKey("LinoApp1"),
			},
			totalReportOfThePost: types.NewCoinFromInt64(0),
			totalUpvoteOfThePost: types.NewCoinFromInt64(100),
			initRewardPool:       types.NewCoinFromInt64(100),
			initRewardWindow:     types.NewCoinFromInt64(100),
			expectPostMeta: postModel.PostMeta{
				TotalUpvoteStake:        types.NewCoinFromInt64(100),
				TotalReportStake:        types.NewCoinFromInt64(0),
				TotalDonateCount:        1,
				TotalReward:             types.NewCoinFromInt64(1),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			expectAppWeight: sdk.OneRat(),
		},
		{
			testName: "reward from different app",
			rewardEvent: RewardEvent{
				PostAuthor: user,
				PostID:     postID,
				Consumer:   user1,
				Evaluate:   types.NewCoinFromInt64(100),
				Original:   types.NewCoinFromInt64(100),
				Friction:   types.NewCoinFromInt64(15),
				FromApp:    types.AccountKey("LinoApp2"),
			},
			totalReportOfThePost: types.NewCoinFromInt64(0),
			totalUpvoteOfThePost: types.NewCoinFromInt64(100),
			initRewardPool:       types.NewCoinFromInt64(100),
			initRewardWindow:     types.NewCoinFromInt64(100),
			expectPostMeta: postModel.PostMeta{
				TotalUpvoteStake:        types.NewCoinFromInt64(100),
				TotalReportStake:        types.NewCoinFromInt64(0),
				TotalDonateCount:        1,
				TotalReward:             types.NewCoinFromInt64(100),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			expectAppWeight: sdk.NewRat(100, 251),
		},
	}

	for _, tc := range testCases {
		gs.SetConsumptionMeta(ctx, &globalModel.ConsumptionMeta{
			ConsumptionRewardPool: tc.initRewardPool,
			ConsumptionWindow:     tc.initRewardWindow,
		})
		pm.postStorage.SetPostMeta(ctx, types.GetPermlink(user, postID),
			&postModel.PostMeta{
				TotalUpvoteStake: tc.totalUpvoteOfThePost,
				TotalReportStake: tc.totalReportOfThePost,
			})

		err := tc.rewardEvent.Execute(ctx, pm, am, gm, dm)
		if err != nil {
			t.Errorf("%s: failed to execute, got err %v", err)
		}
		checkPostMeta(t, ctx, types.GetPermlink(user, postID), tc.expectPostMeta)
		if dm.DoesDeveloperExist(ctx, tc.rewardEvent.FromApp) {
			consumptionWeight, err := dm.GetConsumptionWeight(ctx, tc.rewardEvent.FromApp)
			if err != nil {
				t.Errorf("%s: failed to get consumption weight, got err %v", err)
			}
			if !tc.expectAppWeight.Equal(consumptionWeight) {
				t.Errorf("%s: diff consumption weight, got %v, want %v", tc.testName, consumptionWeight, tc.expectAppWeight)
			}
		}
	}
}
