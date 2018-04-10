package account

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"time"
)

// Construct some global addrs and txs for tests.
var (
	TestAccountKVStoreKey = sdk.NewKVStoreKey("account")
)

func setupTest(t *testing.T, height int64) (sdk.Context, *AccountManager) {
	ctx := getContext(height)
	accManager := NewAccountManager(TestAccountKVStoreKey)
	return ctx, accManager
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{ChainID: "Lino", Height: height, Time: time.Now().Unix()}, false, nil)
}

func createTestAccount(ctx sdk.Context, am *AccountManager, username string) crypto.PrivKey {
	priv := crypto.GenPrivKeyEd25519()
	am.AddCoinToAddress(ctx, priv.PubKey().Address(), types.NewCoin(0))
	am.CreateAccount(ctx, types.AccountKey(username), priv.PubKey(), types.NewCoin(0))
	return priv.Wrap()
}
