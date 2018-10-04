package vote

import (
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/recorder"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	rep "github.com/lino-network/lino/x/reputation"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

// Construct some global addrs and txs for tests.
var (
	testAccountKVStoreKey = sdk.NewKVStoreKey("account")
	testVoteKVStoreKey    = sdk.NewKVStoreKey("vote")
	testGlobalKVStoreKey  = sdk.NewKVStoreKey("global")
	testParamKVStoreKey   = sdk.NewKVStoreKey("param")
	testRepKVStoreKey     = sdk.NewKVStoreKey("reputation")

	c100 = types.NewCoinFromInt64(100 * types.Decimals)
	c500 = types.NewCoinFromInt64(500 * types.Decimals)
)

func initGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest(t *testing.T, height int64) (sdk.Context,
	acc.AccountManager, VoteManager, global.GlobalManager, rep.ReputationManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(testParamKVStoreKey)
	ph.InitParam(ctx)
	recorder := recorder.NewRecorder()
	accManager := acc.NewAccountManager(testAccountKVStoreKey, ph)
	voteManager := NewVoteManager(testVoteKVStoreKey, ph)
	globalManager := global.NewGlobalManager(testGlobalKVStoreKey, ph, recorder)
	repManager := rep.NewReputationManager(testRepKVStoreKey, ph)

	cdc := globalManager.WireCodec()
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(acc.ReturnCoinEvent{}, "1", nil)

	err := initGlobalManager(ctx, globalManager)
	assert.Nil(t, err)
	return ctx, accManager, voteManager, globalManager, repManager
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(testAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testVoteKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testRepKVStoreKey, sdk.StoreTypeIAVL, db)

	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{Height: height, Time: time.Unix(0, 0)}, false, log.NewNopLogger())
}

// helper function to create an account for testing purpose
func createTestAccount(ctx sdk.Context, am acc.AccountManager, username string, initCoin types.Coin) types.AccountKey {
	am.CreateAccount(ctx, "referrer", types.AccountKey(username),
		secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(), initCoin)
	return types.AccountKey(username)
}

func coinToString(coin types.Coin) string {
	coinInInt64, _ := coin.ToInt64()
	return strconv.FormatInt(coinInInt64/types.Decimals, 10)
}
