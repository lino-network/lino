package testsuites

import (
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type CtxTestSuite struct {
	suite.Suite
	ms     sdk.CommitMultiStore
	height int64
	time   time.Time
	Ctx    sdk.Context
}

func (suite *CtxTestSuite) SetupCtx(height int64, t time.Time, keys ...*sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	}
	err := ms.LoadLatestVersion()
	suite.Require().Nil(err)
	suite.ms = ms
	suite.Ctx = sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: height, Time: t}, false, log.NewNopLogger())
	suite.height = height
	suite.time = t
}

func (suite *CtxTestSuite) NextBlock(t time.Time) {
	suite.Ctx = sdk.NewContext(
		suite.ms, abci.Header{
			ChainID: "Lino", Height: suite.height + 1, Time: t}, false, log.NewNopLogger())
	suite.time = t
}
