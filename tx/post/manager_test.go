package post

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/post/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
)

// test create post
func TestCreatePost(t *testing.T) {
	ctx, am, _, pm, _ := setupTest(t, 1)
	user1 := createTestAccount(t, ctx, am, "user1")
	user2 := createTestAccount(t, ctx, am, "user2")

	cases := []struct {
		postID       string
		author       types.AccountKey
		sourcePostID string
		sourceAuthor types.AccountKey
		expectResult sdk.Error
	}{
		{"postID", user1, "", "", nil},
		{"postID", user2, "", "", nil},
		{"postID", user1, "", "", ErrPostExist(types.GetPermLink(user1, "postID"))},
		{"postID", user2, "postID", user1, ErrPostExist(types.GetPermLink(user2, "postID"))},
		{"postID", user2, "postID", user2, ErrPostExist(types.GetPermLink(user2, "postID"))},
		{"postID2", user2, "postID", user1, nil},
		{"postID3", user2, "postID3", user1,
			ErrCreatePostSourceInvalid(types.GetPermLink(user2, "postID3"))},
	}

	for _, cs := range cases {
		// test valid postInfo
		postCreateParams := PostCreateParams{
			PostID:       cs.postID,
			Title:        string(make([]byte, 50)),
			Content:      string(make([]byte, 1000)),
			Author:       cs.author,
			SourceAuthor: cs.sourceAuthor,
			SourcePostID: cs.sourcePostID,
			Links:        nil,
			RedistributionSplitRate: "0",
		}
		err := pm.CreatePost(ctx, &postCreateParams)
		assert.Equal(t, err, cs.expectResult)

		if err != nil {
			continue
		}
		postInfo := model.PostInfo{
			PostID:       postCreateParams.PostID,
			Title:        postCreateParams.Title,
			Content:      postCreateParams.Content,
			Author:       postCreateParams.Author,
			SourceAuthor: postCreateParams.SourceAuthor,
			SourcePostID: postCreateParams.SourcePostID,
			Links:        postCreateParams.Links,
		}

		postMeta := model.PostMeta{
			Created:                 ctx.BlockHeader().Time,
			LastUpdate:              ctx.BlockHeader().Time,
			LastActivity:            ctx.BlockHeader().Time,
			AllowReplies:            true,
			IsDeleted:               false,
			RedistributionSplitRate: sdk.ZeroRat,
		}
		checkPostKVStore(t, ctx,
			types.GetPermLink(postCreateParams.Author, postCreateParams.PostID), postInfo, postMeta)
	}
}

// test get source post
func TestGetSourcePost(t *testing.T) {
	ctx, _, _, pm, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")
	user2 := types.AccountKey("user2")
	user3 := types.AccountKey("user3")
	cases := []struct {
		postID             string
		author             types.AccountKey
		sourcePostID       string
		sourceAuthor       types.AccountKey
		expectSourcePostID string
		expectSourceAuthor types.AccountKey
	}{
		{"postID", user1, "", "", "", ""},
		{"postID1", user1, "postID", user1, "postID", user1},
		{"postID", user2, "postID1", user1, "postID", user1},
		{"postID", user3, "postID", user2, "postID", user1},
	}

	for _, cs := range cases {
		postCreateParams := PostCreateParams{
			PostID:       cs.postID,
			Title:        string(make([]byte, 50)),
			Content:      string(make([]byte, 1000)),
			Author:       cs.author,
			ParentAuthor: "",
			ParentPostID: "",
			SourceAuthor: cs.sourceAuthor,
			SourcePostID: cs.sourcePostID,
			Links:        nil,
			RedistributionSplitRate: "0",
		}
		err := pm.CreatePost(ctx, &postCreateParams)
		assert.Nil(t, err)
		sourceAuthor, sourcePostID, err :=
			pm.GetSourcePost(ctx, types.GetPermLink(cs.author, cs.postID))
		assert.Nil(t, err)
		assert.Equal(t, sourceAuthor, cs.expectSourceAuthor)
		assert.Equal(t, sourcePostID, cs.expectSourcePostID)
	}
}

