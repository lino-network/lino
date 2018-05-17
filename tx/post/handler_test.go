package post

import (
	"fmt"
	"testing"
	"time"

	"github.com/lino-network/lino/tx/post/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
)

func TestHandlerCreatePost(t *testing.T) {
	ctx, am, _, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user := createTestAccount(t, ctx, am, "user1")

	// test valid post
	postCreateParams := PostCreateParams{
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
	msg := NewCreatePostMsg(postCreateParams)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	assert.True(t, pm.IsPostExist(ctx, types.GetPermLink(postCreateParams.Author, postCreateParams.PostID)))

	// test invlaid author
	postCreateParams.Author = types.AccountKey("invalid")
	msg = NewCreatePostMsg(postCreateParams)
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrCreatePostAuthorNotFound(postCreateParams.Author).Result())
}

func TestHandlerUpdatePost(t *testing.T) {
	ctx, am, _, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")
	user1 := createTestAccount(t, ctx, am, "user1")

	cases := []struct {
		TestName     string
		msg          UpdatePostMsg
		expectResult sdk.Result
	}{
		{"normal update",
			NewUpdatePostMsg(string(user), postID, "update title", "update content", []types.IDToURLMapping(nil), "1"),
			sdk.Result{},
		},
		{"update user doesn't exist",
			NewUpdatePostMsg("invalid", postID, "update title", "update content", []types.IDToURLMapping(nil), "1"),
			ErrUpdatePostAuthorNotFound("invalid").Result(),
		},
		{"update post doesn't exist, post ID invalid",
			NewUpdatePostMsg(string(user), "invalid", "update title", "update content", []types.IDToURLMapping(nil), "1"),
			ErrUpdatePostNotFound(types.GetPermLink(user, "invalid")).Result(),
		},
		{"update post doesn't exist, author invalid",
			NewUpdatePostMsg(string(user1), postID, "update title", "update content", []types.IDToURLMapping(nil), "1"),
			ErrUpdatePostNotFound(types.GetPermLink(user1, postID)).Result(),
		},
	}
	for _, cs := range cases {
		splitRate, err := sdk.NewRatFromDecimal(cs.msg.RedistributionSplitRate)
		assert.Nil(t, err)
		result := handler(ctx, cs.msg)
		assert.Equal(t, cs.expectResult, result)
		if cs.expectResult.Code != sdk.CodeOK {
			continue
		}
		postInfo := model.PostInfo{
			PostID:       cs.msg.PostID,
			Title:        cs.msg.Title,
			Content:      cs.msg.Content,
			Author:       cs.msg.Author,
			SourceAuthor: "",
			SourcePostID: "",
			Links:        cs.msg.Links,
		}

		postMeta := model.PostMeta{
			Created:                 ctx.BlockHeader().Time,
			LastUpdate:              ctx.BlockHeader().Time,
			LastActivity:            ctx.BlockHeader().Time,
			AllowReplies:            true,
			IsDeleted:               false,
			RedistributionSplitRate: splitRate,
		}
		checkPostKVStore(t, ctx,
			types.GetPermLink(cs.msg.Author, cs.msg.PostID), postInfo, postMeta)
	}
}

