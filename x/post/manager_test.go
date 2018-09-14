package post

import (
	"math/big"
	"testing"
	"time"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/post/model"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// test create post
func TestCreatePost(t *testing.T) {
	ctx, am, _, pm, _, _, _ := setupTest(t, 1)
	user1 := createTestAccount(t, ctx, am, "user1")
	user2 := createTestAccount(t, ctx, am, "user2")

	testCases := []struct {
		testName     string
		postID       string
		author       types.AccountKey
		sourcePostID string
		sourceAuthor types.AccountKey
		expectResult sdk.Error
	}{
		{
			testName:     "creates (postID, user1) successfully",
			postID:       "postID",
			author:       user1,
			sourcePostID: "",
			sourceAuthor: "",
			expectResult: nil,
		},
		{
			testName:     "creates (postID, user2) successfully",
			postID:       "postID",
			author:       user2,
			sourcePostID: "",
			sourceAuthor: "",
			expectResult: nil,
		},
		{
			testName:     "(postID, user1) already exists",
			postID:       "postID",
			author:       user1,
			sourcePostID: "",
			sourceAuthor: "",
			expectResult: ErrPostAlreadyExist(types.GetPermlink(user1, "postID")),
		},
		{
			testName:     "(postID, user2) already exists case 1",
			postID:       "postID",
			author:       user2,
			sourcePostID: "postID",
			sourceAuthor: user1,
			expectResult: ErrPostAlreadyExist(types.GetPermlink(user2, "postID")),
		},
		{
			testName:     "(postID, user2) already exists case 2",
			postID:       "postID",
			author:       user2,
			sourcePostID: "postID",
			sourceAuthor: user2,
			expectResult: ErrPostAlreadyExist(types.GetPermlink(user2, "postID")),
		},
		{
			testName:     "creates (postID2, user2) successfully",
			postID:       "postID2",
			author:       user2,
			sourcePostID: "postID",
			sourceAuthor: user1,
			expectResult: nil,
		},
		{
			testName:     "source doesn't exist",
			postID:       "postID3",
			author:       user2,
			sourcePostID: "postID3",
			sourceAuthor: user1,
			expectResult: ErrCreatePostSourceInvalid(types.GetPermlink(user2, "postID3")),
		},
	}

	for _, tc := range testCases {
		// test valid postInfo
		msg := CreatePostMsg{
			PostID:       tc.postID,
			Title:        string(make([]byte, 50)),
			Content:      string(make([]byte, 1000)),
			Author:       tc.author,
			SourceAuthor: tc.sourceAuthor,
			SourcePostID: tc.sourcePostID,
			Links:        nil,
		}
		err := pm.CreatePost(
			ctx, msg.Author, msg.PostID, msg.SourceAuthor, msg.SourcePostID,
			msg.ParentAuthor, msg.ParentPostID, msg.Content,
			msg.Title, sdk.ZeroRat(), msg.Links)
		if !assert.Equal(t, err, tc.expectResult) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, err, tc.expectResult)
		}

		if err != nil {
			continue
		}
		postInfo := model.PostInfo{
			PostID:       msg.PostID,
			Title:        msg.Title,
			Content:      msg.Content,
			Author:       msg.Author,
			SourceAuthor: msg.SourceAuthor,
			SourcePostID: msg.SourcePostID,
			Links:        msg.Links,
		}

		postMeta := model.PostMeta{
			CreatedAt:               ctx.BlockHeader().Time.Unix(),
			LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
			LastActivityAt:          ctx.BlockHeader().Time.Unix(),
			AllowReplies:            true,
			IsDeleted:               false,
			RedistributionSplitRate: sdk.ZeroRat(),
			TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
			TotalReportCoinDay:      types.NewCoinFromInt64(0),
			TotalReward:             types.NewCoinFromInt64(0),
		}
		checkPostKVStore(t, ctx,
			types.GetPermlink(msg.Author, msg.PostID), postInfo, postMeta)
	}
}

