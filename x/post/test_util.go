package post

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	dev "github.com/lino-network/lino/x/developer"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/post/model"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestAccountKVStoreKey   = sdk.NewKVStoreKey("account")
	TestPostKVStoreKey      = sdk.NewKVStoreKey("post")
	TestGlobalKVStoreKey    = sdk.NewKVStoreKey("global")
	TestDeveloperKVStoreKey = sdk.NewKVStoreKey("developer")
	TestParamKVStoreKey     = sdk.NewKVStoreKey("param")

	initCoin = types.NewCoinFromInt64(1 * types.Decimals)
	referrer = types.AccountKey("referrer")
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
		ms, abci.Header{ChainID: "Lino", Height: height, Time: time.Now().Unix()}, false, log.NewNopLogger())
}

func checkPostKVStore(
	t *testing.T, ctx sdk.Context, postKey types.Permlink, postInfo model.PostInfo, postMeta model.PostMeta) {
	// check all post related structs in KVStore
	postStorage := model.NewPostStorage(TestPostKVStoreKey)
	postPtr, err := postStorage.GetPostInfo(ctx, postKey)
	assert.Nil(t, err)
	assert.Equal(t, postInfo, *postPtr, "postInfo should be equal")
	checkPostMeta(t, ctx, postKey, postMeta)
}

func checkPostMeta(t *testing.T, ctx sdk.Context, postKey types.Permlink, postMeta model.PostMeta) {
	// check post meta structs in KVStore
	postStorage := model.NewPostStorage(TestPostKVStoreKey)
	postMetaPtr, err := postStorage.GetPostMeta(ctx, postKey)
	assert.Nil(t, err)
	assert.Equal(t, postMeta, *postMetaPtr, "Post meta should be equal")
}

func createTestAccount(
	t *testing.T, ctx sdk.Context, am acc.AccountManager, username string) types.AccountKey {
	priv := crypto.GenPrivKeyEd25519()
	err := am.CreateAccount(ctx, referrer, types.AccountKey(username),
		priv.PubKey(), priv.Generate(0).PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), initCoin)
	assert.Nil(t, err)
	return types.AccountKey(username)
}

func createTestPost(
	t *testing.T, ctx sdk.Context, username, postID string,
	am acc.AccountManager, pm PostManager, redistributionRate string) (types.AccountKey, string) {
	user := createTestAccount(t, ctx, am, username)

	splitRate, err := sdk.NewRatFromDecimal(redistributionRate, types.NewRatFromDecimalPrecision)
	assert.Nil(t, err)
	err = pm.CreatePost(
		ctx, types.AccountKey(user), postID, "", "", "", "",
		string(make([]byte, 1000)), string(make([]byte, 50)),
		splitRate, []types.IDToURLMapping{})
	assert.Nil(t, err)
	return user, postID
}

func createTestRepost(
	t *testing.T, ctx sdk.Context, username, postID string,
	am acc.AccountManager, pm PostManager, sourceUser types.AccountKey,
	sourcePostID string) (types.AccountKey, string) {
	user := createTestAccount(t, ctx, am, username)

	err := pm.CreatePost(
		ctx, types.AccountKey(user), postID, sourceUser, sourcePostID, "", "",
		string(make([]byte, 1000)), string(make([]byte, 50)),
		sdk.ZeroRat(), []types.IDToURLMapping{})
	assert.Nil(t, err)
	return user, postID
}
