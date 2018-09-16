package post

import (
	"testing"
	"time"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/post/model"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	accmodel "github.com/lino-network/lino/x/account/model"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestHandlerCreatePost(t *testing.T) {
	ctx, am, ph, pm, gm, dm, _, rm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm, rm)
	postParam, _ := ph.GetPostParam(ctx)

	user := createTestAccount(t, ctx, am, "user1")

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(postParam.PostIntervalSec, 0)})
	// test valid post
	msg := CreatePostMsg{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        nil,
		RedistributionSplitRate: "0",
	}
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	assert.True(t, pm.DoesPostExist(ctx, types.GetPermlink(msg.Author, msg.PostID)))

	// test invlaid author
	msg.Author = types.AccountKey("invalid")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrAccountNotFound(msg.Author).Result())

	// test duplicate post
	msg.Author = user
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostAlreadyExist(types.GetPermlink(user, msg.PostID)).Result())

	// test post too often
	msg.PostID = "Post too often"
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostTooOften(msg.Author).Result())
}

func TestHandlerUpdatePost(t *testing.T) {
	ctx, am, _, pm, gm, dm, _, rm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm, rm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user1, postID1 := createTestPost(t, ctx, "user1", "postID1", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	err := pm.DeletePost(ctx, types.GetPermlink(user1, postID1))
	assert.Nil(t, err)

	testCases := map[string]struct {
		msg        UpdatePostMsg
		wantResult sdk.Result
	}{
		"normal update": {
			msg:        NewUpdatePostMsg(string(user), postID, "update title", "update content", []types.IDToURLMapping(nil)),
			wantResult: sdk.Result{},
		},
		"update author doesn't exist": {
			msg:        NewUpdatePostMsg("invalid", postID, "update title", "update content", []types.IDToURLMapping(nil)),
			wantResult: ErrAccountNotFound("invalid").Result(),
		},
		"update post doesn't exist - invalid post ID": {
			msg:        NewUpdatePostMsg(string(user), "invalid", "update title", "update content", []types.IDToURLMapping(nil)),
			wantResult: ErrPostNotFound(types.GetPermlink(user, "invalid")).Result(),
		},
		"update post doesn't exist - invalid author": {
			msg:        NewUpdatePostMsg(string(user2), postID, "update title", "update content", []types.IDToURLMapping(nil)),
			wantResult: ErrPostNotFound(types.GetPermlink(user2, postID)).Result(),
		},
		"update deleted post": {
			msg:        NewUpdatePostMsg(string(user1), postID1, "update title", "update content", []types.IDToURLMapping(nil)),
			wantResult: ErrUpdatePostIsDeleted(types.GetPermlink(user1, postID1)).Result(),
		},
	}
	for testName, tc := range testCases {
		result := handler(ctx, tc.msg)
		if !assert.Equal(t, tc.wantResult, result) {
			t.Errorf("%s: diff result, got %v, want %v", testName, result, tc.wantResult)
		}
		if tc.wantResult.Code != sdk.ABCICodeOK {
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
			CreatedAt:               ctx.BlockHeader().Time.Unix(),
			LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
			LastActivityAt:          ctx.BlockHeader().Time.Unix(),
			AllowReplies:            true,
			IsDeleted:               false,
			TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
			TotalReward:             types.NewCoinFromInt64(0),
			TotalReportCoinDay:      types.NewCoinFromInt64(0),
			RedistributionSplitRate: sdk.ZeroRat(),
		}
		checkPostKVStore(t, ctx,
			types.GetPermlink(tc.msg.Author, tc.msg.PostID), postInfo, postMeta)
	}
}

func TestHandlerDeletePost(t *testing.T) {
	ctx, am, _, pm, gm, dm, _, rm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm, rm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user1 := createTestAccount(t, ctx, am, "user1")

	testCases := map[string]struct {
		msg        DeletePostMsg
		wantResult sdk.Result
	}{
		"normal delete": {
			msg: DeletePostMsg{
				Author: user,
				PostID: postID,
			},
			wantResult: sdk.Result{},
		},
		"author doesn't exist": {
			msg: DeletePostMsg{
				Author: types.AccountKey("invalid"),
				PostID: postID,
			},
			wantResult: ErrAccountNotFound("invalid").Result(),
		},
		"post doesn't exist - invalid author": {
			msg: DeletePostMsg{
				Author: user1,
				PostID: "postID",
			},
			wantResult: ErrPostNotFound(types.GetPermlink(user1, postID)).Result(),
		},
		"post doesn't exist - invalid postID": {
			msg: DeletePostMsg{
				Author: user,
				PostID: "invalid",
			},
			wantResult: ErrPostNotFound(types.GetPermlink(user, "invalid")).Result(),
		},
	}
	for testName, tc := range testCases {
		result := handler(ctx, tc.msg)
		if !assert.Equal(t, tc.wantResult, result) {
			t.Errorf("%s: diff result, got %v, want %v", testName, result, tc.wantResult)
		}
	}
}

func TestHandlerCreateComment(t *testing.T) {
	ctx, am, ph, pm, gm, dm, _, rm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm, rm)
	postParam, err := ph.GetPostParam(ctx)
	assert.Nil(t, err)

	baseTime := time.Now()
	baseTime1 := baseTime.Add(time.Duration(postParam.PostIntervalSec) * time.Second)
	baseTime2 := baseTime1.Add(time.Duration(postParam.PostIntervalSec) * time.Second)
	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: baseTime})
	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")

	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: baseTime1})
	// test comment
	msg := CreatePostMsg{
		PostID:       "comment",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: user,
		ParentPostID: postID,
		SourceAuthor: "",
		SourcePostID: "",
		Links:        nil,
		RedistributionSplitRate: "0",
	}
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	postInfo := model.PostInfo{
		PostID:       msg.PostID,
		Title:        msg.Title,
		Content:      msg.Content,
		Author:       msg.Author,
		ParentAuthor: msg.ParentAuthor,
		ParentPostID: msg.ParentPostID,
		SourceAuthor: msg.SourceAuthor,
		SourcePostID: msg.SourcePostID,
		Links:        msg.Links,
	}

	postMeta := model.PostMeta{
		CreatedAt:               baseTime1.Unix(),
		LastUpdatedAt:           baseTime1.Unix(),
		LastActivityAt:          baseTime1.Unix(),
		AllowReplies:            true,
		TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
		TotalReward:             types.NewCoinFromInt64(0),
		TotalReportCoinDay:      types.NewCoinFromInt64(0),
		RedistributionSplitRate: sdk.ZeroRat(),
	}

	checkPostKVStore(t, ctx, types.GetPermlink(user, "comment"), postInfo, postMeta)

	// check parent
	postInfo.PostID = postID
	postInfo.ParentAuthor = ""
	postInfo.ParentPostID = ""
	postMeta.CreatedAt = baseTime.Unix()
	postMeta.LastUpdatedAt = baseTime.Unix()
	checkPostKVStore(t, ctx, types.GetPermlink(user, postID), postInfo, postMeta)

	// test post too often
	msg.PostID = "post too often"

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostTooOften(user).Result())

	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: baseTime2})
	// test invalid parent
	msg.PostID = "invalid post"
	msg.ParentAuthor = user
	msg.ParentPostID = "invalid parent"

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostNotFound(types.GetPermlink(user, msg.ParentPostID)).Result())

	// test duplicate comment
	msg.Author = user
	msg.PostID = "comment"
	msg.ParentAuthor = user
	msg.ParentPostID = "TestPostID"

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostAlreadyExist(types.GetPermlink(msg.Author, msg.PostID)).Result())

	// test cycle comment
	msg.Author = user
	msg.PostID = "newComment"
	msg.ParentAuthor = user
	msg.ParentPostID = "newComment"

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostNotFound(types.GetPermlink(user, msg.PostID)).Result())
}

