package post

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/post/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestHandlerCreatePost(t *testing.T) {
	ctx, am, pm, gm := setupTest(t, 1)
	handler := NewHandler(*pm, *am, *gm)

	user := createTestAccount(ctx, am, "user1")

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
		RedistributionSplitRate: sdk.ZeroRat,
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
	handler := NewHandler(*pm, *am, *gm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, sdk.ZeroRat)

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
		RedistributionSplitRate: sdk.ZeroRat,
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
		Created:                 1,
		LastUpdate:              1,
		LastActivity:            1,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat,
	}

	checkPostKVStore(t, ctx, types.GetPostKey(user, "comment"), postInfo, postMeta)

	// check parent
	postInfo.PostID = postID
	postInfo.ParentAuthor = ""
	postInfo.ParentPostID = ""
	postMeta.Created = 1
	postMeta.LastUpdate = 1
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
	handler := NewHandler(*pm, *am, *gm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, sdk.ZeroRat)

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
		RedistributionSplitRate: sdk.ZeroRat,
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
		Created:                 1,
		LastUpdate:              1,
		LastActivity:            1,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat,
	}

	checkPostKVStore(t, ctx, types.GetPostKey(user, "repost"), postInfo, postMeta)

	// test 2 depth repost
	postCreateParams.PostID = "repost-repost"
	postCreateParams.SourceAuthor = user
	postCreateParams.SourcePostID = "repost"
	msg = NewCreatePostMsg(postCreateParams)
	ctx = ctx.WithBlockHeight(2)
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	// check 2 depth repost
	postInfo.PostID = "repost-repost"
	postMeta = model.PostMeta{
		Created:                 2,
		LastUpdate:              2,
		LastActivity:            2,
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroRat,
	}
	postInfo.SourceAuthor = user
	postInfo.SourcePostID = postID
	checkPostKVStore(t, ctx, types.GetPostKey(user, postInfo.PostID), postInfo, postMeta)
}

func TestHandlerPostLike(t *testing.T) {
	ctx, am, pm, gm := setupTest(t, 1)
	handler := NewHandler(*pm, *am, *gm)

	user, postID := createTestPost(t, ctx, "user", "postID", am, pm, sdk.ZeroRat)

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
		Created:                 1,
		LastUpdate:              1,
		LastActivity:            1,
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
	postMeta.TotalLikeWeight = -10000
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
	handler := NewHandler(*pm, *am, *gm)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, sdk.ZeroRat)
	user2 := createTestAccount(ctx, am, "user2")
	err := am.AddCoin(ctx, user2, types.NewCoin(123*types.Decimals))
	assert.Nil(t, err)

	donateMsg := NewDonateMsg(user2, types.LNO(sdk.NewRat(100)), user1, postID)
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
		Created:                 1,
		LastUpdate:              1,
		LastActivity:            1,
		AllowReplies:            true,
		TotalDonateCount:        1,
		TotalReward:             types.Coin{99 * types.Decimals},
		RedistributionSplitRate: sdk.ZeroRat,
	}

	checkPostKVStore(t, ctx, types.GetPostKey(user1, postID), postInfo, postMeta)

	acc1Balance, _ := am.GetBankBalance(ctx, user1)
	acc2Balance, _ := am.GetBankBalance(ctx, user2)

	assert.Equal(t, true, acc1Balance.IsEqual(types.Coin{99 * types.Decimals}))
	assert.Equal(t, true, acc2Balance.IsEqual(types.Coin{23 * types.Decimals}))
	// test invalid donation target
	donateMsg = NewDonateMsg(user1, types.LNO(sdk.NewRat(100)), user1, "invalid")
	result = handler(ctx, donateMsg)
	assert.Equal(t, result, ErrDonatePostDoesntExist(types.GetPostKey(user1, "invalid")).Result())
	checkPostKVStore(t, ctx, types.GetPostKey(user1, postID), postInfo, postMeta)

	// test invalid user1name
	donateMsg = NewDonateMsg(types.AccountKey("invalid"), types.LNO(sdk.NewRat(100)), user1, postID)
	result = handler(ctx, donateMsg)

	assert.Equal(t, result, ErrDonateUserNotFound(types.AccountKey("invalid")).Result())
	checkPostKVStore(t, ctx, types.GetPostKey(user1, postID), postInfo, postMeta)

	// test insufficient deposit
	donateMsg = NewDonateMsg(user2, types.LNO(sdk.NewRat(100)), user1, postID)
	result = handler(ctx, donateMsg)

	assert.Equal(t, result, ErrDonateFailed(types.GetPostKey(user1, postID)).Result())
}

func TestHandlerRePostDonate(t *testing.T) {
	ctx, am, pm, gm := setupTest(t, 1)
	handler := NewHandler(*pm, *am, *gm)

	user1, postID := createTestPost(t, ctx, "user1", "postID", am, pm, sdk.NewRat(15, 100))
	user2 := createTestAccount(ctx, am, "user2")
	user3 := createTestAccount(ctx, am, "user3")
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
		RedistributionSplitRate: sdk.ZeroRat,
	}
	msg := NewCreatePostMsg(postCreateParams)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	donateMsg := NewDonateMsg(types.AccountKey(user3), types.LNO(sdk.NewRat(100)), user2, "repost")
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
		Created:                 1,
		LastUpdate:              1,
		LastActivity:            1,
		AllowReplies:            true,
		TotalDonateCount:        1,
		TotalReward:             types.RatToCoin(sdk.NewRat(15 * types.Decimals).Mul(sdk.NewRat(99, 100))),
		RedistributionSplitRate: sdk.ZeroRat,
	}

	checkPostKVStore(t, ctx, types.GetPostKey(user2, "repost"), postInfo, postMeta)

	// check source post
	postMeta.TotalReward = types.Coin{sdk.NewRat(85 * types.Decimals).Mul(sdk.NewRat(99, 100)).Evaluate()}
	postInfo.Author = user1
	postInfo.PostID = postID
	postInfo.SourceAuthor = ""
	postInfo.SourcePostID = ""
	postMeta.RedistributionSplitRate = sdk.NewRat(15, 100)

	checkPostKVStore(t, ctx, types.GetPostKey(user1, postID), postInfo, postMeta)

	acc1Balance, _ := am.GetBankBalance(ctx, user1)
	acc2Balance, _ := am.GetBankBalance(ctx, user2)
	acc3Balance, _ := am.GetBankBalance(ctx, user3)
	assert.Equal(t, true, acc1Balance.IsEqual(types.RatToCoin(sdk.NewRat(85*types.Decimals).Mul(sdk.NewRat(99, 100)))))
	assert.Equal(t, true, acc2Balance.IsEqual(types.RatToCoin(sdk.NewRat(15*types.Decimals).Mul(sdk.NewRat(99, 100)))))
	assert.Equal(t, true, acc3Balance.IsEqual(types.NewCoin(23*types.Decimals)))
}
