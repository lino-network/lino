package account

import (
	"testing"
	"time"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tmlibs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestAccountKVStoreKey = sdk.NewKVStoreKey("account")
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

func setupTest(t *testing.T, height int64) (sdk.Context, AccountManager, param.AccountParam) {
	ctx := getContext(height)
	ph := param.NewParamHolder(TestParamKVStoreKey)
	ph.InitParam(ctx)
	accManager := NewAccountManager(TestAccountKVStoreKey, ph)
	accParam, _ := accManager.paramHolder.GetAccountParam(ctx)
	return ctx, accManager, *accParam
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: height, Time: time.Now().Unix()},
		false, nil, log.NewNopLogger())
}

func createTestAccount(ctx sdk.Context, am AccountManager, username string) crypto.PrivKeyEd25519 {
	priv := crypto.GenPrivKeyEd25519()
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	am.CreateAccount(ctx, accountReferrer, types.AccountKey(username),
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), accParam.RegisterFee)
	return priv
}
