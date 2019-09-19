package manager

import (
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	accmn "github.com/lino-network/lino/x/account/manager"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

// Construct some global addrs and txs for tests.
var (
	testAccountKVStoreKey = sdk.NewKVStoreKey("account")
	testVoteKVStoreKey    = sdk.NewKVStoreKey("vote")
	testGlobalKVStoreKey  = sdk.NewKVStoreKey("global")
	testParamKVStoreKey   = sdk.NewKVStoreKey("param")

	c100 = types.NewCoinFromInt64(100 * types.Decimals)
	c500 = types.NewCoinFromInt64(500 * types.Decimals)
)

func initGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest(t *testing.T, height int64) (sdk.Context,
	acc.AccountKeeper, VoteManager, global.GlobalManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(testParamKVStoreKey)
	err := ph.InitParam(ctx)
	if err != nil {
		panic(err)
	}
	gm := global.NewGlobalManager(testGlobalKVStoreKey, ph)
	accManager := accmn.NewAccountManager(testAccountKVStoreKey, ph, &gm)
	voteManager := NewVoteManager(testVoteKVStoreKey, ph, accManager, &gm)

	cdc := gm.WireCodec()
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(accmn.ReturnCoinEvent{}, "1", nil)

	err = initGlobalManager(ctx, gm)
	assert.Nil(t, err)
	return ctx, accManager, voteManager, gm
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(testAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testVoteKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testParamKVStoreKey, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	return sdk.NewContext(ms, abci.Header{Height: height, Time: time.Unix(0, 0)}, false, log.NewNopLogger())
}

// helper function to create an account for testing purpose
func createTestAccount(ctx sdk.Context, am acc.AccountKeeper, username string, initCoin types.Coin) types.AccountKey {
	err := am.CreateAccount(
		ctx, types.AccountKey(username), secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey())
	if err != nil {
		panic(err)
	}
	err = am.AddCoinToUsername(ctx, types.AccountKey(username), initCoin)
	if err != nil {
		panic(err)
	}
	return types.AccountKey(username)
}

func coinToString(coin types.Coin) string {
	coinInInt64, _ := coin.ToInt64()
	return strconv.FormatInt(coinInInt64/types.Decimals, 10)
}
