package model

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("validator")
)

func setup(t *testing.T) (sdk.Context, ValidatorStorage) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)
	vs := NewValidatorStorage(TestKVStoreKey)
	err := vs.InitGenesis(ctx)
	assert.Nil(t, err)
	return ctx, vs
}

func TestValidator(t *testing.T) {
	ctx, vs := setup(t)

	priv := crypto.GenPrivKeyEd25519()
	cases := []struct {
		power   int64
		user    types.AccountKey
		deposit types.Coin
	}{
		{1000, types.AccountKey("user"), types.NewCoin(100)},
		{10000, types.AccountKey("user"), types.NewCoin(0)},
		{1, types.AccountKey("user2"), types.NewCoin(10)},
	}

	for _, cs := range cases {
		validator := Validator{
			ABCIValidator: abci.Validator{PubKey: priv.PubKey().Bytes(), Power: cs.power},
			Username:      cs.user,
			Deposit:       cs.deposit,
		}
		err := vs.SetValidator(ctx, cs.user, &validator)
		assert.Nil(t, err)
		valPtr, err := vs.GetValidator(ctx, cs.user)
		assert.Nil(t, err)
		assert.Equal(t, validator, *valPtr)
	}
}

func TestValidatorList(t *testing.T) {
	ctx, vs := setup(t)

	cases := []struct {
		ValidatorList
	}{
		{ValidatorList{[]types.AccountKey{types.AccountKey("user1")},
			[]types.AccountKey{types.AccountKey("user2")},
			nil,
			types.NewCoin(100), types.AccountKey("user2")}},
	}

	for _, cs := range cases {
		err := vs.SetValidatorList(ctx, &cs.ValidatorList)
		assert.Nil(t, err)
		valListPtr, err := vs.GetValidatorList(ctx)
		assert.Nil(t, err)
		assert.Equal(t, cs.ValidatorList, *valListPtr)
	}
}
