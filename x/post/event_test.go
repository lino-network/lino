package post

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	accModel "github.com/lino-network/lino/x/account/model"
	globalModel "github.com/lino-network/lino/x/global/model"
	postModel "github.com/lino-network/lino/x/post/model"
)

func TestRewardEvent(t *testing.T) {
	ctx, am, _, pm, gm, dm, vm, rm := setupTest(t, 1)
	gs := globalModel.NewGlobalStorage(testGlobalKVStoreKey)
	as := accModel.NewAccountStorage(testAccountKVStoreKey)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user2, deletedPostID := createTestPost(t, ctx, "user2", "deleted", am, pm, "0")
	err := pm.DeletePost(ctx, types.GetPermlink(user2, deletedPostID))
	assert.Nil(t, err)
	user1 := createTestAccount(t, ctx, am, "user1")
	err = dm.RegisterDeveloper(ctx, "LinoApp1", types.NewCoinFromInt64(1000000*types.Decimals), "", "", "")
	assert.Nil(t, err)
	err = dm.RegisterDeveloper(ctx, "LinoApp2", types.NewCoinFromInt64(1000000*types.Decimals), "", "", "")
	assert.Nil(t, err)

	testCases := []struct {
		testName           string
		rewardEvent        RewardEvent
		initRewardPool     types.Coin
		initRewardWindow   types.Coin
		expectPostMeta     postModel.PostMeta
		expectAppWeight    sdk.Rat
		expectAuthorReward accModel.Reward
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
			initRewardPool:   types.NewCoinFromInt64(100),
			initRewardWindow: types.NewCoinFromInt64(100),
			expectPostMeta: postModel.PostMeta{
				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
				TotalReportCoinDay:      types.NewCoinFromInt64(0),
				TotalDonateCount:        1,
				TotalReward:             types.NewCoinFromInt64(100),
				RedistributionSplitRate: sdk.ZeroRat(),
				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
			},
			expectAppWeight: sdk.OneRat(),
			expectAuthorReward: accModel.Reward{
				TotalIncome:     types.NewCoinFromInt64(100),
				OriginalIncome:  types.NewCoinFromInt64(15),
				FrictionIncome:  types.NewCoinFromInt64(15),
				InflationIncome: types.NewCoinFromInt64(100),
				UnclaimReward:   types.NewCoinFromInt64(100),
			},
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
			initRewardPool:   types.NewCoinFromInt64(100),
			initRewardWindow: types.NewCoinFromInt64(100),
			expectPostMeta: postModel.PostMeta{
				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
				TotalReportCoinDay:      types.NewCoinFromInt64(0),
				TotalDonateCount:        1,
				TotalReward:             types.NewCoinFromInt64(1),
				RedistributionSplitRate: sdk.ZeroRat(),
				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
			},
			expectAppWeight: sdk.OneRat(),
			expectAuthorReward: accModel.Reward{
				TotalIncome:     types.NewCoinFromInt64(1),
				OriginalIncome:  types.NewCoinFromInt64(15),
				FrictionIncome:  types.NewCoinFromInt64(15),
				InflationIncome: types.NewCoinFromInt64(1),
				UnclaimReward:   types.NewCoinFromInt64(1),
			},
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
			initRewardPool:   types.NewCoinFromInt64(100),
			initRewardWindow: types.NewCoinFromInt64(100),
			expectPostMeta: postModel.PostMeta{
				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
				TotalReportCoinDay:      types.NewCoinFromInt64(0),
				TotalDonateCount:        1,
				TotalReward:             types.NewCoinFromInt64(100),
				RedistributionSplitRate: sdk.ZeroRat(),
				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
			},
			expectAppWeight: sdk.NewRat(1243781, 2500000),
			expectAuthorReward: accModel.Reward{
				TotalIncome:     types.NewCoinFromInt64(100),
				OriginalIncome:  types.NewCoinFromInt64(15),
				FrictionIncome:  types.NewCoinFromInt64(15),
				InflationIncome: types.NewCoinFromInt64(100),
				UnclaimReward:   types.NewCoinFromInt64(100),
			},
		},
		{
			testName: "deleted post can't get any inflation",
			rewardEvent: RewardEvent{
				PostAuthor: user2,
				PostID:     deletedPostID,
				Consumer:   user1,
				Evaluate:   types.NewCoinFromInt64(33333),
				Original:   types.NewCoinFromInt64(100),
				Friction:   types.NewCoinFromInt64(15),
				FromApp:    types.AccountKey("LinoApp2"),
			},
			initRewardPool:   types.NewCoinFromInt64(5555),
			initRewardWindow: types.NewCoinFromInt64(77777),
			expectPostMeta: postModel.PostMeta{
				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
				TotalReportCoinDay:      types.NewCoinFromInt64(0),
				TotalDonateCount:        1,
				TotalReward:             types.NewCoinFromInt64(0),
				RedistributionSplitRate: sdk.ZeroRat(),
				IsDeleted:               true,
				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
			},
			expectAppWeight: sdk.NewRat(1243781, 2500000),
			expectAuthorReward: accModel.Reward{
				TotalIncome:     types.NewCoinFromInt64(0),
				OriginalIncome:  types.NewCoinFromInt64(15),
				FrictionIncome:  types.NewCoinFromInt64(15),
				InflationIncome: types.NewCoinFromInt64(0),
				UnclaimReward:   types.NewCoinFromInt64(0),
			},
		},
	}

	for _, tc := range testCases {
		gs.SetConsumptionMeta(ctx, &globalModel.ConsumptionMeta{
			ConsumptionRewardPool: tc.initRewardPool,
			ConsumptionWindow:     tc.initRewardWindow,
		})
		pm.postStorage.SetPostMeta(ctx, types.GetPermlink(tc.rewardEvent.PostAuthor, tc.rewardEvent.PostID),
			&postModel.PostMeta{
				TotalUpvoteCoinDay: types.NewCoinFromInt64(0),
				TotalReportCoinDay: types.NewCoinFromInt64(0),
				TotalReward:        types.NewCoinFromInt64(0),
				IsDeleted:          tc.expectPostMeta.IsDeleted,
			})

		as.SetReward(ctx, tc.rewardEvent.PostAuthor, &accModel.Reward{})
		vm.AddVoter(ctx, tc.rewardEvent.PostAuthor, types.NewCoinFromInt64(0))
		err := tc.rewardEvent.Execute(ctx, pm, am, gm, dm, vm, rm)
		if err != nil {
			t.Errorf("%s: failed to execute, got err %v", tc.testName, err)
		}
		checkPostMeta(t, ctx, types.GetPermlink(tc.rewardEvent.PostAuthor, tc.rewardEvent.PostID), tc.expectPostMeta)
		if dm.DoesDeveloperExist(ctx, tc.rewardEvent.FromApp) {
			consumptionWeight, err := dm.GetConsumptionWeight(ctx, tc.rewardEvent.FromApp)
			if err != nil {
				t.Errorf("%s: failed to get consumption weight, got err %v", tc.testName, err)
			}
			if !tc.expectAppWeight.Equal(consumptionWeight) {
				t.Errorf("%s: diff consumption weight, got %v, want %v", tc.testName, consumptionWeight, tc.expectAppWeight)
			}
		}
		reward, err := as.GetReward(ctx, tc.rewardEvent.PostAuthor)
		if err != nil {
			t.Errorf("%s: failed to get reward, got err %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.expectAuthorReward, *reward) {
			t.Errorf("%s: diff reward, got %v, want %v", tc.testName, *reward, tc.expectAuthorReward)
		}
	}
}