func TestHandlerCreateComment(t *testing.T) {
	ctx, am, _, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")

	// test comment
	postCreateParams := PostCreateParams{
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
	msg := NewCreatePostMsg(postCreateParams)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	postInfo := model.PostInfo{
		PostID:       postCreateParams.PostID,
		Title:        postCreateParams.Title,
		Content:      postCreateParams.Content,
		Author:       postCreateParams.Author,
		ParentAuthor: postCreateParams.ParentAuthor,
		ParentPostID: postCreateParams.ParentPostID,
		SourceAuthor: postCreateParams.SourceAuthor,
		SourcePostID: postCreateParams.SourcePostID,
		Links:        postCreateParams.Links,
	}

	postMeta := model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time,
		LastUpdatedAt:           ctx.BlockHeader().Time,
		LastActivityAt:          ctx.BlockHeader().Time,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat,
	}

	checkPostKVStore(t, ctx, types.GetPermLink(user, "comment"), postInfo, postMeta)

	// check parent
	postInfo.PostID = postID
	postInfo.ParentAuthor = ""
	postInfo.ParentPostID = ""
	postMeta.CreatedAt = ctx.BlockHeader().Time
	postMeta.LastUpdatedAt = ctx.BlockHeader().Time
	checkPostKVStore(t, ctx, types.GetPermLink(user, postID), postInfo, postMeta)

	// test invalid parent
	postCreateParams.PostID = "invalid post"
	postCreateParams.ParentAuthor = user
	postCreateParams.ParentPostID = "invalid parent"
	msg = NewCreatePostMsg(postCreateParams)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrCommentInvalidParent(types.GetPermLink(user, postCreateParams.ParentPostID)).Result())

	// test duplicate comment
	postCreateParams.Author = user
	postCreateParams.PostID = "comment"
	postCreateParams.ParentAuthor = user
	postCreateParams.ParentPostID = "TestPostID"
	msg = NewCreatePostMsg(postCreateParams)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrCreateExistPost(types.GetPermLink(postCreateParams.Author, postCreateParams.PostID)).Result())

	// test cycle comment
	postCreateParams.Author = user
	postCreateParams.PostID = "newComment"
	postCreateParams.ParentAuthor = user
	postCreateParams.ParentPostID = "newComment"
	msg = NewCreatePostMsg(postCreateParams)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrCommentInvalidParent(types.GetPermLink(user, postCreateParams.PostID)).Result())
}

func TestHandlerRepost(t *testing.T) {
	ctx, am, _, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")

	// test repost
	postCreateParams := PostCreateParams{
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
	msg := NewCreatePostMsg(postCreateParams)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	postInfo := model.PostInfo{
		PostID:       postCreateParams.PostID,
		Title:        postCreateParams.Title,
		Content:      postCreateParams.Content,
		Author:       postCreateParams.Author,
		ParentAuthor: postCreateParams.ParentAuthor,
		ParentPostID: postCreateParams.ParentPostID,
		SourceAuthor: postCreateParams.SourceAuthor,
		SourcePostID: postCreateParams.SourcePostID,
		Links:        postCreateParams.Links,
	}

	postMeta := model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time,
		LastUpdatedAt:           ctx.BlockHeader().Time,
		LastActivityAt:          ctx.BlockHeader().Time,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat,
	}

	checkPostKVStore(t, ctx, types.GetPermLink(user, "repost"), postInfo, postMeta)

	// test 2 depth repost
	postCreateParams.PostID = "repost-repost"
	postCreateParams.SourceAuthor = user
	postCreateParams.SourcePostID = "repost"
	msg = NewCreatePostMsg(postCreateParams)
	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: time.Now().Unix()})
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	// check 2 depth repost
	postInfo.PostID = "repost-repost"
	postMeta = model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time,
		LastUpdatedAt:           ctx.BlockHeader().Time,
		LastActivityAt:          ctx.BlockHeader().Time,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat,
	}
	postInfo.SourceAuthor = user
	postInfo.SourcePostID = postID
	checkPostKVStore(t, ctx, types.GetPermLink(user, postInfo.PostID), postInfo, postMeta)
}

