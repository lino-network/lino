package account

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestKVStoreKey = sdk.NewKVStoreKey("account")
)

func newLinoAccountManager() LinoAccountManager {
	return NewLinoAccountManager(TestKVStoreKey)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

// helper function to create an account for testing purpose
func privAndBank() (crypto.PrivKey, *types.AccountBank) {
	priv := crypto.GenPrivKeyEd25519()
	accBank := &types.AccountBank{
		Address: priv.PubKey().Address(),
		Coins:   sdk.Coins{sdk.Coin{Denom: "dummy", Amount: 123}},
	}
	return priv.Wrap(), accBank
}

func createTestAccount(ctx sdk.Context, lam LinoAccountManager, username string) {
	priv, bank := privAndBank()
	user := types.AccountKey(username)
	lam.CreateAccount(ctx, user, priv.PubKey(), bank)
}