func TestAddOrUpdateLikeToPost(t *testing.T) {
	ctx, am, _, pm, _ := setupTest(t, 1)
	user1, postID1 := createTestPost(t, ctx, "user1", "postID1", am, pm, "0")
	user2, postID2 := createTestPost(t, ctx, "user2", "postID2", am, pm, "0")
	user3 := types.AccountKey("user3")

	cases := []struct {
		likeUser                 types.AccountKey
		postID                   string
		author                   types.AccountKey
		weight                   int64
		expectTotalLikeCount     int64
		expectTotalLikeWeight    int64
		expectTotalDislikeWeight int64
	}{
		{user3, postID1, user1, 10000, 1, 10000, 0},
		{user3, postID2, user2, 10000, 1, 10000, 0},
		{user1, postID2, user2, 10000, 2, 20000, 0},
		{user2, postID1, user1, -10000, 2, 10000, 10000},
		{user3, postID2, user2, 0, 2, 10000, 0},
		{user3, postID1, user1, -10000, 2, 0, 20000},
	}

	for _, cs := range cases {
		postKey := types.GetPermLink(cs.author, cs.postID)
		err := pm.AddOrUpdateLikeToPost(ctx, postKey, cs.likeUser, cs.weight)
		assert.Nil(t, err)
		postMeta := model.PostMeta{
			Created:                 ctx.BlockHeader().Time,
			LastUpdate:              ctx.BlockHeader().Time,
			LastActivity:            ctx.BlockHeader().Time,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat,
			TotalLikeCount:          cs.expectTotalLikeCount,
			TotalLikeWeight:         cs.expectTotalLikeWeight,
			TotalDislikeWeight:      cs.expectTotalDislikeWeight,
		}
		checkPostMeta(t, ctx, postKey, postMeta)
	}
}

func TestAddOrUpdateViewToPost(t *testing.T) {
	ctx, am, _, pm, _ := setupTest(t, 1)
	createTime := ctx.BlockHeader().Time
	user1, postID1 := createTestPost(t, ctx, "user1", "postID1", am, pm, "0")
	user2, _ := createTestPost(t, ctx, "user2", "postID2", am, pm, "0")
	user3 := types.AccountKey("user3")

	cases := []struct {
		viewUser             types.AccountKey
		postID               string
		author               types.AccountKey
		viewTime             int64
		expectTotalViewCount int64
		expectUserViewCount  int64
	}{
		{user3, postID1, user1, 1, 1, 1},
		{user3, postID1, user1, 2, 2, 2},
		{user2, postID1, user1, 3, 3, 1},
		{user2, postID1, user1, 4, 4, 2},
		{user1, postID1, user1, 5, 5, 1},
	}

	for _, cs := range cases {
		postKey := types.GetPermLink(cs.author, cs.postID)
		ctx = ctx.WithBlockHeader(abci.Header{Time: cs.viewTime})
		err := pm.AddOrUpdateViewToPost(ctx, postKey, cs.viewUser)
		assert.Nil(t, err)
		postMeta := model.PostMeta{
			Created:                 createTime,
			LastUpdate:              createTime,
			LastActivity:            createTime,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat,
			TotalViewCount:          cs.expectTotalViewCount,
		}
		checkPostMeta(t, ctx, postKey, postMeta)
		view, err := pm.postStorage.GetPostView(ctx, postKey, cs.viewUser)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectUserViewCount, view.Times)
		assert.Equal(t, cs.viewTime, view.LastView)
	}
}

