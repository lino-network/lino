package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("post")
)

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func TestPost(t *testing.T) {
	ps := NewPostStorage(TestKVStoreKey)
	ctx := getContext()

	postInfo := PostInfo{
		PostID:       "Test Post",
		Title:        "Test Post",
		Content:      "Test Post",
		Author:       types.AccountKey("author"),
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        nil,
	}
	err := ps.SetPostInfo(ctx, &postInfo)
	assert.Nil(t, err)

	resultPtr, err := ps.GetPostInfo(ctx, types.GetPermLink(postInfo.Author, postInfo.PostID))
	assert.Nil(t, err)
	assert.Equal(t, postInfo, *resultPtr, "postInfo should be equal")
}

func TestPostMeta(t *testing.T) {
	ps := NewPostStorage(TestKVStoreKey)
	ctx := getContext()

	postMeta := PostMeta{
		AllowReplies: true,
	}
	err := ps.SetPostMeta(ctx, types.PermLink("test"), &postMeta)
	assert.Nil(t, err)

	resultPtr, err := ps.GetPostMeta(ctx, types.PermLink("test"))
	assert.Nil(t, err)
	assert.Equal(t, postMeta, *resultPtr, "Post meta should be equal")
}

func TestPostLike(t *testing.T) {
	ps := NewPostStorage(TestKVStoreKey)
	ctx := getContext()
	user := types.AccountKey("test")

	postLike := Like{Username: user, Weight: 10000, Created: 100}
	err := ps.SetPostLike(ctx, types.PermLink("test"), &postLike)
	assert.Nil(t, err)

	resultPtr, err := ps.GetPostLike(ctx, types.PermLink("test"), user)
	assert.Nil(t, err)
	assert.Equal(t, postLike, *resultPtr, "Post like should be equal")
}

func TestPostComment(t *testing.T) {
	ps := NewPostStorage(TestKVStoreKey)
	ctx := getContext()
	user := types.AccountKey("test")

	postComment := Comment{Author: user, PostID: "test", Created: 100}
	err := ps.SetPostComment(ctx, types.PermLink("test"), &postComment)
	assert.Nil(t, err)

	resultPtr, err := ps.GetPostComment(ctx, types.PermLink("test"), types.GetPermLink(user, "test"))
	assert.Nil(t, err)
	assert.Equal(t, postComment, *resultPtr, "Post comment should be equal")
}

func TestPostView(t *testing.T) {
	ps := NewPostStorage(TestKVStoreKey)
	ctx := getContext()
	user := types.AccountKey("test")

	postView := View{Username: user, Created: 100}
	err := ps.SetPostView(ctx, types.PermLink("test"), &postView)
	assert.Nil(t, err)

	resultPtr, err := ps.GetPostView(ctx, types.PermLink("test"), user)
	assert.Nil(t, err)
	assert.Equal(t, postView, *resultPtr, "Post view should be equal")
}

func TestPostDonate(t *testing.T) {
	ps := NewPostStorage(TestKVStoreKey)
	ctx := getContext()
	user := types.AccountKey("test")

	postDonations := Donations{Username: user, DonationList: []Donation{Donation{Created: 100}}}
	err := ps.SetPostDonations(ctx, types.PermLink("test"), &postDonations)
	assert.Nil(t, err)

	resultPtr, err := ps.GetPostDonations(ctx, types.PermLink("test"), user)
	assert.Nil(t, err)
	assert.Equal(t, postDonations, *resultPtr, "Post donation should be equal")
}
