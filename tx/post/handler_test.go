package post

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/post/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
)

func TestHandlerCreatePost(t *testing.T) {
	ctx, am, pm, gm := setupTest(t, 1)
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
		Links:        []types.IDToURLMapping{},
		RedistributionSplitRate: "0",
	}
	msg := NewCreatePostMsg(postCreateParams)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	assert.True(t, pm.IsPostExist(ctx, types.GetPostKey(postCreateParams.Author, postCreateParams.PostID)))

	// test invlaid author
	postCreateParams.Author = types.AccountKey("invalid")
	msg = NewCreatePostMsg(postCreateParams)
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrCreatePostAuthorNotFound(postCreateParams.Author).Result())
}

func TestHandlerCreateComment(t *testing.T) {
	ctx, am, pm, gm := setupTest(t, 1)
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
		Links:        []types.IDToURLMapping{},
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
		Created:                 ctx.BlockHeader().Time,
		LastUpdate:              ctx.BlockHeader().Time,
		LastActivity:            ctx.BlockHeader().Time,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat,
	}

	checkPostKVStore(t, ctx, types.GetPostKey(user, "comment"), postInfo, postMeta)

	// check parent
	postInfo.PostID = postID
	postInfo.ParentAuthor = ""
	postInfo.ParentPostID = ""
	postMeta.Created = ctx.BlockHeader().Time
	postMeta.LastUpdate = ctx.BlockHeader().Time
	checkPostKVStore(t, ctx, types.GetPostKey(user, postID), postInfo, postMeta)

	// test invalid parent
	postCreateParams.PostID = "invalid post"
	postCreateParams.ParentAuthor = user
	postCreateParams.ParentPostID = "invalid parent"
	msg = NewCreatePostMsg(postCreateParams)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrCommentInvalidParent(types.GetPostKey(user, postCreateParams.ParentPostID)).Result())

	// test duplicate comment
	postCreateParams.Author = user
	postCreateParams.PostID = "comment"
	postCreateParams.ParentAuthor = user
	postCreateParams.ParentPostID = "TestPostID"
	msg = NewCreatePostMsg(postCreateParams)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrCreateExistPost(types.GetPostKey(postCreateParams.Author, postCreateParams.PostID)).Result())

	// test cycle comment
	postCreateParams.Author = user
	postCreateParams.PostID = "newComment"
	postCreateParams.ParentAuthor = user
	postCreateParams.ParentPostID = "newComment"
	msg = NewCreatePostMsg(postCreateParams)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrCommentInvalidParent(types.GetPostKey(user, postCreateParams.PostID)).Result())
}

func TestHandlerRepost(t *testing.T) {
	ctx, am, pm, gm := setupTest(t, 1)
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
		Links:        []types.IDToURLMapping{},
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
		Created:                 ctx.BlockHeader().Time,
		LastUpdate:              ctx.BlockHeader().Time,
		LastActivity:            ctx.BlockHeader().Time,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat,
	}

	checkPostKVStore(t, ctx, types.GetPostKey(user, "repost"), postInfo, postMeta)

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
		Created:                 ctx.BlockHeader().Time,
		LastUpdate:              ctx.BlockHeader().Time,
		LastActivity:            ctx.BlockHeader().Time,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat,
	}
	postInfo.SourceAuthor = user
	postInfo.SourcePostID = postID
	checkPostKVStore(t, ctx, types.GetPostKey(user, postInfo.PostID), postInfo, postMeta)
}

func TestHandlerPostLike(t *testing.T) {
	ctx, am, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, "0")

	likeMsg := NewLikeMsg(types.AccountKey(user), 10000, user, postID)
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
		Links:        []types.IDToURLMapping{},
	}
	postMeta := model.PostMeta{
		Created:                 ctx.BlockHeader().Time,
		LastUpdate:              ctx.BlockHeader().Time,
		LastActivity:            ctx.BlockHeader().Time,
		AllowReplies:            true,
		TotalLikeCount:          1,
		TotalLikeWeight:         10000,
		RedistributionSplitRate: sdk.ZeroRat,
	}
	checkPostKVStore(t, ctx, types.GetPostKey(user, postID), postInfo, postMeta)

	// test update like
	likeMsg = NewLikeMsg(user, -10000, user, postID)
	result = handler(ctx, likeMsg)
	assert.Equal(t, result, sdk.Result{})
	postMeta.TotalDislikeWeight = 10000
	postMeta.TotalLikeWeight = 0
	checkPostKVStore(t, ctx, types.GetPostKey(user, postID), postInfo, postMeta)

	// test invalid like target post
	likeMsg = NewLikeMsg(user, -10000, user, "invalid")
	result = handler(ctx, likeMsg)
	assert.Equal(t, result, ErrLikeNonExistPost(types.GetPostKey(user, "invalid")).Result())
	checkPostKVStore(t, ctx, types.GetPostKey(user, postID), postInfo, postMeta)

	// test invalid like username
	likeMsg = NewLikeMsg(types.AccountKey("invalid"), 10000, user, postID)
	result = handler(ctx, likeMsg)

	assert.Equal(t, result, ErrLikePostUserNotFound(likeMsg.Username).Result())
	checkPostKVStore(t, ctx, types.GetPostKey(user, postID), postInfo, postMeta)
}