func TestHandlerRepost(t *testing.T) {
	ctx, am, ph, pm, gm, dm, _, rm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm, rm)
	postParam, err := ph.GetPostParam(ctx)
	assert.Nil(t, err)

	baseTime := time.Now()
	baseTime1 := baseTime.Add(time.Duration(postParam.PostIntervalSec) * time.Second)
	baseTime2 := baseTime1.Add(time.Duration(postParam.PostIntervalSec) * time.Second)
	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: baseTime})
	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")

	// test repost
	msg := CreatePostMsg{
		PostID:       "repost",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: user,
		SourcePostID: postID,
		Links:        nil,
		RedistributionSplitRate: "0",
	}
	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: baseTime1})
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	postInfo := model.PostInfo{
		PostID:       msg.PostID,
		Title:        msg.Title,
		Content:      msg.Content,
		Author:       msg.Author,
		ParentAuthor: msg.ParentAuthor,
		ParentPostID: msg.ParentPostID,
		SourceAuthor: msg.SourceAuthor,
		SourcePostID: msg.SourcePostID,
		Links:        msg.Links,
	}

	postMeta := model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time.Unix(),
		LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
		LastActivityAt:          ctx.BlockHeader().Time.Unix(),
		AllowReplies:            true,
		TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
		TotalReward:             types.NewCoinFromInt64(0),
		TotalReportCoinDay:      types.NewCoinFromInt64(0),
		RedistributionSplitRate: sdk.ZeroRat(),
	}

	checkPostKVStore(t, ctx, types.GetPermlink(user, "repost"), postInfo, postMeta)

	// test 2 depth repost
	msg.PostID = "repost-repost"
	msg.SourceAuthor = user
	msg.SourcePostID = "repost"
	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: baseTime2})
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	// check 2 depth repost
	postInfo.PostID = "repost-repost"
	postMeta = model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time.Unix(),
		LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
		LastActivityAt:          ctx.BlockHeader().Time.Unix(),
		AllowReplies:            true,
		TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
		TotalReward:             types.NewCoinFromInt64(0),
		TotalReportCoinDay:      types.NewCoinFromInt64(0),
		RedistributionSplitRate: sdk.ZeroRat(),
	}
	postInfo.SourceAuthor = user
	postInfo.SourcePostID = postID
	checkPostKVStore(t, ctx, types.GetPermlink(user, postInfo.PostID), postInfo, postMeta)
}

