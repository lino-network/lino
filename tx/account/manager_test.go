package account

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/account/model"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-crypto"
)

func checkBankKVByAddress(t *testing.T, ctx sdk.Context, addr sdk.Address, bank model.AccountBank) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	bankPtr, err := accStorage.GetBankFromAddress(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, bank, *bankPtr, "bank should be equal")
}

func checkPendingStake(t *testing.T, ctx sdk.Context, addr sdk.Address, pendingStakeQueue model.PendingStakeQueue) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	pendingStakeQueuePtr, err := accStorage.GetPendingStakeQueue(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, pendingStakeQueue, *pendingStakeQueuePtr, "pending stake should be equal")
}

func checkAccountInfo(t *testing.T, ctx sdk.Context, accKey types.AccountKey, accInfo model.AccountInfo) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	infoPtr, err := accStorage.GetInfo(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, accInfo, *infoPtr, "accout info should be equal")
}

func checkAccountMeta(t *testing.T, ctx sdk.Context, accKey types.AccountKey, accMeta model.AccountMeta) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	metaPtr, err := accStorage.GetMeta(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, accMeta, *metaPtr, "accout meta should be equal")
}

func checkAccountReward(t *testing.T, ctx sdk.Context, accKey types.AccountKey, reward model.Reward) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	rewardPtr, err := accStorage.GetReward(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, reward, *rewardPtr, "accout reward should be equal")
}

func TestIsAccountExist(t *testing.T) {
	ctx, am := setupTest(t, 1)
	createTestAccount(ctx, am, "user1")
	assert.True(t, am.IsAccountExist(ctx, types.AccountKey("user1")))
}

func TestAddCoinToAddress(t *testing.T) {
	ctx, am := setupTest(t, 1)

	// add coin to non-exist account
	err := am.AddCoinToAddress(ctx, sdk.Address("test"), coin1)
	assert.Nil(t, err)

	bank := model.AccountBank{
		Address: sdk.Address("test"),
		Balance: coin1,
	}
	checkBankKVByAddress(t, ctx, sdk.Address("test"), bank)
	pendingStakeQueue := model.PendingStakeQueue{
		PendingStakeList: []model.PendingStake{model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + TotalCoinDaysSec,
			Coin:      coin1,
		}}}
	checkPendingStake(t, ctx, sdk.Address("test"), pendingStakeQueue)

	// add coin to exist bank
	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: time.Now().Unix()})
	err = am.AddCoinToAddress(ctx, sdk.Address("test"), coin100)
	assert.Nil(t, err)
	bank = model.AccountBank{
		Address: sdk.Address("test"),
		Balance: types.NewCoin(101),
	}
	checkBankKVByAddress(t, ctx, sdk.Address("test"), bank)
	pendingStakeQueue.PendingStakeList = append(pendingStakeQueue.PendingStakeList,
		model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + TotalCoinDaysSec,
			Coin:      coin100,
		})
	checkPendingStake(t, ctx, sdk.Address("test"), pendingStakeQueue)

	// add coin to exist bank after previous coin day
	ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 3, Time: (ctx.BlockHeader().Time + 3600*24*CoinDays + 1)})
	err = am.AddCoinToAddress(ctx, sdk.Address("test"), coin100)
	assert.Nil(t, err)
	bank = model.AccountBank{
		Address: sdk.Address("test"),
		Balance: types.NewCoin(201),
	}
	checkBankKVByAddress(t, ctx, sdk.Address("test"), bank)
	pendingStakeQueue.PendingStakeList = append(pendingStakeQueue.PendingStakeList,
		model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + TotalCoinDaysSec,
			Coin:      coin100,
		})
	checkPendingStake(t, ctx, sdk.Address("test"), pendingStakeQueue)
}