func TestHandlerPostDonate(t *testing.T) {
	ctx, am, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	err := am.AddCoin(ctx, user2, types.NewCoin(123*types.Decimals))
	assert.Nil(t, err)

	donateMsg := NewDonateMsg(user2, types.LNO("100"), user1, postID, "")
	result := handler(ctx, donateMsg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	postInfo := model.PostInfo{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user1,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []types.IDToURLMapping{},
	}
	postMeta := model.PostMeta{
		Created:                 ctx.BlockHeader().Time,
		LastUpdate:              ctx.BlockHeader().Time,
		LastActivity:            ctx.BlockHeader().Time,
		AllowReplies:            true,
		TotalDonateCount:        1,
		TotalReward:             types.NewCoin(95 * types.Decimals),
		RedistributionSplitRate: sdk.ZeroRat,
	}

	checkPostKVStore(t, ctx, types.GetPostKey(user1, postID), postInfo, postMeta)

	acc1Balance, _ := am.GetBankBalance(ctx, user1)
	acc2Balance, _ := am.GetBankBalance(ctx, user2)

	assert.Equal(t, acc1Balance, initCoin.Plus(types.NewCoin(95*types.Decimals)))
	assert.Equal(t, acc2Balance, initCoin.Plus(types.NewCoin(23*types.Decimals)))
	// test invalid donation target
	donateMsg = NewDonateMsg(user1, types.LNO("100"), user1, "invalid", "")
	result = handler(ctx, donateMsg)
	assert.Equal(t, result, ErrDonatePostDoesntExist(types.GetPostKey(user1, "invalid")).Result())
	checkPostKVStore(t, ctx, types.GetPostKey(user1, postID), postInfo, postMeta)

	// test invalid user1name
	donateMsg = NewDonateMsg(types.AccountKey("invalid"), types.LNO("100"), user1, postID, "")
	result = handler(ctx, donateMsg)

	assert.Equal(t, result, ErrDonateUserNotFound(types.AccountKey("invalid")).Result())
	checkPostKVStore(t, ctx, types.GetPostKey(user1, postID), postInfo, postMeta)

	// test insufficient deposit
	donateMsg = NewDonateMsg(user2, types.LNO("100"), user1, postID, "")
	result = handler(ctx, donateMsg)

	assert.Equal(t, result, ErrDonateFailed(types.GetPostKey(user1, postID)).Result())
}

func TestHandlerRePostDonate(t *testing.T) {
	ctx, am, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0.15")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")
	err := am.AddCoin(ctx, user3, types.NewCoin(123*types.Decimals))
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
		Links:        []types.IDToURLMapping{},
		RedistributionSplitRate: "0",
	}
	msg := NewCreatePostMsg(postCreateParams)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	donateMsg := NewDonateMsg(types.AccountKey(user3), types.LNO("100"), user2, "repost", "")
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
		Created:                 ctx.BlockHeader().Time,
		LastUpdate:              ctx.BlockHeader().Time,
		LastActivity:            ctx.BlockHeader().Time,
		AllowReplies:            true,
		TotalDonateCount:        1,
		TotalReward:             types.RatToCoin(sdk.NewRat(15 * types.Decimals).Mul(sdk.NewRat(95, 100))),
		RedistributionSplitRate: sdk.ZeroRat,
	}

	checkPostKVStore(t, ctx, types.GetPostKey(user2, "repost"), postInfo, postMeta)

	// check source post
	postMeta.TotalReward = types.Coin{sdk.NewRat(85 * types.Decimals).Mul(sdk.NewRat(95, 100)).Evaluate()}
	postInfo.Author = user1
	postInfo.PostID = postID
	postInfo.SourceAuthor = ""
	postInfo.SourcePostID = ""
	postMeta.RedistributionSplitRate = sdk.NewRat(15, 100)

	checkPostKVStore(t, ctx, types.GetPostKey(user1, postID), postInfo, postMeta)

	acc1Balance, _ := am.GetBankBalance(ctx, user1)
	acc2Balance, _ := am.GetBankBalance(ctx, user2)
	acc3Balance, _ := am.GetBankBalance(ctx, user3)
	assert.Equal(t, acc1Balance, initCoin.Plus(types.RatToCoin(sdk.NewRat(85*types.Decimals).Mul(sdk.NewRat(95, 100)))))
	assert.Equal(t, acc2Balance, initCoin.Plus(types.RatToCoin(sdk.NewRat(15*types.Decimals).Mul(sdk.NewRat(95, 100)))))
	assert.Equal(t, acc3Balance, initCoin.Plus(types.NewCoin(23*types.Decimals)))
}

