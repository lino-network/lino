package post

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestHandlerCreatePost(t *testing.T) {
	pm := newPostManager()
	lam := acc.NewLinoAccountManager(TestKVStoreKey)
	ctx := getContext()

	handler := NewHandler(pm, lam)

	user := acc.AccountKey("testuser")
	createTestAccount(ctx, lam, string(user))

	// test valid post
	postInfo := PostInfo{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []IDToURLMapping{},
	}
	msg := NewCreatePostMsg(postInfo)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})
	// after handler check KVStore
	postMeta := PostMeta{
		Created:      0,
		LastUpdate:   0,
		LastActivity: 0,
		AllowReplies: true,
	}
	checkPostKVStore(t, ctx, pm, GetPostKey(user, "TestPostID"), postInfo, postMeta)

	// test invlaid author
	postInfo.Author = acc.AccountKey("invalid")
	msg = NewCreatePostMsg(postInfo)
	result = handler(ctx, msg)
	assert.Equal(t, result, acc.ErrUsernameNotFound().Result())
}

func TestHandlerCreateComment(t *testing.T) {
	pm := newPostManager()
	lam := acc.NewLinoAccountManager(TestKVStoreKey)
	ctx := getContext()

	handler := NewHandler(pm, lam)

	user := acc.AccountKey("testuser")
	createTestAccount(ctx, lam, string(user))

	postInfo := PostInfo{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []IDToURLMapping{},
	}
	msg := NewCreatePostMsg(postInfo)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// test comment
	postInfo.Author = user
	postInfo.PostID = "comment"
	postInfo.ParentAuthor = user
	postInfo.ParentPostID = "TestPostID"
	msg = NewCreatePostMsg(postInfo)
	ctx = ctx.WithBlockHeight(1)
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	// check comment
	postMeta := PostMeta{
		Created:      1,
		LastUpdate:   1,
		LastActivity: 1,
		AllowReplies: true,
	}

	checkPostKVStore(t, ctx, pm, GetPostKey(user, "comment"), postInfo, postMeta)

	// check parent
	postInfo.PostID = "TestPostID"
	postInfo.ParentAuthor = ""
	postInfo.ParentPostID = ""
	postMeta.Created = 0
	postMeta.LastUpdate = 0
	checkPostKVStore(t, ctx, pm, GetPostKey(user, "TestPostID"), postInfo, postMeta)

	// test invalid parent
	postInfo.PostID = "invalid post"
	postInfo.ParentAuthor = user
	postInfo.ParentPostID = "invalid parent"
	msg = NewCreatePostMsg(postInfo)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostMetaNotFound(GetPostMetaKey(GetPostKey(user, "invalid parent"))).Result())

	// test duplicate comment
	postInfo.Author = user
	postInfo.PostID = "comment"
	postInfo.ParentAuthor = user
	postInfo.ParentPostID = "TestPostID"
	msg = NewCreatePostMsg(postInfo)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostExist().Result())

	// test cycle comment
	postInfo.Author = user
	postInfo.PostID = "newComment"
	postInfo.ParentAuthor = user
	postInfo.ParentPostID = "newComment"
	msg = NewCreatePostMsg(postInfo)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostMetaNotFound(GetPostMetaKey(GetPostKey(user, "newComment"))).Result())
}

func TestHandlerRepost(t *testing.T) {
	pm := newPostManager()
	lam := acc.NewLinoAccountManager(TestKVStoreKey)
	ctx := getContext()

	handler := NewHandler(pm, lam)

	user := acc.AccountKey("testuser")
	createTestAccount(ctx, lam, string(user))

	postInfo := PostInfo{
		PostID:       "TestPostID",
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []IDToURLMapping{},
	}
	msg := NewCreatePostMsg(postInfo)
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// test 1 depth repost
	postInfo.Author = user
	postInfo.PostID = "repost"
	postInfo.SourceAuthor = user
	postInfo.SourcePostID = "TestPostID"
	msg = NewCreatePostMsg(postInfo)
	ctx = ctx.WithBlockHeight(1)
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	// check 1 depth repost
	postMeta := PostMeta{
		Created:      1,
		LastUpdate:   1,
		LastActivity: 1,
		AllowReplies: true,
	}
	checkPostKVStore(t, ctx, pm, GetPostKey(user, "repost"), postInfo, postMeta)

	// test 2 depth repost
	postInfo.PostID = "repost-repost"
	postInfo.SourceAuthor = user
	postInfo.SourcePostID = "repost"
	msg = NewCreatePostMsg(postInfo)
	ctx = ctx.WithBlockHeight(2)
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	// check 2 depth repost
	postMeta = PostMeta{
		Created:      2,
		LastUpdate:   2,
		LastActivity: 2,
		AllowReplies: true,
	}
	postInfo.SourceAuthor = user
	postInfo.SourcePostID = "TestPostID"
	checkPostKVStore(t, ctx, pm, GetPostKey(user, "repost-repost"), postInfo, postMeta)
}