func TestCreateAccount(t *testing.T) {
	ctx, am := setupTest(t, 1)
	priv := crypto.GenPrivKeyEd25519()
	accKey := types.AccountKey("accKey")

	// normal test
	assert.False(t, am.IsAccountExist(ctx, accKey))
	err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), coin100)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey, priv.PubKey(), coin0)
	assert.Nil(t, err)

	assert.True(t, am.IsAccountExist(ctx, accKey))
	bank := model.AccountBank{
		Address:  priv.PubKey().Address(),
		Balance:  coin100,
		Username: accKey,
	}
	checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)
	pendingStakeQueue := model.PendingStakeQueue{
		PendingStakeList: []model.PendingStake{model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + TotalCoinDaysSec,
			Coin:      coin100,
		}}}
	checkPendingStake(t, ctx, priv.PubKey().Address(), pendingStakeQueue)
	accInfo := model.AccountInfo{
		Username: accKey,
		Created:  1,
		PostKey:  priv.PubKey(),
		OwnerKey: priv.PubKey(),
		Address:  priv.PubKey().Address(),
	}
	checkAccountInfo(t, ctx, accKey, accInfo)
	accMeta := model.AccountMeta{
		LastActivity:   ctx.BlockHeight(),
		ActivityBurden: types.DefaultActivityBurden,
	}
	checkAccountMeta(t, ctx, accKey, accMeta)

	reward := model.Reward{coin0, coin0, coin0}
	checkAccountReward(t, ctx, accKey, reward)

	// username already took
	err = am.CreateAccount(ctx, accKey, priv.PubKey(), coin0)
	assert.Equal(t, ErrAccountAlreadyExists(accKey), err)

	// bank already registered
	err = am.CreateAccount(ctx, types.AccountKey("newKey"), priv.PubKey(), coin0)
	assert.Equal(t, ErrBankAlreadyRegistered(), err)

	// bank doesn't exist
	priv2 := crypto.GenPrivKeyEd25519()
	err = am.CreateAccount(ctx, types.AccountKey("newKey"), priv2.PubKey(), coin0)
	assert.Equal(t, "Error{311:create account newKey failed,Error{310:account bank doesn't exist,<nil>,0},1}", err.Error())

	// register fee doesn't enough
	err = am.AddCoinToAddress(ctx, priv2.PubKey().Address(), coin100)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, types.AccountKey("newKey"), priv2.PubKey(), types.NewCoin(101))
	assert.Equal(t, ErrRegisterFeeInsufficient(), err)
}

func TestCoinDayByAddress(t *testing.T) {
	ctx, am := setupTest(t, 1)
	priv := crypto.GenPrivKeyEd25519()
	accKey := types.AccountKey("accKey")

	// create bank and account
	err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), coin100)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey, priv.PubKey(), coin0)
	assert.Nil(t, err)

	baseTime1 := ctx.BlockHeader().Time
	baseTime2 := baseTime1 + CoinDays*24*5400 + 1000
	cases := []struct {
		AddCoin           types.Coin
		AtWhen            int64
		ExpectBalance     types.Coin
		ExpectStake       types.Coin
		ExpectStakeInBank types.Coin
	}{
		{coin0, baseTime1 + 3456, coin100, coin0, coin0},
		{coin0, baseTime1 + 3457, coin100, coin1, coin0},
		{coin0, baseTime1 + TotalCoinDaysSec/2, coin100, coin50, coin0},
		{coin100, baseTime1 + TotalCoinDaysSec/2, coin200, coin50, coin0},
		{coin0, baseTime1 + TotalCoinDaysSec + 1, coin200, types.NewCoin(150), coin100},
		{coin0, baseTime1 + CoinDays*24*5400 + 1, coin200, coin200, coin200},
		{coin1, baseTime2, types.NewCoin(201), coin200, coin200},
		{coin0, baseTime2 + TotalCoinDaysSec/2, types.NewCoin(201), coin200, coin200},
		{coin0, baseTime2 + TotalCoinDaysSec/2 + 1, types.NewCoin(201), types.NewCoin(201), coin200},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), cs.AddCoin)
		assert.Nil(t, err)
		coin, err := am.GetStake(ctx, accKey)
		assert.Nil(t, err)
		assert.Equal(t, cs.ExpectStake, coin)

		bank := model.AccountBank{
			Address:  priv.PubKey().Address(),
			Balance:  cs.ExpectBalance,
			Stake:    cs.ExpectStakeInBank,
			Username: accKey,
		}
		checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)
	}
}