func TestHandlerReportOrUpvote(t *testing.T) {
	ctx, am, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")

	cases := []struct {
		reportOrUpvoteUser     types.AccountKey
		isReport               bool
		isRevoke               bool
		expectTotalReportStake types.Coin
		expectTotalUpvoteStake types.Coin
	}{
		{user2, true, false, types.NewCoin(100), types.NewCoin(0)},
		{user3, true, false, types.NewCoin(200), types.NewCoin(0)},
		{user2, false, false, types.NewCoin(100), types.NewCoin(100)},
		{user3, false, false, types.NewCoin(0), types.NewCoin(200)},
		{user2, false, true, types.NewCoin(0), types.NewCoin(100)},
		{user3, false, true, types.NewCoin(0), types.NewCoin(0)},
		{user2, true, false, types.NewCoin(100), types.NewCoin(0)},
		{user3, true, false, types.NewCoin(200), types.NewCoin(0)},
		{user2, false, false, types.NewCoin(100), types.NewCoin(100)},
		{user3, false, false, types.NewCoin(0), types.NewCoin(200)},
	}

	for _, cs := range cases {
		newCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: ctx.BlockHeader().Time + acc.TotalCoinDaysSec})
		msg := NewReportOrUpvoteMsg(cs.reportOrUpvoteUser, user1, postID, cs.isReport, cs.isRevoke)
		result := handler(newCtx, msg)
		assert.Equal(t, result, sdk.Result{})
		postMeta := model.PostMeta{
			Created:                 ctx.BlockHeader().Time,
			LastUpdate:              ctx.BlockHeader().Time,
			LastActivity:            newCtx.BlockHeader().Time,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat,
			TotalReportStake:        cs.expectTotalReportStake,
			TotalUpvoteStake:        cs.expectTotalUpvoteStake,
		}
		postKey := types.GetPostKey(user1, postID)
		checkPostMeta(t, ctx, postKey, postMeta)
	}
}

func TestHandlerRepostReportOrUpvote(t *testing.T) {
	ctx, am, pm, gm := setupTest(t, 1)
	handler := NewHandler(pm, am, gm)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, "0")
	user2 := createTestAccount(t, ctx, am, "user2")
	user3 := createTestAccount(t, ctx, am, "user3")

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
		Links:        []types.IDToURLMapping{},
		RedistributionSplitRate: "0",
	}
	msg := NewCreatePostMsg(postCreateParams)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	cases := []struct {
		reportOrUpvoteUser      types.AccountKey
		isReport                bool
		isRevoke                bool
		expectSourceReportStake types.Coin
		expectSourceUpvoteStake types.Coin
	}{
		{user2, true, false, types.NewCoin(100), types.NewCoin(0)},
		{user3, true, false, types.NewCoin(200), types.NewCoin(0)},
		{user2, false, false, types.NewCoin(100), types.NewCoin(100)},
		{user3, false, false, types.NewCoin(0), types.NewCoin(200)},
		{user2, false, true, types.NewCoin(0), types.NewCoin(100)},
		{user3, false, true, types.NewCoin(0), types.NewCoin(0)},
		{user2, true, false, types.NewCoin(100), types.NewCoin(0)},
		{user3, true, false, types.NewCoin(200), types.NewCoin(0)},
		{user2, false, false, types.NewCoin(100), types.NewCoin(100)},
		{user3, false, false, types.NewCoin(0), types.NewCoin(200)},
	}

	for _, cs := range cases {
		newCtx := ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: ctx.BlockHeader().Time + acc.TotalCoinDaysSec})
		msg := NewReportOrUpvoteMsg(cs.reportOrUpvoteUser, user2, repostID, cs.isReport, cs.isRevoke)
		result := handler(newCtx, msg)
		assert.Equal(t, result, sdk.Result{})
		postMeta := model.PostMeta{
			Created:                 ctx.BlockHeader().Time,
			LastUpdate:              ctx.BlockHeader().Time,
			LastActivity:            newCtx.BlockHeader().Time,
			AllowReplies:            true,
			RedistributionSplitRate: sdk.ZeroRat,
			TotalReportStake:        cs.expectSourceReportStake,
			TotalUpvoteStake:        cs.expectSourceUpvoteStake,
		}
		postKey := types.GetPostKey(user1, postID)
		checkPostMeta(t, ctx, postKey, postMeta)
	}
}