func TestHandlerPostLike(t *testing.T) {
	ctx, am, _, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")

	likeMsg := NewLikeMsg(string(user), 10000, string(user), postID)
	result := handler(ctx, likeMsg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	postInfo := model.PostInfo{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        nil,
	}
	postMeta := model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time,
		LastUpdatedAt:           ctx.BlockHeader().Time,
		LastActivityAt:          ctx.BlockHeader().Time,
		AllowReplies:            true,
		TotalLikeCount:          1,
		TotalLikeWeight:         10000,
		RedistributionSplitRate: sdk.ZeroRat,
	}
	checkPostKVStore(t, ctx, types.GetPermLink(user, postID), postInfo, postMeta)

	// test update like
	likeMsg = NewLikeMsg(string(user), -10000, string(user), postID)
	result = handler(ctx, likeMsg)
	assert.Equal(t, result, sdk.Result{})
	postMeta.TotalDislikeWeight = 10000
	postMeta.TotalLikeWeight = 0
	checkPostKVStore(t, ctx, types.GetPermLink(user, postID), postInfo, postMeta)

	// test invalid like target post
	likeMsg = NewLikeMsg(string(user), -10000, string(user), "invalid")
	result = handler(ctx, likeMsg)
	assert.Equal(t, result, ErrLikeNonExistPost(types.GetPermLink(user, "invalid")).Result())
	checkPostKVStore(t, ctx, types.GetPermLink(user, postID), postInfo, postMeta)

	// test invalid like username
	likeMsg = NewLikeMsg("invalid", 10000, string(user), postID)
	result = handler(ctx, likeMsg)

	assert.Equal(t, result, ErrLikePostUserNotFound(likeMsg.Username).Result())
	checkPostKVStore(t, ctx, types.GetPermLink(user, postID), postInfo, postMeta)
}

