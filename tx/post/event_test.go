package post

import (
	"math/big"
	"testing"

	globalModel "github.com/lino-network/lino/tx/global/model"
	postModel "github.com/lino-network/lino/tx/post/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestRewardEvent(t *testing.T) {
	ctx, am, _, pm, gm, dm := setupTest(t, 1)
	gs := globalModel.NewGlobalStorage(TestGlobalKVStoreKey)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user1 := createTestAccount(t, ctx, am, "user1")
	err := dm.RegisterDeveloper(ctx, "LinoApp1", types.NewCoinFromInt64(1000000*types.Decimals))
	assert.Nil(t, err)
	err = dm.RegisterDeveloper(ctx, "LinoApp2", types.NewCoinFromInt64(1000000*types.Decimals))
	assert.Nil(t, err)

	testCases := []struct {
		testName             string
		rewardEvent          RewardEvent
		totalReportOfthePost types.Coin
		totalUpvoteOfthePost types.Coin
		initRewardPool       types.Coin
		initRewardWindow     types.Coin
		expectPostMeta       postModel.PostMeta
		expectAppWeight      *big.Rat
	}{
		{"normal event",
			RewardEvent{
				PostAuthor: user,
				PostID:     postID,
				Consumer:   user1,
				Evaluate:   types.NewCoinFromInt64(100),
				Original:   types.NewCoinFromInt64(100),
				Friction:   types.NewCoinFromInt64(15),
				FromApp:    types.AccountKey("LinoApp1"),
			}, types.NewCoinFromInt64(0), types.NewCoinFromInt64(100),
			types.NewCoinFromInt64(100), types.NewCoinFromInt64(100),
			postModel.PostMeta{
				TotalUpvoteStake: types.NewCoinFromInt64(100),
				TotalReportStake: types.NewCoinFromInt64(0),
				TotalDonateCount: 1,
				TotalReward:      types.NewCoinFromInt64(100),
			}, new(big.Rat).SetFloat64(1.0),
		},
		{"100% panelty reward post",
			RewardEvent{
				PostAuthor: user,
				PostID:     postID,
				Consumer:   user1,
				Evaluate:   types.NewCoinFromInt64(100),
				Original:   types.NewCoinFromInt64(100),
				Friction:   types.NewCoinFromInt64(15),
				FromApp:    types.AccountKey("LinoApp1"),
			}, types.NewCoinFromInt64(100), types.NewCoinFromInt64(100),
			types.NewCoinFromInt64(100), types.NewCoinFromInt64(100),
			postModel.PostMeta{
				TotalUpvoteStake: types.NewCoinFromInt64(100),
				TotalReportStake: types.NewCoinFromInt64(100),
				TotalDonateCount: 1,
				TotalReward:      types.NewCoinFromInt64(0),
			}, new(big.Rat).SetFloat64(1.0),
		},
		{"50% panelty reward post",
			RewardEvent{
				PostAuthor: user,
				PostID:     postID,
				Consumer:   user1,
				Evaluate:   types.NewCoinFromInt64(100),
				Original:   types.NewCoinFromInt64(100),
				Friction:   types.NewCoinFromInt64(15),
				FromApp:    types.AccountKey("LinoApp1"),
			}, types.NewCoinFromInt64(50), types.NewCoinFromInt64(100),
			types.NewCoinFromInt64(100), types.NewCoinFromInt64(100),
			postModel.PostMeta{
				TotalUpvoteStake: types.NewCoinFromInt64(100),
				TotalReportStake: types.NewCoinFromInt64(50),
				TotalDonateCount: 1,
				TotalReward:      types.NewCoinFromInt64(50),
			}, new(big.Rat).SetFloat64(1.0),
		},
		{"evaluate as 1% of total window",
			RewardEvent{
				PostAuthor: user,
				PostID:     postID,
				Consumer:   user1,
				Evaluate:   types.NewCoinFromInt64(1),
				Original:   types.NewCoinFromInt64(100),
				Friction:   types.NewCoinFromInt64(15),
				FromApp:    types.AccountKey("LinoApp1"),
			}, types.NewCoinFromInt64(0), types.NewCoinFromInt64(100),
			types.NewCoinFromInt64(100), types.NewCoinFromInt64(100),
			postModel.PostMeta{
				TotalUpvoteStake: types.NewCoinFromInt64(100),
				TotalReportStake: types.NewCoinFromInt64(0),
				TotalDonateCount: 1,
				TotalReward:      types.NewCoinFromInt64(1),
			}, new(big.Rat).SetFloat64(1.0),
		},
		{"reward from different app",
			RewardEvent{
				PostAuthor: user,
				PostID:     postID,
				Consumer:   user1,
				Evaluate:   types.NewCoinFromInt64(100),
				Original:   types.NewCoinFromInt64(100),
				Friction:   types.NewCoinFromInt64(15),
				FromApp:    types.AccountKey("LinoApp2"),
			}, types.NewCoinFromInt64(0), types.NewCoinFromInt64(100),
			types.NewCoinFromInt64(100), types.NewCoinFromInt64(100),
			postModel.PostMeta{
				TotalUpvoteStake: types.NewCoinFromInt64(100),
				TotalReportStake: types.NewCoinFromInt64(0),
				TotalDonateCount: 1,
				TotalReward:      types.NewCoinFromInt64(100),
			}, new(big.Rat).SetFrac(new(big.Int).SetInt64(100), new(big.Int).SetInt64(251)),
		},
	}

	for _, tc := range testCases {
		gs.SetConsumptionMeta(ctx, &globalModel.ConsumptionMeta{
			ConsumptionRewardPool: tc.initRewardPool,
			ConsumptionWindow:     tc.initRewardWindow,
		})
		pm.postStorage.SetPostMeta(ctx, types.GetPermLink(user, postID),
			&postModel.PostMeta{
				TotalUpvoteStake: tc.totalUpvoteOfthePost,
				TotalReportStake: tc.totalReportOfthePost,
			})
		err := tc.rewardEvent.Execute(ctx, pm, am, gm, dm)
		assert.Nil(t, err)
		postMeta, err := pm.postStorage.GetPostMeta(ctx, types.GetPermLink(user, postID))
		assert.Nil(t, err)
		if *postMeta != tc.expectPostMeta {
			t.Errorf("%s get post meta failed: got %v, want %v",
				tc.testName, *postMeta, tc.expectPostMeta)
		}
		if dm.IsDeveloperExist(ctx, tc.rewardEvent.FromApp) {
			consumptionWeight, err := dm.GetConsumptionWeight(ctx, tc.rewardEvent.FromApp)
			assert.Nil(t, err)
			if tc.expectAppWeight.Cmp(consumptionWeight) != 0 {
				t.Errorf("%s get expect app weight failed: got %v, want %v",
					tc.testName, consumptionWeight, tc.expectAppWeight)
			}
		}
	}
}
