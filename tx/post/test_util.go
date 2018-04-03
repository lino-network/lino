package post

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestKVStoreKey = sdk.NewKVStoreKey("account")
)

func newLinoAccountManager() acc.AccountManager {
	return acc.NewLinoAccountManager(TestKVStoreKey)
}

func newPostManager() PostManager {
	return NewPostMananger(TestKVStoreKey)
}

func newAmount(amount int64) types.Coin {
	return types.NewCoin(amount)
}

func getContext() sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{}, false, nil)
}

func privAndBank() (crypto.PrivKey, *acc.AccountBank) {
	priv := crypto.GenPrivKeyEd25519()
	accBank := &acc.AccountBank{
		Address: priv.PubKey().Address(),
		Balance: types.Coin{123 * types.Decimals},
	}
	return priv.Wrap(), accBank
}

func createTestAccount(ctx sdk.Context, lam acc.AccountManager, username string) *acc.AccountProxy {
	priv, bank := privAndBank()
	account := acc.NewAccountProxy(acc.AccountKey(username), &lam)
	account.CreateAccount(ctx, acc.AccountKey(username), priv.PubKey(), bank)
	account.Apply(ctx)
	return account
}

func createTestPost(ctx sdk.Context, lam acc.AccountManager, pm PostManager, username, postID string) *PostProxy {
	createTestAccount(ctx, lam, username)
	postInfo := PostInfo{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       acc.AccountKey(username),
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []IDToURLMapping{},
	}
	post := NewPostProxy(postInfo.Author, postInfo.PostID, &pm)
	post.CreatePost(ctx, &postInfo)
	post.Apply(ctx)
	return post
}
