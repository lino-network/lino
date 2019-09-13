package reputation

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

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

func (suite *reputationTestSuite) timefies() {
	suite.t = suite.t.Add(25 * time.Hour)
	suite.height += 1
	// donate several rounds
	newctx := sdk.NewContext(
		suite.ms, abci.Header{
			ChainID: "Lino",
			Height:  suite.height,
			Time:    suite.t}, false, log.NewNopLogger())

	_ = suite.rep.Update(newctx)
	suite.ctx = newctx
}

func TestReputationTestSuite(t *testing.T) {
	suite.Run(t, &reputationTestSuite{})
}

func (suite *reputationTestSuite) SetupTest() {
	suite.t = time.Now()
	TestParamKVStoreKey := sdk.NewKVStoreKey("param")
	TestRepv2KVStoreKey := sdk.NewKVStoreKey("repv2")
	// TestReputationKVStoreKey := sdk.NewKVStoreKey("rep")
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestRepv2KVStoreKey, sdk.StoreTypeIAVL, db)
	_ = ms.LoadLatestVersion()
	ctx := sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: 1, Time: suite.t}, false, log.NewNopLogger())

	ph := param.NewParamHolder(TestParamKVStoreKey)
	_ = ph.InitParam(ctx)
	rep := NewReputationManager(TestRepv2KVStoreKey, ph)
	suite.ph = ph
	suite.rep = rep.(ReputationManager)
	suite.ctx = ctx
	suite.ms = ms
	suite.height = 1
}

func (suite *reputationTestSuite) TestGetHandlers() {
	rep := suite.rep
	v2 := rep.getHandlerV2(suite.ctx)
	v2impl := v2.(repv2.ReputationImpl)

	suite.Equal(200, v2impl.BestN)
	suite.Equal(50, v2impl.UserMaxN)
	suite.Equal(int64(25*3600), v2impl.RoundDurationSeconds)
	suite.Equal(int64(10), v2impl.SampleWindowSize)
	suite.Equal(int64(10), v2impl.DecayFactor)
}

func (suite *reputationTestSuite) TestUserPostBasicCheck() {
	suite.Nil(suite.rep.checkUsername("x"))
	suite.NotNil(suite.rep.checkUsername(""))
	suite.Nil(suite.rep.checkPost("x"))
	suite.NotNil(suite.rep.checkPost(""))
	suite.Nil(suite.rep.basicCheck("x", "y"))
	suite.NotNil(suite.rep.basicCheck("", "y"))
	suite.NotNil(suite.rep.basicCheck("x", ""))
	suite.NotNil(suite.rep.basicCheck("", ""))
}

func (suite *reputationTestSuite) TestGetCurrentRound() {
	rep := suite.rep
	ts, err := rep.GetCurrentRound(suite.ctx)
	suite.Nil(err)
	suite.Equal(int64(0), ts)

	suite.timefies()
	ts, err = rep.GetCurrentRound(suite.ctx)
	suite.Nil(err)
	suite.Equal(suite.t.Unix(), ts)

	suite.timefies()
	ts, err = rep.GetCurrentRound(suite.ctx)
	suite.Nil(err)
	suite.Equal(suite.t.Unix(), ts)
}

func (suite *reputationTestSuite) TestDonateInvalid() {
	rep := suite.rep
	// errors
	_, err := rep.DonateAt(suite.ctx, "", "", types.NewMiniDollar(3333))
	suite.NotNil(err)
	_, err = rep.DonateAt(suite.ctx, "yxia", "", types.NewMiniDollar(3333))
	suite.NotNil(err)
	_, err = rep.DonateAt(suite.ctx, "", "xxx", types.NewMiniDollar(3333))
	suite.NotNil(err)
}

