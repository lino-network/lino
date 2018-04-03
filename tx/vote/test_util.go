package vote

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestKVStoreKey = sdk.NewKVStoreKey("vote")
)

func newLinoAccountManager() acc.AccountManager {
	return acc.NewLinoAccountManager(TestKVStoreKey)
}

func newVoteManager() VoteManager {
	return NewVoteMananger(TestKVStoreKey)
}

func newGlobalManager() global.GlobalManager {
	return global.NewGlobalManager(TestKVStoreKey)
}

func newGlobalProxy() global.GlobalProxy {
	var gm global.GlobalManager
	gm = newGlobalManager()
	return *global.NewGlobalProxy(&gm)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

// helper function to create an account for testing purpose
func privAndBank() (crypto.PrivKey, *acc.AccountBank) {
	priv := crypto.GenPrivKeyEd25519()
	accBank := &acc.AccountBank{
		Address: priv.PubKey().Address(),
	}
	return priv.Wrap(), accBank
}

func createTestAccount(ctx sdk.Context, lam acc.AccountManager, username string) *acc.Account {
	priv, bank := privAndBank()
	account := acc.NewProxyAccount(acc.AccountKey(username), &lam)
	account.CreateAccount(ctx, acc.AccountKey(username), priv.PubKey(), bank)
	account.Apply(ctx)
	return account
}
