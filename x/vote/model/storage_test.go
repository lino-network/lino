package model

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/lino-network/lino/testsuites"
	"github.com/lino-network/lino/testutils"
	linotypes "github.com/lino-network/lino/types"

	"github.com/lino-network/lino/x/vote/types"
)

var (
	storeKeyStr = "testVoterStore"
	kvStoreKey  = sdk.NewKVStoreKey(storeKeyStr)
)

type VoteStoreDumper struct{}

func (dumper VoteStoreDumper) NewDumper() *testutils.Dumper {
	return NewVoteDumper(NewVoteStorage(kvStoreKey))
}

type voteStoreTestSuite struct {
	testsuites.GoldenTestSuite
	store VoteStorage
}

func NewVoteStoreTestSuite() *voteStoreTestSuite {
	return &voteStoreTestSuite{
		GoldenTestSuite: testsuites.NewGoldenTestSuite(VoteStoreDumper{}, kvStoreKey),
	}
}

func (suite *voteStoreTestSuite) SetupTest() {
	suite.SetupCtx(0, time.Unix(0, 0), kvStoreKey)
	suite.store = NewVoteStorage(kvStoreKey)
}

func TestVoteStoreSuite(t *testing.T) {
	suite.Run(t, NewVoteStoreTestSuite())
}

func (suite *voteStoreTestSuite) TestGetSetVoter() {
	store := suite.store
	ctx := suite.Ctx
	user1 := linotypes.AccountKey("user1")
	user2 := linotypes.AccountKey("user2")
	voter1 := Voter{
		Username:          user1,
		LinoStake:         linotypes.NewCoinFromInt64(123),
		LastPowerChangeAt: 777,
		Interest:          linotypes.NewCoinFromInt64(234),
		Duty:              types.DutyValidator,
		FrozenAmount:      linotypes.NewCoinFromInt64(9),
	}
	voter2 := Voter{
		Username:          user2,
		LinoStake:         linotypes.NewCoinFromInt64(345),
		LastPowerChangeAt: 888,
		Interest:          linotypes.NewCoinFromInt64(456),
		Duty:              types.DutyValidator,
		FrozenAmount:      linotypes.NewCoinFromInt64(12),
	}

	suite.False(store.DoesVoterExist(ctx, user1))
	_, err := store.GetVoter(ctx, user1)
	suite.Equal(types.ErrVoterNotFound(), err)

	suite.store.SetVoter(ctx, &voter1)
	suite.store.SetVoter(ctx, &voter2)

	v1, err := store.GetVoter(ctx, user1)
	suite.Nil(err)
	suite.Equal(&voter1, v1)

	v2, err := store.GetVoter(ctx, user2)
	suite.Nil(err)
	suite.Equal(&voter2, v2)

	suite.Golden()
}

func (suite *voteStoreTestSuite) TestGetSetLinoStakeStats() {
	store := suite.store
	ctx := suite.Ctx
	stats1 := LinoStakeStat{
		TotalConsumptionFriction: linotypes.NewCoinFromInt64(123),
		UnclaimedFriction:        linotypes.NewCoinFromInt64(456),
		TotalLinoStake:           linotypes.NewCoinFromInt64(789),
		UnclaimedLinoStake:       linotypes.NewCoinFromInt64(999),
	}
	stats2 := LinoStakeStat{
		TotalConsumptionFriction: linotypes.NewCoinFromInt64(888),
		UnclaimedFriction:        linotypes.NewCoinFromInt64(888),
		TotalLinoStake:           linotypes.NewCoinFromInt64(888),
		UnclaimedLinoStake:       linotypes.NewCoinFromInt64(888),
	}

	_, err := store.GetLinoStakeStat(ctx, 1)
	suite.Equal(types.ErrStakeStatNotFound(1), err)

	suite.store.SetLinoStakeStat(ctx, 1, &stats1)
	suite.store.SetLinoStakeStat(ctx, 2, &stats2)

	v1, err := store.GetLinoStakeStat(ctx, 1)
	suite.Nil(err)
	suite.Equal(&stats1, v1)

	v2, err := store.GetLinoStakeStat(ctx, 2)
	suite.Nil(err)
	suite.Equal(&stats2, v2)

	suite.Golden()
}
