package post

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestKVStoreKey = sdk.NewKVStoreKey("post")
)

func newPostManager() PostManager {
	return NewPostMananger(TestKVStoreKey)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func TestPost(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()

	post := types.Post{
		PostID:  "Test Post",
		Title:   "Test Post",
		Content: "Test Post",
		Author:  types.AccountKey("author"),
		Parent:  "",
		Source:  "",
		Created: 0,
		Links:   []types.IDToURLMapping{},
	}
	err := pm.SetPost(ctx, &post)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPost(ctx, types.GetPostKey(post.Author, post.PostID))
	assert.Nil(t, err)
	assert.Equal(t, post, *resultPtr, "post should be equal")
}

func TestPostMeta(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()

	postMeta := types.PostMeta{
		LastUpdate:   0,
		LastActivity: 0,
		AllowReplies: false,
	}
	err := pm.SetPostMeta(ctx, types.PostKey("test"), &postMeta)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPostMeta(ctx, types.PostKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, postMeta, *resultPtr, "Post meta should be equal")
}

func TestPostLikes(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()

	postLikes := types.PostLikes{Likes: []types.Like{}}
	err := pm.SetPostLikes(ctx, types.PostKey("test"), &postLikes)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPostLikes(ctx, types.PostKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, postLikes, *resultPtr, "Post like list should be equal")
}

func TestPostComments(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()

	postComments := types.PostComments{Comments: []types.PostKey{}}
	err := pm.SetPostComments(ctx, types.PostKey("test"), &postComments)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPostComments(ctx, types.PostKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, postComments, *resultPtr, "Post comments should be equal")
}

func TestPostViews(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()

	postViews := types.PostViews{Views: []types.View{}}
	err := pm.SetPostViews(ctx, types.PostKey("test"), &postViews)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPostViews(ctx, types.PostKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, postViews, *resultPtr, "Post views should be equal")
}

func TestPostDonate(t *testing.T) {
	pm := newPostManager()
	ctx := getContext()

	postDonations := types.PostDonations{Donations: []types.Donation{}, Reward: sdk.Coins{}}
	err := pm.SetPostDonations(ctx, types.PostKey("test"), &postDonations)
	assert.Nil(t, err)

	resultPtr, err := pm.GetPostDonations(ctx, types.PostKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, postDonations, *resultPtr, "Post donations should be equal")
}
