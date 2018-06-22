package proposal

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"
	"github.com/lino-network/lino/x/post"
	"github.com/lino-network/lino/x/vote"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/tmlibs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	val "github.com/lino-network/lino/x/validator"
	abci "github.com/tendermint/abci/types"
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

	return sdk.NewContext(ms, abci.Header{Height: height}, false, nil, log.NewNopLogger())
}

// helper function to create an account for testing purpose
func createTestAccount(
	ctx sdk.Context, am acc.AccountManager, username string, initCoin types.Coin) types.AccountKey {
	priv := crypto.GenPrivKeyEd25519()
	am.CreateAccount(ctx, "referrer", types.AccountKey(username),
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), initCoin)
	return types.AccountKey(username)
}

func createTestPost(
	t *testing.T, ctx sdk.Context, username, postID string, initCoin types.Coin,
	am acc.AccountManager, pm post.PostManager, redistributionRate string) (types.AccountKey, string) {
	user := createTestAccount(ctx, am, username, initCoin)
	msg := &post.CreatePostMsg{
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
	splitRate, err := sdk.NewRatFromDecimal(redistributionRate)
	assert.Nil(t, err)

	err = pm.CreatePost(
		ctx, msg.Author, msg.PostID, msg.SourceAuthor, msg.SourcePostID,
		msg.ParentAuthor, msg.ParentPostID, msg.Content,
		msg.Title, splitRate, msg.Links)

	assert.Nil(t, err)
	return user, postID
}

func addProposalInfo(ctx sdk.Context, pm ProposalManager, proposalID types.ProposalKey,
	agreeVotes, disagreeVotes types.Coin) sdk.Error {
	proposal, err := pm.storage.GetProposal(ctx, proposalID)
	if err != nil {
		return err
	}

	proposalInfo := proposal.GetProposalInfo()
	proposalInfo.AgreeVotes = agreeVotes
	proposalInfo.DisagreeVotes = disagreeVotes

	proposal.SetProposalInfo(proposalInfo)

	if err := pm.storage.SetProposal(ctx, proposalID, proposal); err != nil {
		return err
	}
	return nil
}