func TestReportOrUpvoteToPost(t *testing.T) {
	ctx, am, _, pm, _ := setupTest(t, 1)
	user1, postID1 := createTestPost(t, ctx, "user1", "postID1", am, pm, "0")
	user2, _ := createTestPost(t, ctx, "user2", "postID2", am, pm, "0")
	user3 := types.AccountKey("user3")

	cases := []struct {
		user                   types.AccountKey
		stake                  types.Coin
		postID                 string
		author                 types.AccountKey
		isReport               bool
		expectTotalReportStake types.Coin
		expectTotalUpvoteStake types.Coin
	}{
		{user3, types.NewCoin(1), postID1, user1, true, types.NewCoin(1), types.NewCoin(0)},
		{user2, types.NewCoin(100), postID1, user1, false, types.NewCoin(1), types.NewCoin(100)},
	}

	for _, cs := range cases {
		postKey := types.GetPermLink(cs.author, cs.postID)
		err := pm.ReportOrUpvoteToPost(ctx, postKey, cs.user, cs.stake, cs.isReport)
		assert.Nil(t, err)
		postMeta := model.PostMeta{
			Created:                 ctx.BlockHeader().Time,
			LastUpdate:              ctx.BlockHeader().Time,
			LastActivity:            ctx.BlockHeader().Time,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat,
			TotalReportStake:        cs.expectTotalReportStake,
			TotalUpvoteStake:        cs.expectTotalUpvoteStake,
		}
		checkPostMeta(t, ctx, postKey, postMeta)
	}
}

func TestDonation(t *testing.T) {
	ctx, am, _, pm, _ := setupTest(t, 1)
	user1, postID1 := createTestPost(t, ctx, "user1", "postID1", am, pm, "0")
	user2, postID2 := createTestPost(t, ctx, "user2", "postID2", am, pm, "0")
	user3 := types.AccountKey("user3")

	baseTime := ctx.BlockHeader().Time
	cases := []struct {
		user                types.AccountKey
		donateAt            int64
		amount              types.Coin
		postID              string
		author              types.AccountKey
		expectDonateCount   int64
		expectTotalDonation types.Coin
		expectDonationList  model.Donations
	}{
		{user3, baseTime, types.NewCoin(1), postID1, user1, 1, types.NewCoin(1),
			model.Donations{user3, []model.Donation{model.Donation{types.NewCoin(1), baseTime}}}},
		{user3, baseTime, types.NewCoin(1), postID2, user2, 1, types.NewCoin(1),
			model.Donations{user3, []model.Donation{model.Donation{types.NewCoin(1), baseTime}}}},
		{user3, baseTime, types.NewCoin(20), postID2, user2, 2, types.NewCoin(21),
			model.Donations{user3,
				[]model.Donation{model.Donation{types.NewCoin(1), baseTime},
					model.Donation{types.NewCoin(20), baseTime}}}},
	}

	for _, cs := range cases {
		postKey := types.GetPermLink(cs.author, cs.postID)
		err := pm.AddDonation(ctx, postKey, cs.user, cs.amount)
		assert.Nil(t, err)
		postMeta := model.PostMeta{
			Created:                 ctx.BlockHeader().Time,
			LastUpdate:              ctx.BlockHeader().Time,
			LastActivity:            ctx.BlockHeader().Time,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat,
			TotalDonateCount:        cs.expectDonateCount,
			TotalReward:             cs.expectTotalDonation,
		}
		checkPostMeta(t, ctx, postKey, postMeta)
		storage := model.NewPostStorage(TestPostKVStoreKey)
		donations, _ := storage.GetPostDonations(ctx, postKey, cs.user)
		assert.Equal(t, cs.expectDonationList, *donations)
	}
}