func TestHandlerPostDonate(t *testing.T) {
	ctx, am, ph, pm, gm, dm, _, rm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm, rm)

	accParam, err := ph.GetAccountParam(ctx)
	assert.Nil(t, err)

	author, postID := createTestPost(t, ctx, "author", "postID", am, pm, "0")
	author1, deletedPostID := createTestPost(t, ctx, "author1", "delete", am, pm, "0")

	pm.DeletePost(ctx, types.GetPermlink(author1, deletedPostID))

	userWithSufficientSaving := createTestAccount(t, ctx, am, "userWithSufficientSaving")
	err = am.AddSavingCoin(
		ctx, userWithSufficientSaving, types.NewCoinFromInt64(100*types.Decimals),
		referrer, "", types.TransferIn)
	assert.Nil(t, err)

	secondUserWithSufficientSaving := createTestAccount(t, ctx, am, "secondUserWithSufficientSaving")
	err = am.AddSavingCoin(
		ctx, secondUserWithSufficientSaving, types.NewCoinFromInt64(100*types.Decimals),
		referrer, "", types.TransferIn)
	assert.Nil(t, err)

	micropaymentUser := createTestAccount(t, ctx, am, "micropaymentUser")
	err = am.AddSavingCoin(
		ctx, micropaymentUser, types.NewCoinFromInt64(1*types.Decimals),
		referrer, "", types.TransferIn)
	assert.Nil(t, err)
	testCases := []struct {
		testName            string
		donateUser          types.AccountKey
		amount              types.LNO
		toAuthor            types.AccountKey
		toPostID            string
		expectErr           sdk.Result
		expectPostMeta      model.PostMeta
		expectDonatorSaving types.Coin
		expectAuthorSaving  types.Coin
		//https://github.com/lino-network/lino/issues/154
		expectRegisteredEvent             RewardEvent
		expectDonateTimesFromUserToAuthor int64
		expectCumulativeConsumption       types.Coin
		expectAuthorReward                accmodel.Reward
	}{
		{
			testName:   "donate from sufficient saving",
			donateUser: userWithSufficientSaving,
			amount:     types.LNO("100"),
			toAuthor:   author,
			toPostID:   postID,
			expectErr:  sdk.Result{},
			expectPostMeta: model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time.Unix(),
				LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
				AllowReplies:            true,
				TotalDonateCount:        1,
				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
				TotalReportCoinDay:      types.NewCoinFromInt64(0),
				TotalReward:             types.NewCoinFromInt64(95 * types.Decimals),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			expectDonatorSaving: accParam.RegisterFee,
			expectAuthorSaving: accParam.RegisterFee.Plus(
				types.NewCoinFromInt64(95 * types.Decimals)),
			expectRegisteredEvent: RewardEvent{
				PostAuthor: author,
				PostID:     postID,
				Consumer:   userWithSufficientSaving,
				Evaluate:   types.NewCoinFromInt64(59380),
				Original:   types.NewCoinFromInt64(100 * types.Decimals),
				Friction:   types.NewCoinFromInt64(5 * types.Decimals),
				FromApp:    "",
			},
			expectDonateTimesFromUserToAuthor: 1,
			expectCumulativeConsumption:       types.NewCoinFromInt64(100 * types.Decimals),
			expectAuthorReward: accmodel.Reward{
				TotalIncome:    types.NewCoinFromInt64(95 * types.Decimals),
				OriginalIncome: types.NewCoinFromInt64(95 * types.Decimals),
			},
		},
		{
			testName:            "donate from insufficient saving",
			donateUser:          userWithSufficientSaving,
			amount:              types.LNO("100"),
			toAuthor:            author,
			toPostID:            postID,
			expectErr:           acc.ErrAccountSavingCoinNotEnough().Result(),
			expectPostMeta:      model.PostMeta{},
			expectDonatorSaving: accParam.RegisterFee,
			expectAuthorSaving: accParam.RegisterFee.Plus(
				types.NewCoinFromInt64(95 * types.Decimals)),
			expectRegisteredEvent:             RewardEvent{},
			expectDonateTimesFromUserToAuthor: 1,
			expectCumulativeConsumption:       types.NewCoinFromInt64(100 * types.Decimals),
			expectAuthorReward: accmodel.Reward{
				TotalIncome:    types.NewCoinFromInt64(95 * types.Decimals),
				OriginalIncome: types.NewCoinFromInt64(95 * types.Decimals),
			},
		},
		{
			testName:   "donate less money from second user with sufficient saving",
			donateUser: secondUserWithSufficientSaving,
			amount:     types.LNO("50"),
			toAuthor:   author,
			toPostID:   postID,
			expectErr:  sdk.Result{},
			expectPostMeta: model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time.Unix(),
				LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
				AllowReplies:            true,
				TotalDonateCount:        2,
				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
				TotalReportCoinDay:      types.NewCoinFromInt64(0),
				TotalReward:             types.NewCoinFromInt64(14250000),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			expectDonatorSaving: accParam.RegisterFee.Plus(
				types.NewCoinFromInt64(50 * types.Decimals)),
			expectAuthorSaving: accParam.RegisterFee.Plus(types.NewCoinFromInt64(14250000)),
			expectRegisteredEvent: RewardEvent{
				PostAuthor: author,
				PostID:     postID,
				Consumer:   secondUserWithSufficientSaving,
				Evaluate:   types.NewCoinFromInt64(59361),
				Original:   types.NewCoinFromInt64(50 * types.Decimals),
				Friction:   types.NewCoinFromInt64(250000),
				FromApp:    "",
			},
			expectDonateTimesFromUserToAuthor: 1,
			expectCumulativeConsumption:       types.NewCoinFromInt64(150 * types.Decimals),
			expectAuthorReward: accmodel.Reward{
				TotalIncome:    types.NewCoinFromInt64(14250000),
				OriginalIncome: types.NewCoinFromInt64(14250000),
			},
		},
		{
			testName:   "donate second times from second user with sufficient saving (donate stake is zero)",
			donateUser: secondUserWithSufficientSaving,
			amount:     types.LNO("50"),
			toAuthor:   author,
			toPostID:   postID,
			expectErr:  sdk.Result{},
			expectPostMeta: model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time.Unix(),
				LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
				AllowReplies:            true,
				TotalDonateCount:        3,
				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
				TotalReportCoinDay:      types.NewCoinFromInt64(0),
				TotalReward:             types.NewCoinFromInt64(190 * types.Decimals),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			expectDonatorSaving: accParam.RegisterFee,
			expectAuthorSaving:  accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
			expectRegisteredEvent: RewardEvent{
				PostAuthor: author,
				PostID:     postID,
				Consumer:   secondUserWithSufficientSaving,
				Evaluate:   types.NewCoinFromInt64(0),
				Original:   types.NewCoinFromInt64(50 * types.Decimals),
				Friction:   types.NewCoinFromInt64(250000),
				FromApp:    "",
			},
			expectDonateTimesFromUserToAuthor: 2,
			expectCumulativeConsumption:       types.NewCoinFromInt64(200 * types.Decimals),
			expectAuthorReward: accmodel.Reward{
				TotalIncome:    types.NewCoinFromInt64(190 * types.Decimals),
				OriginalIncome: types.NewCoinFromInt64(190 * types.Decimals),
			},
		},
		{
			testName:   "micropayment",
			donateUser: micropaymentUser,
			amount:     types.LNO("0.00001"),
			toAuthor:   author,
			toPostID:   postID,
			expectErr:  sdk.Result{},
			expectPostMeta: model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time.Unix(),
				LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
				AllowReplies:            true,
				TotalDonateCount:        4,
				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
				TotalReportCoinDay:      types.NewCoinFromInt64(0),
				TotalReward:             types.NewCoinFromInt64(19000001),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			expectDonatorSaving: types.NewCoinFromInt64(199999),
			expectAuthorSaving:  accParam.RegisterFee.Plus(types.NewCoinFromInt64(19000001)),
			expectRegisteredEvent: RewardEvent{
				PostAuthor: author,
				PostID:     postID,
				Consumer:   micropaymentUser,
				Evaluate:   types.NewCoinFromInt64(5),
				Original:   types.NewCoinFromInt64(1),
				Friction:   types.NewCoinFromInt64(0),
				FromApp:    "",
			},
			expectDonateTimesFromUserToAuthor: 1,
			expectCumulativeConsumption:       types.NewCoinFromInt64(20000001),
			expectAuthorReward: accmodel.Reward{
				TotalIncome:    types.NewCoinFromInt64(19000001),
				OriginalIncome: types.NewCoinFromInt64(19000001),
			},
		},
		{
			testName:                          "invalid target postID",
			donateUser:                        userWithSufficientSaving,
			amount:                            types.LNO("1"),
			toAuthor:                          author,
			toPostID:                          "invalid",
			expectErr:                         ErrPostNotFound(types.GetPermlink(author, "invalid")).Result(),
			expectPostMeta:                    model.PostMeta{},
			expectDonatorSaving:               accParam.RegisterFee,
			expectAuthorSaving:                accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
			expectRegisteredEvent:             RewardEvent{},
			expectDonateTimesFromUserToAuthor: 1,
			expectCumulativeConsumption:       types.NewCoinFromInt64(200 * types.Decimals),
			expectAuthorReward: accmodel.Reward{
				TotalIncome:    types.NewCoinFromInt64(19000001),
				OriginalIncome: types.NewCoinFromInt64(19000001),
			},
		},
		{
			testName:                          "invalid target author",
			donateUser:                        userWithSufficientSaving,
			amount:                            types.LNO("1"),
			toAuthor:                          types.AccountKey("invalid"),
			toPostID:                          postID,
			expectErr:                         ErrPostNotFound(types.GetPermlink(types.AccountKey("invalid"), postID)).Result(),
			expectPostMeta:                    model.PostMeta{},
			expectDonatorSaving:               accParam.RegisterFee,
			expectAuthorSaving:                accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
			expectRegisteredEvent:             RewardEvent{},
			expectDonateTimesFromUserToAuthor: 0,
			expectCumulativeConsumption:       types.NewCoinFromInt64(200 * types.Decimals),
			expectAuthorReward: accmodel.Reward{
				TotalIncome:    types.NewCoinFromInt64(19000001),
				OriginalIncome: types.NewCoinFromInt64(19000001),
			},
		},
		{
			testName:   "donate to self",
			donateUser: author,
			amount:     types.LNO("100"),
			toAuthor:   author,
			toPostID:   postID,
			expectErr:  ErrCannotDonateToSelf(author).Result(),
			expectPostMeta: model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time.Unix(),
				LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
				LastActivityAt:          ctx.BlockHeader().Time.Unix(),
				AllowReplies:            true,
				TotalDonateCount:        2,
				TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
				TotalReportCoinDay:      types.NewCoinFromInt64(0),
				TotalReward:             types.NewCoinFromInt64(19000001),
				RedistributionSplitRate: sdk.ZeroRat(),
			},
			expectDonatorSaving:               accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
			expectAuthorSaving:                accParam.RegisterFee.Plus(types.NewCoinFromInt64(190 * types.Decimals)),
			expectRegisteredEvent:             RewardEvent{},
			expectDonateTimesFromUserToAuthor: 0,
			expectCumulativeConsumption:       types.NewCoinFromInt64(20000000),
			expectAuthorReward: accmodel.Reward{
				TotalIncome:    types.NewCoinFromInt64(19000001),
				OriginalIncome: types.NewCoinFromInt64(19000001),
			},
		},
		{
			testName:   "donate to deleted post",
			donateUser: userWithSufficientSaving,
			amount:     types.LNO("1"),
			toAuthor:   author1,
			toPostID:   deletedPostID,
			expectErr:  ErrDonatePostIsDeleted(types.GetPermlink(author1, deletedPostID)).Result(),
		},
	}

	for _, tc := range testCases {
		donateMsg := NewDonateMsg(
			string(tc.donateUser), tc.amount, string(tc.toAuthor), tc.toPostID, "", memo1)
		result := handler(ctx, donateMsg)
		if !assert.Equal(t, tc.expectErr, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectErr)
		}
		if tc.expectErr.Code == sdk.ABCICodeOK {
			checkPostMeta(t, ctx, types.GetPermlink(tc.toAuthor, tc.toPostID), tc.expectPostMeta)
		} else {
			continue
		}

		authorSaving, err := am.GetSavingFromBank(ctx, tc.toAuthor)
		if err != nil {
			t.Errorf("%s: failed to get author saving from bank, got err %v", tc.testName, err)
		}
		if !authorSaving.IsEqual(tc.expectAuthorSaving) {
			t.Errorf("%s: diff author saving, got %v, want %v", tc.testName, authorSaving, tc.expectAuthorSaving)
			return
		}

		donatorSaving, err := am.GetSavingFromBank(ctx, tc.donateUser)
		if err != nil {
			t.Errorf("%s: failed to get donator saving from bank, got err %v", tc.testName, err)
		}
		if !donatorSaving.IsEqual(tc.expectDonatorSaving) {
			t.Errorf("%s: diff donator saving %v, got %v", tc.testName, donatorSaving, tc.expectDonatorSaving)
			return
		}

		if tc.expectErr.Code == sdk.ABCICodeOK {
			eventList := gm.GetTimeEventListAtTime(ctx, ctx.BlockHeader().Time.Unix()+3600*7*24)
			if !assert.Equal(t, tc.expectRegisteredEvent, eventList.Events[len(eventList.Events)-1]) {
				t.Errorf("%s: diff event, got %v, want %v", tc.testName,
					eventList.Events[len(eventList.Events)-1], tc.expectRegisteredEvent)
			}

			as := accmodel.NewAccountStorage(testAccountKVStoreKey)
			reward, err := as.GetReward(ctx, tc.toAuthor)
			if err != nil {
				t.Errorf("%s: failed to get reward, got err %v", tc.testName, err)
			}
			tc.expectAuthorReward.FrictionIncome = types.NewCoinFromInt64(0)
			tc.expectAuthorReward.InflationIncome = types.NewCoinFromInt64(0)
			tc.expectAuthorReward.UnclaimReward = types.NewCoinFromInt64(0)
			tc.expectAuthorReward.Interest = types.NewCoinFromInt64(0)
			if !assert.Equal(t, tc.expectAuthorReward, *reward) {
				t.Errorf("%s: diff reward, got %v, want %v", tc.testName, *reward, tc.expectAuthorReward)
			}
		}

		times, err := am.GetDonationRelationship(ctx, tc.toAuthor, tc.donateUser)
		if err != nil {
			t.Errorf("%s: failed to get donation relationship, got err %v", tc.testName, err)
		}
		if tc.expectDonateTimesFromUserToAuthor != times {
			t.Errorf("%s: diff donate times, got %v, want %v", tc.testName, times, tc.expectDonateTimesFromUserToAuthor)
			return
		}

		cumulativeConsumption, err := gm.GetConsumption(ctx)
		if err != nil {
			t.Errorf("%s: failed to get consumption, got err %v", tc.testName, err)
		}
		if !tc.expectCumulativeConsumption.IsEqual(cumulativeConsumption) {
			t.Errorf("%s: diff cumulative consumption, got %v, want %v",
				tc.testName, cumulativeConsumption, tc.expectCumulativeConsumption)
			return
		}
	}
}

