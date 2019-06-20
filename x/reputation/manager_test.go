package reputation

import (
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/reputation/repv2"
)

type reputationTestSuite struct {
	suite.Suite
	ms     sdk.CommitMultiStore
	ph     param.ParamHolder
	rep    ReputationManager
	height int64
	ctx    sdk.Context
	t      time.Time
}

func (suite *reputationTestSuite) timefies(afterUpdate6 bool) {
	suite.t = suite.t.Add(25 * time.Hour)
	suite.height = suite.height + 1
	if afterUpdate6 && suite.height < types.BlockchainUpgrade1Update6Height {
		suite.height = types.BlockchainUpgrade1Update6Height
	}
	// donate several rounds
	newctx := sdk.NewContext(
		suite.ms, abci.Header{
			ChainID: "Lino",
			Height:  suite.height,
			Time:    suite.t}, false, log.NewNopLogger())

	suite.rep.Update(newctx)
	suite.ctx = newctx
}

func TestReputationTestSuite(t *testing.T) {
	suite.Run(t, &reputationTestSuite{})
}

func (suite *reputationTestSuite) SetupTest() {
	suite.t = time.Now()
	TestParamKVStoreKey := sdk.NewKVStoreKey("param")
	TestRepv1KVStoreKey := sdk.NewKVStoreKey("repv1")
	TestRepv2KVStoreKey := sdk.NewKVStoreKey("repv2")
	// TestReputationKVStoreKey := sdk.NewKVStoreKey("rep")
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestRepv1KVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestRepv2KVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: 1, Time: suite.t}, false, log.NewNopLogger())

	ph := param.NewParamHolder(TestParamKVStoreKey)
	ph.InitParam(ctx)
	rep := NewReputationManager(TestRepv1KVStoreKey, TestRepv2KVStoreKey, ph)
	suite.ph = ph
	suite.rep = rep
	suite.ctx = ctx
	suite.ms = ms
	suite.height = 1
}

func (suite *reputationTestSuite) TestGetHandlers() {
	rep := suite.rep
	_, err := rep.getHandler(suite.ctx)
	suite.Nil(err)
	v2 := rep.getHandlerV2(suite.ctx)

	v2impl := v2.(repv2.ReputationImpl)

	suite.Equal(200, v2impl.BestN)
	suite.Equal(50, v2impl.UserMaxN)
	suite.Equal(int64(25*3600), v2impl.RoundDurationSeconds)
	suite.Equal(int64(10), v2impl.SampleWindowSize)
	suite.Equal(int64(10), v2impl.DecayFactor)
}

func (suite *reputationTestSuite) TestMisc() {
	suite.Nil(suite.rep.checkUsername("x"))
	suite.NotNil(suite.rep.checkUsername(""))
	suite.Nil(suite.rep.checkPost("x"))
	suite.NotNil(suite.rep.checkPost(""))
	suite.Nil(suite.rep.basicCheck("x", "y"))
	suite.NotNil(suite.rep.basicCheck("", "y"))
	suite.NotNil(suite.rep.basicCheck("x", ""))
	suite.NotNil(suite.rep.basicCheck("", ""))
}

func (suite *reputationTestSuite) TestMigrate() {
	rep := suite.rep
	v1, _ := rep.getHandler(suite.ctx)
	v2 := rep.getHandlerV2(suite.ctx)
	suite.Require().NotNil(v1)
	suite.Require().NotNil(v2)

	v1.IncFreeScore("yxia", big.NewInt(100000000))
	rep.migrate(v1, v2, "yxia")
	suite.Equal(big.NewInt(100100000), v2.GetReputation("yxia"))
}

func (suite *reputationTestSuite) TestReportAt() {
	rep := suite.rep
	coin, err := rep.ReportAt(suite.ctx, "username", "post")
	suite.Nil(err)
	suite.Equal(types.NewCoinFromInt64(-100000), coin)
	suite.timefies(true)
	coin, err = rep.ReportAt(suite.ctx, "username", "post")
	suite.Nil(err)
	suite.Equal(types.NewCoinFromInt64(0).String(), coin.String())
}

