package account

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/recorder"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"
	"github.com/lino-network/lino/x/global"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
)

var (
	testAccountKVStoreKey = sdk.NewKVStoreKey("account")
	testGlobalKVStoreKey  = sdk.NewKVStoreKey("global")
	testParamKVStoreKey   = sdk.NewKVStoreKey("param")

	accountReferrer = types.AccountKey("referrer")

	l0    = types.LNO("0")
	l100  = types.LNO("100")
	l200  = types.LNO("200")
	l1600 = types.LNO("1600")
	l1800 = types.LNO("1800")
	l1900 = types.LNO("1900")
	l1999 = types.LNO("1999")
	l2000 = types.LNO("2000")
	c0    = types.NewCoinFromInt64(0)
	c100  = types.NewCoinFromInt64(100 * types.Decimals)
	c200  = types.NewCoinFromInt64(200 * types.Decimals)
	c300  = types.NewCoinFromInt64(300 * types.Decimals)
	c400  = types.NewCoinFromInt64(400 * types.Decimals)
	c500  = types.NewCoinFromInt64(500 * types.Decimals)
	c600  = types.NewCoinFromInt64(600 * types.Decimals)
	c1000 = types.NewCoinFromInt64(1000 * types.Decimals)
	c1500 = types.NewCoinFromInt64(1500 * types.Decimals)
	c1600 = types.NewCoinFromInt64(1600 * types.Decimals)
	c1800 = types.NewCoinFromInt64(1800 * types.Decimals)
	c1900 = types.NewCoinFromInt64(1900 * types.Decimals)
	c2000 = types.NewCoinFromInt64(2000 * types.Decimals)

	coin0   = types.NewCoinFromInt64(0)
	coin1   = types.NewCoinFromInt64(1)
	coin2   = types.NewCoinFromInt64(2)
	coin3   = types.NewCoinFromInt64(3)
	coin4   = types.NewCoinFromInt64(4)
	coin10  = types.NewCoinFromInt64(10)
	coin50  = types.NewCoinFromInt64(50)
	coin100 = types.NewCoinFromInt64(100)
	coin200 = types.NewCoinFromInt64(200)
	coin300 = types.NewCoinFromInt64(300)
	coin400 = types.NewCoinFromInt64(400)
)

func initGlobalManager(ctx sdk.Context, gm global.GlobalManager) error {
	return gm.InitGlobalManager(ctx, types.NewCoinFromInt64(10000*types.Decimals))
}

func setupTest(t *testing.T, height int64) (sdk.Context, AccountManager, global.GlobalManager) {
	ctx := getContext(height)
	ph := param.NewParamHolder(testParamKVStoreKey)
	ph.InitParam(ctx)
	recorder := recorder.NewRecorder()
	accManager := NewAccountManager(testAccountKVStoreKey, ph)
	globalManager := global.NewGlobalManager(testGlobalKVStoreKey, ph, recorder)

	cdc := globalManager.WireCodec()
	cdc.RegisterInterface((*types.Event)(nil), nil)
	cdc.RegisterConcrete(ReturnCoinEvent{}, "event/return", nil)

	err := initGlobalManager(ctx, globalManager)
	assert.Nil(t, err)
	return ctx, accManager, globalManager
}