func TestUpdatePost(t *testing.T) {
	ctx, am, _, pm, _, _, _ := setupTest(t, 1)
	baseTime := time.Now().Unix()
	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: time.Unix(baseTime, 0)})
	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")

	testCases := []struct {
		testName   string
		msg        UpdatePostMsg
		expectErr  sdk.Error
		updateTime int64
	}{
		{
			testName: "normal update",
			msg: NewUpdatePostMsg(
				string(user), postID, "update to this title", "update to this content",
				[]types.IDToURLMapping{{Identifier: "#1", URL: "https://lino.network"}}),
			expectErr:  nil,
			updateTime: baseTime + 10,
		},
		{
			testName: "update with invalid post id",
			msg: NewUpdatePostMsg(
				"invalid", postID, "update to this title", "update to this content",
				[]types.IDToURLMapping{{Identifier: "#1", URL: "https://lino.network"}}),
			expectErr:  model.ErrPostNotFound(model.GetPostInfoKey(types.GetPermlink("invalid", postID))),
			updateTime: baseTime + 100,
		},
		{
			testName: "update with invalid author",
			msg: NewUpdatePostMsg(
				string(user), "invalid", "update to this title", "update to this content",
				[]types.IDToURLMapping{{Identifier: "#1", URL: "https://lino.network"}}),
			expectErr:  model.ErrPostNotFound(model.GetPostInfoKey(types.GetPermlink(user, "invalid"))),
			updateTime: baseTime + 1000,
		},
	}

	for _, tc := range testCases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: time.Unix(tc.updateTime, 0)})

		err := pm.UpdatePost(
			ctx, tc.msg.Author, tc.msg.PostID, tc.msg.Title, tc.msg.Content, tc.msg.Links)
		if !assert.Equal(t, err, tc.expectErr) {
			t.Errorf("%s: diff err, got %v, want %v", tc.testName, err, tc.expectErr)
		}
		if tc.expectErr != nil {
			continue
		}

		postInfo := model.PostInfo{
			PostID:       tc.msg.PostID,
			Title:        tc.msg.Title,
			Content:      tc.msg.Content,
			Author:       tc.msg.Author,
			SourceAuthor: "",
			SourcePostID: "",
			Links:        tc.msg.Links,
		}

		postMeta := model.PostMeta{
			CreatedAt:               baseTime,
			LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
			LastActivityAt:          baseTime,
			AllowReplies:            true,
			IsDeleted:               false,
			TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
			TotalReportCoinDay:      types.NewCoinFromInt64(0),
			TotalReward:             types.NewCoinFromInt64(0),
			RedistributionSplitRate: sdk.ZeroRat(),
		}
		checkPostKVStore(t, ctx,
			types.GetPermlink(tc.msg.Author, tc.msg.PostID), postInfo, postMeta)
	}
}

