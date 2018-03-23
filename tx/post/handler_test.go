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
	postViews := PostViews{Views: []View{}}
	postLikes := PostLikes{Likes: []Like{}}
	postComments := PostComments{Comments: []PostKey{}}
	postDonations := PostDonations{Donations: []Donation{}, Reward: sdk.Coins{}}
	checkPostKVStore(t, ctx, pm, GetPostKey(user, "TestPostID"), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// test invlaid author
	postInfo.Author = acc.AccountKey("invalid")
	msg = NewCreatePostMsg(postInfo)
	result = handler(ctx, msg)
	assert.Equal(t, result, acc.ErrUsernameNotFound("invalid").Result())
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
	postViews := PostViews{Views: []View{}}
	postLikes := PostLikes{Likes: []Like{}}
	postComments := PostComments{Comments: []PostKey{}}
	postDonations := PostDonations{Donations: []Donation{}, Reward: sdk.Coins{}}
	checkPostKVStore(t, ctx, pm, GetPostKey(user, "comment"), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// check parent
	postInfo.PostID = "TestPostID"
	postInfo.ParentAuthor = ""
	postInfo.ParentPostID = ""
	postMeta.Created = 0
	postMeta.LastUpdate = 0
	postComments = PostComments{Comments: []PostKey{GetPostKey(user, "comment")}}
	checkPostKVStore(t, ctx, pm, GetPostKey(user, "TestPostID"), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// test invalid parent
	postInfo.PostID = "invalid post"
	postInfo.ParentAuthor = user
	postInfo.ParentPostID = "invalid parent"
	msg = NewCreatePostMsg(postInfo)

	result = handler(ctx, msg)
	assert.Equal(t, result, ErrPostCommentsNotFound(GetPostKey(user, "invalid parent")).Result())

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
	assert.Equal(t, result, ErrPostCommentsNotFound(GetPostKey(user, "newComment")).Result())
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
		Created:      0,
		LastUpdate:   0,
		LastActivity: 0,
		AllowReplies: true,
	}
	postViews := PostViews{Views: []View{}}
	postLikes := PostLikes{Likes: []Like{Like{Username: acc.AccountKey(user), Weight: 10000}}, TotalWeight: 10000}
	postComments := PostComments{Comments: []PostKey{}}
	postDonations := PostDonations{Donations: []Donation{}, Reward: sdk.Coins{}}
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// test update like
	likeMsg = NewLikeMsg(acc.AccountKey(user), -10000, acc.AccountKey(user), postID)
	result = handler(ctx, likeMsg)
	assert.Equal(t, result, sdk.Result{})

	postLikes = PostLikes{Likes: []Like{Like{Username: acc.AccountKey(user), Weight: -10000}}, TotalWeight: -10000}
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// test invalid like target post
	likeMsg = NewLikeMsg(acc.AccountKey(user), -10000, acc.AccountKey(user), "invalid")
	result = handler(ctx, likeMsg)
	assert.Equal(t, result, ErrLikePostDoesntExist().Result())
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// test invalid like username
	likeMsg = NewLikeMsg(acc.AccountKey("invalid"), 10000, acc.AccountKey(user), postID)
	result = handler(ctx, likeMsg)
	assert.Equal(t, result, acc.ErrUsernameNotFound(string(likeMsg.Username)).Result())
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

}

func TestHandlerPostDonate(t *testing.T) {
	pm := newPostManager()
	lam := acc.NewLinoAccountManager(TestKVStoreKey)
	ctx := getContext()

	user := "username"
	postID := "postID"
	handler := NewHandler(pm, lam)
	createTestAccount(ctx, lam, user)
	createTestPost(ctx, lam, pm, user, postID)

	donateMsg := NewDonateMsg(acc.AccountKey(user), newAmount(100), acc.AccountKey(user), postID)
	result := handler(ctx, donateMsg)
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
		Created:      0,
		LastUpdate:   0,
		LastActivity: 0,
		AllowReplies: true,
	}
	postViews := PostViews{Views: []View{}}
	postLikes := PostLikes{Likes: []Like{}}
	postComments := PostComments{Comments: []PostKey{}}
	postDonations := PostDonations{Donations: []Donation{
		Donation{
			Username: donateMsg.Username,
			Amount:   donateMsg.Amount,
			Created:  types.Height(ctx.BlockHeight()),
		}}, Reward: newAmount(100),
	}
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// test invalid donation target
	donateMsg = NewDonateMsg(acc.AccountKey(user), newAmount(100), acc.AccountKey(user), "invalid")
	result = handler(ctx, donateMsg)
	assert.Equal(t, result, ErrDonatePostDoesntExist().Result())
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// test invalid username
	donateMsg = NewDonateMsg(acc.AccountKey("invalid"), newAmount(100), acc.AccountKey(user), postID)
	result = handler(ctx, donateMsg)
	assert.Equal(t, result, acc.ErrUsernameNotFound(string(donateMsg.Username)).Result())
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta, postLikes, postComments, postViews, postDonations)

	// test insufficient deposit
	donateMsg = NewDonateMsg(acc.AccountKey(user), newAmount(100), acc.AccountKey(user), postID)
	result = handler(ctx, donateMsg)
	assert.Equal(t, result, acc.ErrAccountManagerFail("Account bank's coins are not enough").Result())
	checkPostKVStore(t, ctx, pm, GetPostKey(acc.AccountKey(user), postID), postInfo, postMeta, postLikes, postComments, postViews, postDonations)
}
