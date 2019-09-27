package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("validator")
)

func setup(t *testing.T) (sdk.Context, ValidatorStorage) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	vs := NewValidatorStorage(TestKVStoreKey)
	err = vs.InitGenesis(ctx)
	assert.Nil(t, err)
	return ctx, vs
}

func TestValidator(t *testing.T) {
	ctx, vs := setup(t)

	priv := secp256k1.GenPrivKey()
	testCases := []struct {
		testName string
		power    int64
		user     types.AccountKey
		votes    types.Coin
	}{
		{
			testName: "user as validator",
			power:    1000,
			user:     types.AccountKey("user"),
			votes:    types.NewCoinFromInt64(10),
		},
		{
			testName: "user as validator again",
			power:    10000,
			user:     types.AccountKey("user"),
			votes:    types.NewCoinFromInt64(10),
		},
		{
			testName: "user2 as validator",
			power:    1,
			user:     types.AccountKey("user2"),
			votes:    types.NewCoinFromInt64(10),
		},
	}

	for _, tc := range testCases {
		validator := Validator{
			ABCIValidator: abci.Validator{
				Address: priv.PubKey().Address(),
				Power:   1000},
			Username:      tc.user,
			ReceivedVotes: tc.votes,
		}
		err := vs.SetValidator(ctx, tc.user, &validator)
		if err != nil {
			t.Errorf("%s: failed to set validator, got err %v", tc.testName, err)
		}

		valPtr, err := vs.GetValidator(ctx, tc.user)
		if err != nil {
			t.Errorf("%s: failed to get validator, got err %v", tc.testName, err)
		}
		if !assert.Equal(t, validator, *valPtr) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, *valPtr, validator)
		}
	}
}

func TestValidatorList(t *testing.T) {
	ctx, vs := setup(t)

	testCases := []struct {
		testName string
		valList  ValidatorList
	}{
		{
			testName: "normal case",
			valList: ValidatorList{
				Oncall: []types.AccountKey{
					types.AccountKey("user1"),
				},
				Standby: []types.AccountKey{
					types.AccountKey("user2"),
				},
				Candidates: []types.AccountKey{
					types.AccountKey("user2"),
				},
				Jail: []types.AccountKey{
					types.AccountKey("user2"),
				},
				PreBlockValidators: []types.AccountKey{
					types.AccountKey("user2"),
				},
				LowestOncallVotes:  types.NewCoinFromInt64(100),
				LowestOncall:       types.AccountKey("user2"),
				LowestStandbyVotes: types.NewCoinFromInt64(100),
				LowestStandby:      types.AccountKey("user2"),
			},
		},
	}

	for _, tc := range testCases {
		err := vs.SetValidatorList(ctx, &tc.valList)
		if err != nil {
			t.Errorf("%s: failed to set validator list, got err %v", tc.testName, err)
		}

		valListPtr, err := vs.GetValidatorList(ctx)
		if err != nil {
			t.Errorf("%s: failed to get validator list, got err %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.valList, *valListPtr) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, *valListPtr, tc.valList)
		}
	}
}

func TestElectionVoteList(t *testing.T) {
	ctx, vs := setup(t)

	testCases := []struct {
		testName string
		lst      ElectionVoteList
		user     types.AccountKey
	}{
		{
			testName: "normal case",
			user:     types.AccountKey("user"),
			lst: ElectionVoteList{
				ElectionVotes: []ElectionVote{
					{
						ValidatorName: types.AccountKey("test"),
						Vote:          types.NewCoinFromInt64(100),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		err := vs.SetElectionVoteList(ctx, tc.user, &tc.lst)
		if err != nil {
			t.Errorf("%s: failed to set election list, got err %v", tc.testName, err)
		}

		lstPtr, err := vs.GetElectionVoteList(ctx, tc.user)
		if err != nil {
			t.Errorf("%s: failed to get election list, got err %v", tc.testName, err)
		}
		if !assert.Equal(t, tc.lst, *lstPtr) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, *lstPtr, tc.lst)
		}
	}
}
