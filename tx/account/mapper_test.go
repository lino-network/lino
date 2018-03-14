package account

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"testing"
)

// Construct some global addrs and txs for tests.
var (
	KVStoreKey = sdk.NewKVStoreKey("account")
)

func newLinoAccountManager() linoAccountManager {
	return NewLinoAccountManager(KVStoreKey)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(KVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func TestAccountInfo(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	priv := crypto.GenPrivKeyEd25519()
	accInfo := types.AccountInfo{
		Username: types.AccountKey("test"),
		Created:  0,
		PostKey:  priv.PubKey(),
		OwnerKey: priv.PubKey(),
		Address:  priv.PubKey().Address(),
	}
	err := lam.SetInfo(ctx, types.AccountKey("test"), &accInfo)
	assert.Nil(t, err)

	resultPtr, err := lam.GetInfo(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accInfo, *resultPtr, "Account info should be equal")
}

func TestInvalidAccountInfo(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	resultPtr, err := lam.GetInfo(ctx, types.AccountKey("test"))
	assert.Nil(t, resultPtr)
	assert.Equal(t, err, ErrCodeAccountManagerFail("linoAccountManager get info failed: info doesn't exist"))
}

func TestAccountBank(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	priv := crypto.GenPrivKeyEd25519()
	accInfo := types.AccountInfo{
		Username: types.AccountKey("test"),
		Created:  0,
		PostKey:  priv.PubKey(),
		OwnerKey: priv.PubKey(),
		Address:  priv.PubKey().Address(),
	}
	err := lam.SetInfo(ctx, types.AccountKey("test"), &accInfo)
	assert.Nil(t, err)

	accBank := types.AccountBank{
		auth.BaseAccount{
			Address: priv.PubKey().Address(),
			PubKey:  priv.PubKey(),
			Coins:   sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
		},
	}
	err = lam.SetBank(ctx, priv.PubKey().Address(), &accBank)
	assert.Nil(t, err)

	resultPtr, err := lam.GetBankFromAccountKey(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accBank, *resultPtr, "Account bank should be equal")

	resultPtr, err = lam.GetBankFromAddress(ctx, priv.PubKey().Address())
	assert.Nil(t, err)
	assert.Equal(t, accBank, *resultPtr, "Account bank should be equal")
}

func TestAccountMeta(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	accMeta := types.AccountMeta{}
	err := lam.SetMeta(ctx, types.AccountKey("test"), &accMeta)
	assert.Nil(t, err)

	resultPtr, err := lam.GetMeta(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, accMeta, *resultPtr, "Account meta should be equal")
}

func TestAccountFollowers(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	followers := types.Followers{Followers: []types.AccountKey{}}
	err := lam.SetFollowers(ctx, types.AccountKey("test"), &followers)
	assert.Nil(t, err)

	resultPtr, err := lam.GetFollowers(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, followers, *resultPtr, "Account followers should be equal")
}

func TestAccountFollowings(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()

	followings := types.Followings{Followings: []types.AccountKey{}}
	err := lam.SetFollowings(ctx, types.AccountKey("test"), &followings)
	assert.Nil(t, err)

	resultPtr, err := lam.GetFollowings(ctx, types.AccountKey("test"))
	assert.Nil(t, err)
	assert.Equal(t, followings, *resultPtr, "Account followers should be equal")
}
