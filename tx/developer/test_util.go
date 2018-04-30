package developer

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	global "github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestInfraKVStoreKey   = sdk.NewKVStoreKey("infra")
	TestAccountKVStoreKey = sdk.NewKVStoreKey("account")
	TestGlobalKVStoreKey  = sdk.NewKVStoreKey("global")

	initCoin = types.NewCoin(100)
)

func setupTest(t *testing.T, height int64) (
	sdk.Context, acc.AccountManager, DeveloperManager, global.GlobalManager) {
	ctx := getContext(height)
	am := acc.NewAccountManager(TestAccountKVStoreKey)
	dm := NewDeveloperManager(TestInfraKVStoreKey)
	gm := global.NewGlobalManager(TestGlobalKVStoreKey)
	cdc := gm.WireCodec()
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(acc.ReturnCoinEvent{}, "event/return", nil)
	return ctx, am, dm, gm
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestInfraKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{Height: height}, false, nil)
}

// helper function to create an account for testing purpose
func createTestAccount(ctx sdk.Context, am acc.AccountManager, username string) types.AccountKey {
	priv := crypto.GenPrivKeyEd25519()
	am.AddCoinToAddress(ctx, priv.PubKey().Address(), initCoin)
	am.CreateAccount(ctx, types.AccountKey(username), priv.PubKey(), types.NewCoin(0))
	return types.AccountKey(username)
}