func TestHandlerPostLike(t *testing.T) {
	pm := newPostManager()
	lam := acc.NewLinoAccountManager(TestKVStoreKey)
	ctx := getContext()

	user := "username"
	postID := "postID"
	handler := NewHandler(pm, lam)
	createTestAccount(ctx, lam, user)
	createTestPost(ctx, lam, pm, user, postID)

	likeMsg := NewLikeMsg(acc.AccountKey(user), 10000, acc.AccountKey(user), postID)
	result := handler(ctx, likeMsg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	postInfo := PostInfo{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       acc.AccountKey(user),
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []IDToURLMapping{},
	}
	postMeta := PostMeta{
		Created:         0,
		LastUpdate:      0,
		LastActivity:    0,
		AllowReplies:    true,
		TotalLikeCount:  1,
		TotalLikeWeight: 10000,
	}
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta)

	// test update like
	likeMsg = NewLikeMsg(acc.AccountKey(user), -10000, acc.AccountKey(user), postID)
	result = handler(ctx, likeMsg)
	assert.Equal(t, result, sdk.Result{})
	postMeta.TotalLikeWeight = -10000
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta)

	// test invalid like target post
	likeMsg = NewLikeMsg(acc.AccountKey(user), -10000, acc.AccountKey(user), "invalid")
	result = handler(ctx, likeMsg)
	assert.Equal(t, result, ErrLikePostDoesntExist().Result())
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta)

	// test invalid like username
	likeMsg = NewLikeMsg(acc.AccountKey("invalid"), 10000, acc.AccountKey(user), postID)
	result = handler(ctx, likeMsg)

	assert.Equal(t, result, acc.ErrUsernameNotFound().Result())
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta)
}

func TestHandlerPostDonate(t *testing.T) {
	pm := newPostManager()
	lam := acc.NewLinoAccountManager(TestKVStoreKey)
	ctx := getContext()

	user1 := "user1"
	user2 := "user2"
	postID := "postID"
	handler := NewHandler(pm, lam)
	accProxy1 := createTestAccount(ctx, lam, user1)
	accProxy2 := createTestAccount(ctx, lam, user2)
	createTestPost(ctx, lam, pm, user1, postID)

	donateMsg := NewDonateMsg(acc.AccountKey(user2), types.TestLNO(sdk.NewRat(100)), acc.AccountKey(user1), postID)
	result := handler(ctx, donateMsg)
	assert.Equal(t, result, sdk.Result{})

	// after handler check KVStore
	postInfo := PostInfo{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       acc.AccountKey(user1),
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []IDToURLMapping{},
	}
	postMeta := PostMeta{
		Created:          0,
		LastUpdate:       0,
		LastActivity:     0,
		AllowReplies:     true,
		TotalDonateCount: 1,
		TotalReward:      types.Coin{100 * types.Decimals},
	}

	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user1), postID), postInfo, postMeta)

	acc1Balance, _ := accProxy1.GetBankBalance(ctx)
	acc2Balance, _ := accProxy2.GetBankBalance(ctx)
	assert.Equal(t, true, acc1Balance.IsEqual(types.Coin{223 * types.Decimals}))
	assert.Equal(t, true, acc2Balance.IsEqual(types.Coin{23 * types.Decimals}))
	// test invalid donation target
	donateMsg = NewDonateMsg(acc.AccountKey(user1), types.TestLNO(sdk.NewRat(100)), acc.AccountKey(user1), "invalid")
	result = handler(ctx, donateMsg)
	assert.Equal(t, result, ErrDonatePostDoesntExist().Result())
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user1), postID), postInfo, postMeta)

	// test invalid user1name
	donateMsg = NewDonateMsg(acc.AccountKey("invalid"), types.TestLNO(sdk.NewRat(100)), acc.AccountKey(user1), postID)
	result = handler(ctx, donateMsg)

	assert.Equal(t, result, acc.ErrUsernameNotFound().Result())
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user1), postID), postInfo, postMeta)

	// test insufficient deposit
	donateMsg = NewDonateMsg(acc.AccountKey(user2), types.TestLNO(sdk.NewRat(100)), acc.AccountKey(user1), postID)
	result = handler(ctx, donateMsg)

	assert.Equal(t, result, acc.ErrAccountCoinNotEnough().Result())
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user1), postID), postInfo, postMeta)
}
