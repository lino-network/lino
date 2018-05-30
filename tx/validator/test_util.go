package validator

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	vote "github.com/lino-network/lino/tx/vote"
	abci "github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var (
	TestAccountKVStoreKey   = sdk.NewKVStoreKey("account")
	TestValidatorKVStoreKey = sdk.NewKVStoreKey("validator")
	TestGlobalKVStoreKey    = sdk.NewKVStoreKey("global")
	TestVoteKVStoreKey      = sdk.NewKVStoreKey("vote")
	TestParamKVStoreKey     = sdk.NewKVStoreKey("param")
)

func InitGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest(t *testing.T, height int64) (sdk.Context,
	acc.AccountManager, ValidatorManager, vote.VoteManager, global.GlobalManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(TestParamKVStoreKey)
	ph.InitParam(ctx)
	accManager := acc.NewAccountManager(TestAccountKVStoreKey, ph)
	postManager := NewValidatorManager(TestValidatorKVStoreKey, ph)
	globalManager := global.NewGlobalManager(TestGlobalKVStoreKey, ph)
	voteManager := vote.NewVoteManager(TestVoteKVStoreKey, ph)

	cdc := globalManager.WireCodec()
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(acc.ReturnCoinEvent{}, "event/return", nil)

	err := InitGlobalManager(ctx, globalManager)
	assert.Nil(t, err)
	return ctx, accManager, postManager, voteManager, globalManager
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(TestAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestValidatorKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestVoteKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(TestParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(ms, abci.Header{Height: height}, false, nil)
}

// helper function to create an account for testing purpose
func createTestAccount(ctx sdk.Context, am acc.AccountManager, username string, initCoin types.Coin) types.AccountKey {
	priv := crypto.GenPrivKeyEd25519()
	am.AddSavingCoinToAddress(ctx, priv.PubKey().Address(), initCoin, types.TransferIn)
	am.CreateAccount(ctx, types.AccountKey(username),
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey())
	return types.AccountKey(username)
}

func coinToString(coin types.Coin) string {
	return strconv.FormatInt(coin.ToInt64()/types.Decimals, 10)
}