func TestHandlerPostDonate(t *testing.T) {
	ctx, am, ph, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	accParam, err := ph.GetAccountParam(ctx)
	assert.Nil(t, err)

	author, postID := createTestPost(t, ctx, "author", "postID", am, pm, "0")

	postInfo := model.PostInfo{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       author,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        nil,
	}

	userWithSufficientSaving := createTestAccount(t, ctx, am, "userWithSufficientSaving")
	userWithSufficientChecking := createTestAccount(t, ctx, am, "userWithSufficientChecking")
	err = am.AddSavingCoin(ctx, userWithSufficientSaving, types.NewCoin(100*types.Decimals))
	assert.Nil(t, err)
	err = am.AddCheckingCoin(ctx, userWithSufficientChecking, types.NewCoin(100*types.Decimals))
	assert.Nil(t, err)

	cases := []struct {
		TestName              string
		DonateUesr            types.AccountKey
		Amount                types.LNO
		ToAuthor              types.AccountKey
		ToPostID              string
		FromChecking          bool
		ExpectErr             sdk.Result
		ExpectPostMeta        model.PostMeta
		ExpectDonatorSaving   types.Coin
		ExpectDonatorChecking types.Coin
		ExpectAuthorSaving    types.Coin
		ExpectAuthorChecking  types.Coin
	}{
		{"donate from sufficient saving",
			userWithSufficientSaving, types.LNO("100"), author, postID, false, sdk.Result{},
			model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time,
				LastUpdatedAt:           ctx.BlockHeader().Time,
				LastActivityAt:          ctx.BlockHeader().Time,
				AllowReplies:            true,
				TotalDonateCount:        1,
				TotalReward:             types.NewCoin(95 * types.Decimals),
				RedistributionSplitRate: sdk.ZeroRat,
			},
			accParam.RegisterFee, types.NewCoin(0),
			accParam.RegisterFee.Plus(types.NewCoin(95 * types.Decimals)), types.NewCoin(0),
		},
		{"donate from sufficient checking",
			userWithSufficientChecking, types.LNO("100"), author, postID, true, sdk.Result{},
			model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time,
				LastUpdatedAt:           ctx.BlockHeader().Time,
				LastActivityAt:          ctx.BlockHeader().Time,
				AllowReplies:            true,
				TotalDonateCount:        2,
				TotalReward:             types.NewCoin(190 * types.Decimals),
				RedistributionSplitRate: sdk.ZeroRat,
			},
			accParam.RegisterFee, types.NewCoin(0),
			accParam.RegisterFee.Plus(types.NewCoin(190 * types.Decimals)), types.NewCoin(0),
		},
		{"donate from insufficient saving",
			userWithSufficientSaving, types.LNO("100"), author, postID, false,
			ErrAccountSavingCoinNotEnough(types.GetPermLink(author, postID)).Result(),
			model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time,
				LastUpdatedAt:           ctx.BlockHeader().Time,
				LastActivityAt:          ctx.BlockHeader().Time,
				AllowReplies:            true,
				TotalDonateCount:        2,
				TotalReward:             types.NewCoin(190 * types.Decimals),
				RedistributionSplitRate: sdk.ZeroRat,
			},
			accParam.RegisterFee, types.NewCoin(0),
			accParam.RegisterFee.Plus(types.NewCoin(190 * types.Decimals)), types.NewCoin(0),
		},
		{"donate from insufficient checking",
			userWithSufficientChecking, types.LNO("100"), author, postID, true,
			ErrAccountCheckingCoinNotEnough(types.GetPermLink(author, postID)).Result(),
			model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time,
				LastUpdatedAt:           ctx.BlockHeader().Time,
				LastActivityAt:          ctx.BlockHeader().Time,
				AllowReplies:            true,
				TotalDonateCount:        2,
				TotalReward:             types.NewCoin(190 * types.Decimals),
				RedistributionSplitRate: sdk.ZeroRat,
			},
			accParam.RegisterFee, types.NewCoin(0),
			accParam.RegisterFee.Plus(types.NewCoin(190 * types.Decimals)), types.NewCoin(0),
		},
		{"invalid target postID",
			userWithSufficientChecking, types.LNO("100"), author, "invalid", true,
			ErrDonatePostNotFound(types.GetPermLink(author, "invalid")).Result(),
			model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time,
				LastUpdatedAt:           ctx.BlockHeader().Time,
				LastActivityAt:          ctx.BlockHeader().Time,
				AllowReplies:            true,
				TotalDonateCount:        2,
				TotalReward:             types.NewCoin(190 * types.Decimals),
				RedistributionSplitRate: sdk.ZeroRat,
			},
			accParam.RegisterFee, types.NewCoin(0),
			accParam.RegisterFee.Plus(types.NewCoin(190 * types.Decimals)), types.NewCoin(0),
		},
		{"invalid target author",
			userWithSufficientChecking, types.LNO("100"), types.AccountKey("invalid"), postID, true,
			ErrDonatePostNotFound(types.GetPermLink(types.AccountKey("invalid"), postID)).Result(),
			model.PostMeta{
				CreatedAt:               ctx.BlockHeader().Time,
				LastUpdatedAt:           ctx.BlockHeader().Time,
				LastActivityAt:          ctx.BlockHeader().Time,
				AllowReplies:            true,
				TotalDonateCount:        2,
				TotalReward:             types.NewCoin(190 * types.Decimals),
				RedistributionSplitRate: sdk.ZeroRat,
			},
			accParam.RegisterFee, types.NewCoin(0),
			accParam.RegisterFee.Plus(types.NewCoin(190 * types.Decimals)), types.NewCoin(0),
		},
	}

	for _, cs := range cases {
		donateMsg := NewDonateMsg(
			string(cs.DonateUesr), cs.Amount, string(cs.ToAuthor), cs.ToPostID, "", cs.FromChecking, memo1)
		result := handler(ctx, donateMsg)
		assert.Equal(t, cs.ExpectErr, result)
		if cs.ExpectErr.Code == sdk.CodeOK {
			checkPostKVStore(t, ctx, types.GetPermLink(cs.ToAuthor, cs.ToPostID), postInfo, cs.ExpectPostMeta)
		}
		authorSaving, err := am.GetSavingFromBank(ctx, author)
		assert.Nil(t, err)
		if !authorSaving.IsEqual(cs.ExpectAuthorSaving) {
			t.Errorf(
				"%s: expect author saving %v, got %v",
				cs.TestName, cs.ExpectAuthorSaving, authorSaving)
			return
		}
		donatorSaving, err := am.GetSavingFromBank(ctx, cs.DonateUesr)
		assert.Nil(t, err)
		if !donatorSaving.IsEqual(cs.ExpectDonatorSaving) {
			t.Errorf(
				"%s: expect donator saving %v, got %v",
				cs.TestName, cs.ExpectDonatorSaving, donatorSaving)
			return
		}
		authorChecking, err := am.GetCheckingFromBank(ctx, author)
		assert.Nil(t, err)
		if !authorChecking.IsEqual(cs.ExpectAuthorChecking) {
			t.Errorf(
				"%s: expect author checking %v, got %v",
				cs.TestName, cs.ExpectAuthorChecking, authorChecking)
			return
		}
		donatorChecking, err := am.GetCheckingFromBank(ctx, cs.DonateUesr)
		assert.Nil(t, err)
		if !donatorChecking.IsEqual(cs.ExpectDonatorChecking) {
			t.Errorf(
				"%s: expect donator checking %v, got %v",
				cs.TestName, cs.ExpectDonatorChecking, donatorChecking)
			return
		}
	}
}