// test get source post
func TestGetSourcePost(t *testing.T) {
	ctx, _, _, pm, _, _, _ := setupTest(t, 1)
	user1 := types.AccountKey("user1")
	user2 := types.AccountKey("user2")
	user3 := types.AccountKey("user3")

	testCases := []struct {
		testName           string
		postID             string
		author             types.AccountKey
		sourcePostID       string
		sourceAuthor       types.AccountKey
		expectSourcePostID string
		expectSourceAuthor types.AccountKey
	}{
		{
			testName:           "create post without source",
			postID:             "postID",
			author:             user1,
			sourcePostID:       "",
			sourceAuthor:       "",
			expectSourcePostID: "",
			expectSourceAuthor: "",
		},
		{
			testName:           "creat post with original source",
			postID:             "postID1",
			author:             user1,
			sourcePostID:       "postID",
			sourceAuthor:       user1,
			expectSourcePostID: "postID",
			expectSourceAuthor: user1,
		},
		{
			testName:           "create post with secondary source, but expect original source",
			postID:             "postID",
			author:             user2,
			sourcePostID:       "postID1",
			sourceAuthor:       user1,
			expectSourcePostID: "postID",
			expectSourceAuthor: user1,
		},
		{
			testName:           "create post with secodary source again, but expect orignal source",
			postID:             "postID",
			author:             user3,
			sourcePostID:       "postID",
			sourceAuthor:       user2,
			expectSourcePostID: "postID",
			expectSourceAuthor: user1,
		},
	}

	for _, tc := range testCases {
		msg := CreatePostMsg{
			PostID:       tc.postID,
			Title:        string(make([]byte, 50)),
			Content:      string(make([]byte, 1000)),
			Author:       tc.author,
			ParentAuthor: "",
			ParentPostID: "",
			SourceAuthor: tc.sourceAuthor,
			SourcePostID: tc.sourcePostID,
			Links:        nil,
			RedistributionSplitRate: "0",
		}
		err := pm.CreatePost(
			ctx, msg.Author, msg.PostID, msg.SourceAuthor, msg.SourcePostID,
			msg.ParentAuthor, msg.ParentPostID, msg.Content,
			msg.Title, sdk.ZeroRat(), msg.Links)
		if err != nil {
			t.Errorf("%s: failed to create post, got err %v", tc.testName, err)
		}

		sourceAuthor, sourcePostID, err :=
			pm.GetSourcePost(ctx, types.GetPermlink(tc.author, tc.postID))
		if err != nil {
			t.Errorf("%s: failed to get source post, got err %v", tc.testName, err)
		}
		if sourceAuthor != tc.expectSourceAuthor {
			t.Errorf("%s: diff source author, got %v, want %v", tc.testName, sourceAuthor, tc.expectSourceAuthor)
		}
		if sourcePostID != tc.expectSourcePostID {
			t.Errorf("%s: diff source post id, got %v, want %v", tc.testName, sourcePostID, tc.expectSourcePostID)
		}
	}
}

func TestAddOrUpdateViewToPost(t *testing.T) {
	ctx, am, _, pm, _, _, _ := setupTest(t, 1)
	createTime := ctx.BlockHeader().Time
	user1, postID1 := createTestPost(t, ctx, "user1", "postID1", am, pm, "0")
	user2, _ := createTestPost(t, ctx, "user2", "postID2", am, pm, "0")
	user3 := types.AccountKey("user3")

	testCases := []struct {
		testName             string
		viewUser             types.AccountKey
		postID               string
		author               types.AccountKey
		viewTime             int64
		expectTotalViewCount int64
		expectUserViewCount  int64
	}{
		{
			testName:             "user3 views (postID1, user1)",
			viewUser:             user3,
			postID:               postID1,
			author:               user1,
			viewTime:             1,
			expectTotalViewCount: 1,
			expectUserViewCount:  1,
		},
		{
			testName:             "user3 views (postID1, user1) again",
			viewUser:             user3,
			postID:               postID1,
			author:               user1,
			viewTime:             2,
			expectTotalViewCount: 2,
			expectUserViewCount:  2,
		},
		{
			testName:             "user2 views (postID1, user1)",
			viewUser:             user2,
			postID:               postID1,
			author:               user1,
			viewTime:             3,
			expectTotalViewCount: 3,
			expectUserViewCount:  1,
		},
		{
			testName:             "user2 views (postID1, user1) again",
			viewUser:             user2,
			postID:               postID1,
			author:               user1,
			viewTime:             4,
			expectTotalViewCount: 4,
			expectUserViewCount:  2,
		},
		{
			testName:             "user1 views (postID1, user1)",
			viewUser:             user1,
			postID:               postID1,
			author:               user1,
			viewTime:             5,
			expectTotalViewCount: 5,
			expectUserViewCount:  1,
		},
	}

	for _, tc := range testCases {
		postKey := types.GetPermlink(tc.author, tc.postID)
		ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(tc.viewTime, 0)})
		err := pm.AddOrUpdateViewToPost(ctx, postKey, tc.viewUser)
		if err != nil {
			t.Errorf("%s: failed to add or update view to post, got err %v", tc.testName, err)
		}

		postMeta := model.PostMeta{
			CreatedAt:               createTime.Unix(),
			LastUpdatedAt:           createTime.Unix(),
			LastActivityAt:          createTime.Unix(),
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat(),
			TotalViewCount:          tc.expectTotalViewCount,
			TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
			TotalReportCoinDay:      types.NewCoinFromInt64(0),
			TotalReward:             types.NewCoinFromInt64(0),
		}
		checkPostMeta(t, ctx, postKey, postMeta)
		view, err := pm.postStorage.GetPostView(ctx, postKey, tc.viewUser)
		if err != nil {
			t.Errorf("%s: failed to get post view, got err %v", tc.testName, err)
		}
		if view.Times != tc.expectUserViewCount {
			t.Errorf("%s: diff user view count, got %v, want %v", tc.testName, view.Times, tc.expectUserViewCount)
		}
		if view.LastViewAt != tc.viewTime {
			t.Errorf("%s: diff view time, got %v, want %v", tc.testName, view.LastViewAt, tc.viewTime)
		}
	}
}

