package account

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestKVStoreKey = sdk.NewKVStoreKey("account")
)

func newLinoAccountManager() AccountManager {
	return NewLinoAccountManager(TestKVStoreKey)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

// helper function to create an account for testing purpose
func privAndBank() (crypto.PrivKey, *AccountBank) {
	priv := crypto.GenPrivKeyEd25519()
	accBank := &AccountBank{
		Address: priv.PubKey().Address(),
	}
	return priv.Wrap(), accBank
}

func createTestAccount(ctx sdk.Context, lam AccountManager, username string) *AccountProxy {
	priv, bank := privAndBank()
	acc := NewAccountProxy(AccountKey(username), &lam)
	acc.CreateAccount(ctx, AccountKey(username), priv.PubKey(), bank)
	acc.Apply(ctx)
	return acc
}