func TestHandlerRePostDonate(t *testing.T) {
	ctx, am, _, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0.15")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")
	err := am.AddSavingCoin(ctx, user3, types.NewCoin(123*types.Decimals))
	assert.Nil(t, err)
	// repost
	postCreateParams := PostCreateParams{
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
	msg := NewCreatePostMsg(postCreateParams)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	donateMsg := NewDonateMsg(
		string(user3), types.LNO("100"), string(user2), "repost", "", false, memo1)
	result = handler(ctx, donateMsg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	// check repost first
	postInfo := model.PostInfo{
		PostID:       postCreateParams.PostID,
		Title:        postCreateParams.Title,
		Content:      postCreateParams.Content,
		Author:       postCreateParams.Author,
		ParentAuthor: postCreateParams.ParentAuthor,
		ParentPostID: postCreateParams.ParentPostID,
		SourceAuthor: postCreateParams.SourceAuthor,
		SourcePostID: postCreateParams.SourcePostID,
		Links:        postCreateParams.Links,
	}

	postMeta := model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time,
		LastUpdatedAt:           ctx.BlockHeader().Time,
		LastActivityAt:          ctx.BlockHeader().Time,
		AllowReplies:            true,
		TotalDonateCount:        1,
		TotalReward:             types.RatToCoin(sdk.NewRat(15 * types.Decimals).Mul(sdk.NewRat(95, 100))),
		RedistributionSplitRate: sdk.ZeroRat,
	}

	checkPostKVStore(t, ctx, types.GetPermLink(user2, "repost"), postInfo, postMeta)

	// check source post
	postMeta.TotalReward = types.Coin{sdk.NewRat(85 * types.Decimals).Mul(sdk.NewRat(95, 100)).Evaluate()}
	postInfo.Author = user1
	postInfo.PostID = postID
	postInfo.SourceAuthor = ""
	postInfo.SourcePostID = ""
	postMeta.RedistributionSplitRate = sdk.NewRat(15, 100)

	checkPostKVStore(t, ctx, types.GetPermLink(user1, postID), postInfo, postMeta)

	acc1Saving, _ := am.GetSavingFromBank(ctx, user1)
	acc2Saving, _ := am.GetSavingFromBank(ctx, user2)
	acc3Saving, _ := am.GetSavingFromBank(ctx, user3)
	assert.Equal(t, acc1Saving, initCoin.Plus(types.RatToCoin(sdk.NewRat(85*types.Decimals).Mul(sdk.NewRat(95, 100)))))
	assert.Equal(t, acc2Saving, initCoin.Plus(types.RatToCoin(sdk.NewRat(15*types.Decimals).Mul(sdk.NewRat(95, 100)))))
	assert.Equal(t, acc3Saving, initCoin.Plus(types.NewCoin(23*types.Decimals)))
}