func TestReportOrUpvoteToPost(t *testing.T) {
	ctx, am, _, pm, _, _, _ := setupTest(t, 1)
	user1, postID1 := createTestPost(t, ctx, "user1", "postID1", am, pm, "0")
	user2, _ := createTestPost(t, ctx, "user2", "postID2", am, pm, "0")
	user3 := types.AccountKey("user3")
	user4 := types.AccountKey("user4")

	permlink := types.GetPermlink(user1, postID1)

	testCases := []struct {
		testName                 string
		user                     types.AccountKey
		coinDay                  types.Coin
		isReport                 bool
		expectResult             sdk.Error
		expectTotalReportCoinDay types.Coin
		expectTotalUpvoteCoinDay types.Coin
	}{
		{
			testName:                 "user3 reports with 1 coin day",
			user:                     user3,
			coinDay:                  types.NewCoinFromInt64(1),
			isReport:                 true,
			expectResult:             nil,
			expectTotalReportCoinDay: types.NewCoinFromInt64(1),
			expectTotalUpvoteCoinDay: types.NewCoinFromInt64(0),
		},
		{
			testName:                 "user2 upvotes with 100 coin day",
			user:                     user2,
			coinDay:                  types.NewCoinFromInt64(100),
			isReport:                 false,
			expectResult:             nil,
			expectTotalReportCoinDay: types.NewCoinFromInt64(1),
			expectTotalUpvoteCoinDay: types.NewCoinFromInt64(100),
		},
		{
			testName:                 "user3 upvotes with 100 coin day and override previous report",
			user:                     user3,
			coinDay:                  types.NewCoinFromInt64(100),
			isReport:                 false,
			expectResult:             nil,
			expectTotalReportCoinDay: types.NewCoinFromInt64(0),
			expectTotalUpvoteCoinDay: types.NewCoinFromInt64(200),
		},
		{
			testName:                 "user4 upvotes with 100 coin day",
			user:                     user4,
			coinDay:                  types.NewCoinFromInt64(100),
			isReport:                 false,
			expectResult:             nil,
			expectTotalReportCoinDay: types.NewCoinFromInt64(0),
			expectTotalUpvoteCoinDay: types.NewCoinFromInt64(300),
		},
		{
			testName:                 "user3 report with 2 coin day which overrides previous upvote",
			user:                     user4,
			coinDay:                  types.NewCoinFromInt64(2),
			isReport:                 true,
			expectResult:             nil,
			expectTotalReportCoinDay: types.NewCoinFromInt64(2),
			expectTotalUpvoteCoinDay: types.NewCoinFromInt64(200),
		},
	}

	for _, tc := range testCases {
		err := pm.ReportOrUpvoteToPost(ctx, permlink, tc.user, tc.coinDay, tc.isReport)
		if !assert.Equal(t, tc.expectResult, err) {
			t.Errorf("%s: diff err, got %v, want %v", tc.testName, err, tc.expectResult)
		}

		postMeta := model.PostMeta{
			CreatedAt:               ctx.BlockHeader().Time.Unix(),
			LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
			LastActivityAt:          ctx.BlockHeader().Time.Unix(),
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat(),
			TotalReportCoinDay:      tc.expectTotalReportCoinDay,
			TotalUpvoteCoinDay:      tc.expectTotalUpvoteCoinDay,
			TotalReward:             types.NewCoinFromInt64(0),
		}
		checkPostMeta(t, ctx, permlink, postMeta)
	}
}