func TestCoinDayByAccountKey(t *testing.T) {
	ctx, am := setupTest(t, 1)
	priv := crypto.GenPrivKeyEd25519()
	accKey := types.AccountKey("accKey")

	// create bank and account
	err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), coin100)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey, priv.PubKey(), coin0)
	assert.Nil(t, err)

	baseTime := ctx.BlockHeader().Time
	baseTime2 := baseTime + TotalCoinDaysSec + 1000
	baseTime3 := baseTime2 + TotalCoinDaysSec + 1000
	baseTime4 := baseTime3 + TotalCoinDaysSec*3/2 + 3

	cases := []struct {
		IsAdd             bool
		Coin              types.Coin
		AtWhen            int64
		ExpectBalance     types.Coin
		ExpectStake       types.Coin
		ExpectStakeInBank types.Coin
	}{
		{true, coin0, baseTime + 3456, coin100, coin0, coin0},
		{true, coin0, baseTime + 3457, coin100, coin1, coin0},
		{false, coin100, baseTime + 3457, coin0, coin0, coin0},
		{true, coin0, baseTime + TotalCoinDaysSec + 1, coin0, coin0, coin0},

		{true, coin100, baseTime2, coin100, coin0, coin0},
		{false, coin50, baseTime2 + TotalCoinDaysSec/2 + 1, coin50, coin50, coin50},
		{true, coin0, baseTime2 + TotalCoinDaysSec + 1, coin50, coin50, coin50},

		{true, coin100, baseTime3, types.NewCoin(150), coin50, coin50},
		{true, coin100, baseTime3 + TotalCoinDaysSec/2 + 1, types.NewCoin(250), coin100, coin50},
		{false, coin50, baseTime3 + TotalCoinDaysSec*3/4 + 2, coin200, types.NewCoin(150), types.NewCoin(125)},
		{true, coin0, baseTime3 + TotalCoinDaysSec + 2, coin200, types.NewCoin(175), types.NewCoin(125)},
		{true, coin0, baseTime3 + TotalCoinDaysSec*3/2 + 2, coin200, coin200, coin200},

		{true, coin1, baseTime4, types.NewCoin(201), coin200, coin200},
		{true, coin0, baseTime4 + TotalCoinDaysSec/2 + 1, types.NewCoin(201), types.NewCoin(201), coin200},
		{false, coin1, baseTime4 + TotalCoinDaysSec/2 + 1, coin200, coin200, coin200},
		{true, coin0, baseTime4 + TotalCoinDaysSec + 1, coin200, coin200, coin200},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.AtWhen})
		if cs.IsAdd {
			err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), cs.Coin)
			assert.Nil(t, err)
		} else {
			err := am.MinusCoin(ctx, accKey, cs.Coin)
			assert.Nil(t, err)
		}
		coin, err := am.GetStake(ctx, accKey)
		assert.Nil(t, err)
		assert.Equal(t, cs.ExpectStake, coin)

		bank := model.AccountBank{
			Address:  priv.PubKey().Address(),
			Balance:  cs.ExpectBalance,
			Stake:    cs.ExpectStakeInBank,
			Username: accKey,
		}
		checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)
	}
}

func TestAccountReward(t *testing.T) {
	ctx, am := setupTest(t, 1)
	priv := crypto.GenPrivKeyEd25519()
	accKey := types.AccountKey("accKey")

	// create bank and account
	err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), c100)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey, priv.PubKey(), coin0)
	assert.Nil(t, err)

	err = am.AddIncomeAndReward(ctx, accKey, c200, c300)
	assert.Nil(t, err)
	reward := model.Reward{c200, c300, c300}
	checkAccountReward(t, ctx, accKey, reward)
	err = am.AddIncomeAndReward(ctx, accKey, c300, c200)
	assert.Nil(t, err)
	reward = model.Reward{c500, c500, c500}
	checkAccountReward(t, ctx, accKey, reward)

	bank := model.AccountBank{
		Address:  priv.PubKey().Address(),
		Balance:  c100,
		Stake:    c0,
		Username: accKey,
	}
	checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)

	err = am.ClaimReward(ctx, accKey)
	assert.Nil(t, err)
	bank.Balance = c600
	checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)
	reward = model.Reward{c500, c500, c0}
	checkAccountReward(t, ctx, accKey, reward)
}
