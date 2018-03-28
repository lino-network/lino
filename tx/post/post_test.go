package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPost(t *testing.T) {
	pm := newPostManager()
	user := acc.AccountKey("user")
	postID := "post ID"

	post := NewProxyPost(user, postID, &pm)
	assert.Equal(t, user, post.GetAuthor())
	assert.Equal(t, postID, post.GetPostID())
	assert.Equal(t, GetPostKey(user, postID), post.GetPostKey())
	assert.NotNil(t, post.postManager)
	assert.Nil(t, post.postInfo)
	assert.Nil(t, post.postMeta)
	assert.False(t, post.writePostInfo)
	assert.False(t, post.writePostMeta)
}

// checkPostKVStore checks all post related structs in the post manager
func checkPostKVStore(t *testing.T, ctx sdk.Context, pm PostManager,
	postKey PostKey, postInfo PostInfo, postMeta PostMeta) {
	// check all post related structs in KVStore
	postPtr, err := pm.GetPostInfo(ctx, postKey)
	assert.Nil(t, err)
	assert.Equal(t, postInfo, *postPtr, "postInfo should be equal")
	postMetaPtr, err := pm.GetPostMeta(ctx, postKey)
	assert.Nil(t, err)
	assert.Equal(t, postMeta, *postMetaPtr, "Post meta should be equal")
}

// test create post
func TestCreatePost(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()
	ctx = ctx.WithBlockHeight(1)
	author := acc.AccountKey("author")
	postID := "TestPostID"
	post := NewProxyPost(author, postID, &pm)
	assert.False(t, post.IsPostExist(ctx))
	// test valid postInfo
	postInfo := PostInfo{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       author,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []IDToURLMapping{},
	}
	err := post.CreatePost(ctx, &postInfo)
	assert.Nil(t, err)

	// test created struct before apply
	assert.Equal(t, postInfo, *post.postInfo, "postInfo should be equal")
	postMeta := PostMeta{
		Created:      1,
		LastUpdate:   1,
		LastActivity: 1,
		AllowReplies: true,
	}
	assert.Equal(t, postMeta, *post.postMeta, "Post meta should be equal")

	// after apply the post proxy should be cleared
	post.Apply(ctx)
	assert.Nil(t, post.postInfo)
	assert.Nil(t, post.postMeta)

	// after apply check KVStore
	postMeta.TotalReward = sdk.Coins{}
	checkPostKVStore(t, ctx, pm, post.GetPostKey(), postInfo, postMeta)
	// test recreate post
	err = post.CreatePost(ctx, &postInfo)
	assert.Equal(t, err, ErrPostExist())
}

func TestComment(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()
	ctx = ctx.WithBlockHeight(1)
	author := acc.AccountKey("author")
	postID := "TestPostID"
	post := NewProxyPost(author, postID, &pm)
	assert.False(t, post.IsPostExist(ctx))

	// test valid postInfo
	postInfo := PostInfo{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       author,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []IDToURLMapping{},
	}
	err := post.CreatePost(ctx, &postInfo)
	assert.Nil(t, err)
	post.Apply(ctx)

	ctx = ctx.WithBlockHeight(2)

	postComment := Comment{Author: author, PostID: "test", Created: types.Height(100)}
	err = post.AddComment(ctx, postComment)
	assert.Nil(t, err)
	post.Apply(ctx)

	// after apply check KVStore
	postMeta := PostMeta{
		Created:      1,
		LastUpdate:   1,
		LastActivity: 2,
		AllowReplies: true,
		TotalReward:  sdk.Coins{},
	}
	checkPostKVStore(t, ctx, pm, post.GetPostKey(), postInfo, postMeta)
}
