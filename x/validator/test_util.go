package validator

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/x/account"
	accmn "github.com/lino-network/lino/x/account/manager"
	vote "github.com/lino-network/lino/x/vote"
	votemn "github.com/lino-network/lino/x/vote/manager"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	testAccountKVStoreKey   = sdk.NewKVStoreKey("account")
	testValidatorKVStoreKey = sdk.NewKVStoreKey("validator")
	testGlobalKVStoreKey    = sdk.NewKVStoreKey("global")
	testVoteKVStoreKey      = sdk.NewKVStoreKey("vote")
	testParamKVStoreKey     = sdk.NewKVStoreKey("param")
)

func initGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest(t *testing.T, height int64) (sdk.Context,
	acc.AccountKeeper, ValidatorManager, vote.VoteKeeper, global.GlobalManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(testParamKVStoreKey)
	err := ph.InitParam(ctx)
	if err != nil {
		panic(err)
	}
	gm := global.NewGlobalManager(testGlobalKVStoreKey, ph)
	am := accmn.NewAccountManager(testAccountKVStoreKey, ph, &gm)
	postManager := NewValidatorManager(testValidatorKVStoreKey, ph)
	voteManager := votemn.NewVoteManager(testVoteKVStoreKey, ph, am, &gm)

	cdc := gm.WireCodec()
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(accmn.ReturnCoinEvent{}, "event/return", nil)

	err = initGlobalManager(ctx, gm)
	assert.Nil(t, err)
	return ctx, am, postManager, voteManager, gm
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(testAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testValidatorKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testVoteKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testParamKVStoreKey, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	return sdk.NewContext(ms, abci.Header{Height: height}, false, log.NewNopLogger())
}

// helper function to create an account for testing purpose
func createTestAccount(ctx sdk.Context, am acc.AccountKeeper, username string, initCoin types.Coin) types.AccountKey {
	err := am.CreateAccount(ctx, types.AccountKey(username), secp256k1.GenPrivKey().PubKey(), secp256k1.GenPrivKey().PubKey())
	if err != nil {
		panic(err)
	}
	err = am.AddCoinToUsername(ctx, types.AccountKey(username), initCoin)
	if err != nil {
		panic(err)
	}
	return types.AccountKey(username)
}

func coinToString(coin types.Coin) string {
	coinInInt64, _ := coin.ToInt64()
	return strconv.FormatInt(coinInInt64/types.Decimals, 10)
}
