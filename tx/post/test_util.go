package post

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/genesis"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/post/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestAccountKVStoreKey = sdk.NewKVStoreKey("account")
	TestPostKVStoreKey    = sdk.NewKVStoreKey("post")
	TestGlobalKVStoreKey  = sdk.NewKVStoreKey("global")

	initCoin = types.NewCoin(100)
)

func InitGlobalManager(ctx sdk.Context, gm *global.GlobalManager) error {
	globalState := genesis.GlobalState{
		TotalLino:                10000,
		GrowthRate:               sdk.Rat{98, 1000},
		InfraAllocation:          sdk.Rat{20, 100},
		ContentCreatorAllocation: sdk.Rat{55, 100},
		DeveloperAllocation:      sdk.Rat{20, 100},
		ValidatorAllocation:      sdk.Rat{5, 100},
		ConsumptionFrictionRate:  sdk.Rat{1, 100},
		FreezingPeriodHr:         24 * 7,
	}
	return gm.InitGlobalManager(ctx, globalState)
}

func setupTest(t *testing.T, height int64) (sdk.Context, *acc.AccountManager, *PostManager, *global.GlobalManager) {
	ctx := getContext(height)
	accManager := acc.NewAccountManager(TestAccountKVStoreKey)
	postManager := NewPostManager(TestPostKVStoreKey)
	globalManager := global.NewGlobalManager(TestGlobalKVStoreKey)
	err := InitGlobalManager(ctx, globalManager)
	assert.Nil(t, err)
	return ctx, accManager, postManager, globalManager
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestPostKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{ChainID: "Lino", Height: height, Time: time.Now().Unix()}, false, nil)
}

func checkPostKVStore(t *testing.T, ctx sdk.Context, postKey types.PostKey, postInfo model.PostInfo, postMeta model.PostMeta) {
	// check all post related structs in KVStore
	postStorage := model.NewPostStorage(TestPostKVStoreKey)
	postPtr, err := postStorage.GetPostInfo(ctx, postKey)
	assert.Nil(t, err)
	assert.Equal(t, postInfo, *postPtr, "postInfo should be equal")
	checkPostMeta(t, ctx, postKey, postMeta)
}

func checkPostMeta(t *testing.T, ctx sdk.Context, postKey types.PostKey, postMeta model.PostMeta) {
	// check post meta structs in KVStore
	postStorage := model.NewPostStorage(TestPostKVStoreKey)
	postMetaPtr, err := postStorage.GetPostMeta(ctx, postKey)
	assert.Nil(t, err)
	assert.Equal(t, postMeta, *postMetaPtr, "Post meta should be equal")
}

func createTestAccount(t *testing.T, ctx sdk.Context, am *acc.AccountManager, username string) types.AccountKey {
	priv := crypto.GenPrivKeyEd25519()
	err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), initCoin)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, types.AccountKey(username), priv.PubKey(), types.NewCoin(0))
	assert.Nil(t, err)
	return types.AccountKey(username)
}

func createTestPost(
	t *testing.T, ctx sdk.Context, username, postID string,
	am *acc.AccountManager, pm *PostManager, redistributionRate sdk.Rat) (types.AccountKey, string) {
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