func getContext(height int64) sdk.Context {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(testAccountKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testGlobalKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(testParamKVStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return sdk.NewContext(
		ms, abci.Header{ChainID: "Lino", Height: height, Time: time.Now()},
		false, log.NewNopLogger())
}

func createTestAccount(ctx sdk.Context, am AccountManager, username string) (secp256k1.PrivKeySecp256k1,
	secp256k1.PrivKeySecp256k1, secp256k1.PrivKeySecp256k1) {
	resetPriv := secp256k1.GenPrivKey()
	txPriv := secp256k1.GenPrivKey()
	appPriv := secp256k1.GenPrivKey()

	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	am.CreateAccount(ctx, accountReferrer, types.AccountKey(username),
		resetPriv.PubKey(), txPriv.PubKey(), appPriv.PubKey(), accParam.RegisterFee)
	return resetPriv, txPriv, appPriv
}

func checkBankKVByUsername(
	t *testing.T, ctx sdk.Context, testName string, username types.AccountKey, bank model.AccountBank) {
	accStorage := model.NewAccountStorage(testAccountKVStoreKey)
	bankPtr, err := accStorage.GetBankFromAccountKey(ctx, username)
	if err != nil {
		t.Errorf("%s, failed to get bank, got err %v", testName, err)
	}
	if !assert.Equal(t, bank, *bankPtr) {
		t.Errorf("%s: diff bank, got %v, want %v", testName, *bankPtr, bank)
	}
}

func checkBalanceHistory(
	t *testing.T, ctx sdk.Context, testName string, username types.AccountKey,
	timeSlot int64, balanceHistory model.BalanceHistory) {
	accStorage := model.NewAccountStorage(testAccountKVStoreKey)
	balanceHistoryPtr, err := accStorage.GetBalanceHistory(ctx, username, timeSlot)
	if err != nil {
		t.Errorf("%s, failed to get balance history, got err %v", testName, err)
	}
	if !assert.Equal(t, balanceHistory, *balanceHistoryPtr) {
		t.Errorf("%s: diff balance history, got %v, want %v", testName, *balanceHistoryPtr, balanceHistory)
	}
}

func checkPendingCoinDay(
	t *testing.T, ctx sdk.Context, testName string, username types.AccountKey, pendingCoinDayQueue model.PendingCoinDayQueue) {
	accStorage := model.NewAccountStorage(testAccountKVStoreKey)
	pendingCoinDayQueuePtr, err := accStorage.GetPendingCoinDayQueue(ctx, username)
	if err != nil {
		t.Errorf("%s, failed to get pending coin day queue, got err %v", testName, err)
	}
	if !assert.Equal(t, pendingCoinDayQueue, *pendingCoinDayQueuePtr) {
		t.Errorf("%s: diff pending coin day queue, got %v, want %v", testName, *pendingCoinDayQueuePtr, pendingCoinDayQueue)
	}
}

func checkAccountInfo(
	t *testing.T, ctx sdk.Context, testName string, accKey types.AccountKey, accInfo model.AccountInfo) {
	accStorage := model.NewAccountStorage(testAccountKVStoreKey)
	info, err := accStorage.GetInfo(ctx, accKey)
	if err != nil {
		t.Errorf("%s, failed to get account info, got err %v", testName, err)
	}
	if !assert.Equal(t, accInfo, *info) {
		t.Errorf("%s: diff account info, got %v, want %v", testName, *info, accInfo)
	}
}

func checkAccountMeta(
	t *testing.T, ctx sdk.Context, testName string, accKey types.AccountKey, accMeta model.AccountMeta) {
	accStorage := model.NewAccountStorage(testAccountKVStoreKey)
	metaPtr, err := accStorage.GetMeta(ctx, accKey)
	if err != nil {
		t.Errorf("%s, failed to get account meta, got err %v", testName, err)
	}
	if !assert.Equal(t, accMeta, *metaPtr) {
		t.Errorf("%s: diff account meta, got %v, want %v", testName, *metaPtr, accMeta)
	}
}

func checkAccountReward(
	t *testing.T, ctx sdk.Context, testName string, accKey types.AccountKey, reward model.Reward) {
	accStorage := model.NewAccountStorage(testAccountKVStoreKey)
	rewardPtr, err := accStorage.GetReward(ctx, accKey)
	if err != nil {
		t.Errorf("%s, failed to get reward, got err %v", testName, err)
	}
	if !assert.Equal(t, reward, *rewardPtr) {
		t.Errorf("%s: diff reward, got %v, want %v", testName, *rewardPtr, reward)
	}
}

func checkRewardHistory(
	t *testing.T, ctx sdk.Context, testName string, accKey types.AccountKey, bucketSlot int64, wantNumOfReward int) {
	accStorage := model.NewAccountStorage(testAccountKVStoreKey)
	rewardHistoryPtr, err := accStorage.GetRewardHistory(ctx, accKey, bucketSlot)
	if err != nil {
		t.Errorf("%s, failed to get reward history, got err %v", testName, err)
	}
	if wantNumOfReward != len(rewardHistoryPtr.Details) {
		t.Errorf("%s: diff account rewards, got %v, want %v", testName, len(rewardHistoryPtr.Details), wantNumOfReward)
	}
}
