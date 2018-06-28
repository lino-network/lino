package developer

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tmlibs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	global "github.com/lino-network/lino/x/global"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestInfraKVStoreKey   = sdk.NewKVStoreKey("infra")
	TestAccountKVStoreKey = sdk.NewKVStoreKey("account")
	TestGlobalKVStoreKey  = sdk.NewKVStoreKey("global")
	TestParamKVStoreKey   = sdk.NewKVStoreKey("param")
)

func InitGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest(t *testing.T, height int64) (
	sdk.Context, acc.AccountManager, DeveloperManager, global.GlobalManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(TestParamKVStoreKey)
	ph.InitParam(ctx)
	am := acc.NewAccountManager(TestAccountKVStoreKey, ph)
	dm := NewDeveloperManager(TestInfraKVStoreKey, ph)
	gm := global.NewGlobalManager(TestGlobalKVStoreKey, ph)
	cdc := gm.WireCodec()
	err := InitGlobalManager(ctx, gm)
	assert.Nil(t, err)
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
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{Height: height}, false, nil, log.NewNopLogger())
}

// helper function to create an account for testing purpose
func createTestAccount(ctx sdk.Context, am acc.AccountManager, username string, initCoin types.Coin) crypto.PrivKeyEd25519 {
	priv := crypto.GenPrivKeyEd25519()
	am.CreateAccount(ctx, "referrer", types.AccountKey(username),
		priv.PubKey(), priv.Generate(0).PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), initCoin)
	return priv
}
