package post

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/param"
	acc "github.com/lino-network/lino/tx/account"
	dev "github.com/lino-network/lino/tx/developer"
	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/tx/post/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestAccountKVStoreKey   = sdk.NewKVStoreKey("account")
	TestPostKVStoreKey      = sdk.NewKVStoreKey("post")
	TestGlobalKVStoreKey    = sdk.NewKVStoreKey("global")
	TestDeveloperKVStoreKey = sdk.NewKVStoreKey("developer")
	TestParamKVStoreKey     = sdk.NewKVStoreKey("param")

	initCoin = types.NewCoinFromInt64(1 * types.Decimals)
)

func InitGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest(
	t *testing.T, height int64) (
	sdk.Context, acc.AccountManager, param.ParamHolder,
	PostManager, global.GlobalManager, dev.DeveloperManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(TestParamKVStoreKey)
	ph.InitParam(ctx)
	accManager := acc.NewAccountManager(TestAccountKVStoreKey, ph)
	postManager := NewPostManager(TestPostKVStoreKey, ph)
	globalManager := global.NewGlobalManager(TestGlobalKVStoreKey, ph)
	devManager := dev.NewDeveloperManager(TestDeveloperKVStoreKey, ph)
	devManager.InitGenesis(ctx)

	cdc := globalManager.WireCodec()
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(RewardEvent{}, "event/reward", nil)

	err := InitGlobalManager(ctx, globalManager)
	assert.Nil(t, err)
	return ctx, accManager, ph, postManager, globalManager, devManager
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestPostKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestDeveloperKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: height, Time: time.Now().Unix()}, false, nil)
}

func checkPostKVStore(
	t *testing.T, ctx sdk.Context, postKey types.PermLink, postInfo model.PostInfo, postMeta model.PostMeta) {
	// check all post related structs in KVStore
	postStorage := model.NewPostStorage(TestPostKVStoreKey)
	postPtr, err := postStorage.GetPostInfo(ctx, postKey)
	assert.Nil(t, err)
	assert.Equal(t, postInfo, *postPtr, "postInfo should be equal")
	checkPostMeta(t, ctx, postKey, postMeta)
}

func checkPostMeta(t *testing.T, ctx sdk.Context, postKey types.PermLink, postMeta model.PostMeta) {
	// check post meta structs in KVStore
	postStorage := model.NewPostStorage(TestPostKVStoreKey)
	postMetaPtr, err := postStorage.GetPostMeta(ctx, postKey)
	assert.Nil(t, err)
	assert.Equal(t, postMeta, *postMetaPtr, "Post meta should be equal")
}

func createTestAccount(
	t *testing.T, ctx sdk.Context, am acc.AccountManager, username string) types.AccountKey {
	priv := crypto.GenPrivKeyEd25519()
	err := am.AddSavingCoinToAddress(ctx, priv.PubKey().Address(), initCoin)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, types.AccountKey(username),
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	assert.Nil(t, err)
	return types.AccountKey(username)
}

func createTestPost(
	t *testing.T, ctx sdk.Context, username, postID string,
	am acc.AccountManager, pm PostManager, redistributionRate string) (types.AccountKey, string) {
	user := createTestAccount(t, ctx, am, username)
	postCreateParams := &PostCreateParams{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: "",
		SourcePostID: "",
		Links:        []types.IDToURLMapping{},
		RedistributionSplitRate: redistributionRate,
	}
	err := pm.CreatePost(ctx, postCreateParams)
	assert.Nil(t, err)
	return user, postID
}

func createTestRepost(
	t *testing.T, ctx sdk.Context, username, postID string,
	am acc.AccountManager, pm PostManager, sourceUser types.AccountKey,
	sourcePostID string) (types.AccountKey, string) {
	user := createTestAccount(t, ctx, am, username)
	postCreateParams := &PostCreateParams{
		PostID:       postID,
		Title:        string(make([]byte, 50)),
		Content:      string(make([]byte, 1000)),
		Author:       user,
		ParentAuthor: "",
		ParentPostID: "",
		SourceAuthor: sourceUser,
		SourcePostID: sourcePostID,
		Links:        []types.IDToURLMapping{},
		RedistributionSplitRate: "0",
	}
	err := pm.CreatePost(ctx, postCreateParams)
	assert.Nil(t, err)
	return user, postID
}