func TestGetPenaltyScore(t *testing.T) {
	ctx, am, _, pm, _ := setupTest(t, 1)
	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	postKey := types.GetPermLink(user, postID)
	cases := []struct {
		totalReportStake types.Coin
		totalUpvoteStake types.Coin
		expectRat        sdk.Rat
	}{
		{types.NewCoin(1), types.NewCoin(0), sdk.OneRat},
		{types.NewCoin(0), types.NewCoin(1), sdk.ZeroRat},
		{types.NewCoin(0), types.NewCoin(0), sdk.ZeroRat},
		{types.NewCoin(100), types.NewCoin(100), sdk.OneRat},
		{types.NewCoin(1000), types.NewCoin(100), sdk.OneRat},
		{types.NewCoin(50), types.NewCoin(100), sdk.NewRat(1, 2)},
	}

	for _, cs := range cases {
		postMeta := &model.PostMeta{
			Created:                 ctx.BlockHeader().Time,
			LastUpdate:              ctx.BlockHeader().Time,
			LastActivity:            ctx.BlockHeader().Time,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat,
			TotalReportStake:        cs.totalReportStake,
			TotalUpvoteStake:        cs.totalUpvoteStake,
		}
		err := pm.postStorage.SetPostMeta(ctx, postKey, postMeta)
		assert.Nil(t, err)
		penaltyScore, err := pm.GetPenaltyScore(ctx, postKey)
		assert.Nil(t, err)
		assert.Equal(t, penaltyScore, cs.expectRat)
	}
}

func TestGetRepostPenaltyScore(t *testing.T) {
	ctx, am, _, pm, _ := setupTest(t, 1)
	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user2, postID2 := createTestRepost(t, ctx, "user2", "repost", am, pm, user, postID)

	postKey := types.GetPermLink(user, postID)
	repostKey := types.GetPermLink(user2, postID2)
	cases := []struct {
		totalReportStake types.Coin
		totalUpvoteStake types.Coin
		expectRat        sdk.Rat
	}{
		{types.NewCoin(1), types.NewCoin(0), sdk.OneRat},
		{types.NewCoin(0), types.NewCoin(1), sdk.ZeroRat},
		{types.NewCoin(0), types.NewCoin(0), sdk.ZeroRat},
		{types.NewCoin(100), types.NewCoin(100), sdk.OneRat},
		{types.NewCoin(1000), types.NewCoin(100), sdk.OneRat},
		{types.NewCoin(50), types.NewCoin(100), sdk.NewRat(1, 2)},
	}

	for _, cs := range cases {
		postMeta := &model.PostMeta{
			Created:                 ctx.BlockHeader().Time,
			LastUpdate:              ctx.BlockHeader().Time,
			LastActivity:            ctx.BlockHeader().Time,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat,
			TotalReportStake:        cs.totalReportStake,
			TotalUpvoteStake:        cs.totalUpvoteStake,
		}
		err := pm.postStorage.SetPostMeta(ctx, postKey, postMeta)
		assert.Nil(t, err)
		penaltyScore, err := pm.GetPenaltyScore(ctx, repostKey)
		assert.Nil(t, err)
		assert.Equal(t, penaltyScore, cs.expectRat)
	}
}

func checkIsDelete(t *testing.T, ctx sdk.Context, pm PostManager, permLink types.PermLink) {
	isDeleted, err := pm.IsDeleted(ctx, permLink)
	assert.Nil(t, err)
	assert.Equal(t, true, isDeleted)
	postInfo, err := pm.postStorage.GetPostInfo(ctx, permLink)
	assert.Nil(t, err)
	assert.Equal(t, "", postInfo.Title)
	assert.Equal(t, "", postInfo.Content)
}

func TestDeletePost(t *testing.T) {
	ctx, am, _, pm, _ := setupTest(t, 1)
	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user2, postID2 := createTestRepost(t, ctx, "user2", "repost", am, pm, user, postID)

	err := pm.DeletePost(ctx, types.GetPermLink(user2, postID2))
	assert.Nil(t, err)
	checkIsDelete(t, ctx, pm, types.GetPermLink(user2, postID2))
	postMeta, err := pm.postStorage.GetPostMeta(ctx, types.GetPermLink(user, postID))
	assert.Nil(t, err)
	assert.Equal(t, false, postMeta.IsDeleted)
	err = pm.DeletePost(ctx, types.GetPermLink(user, postID))
	assert.Nil(t, err)
	checkIsDelete(t, ctx, pm, types.GetPermLink(user, postID))
}
