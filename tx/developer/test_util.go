package developer

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestInfraKVStoreKey   = sdk.NewKVStoreKey("infra")
	TestAccountKVStoreKey = sdk.NewKVStoreKey("account")
	initCoin              = types.NewCoin(100)
)

func setupTest(t *testing.T, height int64) (sdk.Context, *acc.AccountManager, *DeveloperManager) {
	ctx := getContext(height)
	am := acc.NewAccountManager(TestAccountKVStoreKey)
	dm := NewDeveloperManager(TestInfraKVStoreKey)
	return ctx, am, dm
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestInfraKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{Height: height}, false, nil)
}

func createTestAccount(t *testing.T, ctx sdk.Context, am *acc.AccountManager, username string) types.AccountKey {
	priv := crypto.GenPrivKeyEd25519()
	err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), initCoin)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, types.AccountKey(username), priv.PubKey(), types.NewCoin(0))
	assert.Nil(t, err)
	return types.AccountKey(username)
}
