package account

import (
	"testing"
	"time"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"
	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

var (
	TestAccountKVStoreKey = sdk.NewKVStoreKey("account")
	TestGlobalKVStoreKey  = sdk.NewKVStoreKey("global")
	TestParamKVStoreKey   = sdk.NewKVStoreKey("param")

	accountReferrer = types.AccountKey("referrer")

	l0    = types.LNO("0")
	l100  = types.LNO("100")
	l200  = types.LNO("200")
	l1600 = types.LNO("1600")
	l1800 = types.LNO("1800")
	l1900 = types.LNO("1900")
	l1999 = types.LNO("1999")
	l2000 = types.LNO("2000")
	c0    = types.NewCoinFromInt64(0)
	c100  = types.NewCoinFromInt64(100 * types.Decimals)
	c200  = types.NewCoinFromInt64(200 * types.Decimals)
	c300  = types.NewCoinFromInt64(300 * types.Decimals)
	c400  = types.NewCoinFromInt64(400 * types.Decimals)
	c500  = types.NewCoinFromInt64(500 * types.Decimals)
	c600  = types.NewCoinFromInt64(600 * types.Decimals)
	c1000 = types.NewCoinFromInt64(1000 * types.Decimals)
	c1500 = types.NewCoinFromInt64(1500 * types.Decimals)
	c1600 = types.NewCoinFromInt64(1600 * types.Decimals)
	c1800 = types.NewCoinFromInt64(1800 * types.Decimals)
	c1900 = types.NewCoinFromInt64(1900 * types.Decimals)
	c2000 = types.NewCoinFromInt64(2000 * types.Decimals)

	coin0   = types.NewCoinFromInt64(0)
	coin1   = types.NewCoinFromInt64(1)
	coin2   = types.NewCoinFromInt64(2)
	coin3   = types.NewCoinFromInt64(3)
	coin4   = types.NewCoinFromInt64(4)
	coin50  = types.NewCoinFromInt64(50)
	coin100 = types.NewCoinFromInt64(100)
	coin200 = types.NewCoinFromInt64(200)
	coin300 = types.NewCoinFromInt64(300)
	coin400 = types.NewCoinFromInt64(400)
)

func InitGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest(t *testing.T, height int64) (sdk.Context, AccountManager, global.GlobalManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(TestParamKVStoreKey)
	ph.InitParam(ctx)
	accManager := NewAccountManager(TestAccountKVStoreKey, ph)
	globalManager := global.NewGlobalManager(TestGlobalKVStoreKey, ph)

	cdc := globalManager.WireCodec()
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(ReturnCoinEvent{}, "event/return", nil)

	err := InitGlobalManager(ctx, globalManager)
	assert.Nil(t, err)
	return ctx, accManager, globalManager
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: height, Time: time.Now().Unix()},
		false, log.NewNopLogger())
}

func createTestAccount(ctx sdk.Context, am AccountManager, username string) (secp256k1.PrivKeySecp256k1,
	secp256k1.PrivKeySecp256k1, secp256k1.PrivKeySecp256k1) {
	resetPriv := secp256k1.GenPrivKey()
	txPriv := secp256k1.GenPrivKey()
	appPriv := secp256k1.GenPrivKey()

	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	am.CreateAccount(ctx, accountReferrer, types.AccountKey(username),
		resetPriv.PubKey(), txPriv.PubKey(), appPriv.PubKey(), accParam.RegisterFee)
	return resetPriv, txPriv, appPriv
}
