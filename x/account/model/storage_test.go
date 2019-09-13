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
	dbm "github.com/tendermint/tm-db"
)

var (
	TestKVStoreKey = sdk.NewKVStoreKey("account")
)

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	return sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
}

func TestAccountInfo(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	accInfo := AccountInfo{
		Username:       types.AccountKey("test"),
		CreatedAt:      0,
		SigningKey:     secp256k1.GenPrivKey().PubKey(),
		TransactionKey: secp256k1.GenPrivKey().PubKey(),
		Address:        sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
	}
	as.SetInfo(ctx, types.AccountKey("test"), &accInfo)

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

	pubKey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubKey.Address())
	accBank := AccountBank{
		Saving: types.NewCoinFromInt64(int64(123)),
	}
	as.SetBank(ctx, addr, &accBank)

	resultPtr, err := as.GetBank(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, accBank, *resultPtr, "Account bank should be equal")

	accBank.PubKey = pubKey

	as.SetBank(ctx, addr, &accBank)

	resultPtr, err = as.GetBank(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, accBank, *resultPtr, "Account bank should be equal")

}

func TestAccountMeta(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	accMeta := AccountMeta{JSONMeta: "{'test':1}"}
	as.SetMeta(ctx, types.AccountKey("test"), &accMeta)

	resultPtr := as.GetMeta(ctx, types.AccountKey("test"))
	assert.Equal(t, accMeta, *resultPtr, "Account meta should be equal")
}

func TestAccountGrantPubkey(t *testing.T) {
	as := NewAccountStorage(TestKVStoreKey)
	ctx := getContext()

	grantPubKey := GrantPermission{GrantTo: types.AccountKey("grantTo"), Permission: types.AppPermission, Amount: types.NewCoinFromInt64(0)}
	grantPubKey2 := GrantPermission{GrantTo: types.AccountKey("grantTo"), Permission: types.PreAuthorizationPermission, Amount: types.NewCoinFromInt64(10)}
	as.SetGrantPermissions(ctx, types.AccountKey("test"), types.AccountKey("grantTo"), []*GrantPermission{&grantPubKey, &grantPubKey2})

	resultList, err := as.GetGrantPermissions(ctx, types.AccountKey("test"), types.AccountKey("grantTo"))
	assert.Nil(t, err)
	assert.Equal(t, []*GrantPermission{&grantPubKey, &grantPubKey2}, resultList, "Account grant user should be equal")

	resultList, err = as.GetAllGrantPermissions(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, []*GrantPermission{&grantPubKey, &grantPubKey2}, resultList, "Account grant user should be equal")

	as.DeleteAllGrantPermissions(ctx, types.AccountKey("test"), types.AccountKey("grantTo"))
	resultList, err = as.GetGrantPermissions(ctx, types.AccountKey("test"), types.AccountKey("grantTo"))
	assert.NotNil(t, err)
	assert.Nil(t, resultList)
}