func TestHandlerRePostDonate(t *testing.T) {
	ctx, am, ph, pm, gm, dm, _, rm := setupTest(t, 1)
	postParam, _ := ph.GetPostParam(ctx)
	handler := NewHandler(pm, am, gm, dm, rm)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0.15")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")
	err := am.AddSavingCoin(
		ctx, user3, types.NewCoinFromInt64(123*types.Decimals),
		referrer, "", types.TransferIn)
	assert.Nil(t, err)
	// repost
	msg := CreatePostMsg{
		PostID:       "repost",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user2,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: user1,
		SourcePostID: postID,
		Links:        nil,
		RedistributionSplitRate: "0",
	}
	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(postParam.PostIntervalSec, 0)})
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	donateMsg := NewDonateMsg(
		string(user3), types.LNO("100"), string(user2), "repost", "", memo1)
	result = handler(ctx, donateMsg)
	assert.Equal(t, sdk.Result{}, result)
	eventList :=
		gm.GetTimeEventListAtTime(ctx, ctx.BlockHeader().Time.Unix()+3600*7*24)

	// after handler check KVStore
	// check repost first
	postInfo := model.PostInfo{
		PostID:       msg.PostID,
		Title:        msg.Title,
		Content:      msg.Content,
		Author:       msg.Author,
		ParentAuthor: msg.ParentAuthor,
		ParentPostID: msg.ParentPostID,
		SourceAuthor: msg.SourceAuthor,
		SourcePostID: msg.SourcePostID,
		Links:        msg.Links,
	}
	totalReward := types.RatToCoin(sdk.NewRat(15 * types.Decimals).Mul(sdk.NewRat(95, 100)))
	postMeta := model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time.Unix(),
		LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
		LastActivityAt:          ctx.BlockHeader().Time.Unix(),
		AllowReplies:            true,
		TotalDonateCount:        1,
		TotalReward:             totalReward,
		TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
		TotalReportCoinDay:      types.NewCoinFromInt64(0),
		RedistributionSplitRate: sdk.ZeroRat(),
	}
	checkPostKVStore(t, ctx, types.GetPermlink(user2, "repost"), postInfo, postMeta)
	repostRewardEvent := RewardEvent{
		PostAuthor: user2,
		PostID:     "repost",
		Consumer:   user3,
		Evaluate:   types.NewCoinFromInt64(13017),
		Original:   types.NewCoinFromInt64(15 * types.Decimals),
		Friction:   types.NewCoinFromInt64(75000),
		FromApp:    "",
	}
	assert.Equal(t, repostRewardEvent, eventList.Events[1])

	// check source post
	postMeta.TotalReward = types.RatToCoin(sdk.NewRat(85 * types.Decimals).Mul(sdk.NewRat(95, 100)))
	postMeta.CreatedAt = 0
	postMeta.LastUpdatedAt = 0
	postInfo.Author = user1
	postInfo.PostID = postID
	postInfo.SourceAuthor = ""
	postInfo.SourcePostID = ""
	postMeta.RedistributionSplitRate = sdk.NewRat(3, 20)
	postMeta.TotalUpvoteCoinDay = types.NewCoinFromInt64(0)

	checkPostKVStore(t, ctx, types.GetPermlink(user1, postID), postInfo, postMeta)

	acc1Saving, _ := am.GetSavingFromBank(ctx, user1)
	acc2Saving, _ := am.GetSavingFromBank(ctx, user2)
	acc3Saving, _ := am.GetSavingFromBank(ctx, user3)
	acc1SavingCoin := types.RatToCoin(sdk.NewRat(85 * types.Decimals).Mul(sdk.NewRat(95, 100)))
	acc2SavingCoin := types.RatToCoin(sdk.NewRat(15 * types.Decimals).Mul(sdk.NewRat(95, 100)))
	assert.Equal(t, acc1Saving, initCoin.Plus(acc1SavingCoin))
	assert.Equal(t, acc2Saving, initCoin.Plus(acc2SavingCoin))
	assert.Equal(t, acc3Saving, initCoin.Plus(types.NewCoinFromInt64(23*types.Decimals)))

	sourceRewardEvent := RewardEvent{
		PostAuthor: user1,
		PostID:     postID,
		Consumer:   user3,
		Evaluate:   types.NewCoinFromInt64(52141),
		Original:   types.NewCoinFromInt64(85 * types.Decimals),
		Friction:   types.NewCoinFromInt64(425000),
		FromApp:    "",
	}
	assert.Equal(t, sourceRewardEvent, eventList.Events[0])
}

