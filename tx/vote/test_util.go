package vote

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestAccountKVStoreKey = sdk.NewKVStoreKey("account")
	TestVoteKVStoreKey    = sdk.NewKVStoreKey("vote")
	TestGlobalKVStoreKey  = sdk.NewKVStoreKey("global")

	initCoin = types.NewCoin(100)
)

func InitGlobalManager(ctx sdk.Context, gm *global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoin(10000*types.Decimals))
}

func setupTest(t *testing.T, height int64) (sdk.Context, *acc.AccountManager, *VoteManager, *global.GlobalManager) {
	ctx := getContext(height)
	accManager := acc.NewAccountManager(TestAccountKVStoreKey)
	postManager := NewVoteManager(TestVoteKVStoreKey)
	globalManager := global.NewGlobalManager(TestGlobalKVStoreKey)
	err := InitGlobalManager(ctx, globalManager)
	assert.Nil(t, err)
	return ctx, accManager, postManager, globalManager
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestVoteKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{Height: height}, false, nil)
}

// helper function to create an account for testing purpose
func createTestAccount(ctx sdk.Context, am *acc.AccountManager, username string) types.AccountKey {
	priv := crypto.GenPrivKeyEd25519()
	am.AddCoinToAddress(ctx, priv.PubKey().Address(), initCoin)
	am.CreateAccount(ctx, types.AccountKey(username), priv.PubKey(), types.NewCoin(0))
	return types.AccountKey(username)
}