func TestDonation(t *testing.T) {
	ctx, am, _, pm, _, _, _ := setupTest(t, 1)
	user1, postID1 := createTestPost(t, ctx, "user1", "postID1", am, pm, "0")
	user2, postID2 := createTestPost(t, ctx, "user2", "postID2", am, pm, "0")
	user3 := types.AccountKey("user3")

	baseTime := ctx.BlockHeader().Time.Unix()
	testCases := []struct {
		testName            string
		user                types.AccountKey
		donateAt            int64
		amount              types.Coin
		donationType        types.DonationType
		postID              string
		author              types.AccountKey
		expectDonateCount   int64
		expectTotalDonation types.Coin
		expectDonationList  model.Donations
	}{
		{
			testName:            "user3 donates to (postID1, user1)",
			user:                user3,
			donateAt:            baseTime,
			amount:              types.NewCoinFromInt64(1),
			donationType:        types.DirectDeposit,
			postID:              postID1,
			author:              user1,
			expectDonateCount:   1,
			expectTotalDonation: types.NewCoinFromInt64(1),
			expectDonationList: model.Donations{
				Username: user3,
				Amount:   types.NewCoinFromInt64(1),
				Times:    1,
			},
		},
		{
			testName:            "user3 donates to (postID2, user2)",
			user:                user3,
			donateAt:            baseTime,
			amount:              types.NewCoinFromInt64(1),
			donationType:        types.Inflation,
			postID:              postID2,
			author:              user2,
			expectDonateCount:   1,
			expectTotalDonation: types.NewCoinFromInt64(1),
			expectDonationList: model.Donations{
				Username: user3,
				Amount:   types.NewCoinFromInt64(1),
				Times:    1,
			},
		},
		{
			testName:            "user3 donates to (postID2, user2) again",
			user:                user3,
			donateAt:            baseTime,
			amount:              types.NewCoinFromInt64(20),
			donationType:        types.DirectDeposit,
			postID:              postID2,
			author:              user2,
			expectDonateCount:   2,
			expectTotalDonation: types.NewCoinFromInt64(21),
			expectDonationList: model.Donations{
				Username: user3,
				Amount:   types.NewCoinFromInt64(21),
				Times:    2,
			},
		},
	}

	for _, tc := range testCases {
		postKey := types.GetPermlink(tc.author, tc.postID)
		err := pm.AddDonation(ctx, postKey, tc.user, tc.amount, tc.donationType)
		if err != nil {
			t.Errorf("%s: failed to add donation, got err %v", tc.testName, err)
		}

		postMeta := model.PostMeta{
			CreatedAt:               ctx.BlockHeader().Time.Unix(),
			LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
			LastActivityAt:          ctx.BlockHeader().Time.Unix(),
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat(),
			TotalDonateCount:        tc.expectDonateCount,
			TotalReward:             tc.expectTotalDonation,
			TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
			TotalReportCoinDay:      types.NewCoinFromInt64(0),
		}
		checkPostMeta(t, ctx, postKey, postMeta)
		storage := model.NewPostStorage(testPostKVStoreKey)
		donations, _ := storage.GetPostDonations(ctx, postKey, tc.user)
		if !assert.Equal(t, tc.expectDonationList, *donations) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, *donations, tc.expectDonationList)
		}
	}
}

