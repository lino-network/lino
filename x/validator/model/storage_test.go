package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("validator")
)

func setup(t *testing.T) (sdk.Context, ValidatorStorage) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	vs := NewValidatorStorage(TestKVStoreKey)
	err := vs.InitGenesis(ctx)
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
		deposit  types.Coin
	}{
		{
			testName: "user as validator",
			power:    1000,
			user:     types.AccountKey("user"),
			deposit:  types.NewCoinFromInt64(100),
		},
		{
			testName: "user as validator again",
			power:    10000,
			user:     types.AccountKey("user"),
			deposit:  types.NewCoinFromInt64(0),
		},
		{
			testName: "user2 as validator",
			power:    1,
			user:     types.AccountKey("user2"),
			deposit:  types.NewCoinFromInt64(10),
		},
	}

	for _, tc := range testCases {
		validator := Validator{
			ABCIValidator: abci.Validator{
				Address: priv.PubKey().Address(),
				PubKey:  tmtypes.TM2PB.PubKey(priv.PubKey()),
				Power:   1000},
			Username: tc.user,
			Deposit:  tc.deposit,
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
				OncallValidators: []types.AccountKey{
					types.AccountKey("user1"),
				},
				AllValidators: []types.AccountKey{
					types.AccountKey("user2"),
				},
				LowestPower:     types.NewCoinFromInt64(100),
				LowestValidator: types.AccountKey("user2"),
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