// reputation check should be added later
func TestHandlerReportOrUpvote(t *testing.T) {
	ctx, am, ph, pm, gm, dm, _, rm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm, rm)
	coinDayParam, _ := ph.GetCoinDayParam(ctx)
	postParam, _ := ph.GetPostParam(ctx)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")
	user4 := createTestAccount(t, ctx, am, "user4")

	baseTime := ctx.BlockHeader().Time.Unix() + coinDayParam.SecondsToRecoverCoinDay
	invalidPermlink := types.GetPermlink("invalid", "invalid")

	testCases := []struct {
		testName             string
		reportOrUpvoteUser   string
		isReport             bool
		targetPostAuthor     string
		targetPostID         string
		lastReportOrUpvoteAt int64
		expectResult         sdk.Result
	}{
		{
			testName:             "user1 report",
			reportOrUpvoteUser:   string(user1),
			isReport:             true,
			targetPostAuthor:     string(user1),
			targetPostID:         postID,
			lastReportOrUpvoteAt: baseTime - postParam.ReportOrUpvoteIntervalSec,
			expectResult:         sdk.Result{},
		},
		{
			testName:             "user2 report",
			reportOrUpvoteUser:   string(user2),
			isReport:             true,
			targetPostAuthor:     string(user1),
			targetPostID:         postID,
			lastReportOrUpvoteAt: baseTime - postParam.ReportOrUpvoteIntervalSec,
			expectResult:         sdk.Result{},
		},
		{
			testName:             "user3 upvote",
			reportOrUpvoteUser:   string(user3),
			isReport:             false,
			targetPostAuthor:     string(user1),
			targetPostID:         postID,
			lastReportOrUpvoteAt: baseTime - postParam.ReportOrUpvoteIntervalSec,
			expectResult:         sdk.Result{},
		},
		{
			testName:             "user1 wanna change report to upvote",
			reportOrUpvoteUser:   string(user1),
			isReport:             false,
			targetPostAuthor:     string(user1),
			targetPostID:         postID,
			lastReportOrUpvoteAt: baseTime - postParam.ReportOrUpvoteIntervalSec,
			expectResult:         sdk.Result{},
		},
		{
			testName:             "user1 report too often",
			reportOrUpvoteUser:   string(user1),
			isReport:             false,
			targetPostAuthor:     string(user1),
			targetPostID:         postID,
			lastReportOrUpvoteAt: baseTime - postParam.ReportOrUpvoteIntervalSec + 1,
			expectResult:         ErrReportOrUpvoteTooOften().Result(),
		},
		{
			testName:             "user4 report to an invalid post",
			reportOrUpvoteUser:   string(user4),
			isReport:             true,
			targetPostAuthor:     "invalid",
			targetPostID:         "invalid",
			lastReportOrUpvoteAt: baseTime - postParam.ReportOrUpvoteIntervalSec,
			expectResult:         ErrPostNotFound(invalidPermlink).Result(),
		},
	}

	for _, tc := range testCases {
		lastReportOrUpvoteAtCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: time.Unix(tc.lastReportOrUpvoteAt, 0)})
		am.UpdateLastReportOrUpvoteAt(lastReportOrUpvoteAtCtx, types.AccountKey(tc.reportOrUpvoteUser))

		newCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: time.Unix(baseTime, 0)})
		msg := NewReportOrUpvoteMsg(tc.reportOrUpvoteUser, tc.targetPostAuthor, tc.targetPostID, tc.isReport)

		result := handler(newCtx, msg)
		if !assert.Equal(t, tc.expectResult, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectResult)
		}
		if tc.expectResult.Code != sdk.ABCICodeOK {
			continue
		}

		postMeta := model.PostMeta{
			CreatedAt:               ctx.BlockHeader().Time.Unix(),
			LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
			LastActivityAt:          newCtx.BlockHeader().Time.Unix(),
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat(),
			TotalReportCoinDay:      types.NewCoinFromInt64(0),
			TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
			TotalReward:             types.NewCoinFromInt64(0),
		}
		targetPost := types.GetPermlink(types.AccountKey(tc.targetPostAuthor), tc.targetPostID)
		checkPostMeta(t, ctx, targetPost, postMeta)

		lastReportOrUpvoteAt, _ := am.GetLastReportOrUpvoteAt(ctx, types.AccountKey(tc.reportOrUpvoteUser))
		// assert.Equal(t, baseTime, lastReportOrUpvoteAt)
		if baseTime != lastReportOrUpvoteAt {
			t.Errorf("%s: diff time, got %v, want %v", tc.testName, lastReportOrUpvoteAt, baseTime)
		}
	}
}