func TestGetPenaltyScore(t *testing.T) {
	ctx, am, _, pm, _, _, _ := setupTest(t, 1)
	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	postKey := types.GetPermlink(user, postID)
	bigString1 := "1000000000000000000000000"
	bigString2 := "7777777777777777777777777"
	bigStringInt1, _ := new(big.Int).SetString(bigString1, 10)
	bigStringInt2, _ := new(big.Int).SetString(bigString2, 10)
	testCases := []struct {
		testName           string
		totalReportCoinDay types.Coin
		totalUpvoteCoinDay types.Coin
		expectPenaltyScore sdk.Rat
	}{
		{
			testName:           "1 report and 0 upvote expects 1 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(1),
			totalUpvoteCoinDay: types.NewCoinFromInt64(0),
			expectPenaltyScore: sdk.NewRat(1, 1),
		},
		{
			testName:           "0 report and 1 upvote expects 0 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(0),
			totalUpvoteCoinDay: types.NewCoinFromInt64(1),
			expectPenaltyScore: sdk.NewRat(0, 1),
		},
		{
			testName:           "0 report and 0 upvote expects 0 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(0),
			totalUpvoteCoinDay: types.NewCoinFromInt64(0),
			expectPenaltyScore: sdk.NewRat(0, 1),
		},
		{
			testName:           "100 report and 100 upvote expects 1 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(100),
			totalUpvoteCoinDay: types.NewCoinFromInt64(100),
			expectPenaltyScore: sdk.NewRat(1, 1),
		},
		{
			testName:           "1000 report and 100 upvote expects 1 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(1000),
			totalUpvoteCoinDay: types.NewCoinFromInt64(100),
			expectPenaltyScore: sdk.NewRat(1, 1),
		},
		{
			testName:           "50 report and 100 upvote expects 1/2 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(50),
			totalUpvoteCoinDay: types.NewCoinFromInt64(100),
			expectPenaltyScore: sdk.NewRat(1, 2),
		},
		// issue https://github.com/lino-network/lino/issues/150
		{
			testName:           "3333 report and 7777 upvote expects 3/7 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(3333),
			totalUpvoteCoinDay: types.NewCoinFromInt64(7777),
			expectPenaltyScore: sdk.NewRat(2142857, 5000000),
		},
		{
			testName:           "big string report and big string upvote, report is much than upvote",
			totalReportCoinDay: types.NewCoinFromBigInt(bigStringInt2),
			totalUpvoteCoinDay: types.NewCoinFromBigInt(bigStringInt1),
			expectPenaltyScore: sdk.NewRat(1),
		},
		{
			testName:           "big string report and big string upvote, report is less than upvote",
			totalReportCoinDay: types.NewCoinFromBigInt(bigStringInt1),
			totalUpvoteCoinDay: types.NewCoinFromBigInt(bigStringInt2),
			expectPenaltyScore: sdk.NewRat(642857, 5000000),
		},
	}

	for _, tc := range testCases {
		postMeta := &model.PostMeta{
			CreatedAt:               ctx.BlockHeader().Time.Unix(),
			LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
			LastActivityAt:          ctx.BlockHeader().Time.Unix(),
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat(),
			TotalReportCoinDay:      tc.totalReportCoinDay,
			TotalUpvoteCoinDay:      tc.totalUpvoteCoinDay,
		}
		err := pm.postStorage.SetPostMeta(ctx, postKey, postMeta)
		if err != nil {
			t.Errorf("%s: failed to set post meta, got err %v", tc.testName, err)
		}

		penaltyScore, err := pm.GetPenaltyScore(ctx, postKey)
		if err != nil {
			t.Errorf("%s: failed to get penalty score, got err %v", tc.testName, err)
		}
		if !penaltyScore.Equal(tc.expectPenaltyScore) {
			t.Errorf("%s: diff penalty score, got %v, want %v", tc.testName, penaltyScore, tc.expectPenaltyScore)
		}
		assert.Equal(t, penaltyScore, tc.expectPenaltyScore)
	}
}

