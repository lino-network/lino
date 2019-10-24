package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

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
	postInfo1 := &Post{
		PostID:    "Test Post",
		Title:     "Test Post",
		Content:   "Test Post",
		Author:    linotypes.AccountKey("author1"),
		CreatedBy: linotypes.AccountKey("app"),
		CreatedAt: 1,
		UpdatedAt: 2,
		IsDeleted: true,
	}
	postInfo2 := &Post{
		PostID:    "Test Post",
		Title:     "Test Post",
		Content:   "Test Post",
		Author:    linotypes.AccountKey("author"),
		CreatedBy: linotypes.AccountKey("app"),
		CreatedAt: 1,
		UpdatedAt: 2,
		IsDeleted: true,
	}
	permlink1 := linotypes.GetPermlink(postInfo1.Author, postInfo1.PostID)
	permlink2 := linotypes.GetPermlink(postInfo2.Author, postInfo2.PostID)

	// add post1
	suite.False(suite.ps.HasPost(suite.ctx, permlink1))
	suite.False(suite.ps.HasPost(suite.ctx, permlink2))
	suite.ps.SetPost(suite.ctx, postInfo1)
	suite.True(suite.ps.HasPost(suite.ctx, permlink1))
	suite.False(suite.ps.HasPost(suite.ctx, permlink2))
	rst1, err := suite.ps.GetPost(suite.ctx, permlink1)
	suite.Nil(err)
	suite.Equal(postInfo1, rst1)
	_, err = suite.ps.GetPost(suite.ctx, permlink2)
	suite.NotNil(err)

	// add post2
	suite.ps.SetPost(suite.ctx, postInfo2)
	suite.True(suite.ps.HasPost(suite.ctx, permlink1))
	suite.True(suite.ps.HasPost(suite.ctx, permlink2))
	rst2, err := suite.ps.GetPost(suite.ctx, permlink2)
	suite.Nil(err)
	suite.Equal(postInfo2, rst2)
	rst1, err = suite.ps.GetPost(suite.ctx, permlink1)
	suite.Nil(err)
	suite.Equal(postInfo1, rst1)
}

func (suite *postStoreTestSuite) TestConsumptionWIndowGetSet() {
	suite.ps.SetConsumptionWindow(suite.ctx, linotypes.NewMiniDollar(1000))
	suite.Equal(linotypes.NewMiniDollar(1000), suite.ps.GetConsumptionWindow(suite.ctx))

	suite.ps.SetConsumptionWindow(suite.ctx, linotypes.NewMiniDollar(2000))
	suite.Equal(linotypes.NewMiniDollar(2000), suite.ps.GetConsumptionWindow(suite.ctx))
}