func TestHandlerView(t *testing.T) {
	ctx, am, _, pm, gm, dm, _, rm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm, dm, rm)

	createTime := ctx.BlockHeader().Time.Unix()
	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")
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
			testName:             "user3 views (postID, user1)",
			viewUser:             user3,
			postID:               postID,
			author:               user1,
			viewTime:             1,
			expectTotalViewCount: 1,
			expectUserViewCount:  1,
		},
		{
			testName:             "user3 views (postID, user1) again",
			viewUser:             user3,
			postID:               postID,
			author:               user1,
			viewTime:             2,
			expectTotalViewCount: 2,
			expectUserViewCount:  2,
		},
		{
			testName:             "user2 views (postID, user1)",
			viewUser:             user2,
			postID:               postID,
			author:               user1,
			viewTime:             3,
			expectTotalViewCount: 3,
			expectUserViewCount:  1,
		},
		{
			testName:             "user2 views (postID, user1) again",
			viewUser:             user2,
			postID:               postID,
			author:               user1,
			viewTime:             4,
			expectTotalViewCount: 4,
			expectUserViewCount:  2,
		},
		{
			testName:             "user1 views (postID, user1)",
			viewUser:             user1,
			postID:               postID,
			author:               user1,
			viewTime:             5,
			expectTotalViewCount: 5,
			expectUserViewCount:  1,
		},
	}

	for _, tc := range testCases {
		postKey := types.GetPermlink(tc.author, tc.postID)
		ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(tc.viewTime, 0)})
		msg := NewViewMsg(string(tc.viewUser), string(tc.author), tc.postID)
		result := handler(ctx, msg)
		if !assert.Equal(t, result, sdk.Result{}) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, sdk.Result{})
		}

		postMeta := model.PostMeta{
			CreatedAt:               createTime,
			LastUpdatedAt:           createTime,
			LastActivityAt:          createTime,
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
			t.Errorf("%s: diff view times, got %v, want %v", tc.testName, view.Times, tc.expectUserViewCount)
		}
		if view.LastViewAt != tc.viewTime {
			t.Errorf("%s: diff last view at, got %v, want %v", tc.testName, view.LastViewAt, tc.viewTime)
		}
	}
}
