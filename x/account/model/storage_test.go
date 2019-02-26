package model

import (
	"testing"

	"github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("account")
)

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
}

func TestAccountInfo(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	accInfo := AccountInfo{
		Username:       types.AccountKey("test"),
		CreatedAt:      0,
		ResetKey:       secp256k1.GenPrivKey().PubKey(),
		TransactionKey: secp256k1.GenPrivKey().PubKey(),
		AppKey:         secp256k1.GenPrivKey().PubKey(),
	}
	err := as.SetInfo(ctx, types.AccountKey("test"), &accInfo)
	assert.Nil(t, err)

	resultPtr, err := as.GetInfo(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accInfo, *resultPtr, "Account info should be equal")
}

func TestInvalidAccountInfo(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	resultPtr, err := as.GetInfo(ctx, types.AccountKey("test"))
	assert.Nil(t, resultPtr)
	assert.Equal(t, err, ErrAccountInfoNotFound())
}

func TestAccountBank(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	accBank := AccountBank{
		Saving:  types.NewCoinFromInt64(int64(123)),
		CoinDay: types.NewCoinFromInt64(0),
	}
	err := as.SetBankFromAccountKey(ctx, types.AccountKey("test"), &accBank)
	assert.Nil(t, err)

	resultPtr, err := as.GetBankFromAccountKey(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accBank, *resultPtr, "Account bank should be equal")
}

func TestAccountMeta(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	accMeta := AccountMeta{TransactionCapacity: types.NewCoinFromInt64(0)}
	err := as.SetMeta(ctx, types.AccountKey("test"), &accMeta)
	assert.Nil(t, err)

	resultPtr, err := as.GetMeta(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accMeta, *resultPtr, "Account meta should be equal")
}

func TestAccountReward(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	reward := Reward{
		TotalIncome:     types.NewCoinFromInt64(0),
		OriginalIncome:  types.NewCoinFromInt64(0),
		FrictionIncome:  types.NewCoinFromInt64(0),
		InflationIncome: types.NewCoinFromInt64(0),
		UnclaimReward:   types.NewCoinFromInt64(0),
	}
	err := as.SetReward(ctx, types.AccountKey("test"), &reward)
	assert.Nil(t, err)

	resultPtr, err := as.GetReward(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, reward, *resultPtr, "Account reward should be equal")
}

func TestAccountGrantPubkey(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()
	priv := secp256k1.GenPrivKey()

	grantPubKey := GrantPubKey{Amount: types.NewCoinFromInt64(0)}
	err := as.SetGrantPubKey(ctx, types.AccountKey("test"), priv.PubKey(), &grantPubKey)
	assert.Nil(t, err)

	resultPtr, err := as.GetGrantPubKey(ctx, types.AccountKey("test"), priv.PubKey())
	assert.Nil(t, err)
	assert.Equal(t, grantPubKey, *resultPtr, "Account grant user should be equal")

	as.DeleteGrantPubKey(ctx, types.AccountKey("test"), priv.PubKey())
	resultPtr, err = as.GetGrantPubKey(ctx, types.AccountKey("test"), priv.PubKey())
	assert.NotNil(t, err)
	assert.Nil(t, resultPtr)
}
