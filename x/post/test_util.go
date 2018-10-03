package post

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/recorder"
	"github.com/lino-network/lino/types"
	acc "github.com/lino-network/lino/x/account"
	dev "github.com/lino-network/lino/x/developer"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/post/model"
	rep "github.com/lino-network/lino/x/reputation"
	vote "github.com/lino-network/lino/x/vote"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

// Construct some global addrs and txs for tests.
var (
	testAccountKVStoreKey   = sdk.NewKVStoreKey("account")
	testPostKVStoreKey      = sdk.NewKVStoreKey("post")
	testGlobalKVStoreKey    = sdk.NewKVStoreKey("global")
	testDeveloperKVStoreKey = sdk.NewKVStoreKey("developer")
	testVoteKVStoreKey      = sdk.NewKVStoreKey("vote")
	testParamKVStoreKey     = sdk.NewKVStoreKey("param")
	testRepKVStoreKey       = sdk.NewKVStoreKey("reputation")

	initCoin = types.NewCoinFromInt64(1 * types.Decimals)
	referrer = types.AccountKey("referrer")
)

// InitGlobalManager - init global manager
func InitGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest(
	t *testing.T, height int64) (
	sdk.Context, acc.AccountManager, param.ParamHolder, PostManager,
	global.GlobalManager, dev.DeveloperManager, vote.VoteManager, rep.ReputationManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(testParamKVStoreKey)
	recorder := recorder.NewRecorder()
	ph.InitParam(ctx)
	accManager := acc.NewAccountManager(testAccountKVStoreKey, ph)
	postManager := NewPostManager(testPostKVStoreKey, ph, recorder)
	globalManager := global.NewGlobalManager(testGlobalKVStoreKey, ph)
	devManager := dev.NewDeveloperManager(testDeveloperKVStoreKey, ph)
	devManager.InitGenesis(ctx)
	voteManager := vote.NewVoteManager(testVoteKVStoreKey, ph)
	voteManager.InitGenesis(ctx)
	repManager := rep.NewReputationManager(testRepKVStoreKey, ph)

	cdc := globalManager.WireCodec()
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(RewardEvent{}, "event/reward", nil)

	err := InitGlobalManager(ctx, globalManager)
	assert.Nil(t, err)
	return ctx, accManager, ph, postManager, globalManager, devManager, voteManager, repManager
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(testAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testPostKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testDeveloperKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testVoteKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testRepKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: height, Time: time.Unix(0, 0)}, false, log.NewNopLogger())
}

func checkPostKVStore(
	t *testing.T, ctx sdk.Context, postKey types.Permlink, postInfo model.PostInfo, postMeta model.PostMeta) {
	// check all post related structs in KVStore
	postStorage := model.NewPostStorage(testPostKVStoreKey)
	postPtr, err := postStorage.GetPostInfo(ctx, postKey)
	assert.Nil(t, err)
	assert.Equal(t, postInfo, *postPtr, "postInfo should be equal")
	checkPostMeta(t, ctx, postKey, postMeta)
}

func checkPostMeta(t *testing.T, ctx sdk.Context, postKey types.Permlink, postMeta model.PostMeta) {
	// check post meta structs in KVStore
	postStorage := model.NewPostStorage(testPostKVStoreKey)
	postMetaPtr, err := postStorage.GetPostMeta(ctx, postKey)
	assert.Nil(t, err)
	assert.Equal(t, postMeta, *postMetaPtr, "Post meta should be equal")
}

func createTestAccount(
	t *testing.T, ctx sdk.Context, am acc.AccountManager, username string) types.AccountKey {
	err := am.CreateAccount(ctx, referrer, types.AccountKey(username),
		secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(), initCoin)
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
