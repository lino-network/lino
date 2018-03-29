package post

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

func TestPost(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()

	postInfo := PostInfo{
		PostID:       "Test Post",
		Title:        "Test Post",
		Content:      "Test Post",
		Author:       acc.AccountKey("author"),
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []IDToURLMapping{},
	}
	err := pm.SetPostInfo(ctx, &postInfo)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPostInfo(ctx, GetPostKey(postInfo.Author, postInfo.PostID))
	assert.Nil(t, err)
	assert.Equal(t, postInfo, *resultPtr, "postInfo should be equal")
}

func TestPostMeta(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()

	postMeta := PostMeta{
		AllowReplies: true,
		TotalReward:  sdk.Coins{},
	}
	err := pm.SetPostMeta(ctx, PostKey("test"), &postMeta)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPostMeta(ctx, PostKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, postMeta, *resultPtr, "Post meta should be equal")
}

func TestPostLike(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()
	user := acc.AccountKey("test")

	postLike := Like{Username: user, Weight: 10000, Created: types.Height(100)}
	err := pm.SetPostLike(ctx, PostKey("test"), &postLike)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPostLike(ctx, PostKey("test"), user)
	assert.Nil(t, err)
	assert.Equal(t, postLike, *resultPtr, "Post like should be equal")
}

func TestPostComment(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()
	user := acc.AccountKey("test")

	postComment := Comment{Author: user, PostID: "test", Created: types.Height(100)}
	err := pm.SetPostComment(ctx, PostKey("test"), &postComment)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPostComment(ctx, PostKey("test"), GetPostKey(user, "test"))
	assert.Nil(t, err)
	assert.Equal(t, postComment, *resultPtr, "Post comment should be equal")
}

func TestPostView(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()
	user := acc.AccountKey("test")

	postView := View{Username: user, Created: types.Height(100)}
	err := pm.SetPostView(ctx, PostKey("test"), &postView)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPostView(ctx, PostKey("test"), user)
	assert.Nil(t, err)
	assert.Equal(t, postView, *resultPtr, "Post view should be equal")
}

func TestPostDonate(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()
	user := acc.AccountKey("test")

	postDonation := Donation{Username: user, Created: types.Height(100), Amount: sdk.Coins{}}
	err := pm.SetPostDonation(ctx, PostKey("test"), &postDonation)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPostDonation(ctx, PostKey("test"), user)
	assert.Nil(t, err)
	assert.Equal(t, postDonation, *resultPtr, "Post donation should be equal")
}
