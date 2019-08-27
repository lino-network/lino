package developer

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	accmn "github.com/lino-network/lino/x/account/manager"
	global "github.com/lino-network/lino/x/global"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

var (
	testInfraKVStoreKey   = sdk.NewKVStoreKey("infra")
	testAccountKVStoreKey = sdk.NewKVStoreKey("account")
	testGlobalKVStoreKey  = sdk.NewKVStoreKey("global")
	testParamKVStoreKey   = sdk.NewKVStoreKey("param")
)

// InitGlobalManager - init global manager
func InitGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest(t *testing.T, height int64) (
	sdk.Context, acc.AccountKeeper, DeveloperManager, global.GlobalManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(testParamKVStoreKey)
	ph.InitParam(ctx)
	gm := global.NewGlobalManager(testGlobalKVStoreKey, ph)
	am := accmn.NewAccountManager(testAccountKVStoreKey, ph, gm)
	dm := NewDeveloperManager(testInfraKVStoreKey, ph)
	cdc := gm.WireCodec()
	err := InitGlobalManager(ctx, gm)
	assert.Nil(t, err)
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(accmn.ReturnCoinEvent{}, "event/return", nil)
	return ctx, am, dm, gm
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(testInfraKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{Height: height}, false, log.NewNopLogger())
}

// helper function to create an account for testing purpose
func createTestAccount(ctx sdk.Context, am acc.AccountKeeper, username string, initCoin types.Coin) (secp256k1.PrivKeySecp256k1,
	secp256k1.PrivKeySecp256k1, secp256k1.PrivKeySecp256k1) {
	resetPriv := secp256k1.GenPrivKey()
	txPriv := secp256k1.GenPrivKey()
	appPriv := secp256k1.GenPrivKey()
	am.CreateAccount(ctx, types.AccountKey(username), txPriv.PubKey(), resetPriv.PubKey())
	am.AddCoinToUsername(ctx, types.AccountKey(username), initCoin)
	return resetPriv, txPriv, appPriv
}
