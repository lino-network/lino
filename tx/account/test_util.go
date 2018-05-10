package account

import (
	"testing"
	"time"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestAccountKVStoreKey = sdk.NewKVStoreKey("account")
	TestParamKVStoreKey   = sdk.NewKVStoreKey("param")

	l0    = types.LNO("0")
	l100  = types.LNO("100")
	l200  = types.LNO("200")
	l1600 = types.LNO("1600")
	l1800 = types.LNO("1800")
	l1900 = types.LNO("1900")
	l1999 = types.LNO("1999")
	l2000 = types.LNO("2000")
	c0    = types.NewCoin(0)
	c100  = types.NewCoin(100 * types.Decimals)
	c200  = types.NewCoin(200 * types.Decimals)
	c300  = types.NewCoin(300 * types.Decimals)
	c500  = types.NewCoin(500 * types.Decimals)
	c600  = types.NewCoin(600 * types.Decimals)
	c1000 = types.NewCoin(1000 * types.Decimals)
	c1500 = types.NewCoin(1500 * types.Decimals)
	c1600 = types.NewCoin(1600 * types.Decimals)
	c1800 = types.NewCoin(1800 * types.Decimals)
	c1900 = types.NewCoin(1900 * types.Decimals)
	c2000 = types.NewCoin(2000 * types.Decimals)

	coin0   = types.NewCoin(0)
	coin1   = types.NewCoin(1)
	coin50  = types.NewCoin(50)
	coin100 = types.NewCoin(100)
	coin200 = types.NewCoin(200)
	coin300 = types.NewCoin(300)
	coin400 = types.NewCoin(400)
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

	return sdk.NewContext(ms, abci.Header{ChainID: "Lino", Height: height, Time: time.Now().Unix()}, false, nil)
}

func createTestAccount(ctx sdk.Context, am AccountManager, username string) crypto.PrivKeyEd25519 {
	priv := crypto.GenPrivKeyEd25519()
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	am.AddSavingCoinToAddress(ctx, priv.PubKey().Address(), accParam.RegisterFee)
	am.CreateAccount(ctx, types.AccountKey(username),
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	return priv
}
