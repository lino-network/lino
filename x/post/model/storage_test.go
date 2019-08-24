package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	linotypes "github.com/lino-network/lino/types"
)

type postStoreTestSuite struct {
	suite.Suite
	ctx sdk.Context
	ps  PostStorage
}

func TestPostStoreTestSuite(t *testing.T) {
	suite.Run(t, &postStoreTestSuite{})
}

func (suite *postStoreTestSuite) SetupTest() {
	TestKVStoreKey := sdk.NewKVStoreKey("post")
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()
	suite.ctx = sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	suite.ps = NewPostStorage(TestKVStoreKey)
}

func (suite *postStoreTestSuite) TestPostGetSetHas() {
	postInfo := &Post{
		PostID:    "Test Post",
		Title:     "Test Post",
		Content:   "Test Post",
		Author:    linotypes.AccountKey("author"),
		CreatedBy: linotypes.AccountKey("app"),
		CreatedAt: 1,
		UpdatedAt: 2,
		IsDeleted: true,
	}
	permlink := linotypes.GetPermlink(postInfo.Author, postInfo.PostID)

	suite.False(suite.ps.HasPost(suite.ctx, permlink))
	suite.ps.SetPost(suite.ctx, postInfo)
	suite.True(suite.ps.HasPost(suite.ctx, permlink))

	rst, err := suite.ps.GetPost(suite.ctx, permlink)
	suite.Nil(err)
	suite.Equal(postInfo, rst)
}
