package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("post")
)

func TestPost(t *testing.T) {
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

	runTest(t, func(env TestEnv) {
		err := env.ps.SetPostInfo(env.ctx, &postInfo)
		assert.Nil(t, err)

		resultPtr, err := env.ps.GetPostInfo(env.ctx, types.GetPermLink(postInfo.Author, postInfo.PostID))
		assert.Nil(t, err)
		assert.Equal(t, postInfo, *resultPtr, "postInfo should be equal")
	})

}

func TestPostMeta(t *testing.T) {
	postMeta := PostMeta{
		AllowReplies: true,
	}

	runTest(t, func(env TestEnv) {
		err := env.ps.SetPostMeta(env.ctx, types.PermLink("test"), &postMeta)
		assert.Nil(t, err)

		resultPtr, err := env.ps.GetPostMeta(env.ctx, types.PermLink("test"))
		assert.Nil(t, err)
		assert.Equal(t, postMeta, *resultPtr, "Post meta should be equal")
	})
}

func TestPostLike(t *testing.T) {
	user := types.AccountKey("test")
	postLike := Like{Username: user, Weight: 10000, CreatedAt: 100}

	runTest(t, func(env TestEnv) {
		err := env.ps.SetPostLike(env.ctx, types.PermLink("test"), &postLike)
		assert.Nil(t, err)

		resultPtr, err := env.ps.GetPostLike(env.ctx, types.PermLink("test"), user)
		assert.Nil(t, err)
		assert.Equal(t, postLike, *resultPtr, "Post like should be equal")
	})
}

func TestPostComment(t *testing.T) {
	user := types.AccountKey("test")
	postComment := Comment{Author: user, PostID: "test", CreatedAt: 100}

	runTest(t, func(env TestEnv) {
		err := env.ps.SetPostComment(env.ctx, types.PermLink("test"), &postComment)
		assert.Nil(t, err)

		resultPtr, err := env.ps.GetPostComment(env.ctx, types.PermLink("test"), types.GetPermLink(user, "test"))
		assert.Nil(t, err)
		assert.Equal(t, postComment, *resultPtr, "Post comment should be equal")
	})
}

func TestPostView(t *testing.T) {
	user := types.AccountKey("test")
	postView := View{Username: user, LastViewAt: 100, Times: 1}

	runTest(t, func(env TestEnv) {
		err := env.ps.SetPostView(env.ctx, types.PermLink("test"), &postView)
		assert.Nil(t, err)

		resultPtr, err := env.ps.GetPostView(env.ctx, types.PermLink("test"), user)
		assert.Nil(t, err)
		assert.Equal(t, postView, *resultPtr, "Post view should be equal")
	})
}

func TestPostDonate(t *testing.T) {
	user := types.AccountKey("test")
	postDonations := Donations{Username: user, DonationList: []Donation{Donation{CreatedAt: 100}}}

	runTest(t, func(env TestEnv) {
		err := env.ps.SetPostDonations(env.ctx, types.PermLink("test"), &postDonations)
		assert.Nil(t, err)

		resultPtr, err := env.ps.GetPostDonations(env.ctx, types.PermLink("test"), user)
		assert.Nil(t, err)
		assert.Equal(t, postDonations, *resultPtr, "Post donation should be equal")
	})
}

//
// Test Environment setup
//

type TestEnv struct {
	ps  PostStorage
	ctx sdk.Context
}

func runTest(t *testing.T, fc func(env TestEnv)) {
	env := TestEnv{
		ps:  NewPostStorage(TestKVStoreKey),
		ctx: getContext(),
	}
	fc(env)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}
