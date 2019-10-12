package manager

// import (
// 	"testing"
// 	"time"

// 	"github.com/cosmos/cosmos-sdk/store"
// 	"github.com/lino-network/lino/param"
// 	"github.com/lino-network/lino/types"
// 	acc "github.com/lino-network/lino/x/account"
// 	accmn "github.com/lino-network/lino/x/account/manager"
// 	"github.com/lino-network/lino/x/global"
// 	"github.com/lino-network/lino/x/post"
// 	postmn "github.com/lino-network/lino/x/post/manager"
// 	val "github.com/lino-network/lino/x/validator"
// 	"github.com/lino-network/lino/x/vote"
// 	votemn "github.com/lino-network/lino/x/vote/manager"
// 	"github.com/stretchr/testify/assert"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	abci "github.com/tendermint/tendermint/abci/types"
// 	"github.com/tendermint/tendermint/crypto/secp256k1"
// 	"github.com/tendermint/tendermint/libs/log"
// 	dbm "github.com/tendermint/tm-db"
// )

// // Construct some global addrs and txs for tests.
// var (
// 	testAccountKVStoreKey   = sdk.NewKVStoreKey("account")
// 	testGlobalKVStoreKey    = sdk.NewKVStoreKey("global")
// 	testProposalKVStoreKey  = sdk.NewKVStoreKey("proposal")
// 	testVoteKVStoreKey      = sdk.NewKVStoreKey("vote")
// 	testParamKVStoreKey     = sdk.NewKVStoreKey("param")
// 	testValidatorKVStoreKey = sdk.NewKVStoreKey("validator")
// 	testPostKVStoreKey      = sdk.NewKVStoreKey("post")
// )

// func initGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
// 	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
// }

// func setupTest(t *testing.T, height int64) (
// 	sdk.Context, acc.AccountKeeper, ProposalManager, post.PostKeeper, vote.VoteKeeper,
// 	val.ValidatorManager, global.GlobalManager) {
// 	ctx := getContext(height)
// 	ph := param.NewParamHolder(testParamKVStoreKey)
// 	err := ph.InitParam(ctx)
// 	if err != nil {
// 		panic(err)
// 	}

// 	proposalManager := NewProposalManager(testProposalKVStoreKey, ph)
// 	gm := global.NewGlobalManager(testGlobalKVStoreKey, ph)
// 	am := accmn.NewAccountManager(testAccountKVStoreKey, ph, &gm)
// 	voteManager := votemn.NewVoteManager(testGlobalKVStoreKey, ph, am, &gm)
// 	valManager := val.NewValidatorManager(testValidatorKVStoreKey, ph)
// 	postManager := postmn.NewPostManager(testPostKVStoreKey, am, &gm, nil, nil, nil)

// 	cdc := gm.WireCodec()
// 	cdc.RegisterInterface((*types.Event)(nil), nil)
// 	cdc.RegisterConcrete(accmn.ReturnCoinEvent{}, "1", nil)
// 	cdc.RegisterConcrete(param.ChangeParamEvent{}, "2", nil)
// 	cdc.RegisterConcrete(DecideProposalEvent{}, "3", nil)

// 	err = initGlobalManager(ctx, gm)
// 	assert.Nil(t, err)
// 	return ctx, am, proposalManager, postManager, voteManager, valManager, gm
// }

// func getContext(height int64) sdk.Context {
// 	db := dbm.NewMemDB()
// 	ms := store.NewCommitMultiStore(db)
// 	ms.MountStoreWithDB(testAccountKVStoreKey, sdk.StoreTypeIAVL, db)
// 	ms.MountStoreWithDB(testProposalKVStoreKey, sdk.StoreTypeIAVL, db)
// 	ms.MountStoreWithDB(testGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
// 	ms.MountStoreWithDB(testParamKVStoreKey, sdk.StoreTypeIAVL, db)
// 	ms.MountStoreWithDB(testVoteKVStoreKey, sdk.StoreTypeIAVL, db)
// 	ms.MountStoreWithDB(testValidatorKVStoreKey, sdk.StoreTypeIAVL, db)
// 	ms.MountStoreWithDB(testPostKVStoreKey, sdk.StoreTypeIAVL, db)

// 	err := ms.LoadLatestVersion()
// 	if err != nil {
// 		panic(err)
// 	}

// 	return sdk.NewContext(ms, abci.Header{Height: height, Time: time.Now()}, false, log.NewNopLogger())
// }

// // helper function to create an account for testing purpose
// func createTestAccount(
// 	ctx sdk.Context, am acc.AccountKeeper, username string, initCoin types.Coin) types.AccountKey {
// 	err := am.CreateAccount(ctx, types.AccountKey(username), secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey())
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = am.AddCoinToUsername(ctx, types.AccountKey(username), initCoin)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return types.AccountKey(username)
// }

// func createTestPost(
// 	t *testing.T, ctx sdk.Context, username, postID string, initCoin types.Coin,
// 	am acc.AccountKeeper, pm post.PostKeeper, redistributionRate string) (types.AccountKey, string) {
// 	user := createTestAccount(ctx, am, username, initCoin)
// 	msg := &post.CreatePostMsg{
// 		PostID:    postID,
// 		Title:     string(make([]byte, 50)),
// 		Content:   string(make([]byte, 1000)),
// 		Author:    user,
// 		CreatedBy: user,
// 	}

// 	err := pm.CreatePost(ctx, msg.Author, msg.PostID, msg.CreatedBy, msg.Content, msg.Title)
// 	assert.Nil(t, err)
// 	return user, postID
// }

// func addProposalInfo(ctx sdk.Context, pm ProposalManager, proposalID types.ProposalKey,
// 	agreeVotes, disagreeVotes types.Coin) sdk.Error {
// 	proposal, err := pm.storage.GetOngoingProposal(ctx, proposalID)
// 	if err != nil {
// 		return err
// 	}

// 	proposalInfo := proposal.GetProposalInfo()
// 	proposalInfo.AgreeVotes = agreeVotes
// 	proposalInfo.DisagreeVotes = disagreeVotes

// 	proposal.SetProposalInfo(proposalInfo)

// 	if err := pm.storage.SetOngoingProposal(ctx, proposalID, proposal); err != nil {
// 		return err
// 	}
// 	return nil
// }
