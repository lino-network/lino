package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
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
		Links:        []types.IDToURLMapping{types.IDToURLMapping{Identifier: "test", URL: "https://lino.network"}},
	}

	runTest(t, func(env TestEnv) {
		err := env.ps.SetPostInfo(env.ctx, &postInfo)
		assert.Nil(t, err)

		resultPtr, err := env.ps.GetPostInfo(env.ctx, types.GetPermlink(postInfo.Author, postInfo.PostID))
		assert.Nil(t, err)
		assert.Equal(t, postInfo, *resultPtr, "postInfo should be equal")
	})
}

func TestUTF8(t *testing.T) {
	postInfo := PostInfo{
		PostID:       "Test Post",
		Title:        "Test Post",
		Content:      "12345 67890 你好👌12345 67890 你好👌",
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

		resultPtr, err := env.ps.GetPostInfo(env.ctx, types.GetPermlink(postInfo.Author, postInfo.PostID))
		assert.Nil(t, err)
		assert.Equal(t, postInfo, *resultPtr, "postInfo should be equal")
	})
}

func TestPostMeta(t *testing.T) {
	postMeta := PostMeta{
		AllowReplies:            true,
		RedistributionSplitRate: sdk.ZeroDec(),
		TotalUpvoteCoinDay:      types.NewCoinFromInt64(0),
		TotalReportCoinDay:      types.NewCoinFromInt64(0),
		TotalReward:             types.NewCoinFromInt64(0),
	}

	runTest(t, func(env TestEnv) {
		err := env.ps.SetPostMeta(env.ctx, types.Permlink("test"), &postMeta)
		assert.Nil(t, err)

		resultPtr, err := env.ps.GetPostMeta(env.ctx, types.Permlink("test"))
		assert.Nil(t, err)
		assert.Equal(t, postMeta, *resultPtr, "Post meta should be equal")
	})
}

func TestPostComment(t *testing.T) {
	user := types.AccountKey("test")
	postComment := Comment{Author: user, PostID: "test", CreatedAt: 100}

	runTest(t, func(env TestEnv) {
		err := env.ps.SetPostComment(env.ctx, types.Permlink("test"), &postComment)
		assert.Nil(t, err)

		resultPtr, err := env.ps.GetPostComment(env.ctx, types.Permlink("test"), types.GetPermlink(user, "test"))
		assert.Nil(t, err)
		assert.Equal(t, postComment, *resultPtr, "Post comment should be equal")
	})
}

func TestPostView(t *testing.T) {
	user := types.AccountKey("test")
	postView := View{Username: user, LastViewAt: 100, Times: 1}

	runTest(t, func(env TestEnv) {
		err := env.ps.SetPostView(env.ctx, types.Permlink("test"), &postView)
		assert.Nil(t, err)

		resultPtr, err := env.ps.GetPostView(env.ctx, types.Permlink("test"), user)
		assert.Nil(t, err)
		assert.Equal(t, postView, *resultPtr, "Post view should be equal")
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

	return sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
}