func TestHandlerReportOrUpvote(t *testing.T) {
	ctx, am, ph, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)
	coinDayParam, _ := ph.GetCoinDayParam(ctx)
	accParam, _ := ph.GetAccountParam(ctx)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")

	testCases := []struct {
		testName               string
		reportOrUpvoteUser     string
		isReport               bool
		expectTotalReportStake types.Coin
		expectTotalUpvoteStake types.Coin
	}{
		{"user1 report", string(user1), true, accParam.RegisterFee, types.NewCoin(0)},
		{"user2 report", string(user2), true,
			accParam.RegisterFee.Plus(accParam.RegisterFee), types.NewCoin(0)},
		{"user3 upvote", string(user3), false,
			accParam.RegisterFee.Plus(accParam.RegisterFee), accParam.RegisterFee},
	}

	for _, tc := range testCases {
		newCtx := ctx.WithBlockHeader(
			abci.Header{ChainID: "Lino", Time: ctx.BlockHeader().Time + coinDayParam.SecondsToRecoverCoinDayStake})
		msg := NewReportOrUpvoteMsg(tc.reportOrUpvoteUser, string(user1), postID, tc.isReport)
		result := handler(newCtx, msg)
		assert.Equal(t, result, sdk.Result{}, fmt.Sprintf("%s: got %v, want %v", tc.testName, result, sdk.Result{}))
		postMeta := model.PostMeta{
			CreatedAt:               ctx.BlockHeader().Time,
			LastUpdatedAt:           ctx.BlockHeader().Time,
			LastActivityAt:          newCtx.BlockHeader().Time,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat,
			TotalReportStake:        tc.expectTotalReportStake,
			TotalUpvoteStake:        tc.expectTotalUpvoteStake,
		}
		postKey := types.GetPermLink(user1, postID)
		checkPostMeta(t, ctx, postKey, postMeta)
	}
}

func TestHandlerView(t *testing.T) {
	ctx, am, _, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	createTime := ctx.BlockHeader().Time
	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")
	cases := []struct {
		viewUser             types.AccountKey
		postID               string
		author               types.AccountKey
		viewTime             int64
		expectTotalViewCount int64
		expectUserViewCount  int64
	}{
		{user3, postID, user1, 1, 1, 1},
		{user3, postID, user1, 2, 2, 2},
		{user2, postID, user1, 3, 3, 1},
		{user2, postID, user1, 4, 4, 2},
		{user1, postID, user1, 5, 5, 1},
	}

	for _, cs := range cases {
		postKey := types.GetPermLink(cs.author, cs.postID)
		ctx = ctx.WithBlockHeader(abci.Header{Time: cs.viewTime})
		msg := NewViewMsg(string(cs.viewUser), string(cs.author), cs.postID)
		result := handler(ctx, msg)
		assert.Equal(t, result, sdk.Result{})
		postMeta := model.PostMeta{
			CreatedAt:               createTime,
			LastUpdatedAt:           createTime,
			LastActivityAt:          createTime,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat,
			TotalViewCount:          cs.expectTotalViewCount,
		}
		checkPostMeta(t, ctx, postKey, postMeta)
		view, err := pm.postStorage.GetPostView(ctx, postKey, cs.viewUser)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectUserViewCount, view.Times)
		assert.Equal(t, cs.viewTime, view.LastViewAt)
	}
}

func TestHandlerRepostReportOrUpvote(t *testing.T) {
	ctx, am, ph, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")

	accParam, _ := ph.GetAccountParam(ctx)

	// repost
	repostID := "repost"
	postCreateParams := PostCreateParams{
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
	msg := NewCreatePostMsg(postCreateParams)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	cases := []struct {
		reportOrUpvoteUser      types.AccountKey
		isReport                bool
		expectSourceReportStake types.Coin
		expectSourceUpvoteStake types.Coin
	}{
		{user2, true, accParam.RegisterFee, types.NewCoin(0)},
		{user3, false, accParam.RegisterFee, accParam.RegisterFee},
	}

	for _, cs := range cases {
		newCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: ctx.BlockHeader().Time + +7*3600*24})
		msg := NewReportOrUpvoteMsg(string(cs.reportOrUpvoteUser), string(user2), repostID, cs.isReport)
		result := handler(newCtx, msg)
		assert.Equal(t, result, sdk.Result{})
		postMeta := model.PostMeta{
			CreatedAt:               ctx.BlockHeader().Time,
			LastUpdatedAt:           ctx.BlockHeader().Time,
			LastActivityAt:          newCtx.BlockHeader().Time,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat,
			TotalReportStake:        cs.expectSourceReportStake,
			TotalUpvoteStake:        cs.expectSourceUpvoteStake,
		}
		postKey := types.GetPermLink(user1, postID)
		checkPostMeta(t, ctx, postKey, postMeta)
	}
}
