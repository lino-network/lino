package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("vote")
)

func setup(t *testing.T) (sdk.Context, VoteStorage) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	vs := NewVoteStorage(TestKVStoreKey)
	return ctx, vs
}

func TestDoesVoterExist(t *testing.T) {
	ctx, vs := setup(t)
	user := linotypes.AccountKey("user")
	voter := &Voter{
		Username:  user,
		LinoStake: linotypes.NewCoinFromInt64(1000),
	}
	vs.SetVoter(ctx, user, voter)

	testCases := []struct {
		testName  string
		accKey    linotypes.AccountKey
		wantExist bool
	}{
		{
			testName:  "voter exist",
			accKey:    user,
			wantExist: true,
		},
		{
			testName:  "voter doesn't exist",
			accKey:    linotypes.AccountKey("acc"),
			wantExist: false,
		},
	}
	for _, tc := range testCases {
		gotExist := vs.DoesVoterExist(ctx, tc.accKey)
		if gotExist != tc.wantExist {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, gotExist, tc.wantExist)
		}
	}
}

func TestVoter(t *testing.T) {
	ctx, vs := setup(t)

	user := linotypes.AccountKey("user")
	voter := Voter{
		Username:          user,
		LinoStake:         linotypes.NewCoinFromInt64(1000),
		LastPowerChangeAt: 0,
		Duty:              types.DutyValidator,
		Interest:          linotypes.NewCoinFromInt64(0),
		FrozenAmount:      linotypes.NewCoinFromInt64(10),
	}
	vs.SetVoter(ctx, user, &voter)

	voterPtr, err := vs.GetVoter(ctx, user)
	assert.Nil(t, err)
	assert.Equal(t, voter, *voterPtr, "voter should be equal")
}