func (suite *reputationTestSuite) TestOnStake() {
	rep := suite.rep
	rep.OnStakeIn(suite.ctx, "yxia", types.NewCoinFromInt64(10000))
	rv, err := rep.GetReputation(suite.ctx, "yxia")
	suite.Nil(err)
	suite.Equal(types.NewCoinFromInt64(100015), rv)
	suite.timefies(false)
	rv, _ = rep.GetReputation(suite.ctx, "yxia")
	suite.Equal(types.NewCoinFromInt64(100015), rv)
	rep.OnStakeOut(suite.ctx, "yxia", types.NewCoinFromInt64(10000))
	rv, _ = rep.GetReputation(suite.ctx, "yxia")
	suite.Equal(types.NewCoinFromInt64(100000).String(), rv.String())

	suite.timefies(true)
	rv, _ = rep.GetReputation(suite.ctx, "yxia")
	suite.Equal(types.NewCoinFromInt64(1).String(), rv.String())

	suite.timefies(true)
	rep.OnStakeIn(suite.ctx, "yxia", types.NewCoinFromInt64(10000000))
	rv, _ = rep.GetReputation(suite.ctx, "yxia")
	suite.Equal(types.NewCoinFromInt64(1).String(), rv.String())
	rep.OnStakeOut(suite.ctx, "yxia", types.NewCoinFromInt64(10000))
	rv, _ = rep.GetReputation(suite.ctx, "yxia")
	suite.Equal(types.NewCoinFromInt64(1).String(), rv.String())

	suite.timefies(true)
	rv, _ = rep.GetReputation(suite.ctx, "yxia")
	suite.Equal(types.NewCoinFromInt64(1).String(), rv.String())
}

func (suite *reputationTestSuite) TestGetCurrentRound() {
	rep := suite.rep
	ts, err := rep.GetCurrentRound(suite.ctx)
	suite.Nil(err)
	suite.Equal(int64(0), ts)

	suite.timefies(false)
	ts, err = rep.GetCurrentRound(suite.ctx)
	suite.Nil(err)
	suite.Equal(suite.t.Unix(), ts)

	suite.timefies(true)
	ts, err = rep.GetCurrentRound(suite.ctx)
	suite.Nil(err)
	suite.Equal(suite.t.Unix(), ts)
}

func (suite *reputationTestSuite) TestGetSumRep() {
	rep := suite.rep
	rep.DonateAt(suite.ctx, "yxia", "sp", types.NewCoinFromInt64(1))
	sumrep, err := rep.GetSumRep(suite.ctx, "sp")
	suite.Nil(err)
	suite.Equal(types.NewCoinFromInt64(100000), sumrep)

	suite.timefies(true)
	sumrep, err = rep.GetSumRep(suite.ctx, "sp")
	suite.Nil(err)
	suite.Equal(types.NewCoinFromInt64(0), sumrep)
}

func (suite *reputationTestSuite) TestDonateAtUpdate() {
	rep := suite.rep
	// errors
	_, err := rep.DonateAt(suite.ctx, "", "", types.NewCoinFromInt64(3333))
	suite.NotNil(err)
	_, err = rep.DonateAt(suite.ctx, "yxia", "", types.NewCoinFromInt64(3333))
	suite.NotNil(err)
	_, err = rep.DonateAt(suite.ctx, "", "xxx", types.NewCoinFromInt64(3333))
	suite.NotNil(err)

	cases := []struct {
		username types.AccountKey
		post     types.Permlink
		amount   int64
		dp       int64
		rep      int64
	}{
		{"user1", "post1", (100 * 100000), (1 * 100000), (1416175)},
		{"user2", "post1", (100 * 100000), (1 * 100000), (763824)},
		{"user3", "post2", (100 * 100000), (1 * 100000), (1416175)},
		{"user4", "post2", (100 * 100000), (1 * 100000), (763824)},
		{"user5", "post3", (100 * 100000), (1 * 100000), (1416175)},
		{"user6", "post3", (100 * 100000), (1 * 100000), (763824)},
		{"user7", "post4", (100 * 100000), (1 * 100000), (1416175)},
		{"user8", "post4", (100 * 100000), (1 * 100000), (763824)},
		{"user9", "post5", (0.5 * 100000), (0.5 * 100000), 100000},
		{"user10", "post6", (0.5 * 100000), (0.5 * 100000), 100000},
		{"user11", "post7", (0.5 * 100000), (0.5 * 100000), 100000},
		{"user12", "post8", (0.5 * 100000), (0.5 * 100000), 100000},
	}
	for _, v := range cases {
		dp, err := rep.DonateAt(suite.ctx, v.username, v.post, types.NewCoinFromInt64(v.amount))
		suite.Nil(err)
		suite.Equal(dp, types.NewCoinFromInt64(v.dp))
	}
	suite.timefies(false)

	// check reputation.
	for i, v := range cases {
		rv, err := rep.GetReputation(suite.ctx, v.username)
		suite.Nil(err)
		suite.Equal(types.NewCoinFromInt64(v.rep).String(), rv.String(), "%d", i)
	}

	// update6
	suite.timefies(true)
	// check reputation.
	for i, v := range cases {
		rv, err := rep.GetReputation(suite.ctx, v.username)
		suite.Nil(err)
		if v.rep > 100000 {
			suite.Equal(types.NewCoinFromInt64(v.rep).String(), rv.String(), "%d", i)
		} else {
			suite.Equal(types.NewCoinFromInt64(1).String(), rv.String(), "%d", i)
		}
	}

	suite.timefies(true)
	// check reputation.
	for i, v := range cases {
		rv, err := rep.GetReputation(suite.ctx, v.username)
		suite.Nil(err)
		if v.rep > 100000 {
			suite.Equal(types.NewCoinFromInt64(v.rep).String(), rv.String(), "%d", i)
		} else {
			suite.Equal(types.NewCoinFromInt64(1).String(), rv.String(), "%d", i)
		}
	}

	// donate
	for _, v := range cases {
		dp, err := rep.DonateAt(suite.ctx, v.username, v.post, types.NewCoinFromInt64(v.amount))
		suite.Nil(err)
		r := v.rep
		if v.rep <= 100000 {
			r = 1
		}
		suite.Equal(
			types.NewCoinFromInt64(min(r, v.amount)), dp, "ne: %s, %d", dp.String(), v.rep)
	}

	suite.timefies(true)
	// donate
	for _, v := range cases {
		rep.DonateAt(suite.ctx, v.username, v.post, types.NewCoinFromInt64(v.amount))
	}

	suite.timefies(true)
	// check reputation.
	for i, v := range cases {
		rv, err := rep.GetReputation(suite.ctx, v.username)
		suite.Nil(err)
		if i >= 8 {
			suite.Equal(types.NewCoinFromInt64(0).String(), rv.String(), "%d", i)
		} else {
			if i%2 == 0 {
				suite.Equal(types.NewCoinFromInt64(1502021).String(), rv.String(), "%d", i)
			} else {
				suite.Equal(types.NewCoinFromInt64(856196).String(), rv.String(), "%d", i)
			}
		}
	}
}