func TestGetRepostPenaltyScore(t *testing.T) {
	ctx, am, _, pm, _, _, _ := setupTest(t, 1)
	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user2, postID2 := createTestRepost(t, ctx, "user2", "repost", am, pm, user, postID)

	postKey := types.GetPermlink(user, postID)
	repostKey := types.GetPermlink(user2, postID2)
	testCases := []struct {
		testName           string
		totalReportCoinDay types.Coin
		totalUpvoteCoinDay types.Coin
		expectPenaltyScore sdk.Rat
	}{
		{
			testName:           "1 report and 0 upvote expects 1 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(1),
			totalUpvoteCoinDay: types.NewCoinFromInt64(0),
			expectPenaltyScore: sdk.NewRat(1, 1),
		},
		{
			testName:           "0 report and 1 upvote expects 0 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(0),
			totalUpvoteCoinDay: types.NewCoinFromInt64(1),
			expectPenaltyScore: sdk.NewRat(0, 1),
		},
		{
			testName:           "0 report and 0 upvote expects 0 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(0),
			totalUpvoteCoinDay: types.NewCoinFromInt64(0),
			expectPenaltyScore: sdk.NewRat(0, 1),
		},
		{
			testName:           "100 report and 100 upvote expects 1 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(100),
			totalUpvoteCoinDay: types.NewCoinFromInt64(100),
			expectPenaltyScore: sdk.NewRat(1, 1),
		},
		{
			testName:           "1000 report and 100 upvote expects 1 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(1000),
			totalUpvoteCoinDay: types.NewCoinFromInt64(100),
			expectPenaltyScore: sdk.NewRat(1, 1),
		},
		{
			testName:           "50 report and 100 upvote expects 1/2 penalty score",
			totalReportCoinDay: types.NewCoinFromInt64(50),
			totalUpvoteCoinDay: types.NewCoinFromInt64(100),
			expectPenaltyScore: sdk.NewRat(1, 2),
		},
	}

	for _, tc := range testCases {
		postMeta := &model.PostMeta{
			CreatedAt:               ctx.BlockHeader().Time.Unix(),
			LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
			LastActivityAt:          ctx.BlockHeader().Time.Unix(),
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat(),
			TotalReportCoinDay:      tc.totalReportCoinDay,
			TotalUpvoteCoinDay:      tc.totalUpvoteCoinDay,
		}
		err := pm.postStorage.SetPostMeta(ctx, postKey, postMeta)
		if err != nil {
			t.Errorf("%s: failed to set post meta, got err %v", tc.testName, err)
		}

		penaltyScore, err := pm.GetPenaltyScore(ctx, repostKey)
		if err != nil {
			t.Errorf("%s: failed to get penalty score, got err %v", tc.testName, err)
		}
		if !penaltyScore.Equal(tc.expectPenaltyScore) {
			t.Errorf("%s: diff penalty score, got %v, want %v", tc.testName, penaltyScore, tc.expectPenaltyScore)
		}
	}
}

func checkIsDelete(t *testing.T, ctx sdk.Context, pm PostManager, permlink types.Permlink) {
	isDeleted, err := pm.IsDeleted(ctx, permlink)
	assert.Nil(t, err)
	assert.Equal(t, true, isDeleted)
	postInfo, err := pm.postStorage.GetPostInfo(ctx, permlink)
	assert.Nil(t, err)
	assert.Equal(t, "", postInfo.Title)
	assert.Equal(t, "", postInfo.Content)
}

func TestDeletePost(t *testing.T) {
	ctx, am, _, pm, _, _, _ := setupTest(t, 1)
	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user2, postID2 := createTestRepost(t, ctx, "user2", "repost", am, pm, user, postID)

	err := pm.DeletePost(ctx, types.GetPermlink(user2, postID2))
	assert.Nil(t, err)
	checkIsDelete(t, ctx, pm, types.GetPermlink(user2, postID2))
	postMeta, err := pm.postStorage.GetPostMeta(ctx, types.GetPermlink(user, postID))
	assert.Nil(t, err)
	assert.Equal(t, false, postMeta.IsDeleted)
	err = pm.DeletePost(ctx, types.GetPermlink(user, postID))
	assert.Nil(t, err)
	checkIsDelete(t, ctx, pm, types.GetPermlink(user, postID))
}