func (suite *reputationTestSuite) TestDonateAtUpdate() {
	rep := suite.rep
	cases := []struct {
		username types.AccountKey
		post     types.Permlink
		amount   int64
	}{
		{"user1", "post1", (100 * 100000)},
		{"user2", "post1", (100 * 100000)},
		{"user3", "post2", (100 * 100000)},
		{"user4", "post2", (100 * 100000)},
		{"user5", "post3", (100 * 100000)},
		{"user6", "post3", (100 * 100000)},
		{"user7", "post4", (100 * 100000)},
		{"user8", "post4", (100 * 100000)},
		{"user9", "post5", (0.5 * 100000)},
		{"user10", "post6", (0.5 * 100000)},
	}
	for _, v := range cases {
		dp, err := rep.DonateAt(suite.ctx, v.username, v.post, types.NewMiniDollar(v.amount))
		suite.Nil(err)
		suite.Equal(dp, types.NewMiniDollar(repv2.DefaultInitialReputation))
	}
	suite.timefies()

	// check reputation.
	for i, v := range cases {
		rv, err := rep.GetReputation(suite.ctx, v.username)
		suite.Nil(err)
		if i < 8 {
			suite.Equal(types.NewMiniDollar(10), rv, "%d", i)
		} else {
			suite.Equal(types.NewMiniDollar(0), rv, "%d", i)
		}
	}
	suite.timefies()

	for i := 0; i < 30; i++ {
		for _, v := range cases {
			_, err := rep.DonateAt(suite.ctx, v.username, v.post, types.NewMiniDollar(v.amount))
			suite.Nil(err)
		}
		suite.timefies()
	}

	// check reputation.
	for i, v := range cases {
		rv, err := rep.GetReputation(suite.ctx, v.username)
		suite.Nil(err)
		if i < 8 {
			suite.Equal(types.NewMiniDollar(8304386).String(), rv.String(), "%d", i)
		} else {
			suite.Equal(types.NewMiniDollar(0).String(), rv.String(), "%d", i)
		}
	}
}

func (suite *reputationTestSuite) TestExportImport() {
	rep := suite.rep
	cases := []struct {
		username types.AccountKey
		post     types.Permlink
		amount   int64
	}{
		{"user1", "post1", (100 * 100000)},
		{"user2", "post1", (100 * 100000)},
		{"user3", "post2", (100 * 100000)},
		{"user4", "post2", (100 * 100000)},
		{"user5", "post3", (100 * 100000)},
		{"user6", "post3", (100 * 100000)},
		{"user7", "post4", (100 * 100000)},
		{"user8", "post4", (100 * 100000)},
		{"user9", "post5", (0.5 * 100000)},
		{"user10", "post6", (0.5 * 100000)},
	}

	for i := 0; i < 30+1; i++ {
		for _, v := range cases {
			_, err := rep.DonateAt(suite.ctx, v.username, v.post, types.NewMiniDollar(v.amount))
			suite.Nil(err)
		}
		suite.timefies()
	}

	// check reputation
	for i, v := range cases {
		rv, err := rep.GetReputation(suite.ctx, v.username)
		suite.Nil(err)
		if i < 8 {
			suite.Equal(types.NewMiniDollar(8304386).String(), rv.String(), "%d", i)
		} else {
			suite.Equal(types.NewMiniDollar(0).String(), rv.String(), "%d", i)
		}
	}

	dir, err2 := ioutil.TempDir("", "test")
	suite.Require().Nil(err2)
	defer os.RemoveAll(dir) // clean up

	tmpfn := filepath.Join(dir, "tmpfile")
	err2 = rep.ExportToFile(suite.ctx, tmpfn)
	suite.Nil(err2)

	// clear everything
	suite.SetupTest()
	suite.timefies() // start to use repv2
	rep = suite.rep

	// check reputation, should all be initial reputation
	for i, v := range cases {
		rv, err := rep.GetReputation(suite.ctx, v.username)
		suite.Nil(err)
		suite.Require().Equal(
			types.NewMiniDollar(repv2.DefaultInitialReputation).String(), rv.String(), "%d", i)
	}

	err3 := rep.ImportFromFile(suite.ctx, tmpfn)
	suite.Nil(err3)

	// check reputation
	for i, v := range cases {
		rv, err := rep.GetReputation(suite.ctx, v.username)
		suite.Nil(err)
		if i < 8 {
			suite.Equal(types.NewMiniDollar(8304386).String(), rv.String(), "%d", i)
		} else {
			suite.Equal(types.NewMiniDollar(0).String(), rv.String(), "%d", i)
		}
	}
}
