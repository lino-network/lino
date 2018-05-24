package proposal

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/tx/post"
	val "github.com/lino-network/lino/tx/validator"
	"github.com/lino-network/lino/tx/vote"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
)

// Construct some global addrs and txs for tests.
var (
	TestAccountKVStoreKey   = sdk.NewKVStoreKey("account")
	TestGlobalKVStoreKey    = sdk.NewKVStoreKey("global")
	TestProposalKVStoreKey  = sdk.NewKVStoreKey("proposal")
	TestVoteKVStoreKey      = sdk.NewKVStoreKey("vote")
	TestParamKVStoreKey     = sdk.NewKVStoreKey("param")
	TestValidatorKVStoreKey = sdk.NewKVStoreKey("validator")
	TestPostKVStoreKey      = sdk.NewKVStoreKey("post")
)

func InitGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest(t *testing.T, height int64) (
	sdk.Context, acc.AccountManager, ProposalManager, post.PostManager, vote.VoteManager,
	val.ValidatorManager, global.GlobalManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(TestParamKVStoreKey)
	ph.InitParam(ctx)

	accManager := acc.NewAccountManager(TestAccountKVStoreKey, ph)
	proposalManager := NewProposalManager(TestProposalKVStoreKey, ph)
	globalManager := global.NewGlobalManager(TestGlobalKVStoreKey, ph)
	voteManager := vote.NewVoteManager(TestGlobalKVStoreKey, ph)
	valManager := val.NewValidatorManager(TestValidatorKVStoreKey, ph)
	postManager := post.NewPostManager(TestPostKVStoreKey, ph)

	cdc := globalManager.WireCodec()
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(acc.ReturnCoinEvent{}, "1", nil)
	cdc.RegisterConcrete(param.ChangeParamEvent{}, "2", nil)
	cdc.RegisterConcrete(DecideProposalEvent{}, "3", nil)

	err := InitGlobalManager(ctx, globalManager)
	assert.Nil(t, err)
	return ctx, accManager, proposalManager, postManager, voteManager, valManager, globalManager
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestProposalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestVoteKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestValidatorKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestPostKVStoreKey, sdk.StoreTypeIAVL, db)

	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{Height: height}, false, nil)
}

// helper function to create an account for testing purpose
func createTestAccount(
	ctx sdk.Context, am acc.AccountManager, username string, initCoin types.Coin) types.AccountKey {
	priv := crypto.GenPrivKeyEd25519()
	am.AddSavingCoinToAddress(ctx, priv.PubKey().Address(), initCoin)
	am.CreateAccount(ctx, types.AccountKey(username),
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	return types.AccountKey(username)
}

func createTestPost(
	t *testing.T, ctx sdk.Context, username, postID string, initCoin types.Coin,
	am acc.AccountManager, pm post.PostManager, redistributionRate string) (types.AccountKey, string) {
	user := createTestAccount(ctx, am, username, initCoin)
	postCreateParams := &post.PostCreateParams{
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