func (suite *reputationTestSuite) TestExportImport() {
	rep := suite.rep
	cases := []struct {
		username types.AccountKey
		post     types.Permlink
		amount   int64
		dp       int64
		rep      int64
	}{
		{"user1", "post1", (100 * 100000), (1 * 100000), (1416175)},
		{"user2", "post1", (100 * 100000), (1 * 100000), (763824)},
		{"user3", "post2", (100 * 100000), (1 * 100000), (1416175)},
		{"user4", "post2", (100 * 100000), (1 * 100000), (763824)},
		{"user5", "post3", (100 * 100000), (1 * 100000), (1416175)},
		{"user6", "post3", (100 * 100000), (1 * 100000), (763824)},
		{"user7", "post4", (100 * 100000), (1 * 100000), (1416175)},
		{"user8", "post4", (100 * 100000), (1 * 100000), (763824)},
		{"user9", "post5", (0.5 * 100000), (0.5 * 100000), 100000},
		{"user10", "post6", (0.5 * 100000), (0.5 * 100000), 100000},
		{"user11", "post7", (0.5 * 100000), (0.5 * 100000), 100000},
		{"user12", "post8", (0.5 * 100000), (0.5 * 100000), 100000},
	}
	for _, v := range cases {
		dp, err := rep.DonateAt(suite.ctx, v.username, v.post, types.NewCoinFromInt64(v.amount))
		suite.Nil(err)
		suite.Require().Equal(dp, types.NewCoinFromInt64(v.dp))
	}
	suite.timefies(false)
	// check reputation
	for i, v := range cases {
		rv, err := rep.GetReputation(suite.ctx, v.username)
		suite.Nil(err)
		suite.Require().Equal(types.NewCoinFromInt64(v.rep).String(), rv.String(), "%d", i)
	}

	suite.timefies(true)
	rep.DonateAt(suite.ctx, "yxia", "sp", types.NewCoinFromInt64(100000000))
	suite.timefies(true)
	rv, _ := rep.GetReputation(suite.ctx, "yxia")
	suite.Equal(types.NewCoinFromInt64(10), rv)

	dir, err2 := ioutil.TempDir("", "test")
	suite.Require().Nil(err2)
	defer os.RemoveAll(dir) // clean up

	tmpfn := filepath.Join(dir, "tmpfile")
	rep.ExportToFile(suite.ctx, tmpfn)

	// clear everything
	suite.SetupTest()
	suite.timefies(true) // start to use repv2
	rep = suite.rep

	rep.ImportFromFile(suite.ctx, tmpfn)
	rv, err := rep.GetReputation(suite.ctx, "yxia")
	suite.Nil(err)
	suite.Equal(types.NewCoinFromInt64(10).String(), rv.String())

	// check reputation of users that have not donated after update6
	for i, v := range cases {
		rv, err := rep.GetReputation(suite.ctx, v.username)
		suite.Nil(err)
		if v.rep > 100000 {
			suite.Equal(types.NewCoinFromInt64(v.rep).String(), rv.String(), "%d", i)
		} else {
			suite.Equal(types.NewCoinFromInt64(repv2.DefaultInitialReputation).String(),
				rv.String(), "%d", i)
		}
	}

}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
