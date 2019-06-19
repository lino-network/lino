package reputation

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/x/reputation/repv2"
	// "github.com/lino-network/lino/types"
)

type reputationTestSuite struct {
	suite.Suite
	ph  param.ParamHolder
	rep ReputationManager
	ctx sdk.Context
	t   time.Time
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
}

func (suite *reputationTestSuite) TestGetHandlers() {
	rep := suite.rep
	_, err := rep.getHandler(suite.ctx)
	suite.Nil(err)
	v2, err := rep.getHandlerV2(suite.ctx)
	suite.Nil(err)

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
	v2, _ := rep.getHandlerV2(suite.ctx)
	suite.Require().NotNil(v1)
	suite.Require().NotNil(v2)

	v1.IncFreeScore("yxia", big.NewInt(100000000))
	rep.migrate(v1, v2, "yxia")
	suite.Equal(big.NewInt(100100000), v2.GetReputation("yxia"))
}

func (suite *reputationTestSuite) TestDonateAtUpdate() {
	// TODO
}
