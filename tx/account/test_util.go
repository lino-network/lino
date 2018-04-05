package account

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestKVStoreKey = sdk.NewKVStoreKey("account")
)

func newLinoAccountManager() *AccountManager {
	return NewAccountManager(TestKVStoreKey)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func createTestAccount(ctx sdk.Context, am *AccountManager, username string) crypto.PrivKey {
	priv := crypto.GenPrivKeyEd25519()
	am.AddCoinToAddress(ctx, priv.PubKey().Address(), types.NewCoin(0))
	am.CreateAccount(ctx, types.AccountKey(username), priv.PubKey(), types.NewCoin(0))
	return priv.Wrap()
}
