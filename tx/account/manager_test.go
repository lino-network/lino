package account

import (
	"testing"
	"time"

	"github.com/lino-network/lino/tx/account/model"
	"github.com/lino-network/lino/types"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
)

func checkBankKVByAddress(
	t *testing.T, ctx sdk.Context, addr sdk.Address, bank model.AccountBank) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	bankPtr, err := accStorage.GetBankFromAddress(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, bank, *bankPtr, "bank should be equal")
}

func checkPendingStake(
	t *testing.T, ctx sdk.Context, addr sdk.Address, pendingStakeQueue model.PendingStakeQueue) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	pendingStakeQueuePtr, err := accStorage.GetPendingStakeQueue(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, pendingStakeQueue, *pendingStakeQueuePtr, "pending stake should be equal")
}

func checkAccountInfo(
	t *testing.T, ctx sdk.Context, accKey types.AccountKey, accInfo model.AccountInfo) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	infoPtr, err := accStorage.GetInfo(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, accInfo, *infoPtr, "accout info should be equal")
}

func checkAccountMeta(
	t *testing.T, ctx sdk.Context, accKey types.AccountKey, accMeta model.AccountMeta) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	metaPtr, err := accStorage.GetMeta(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, accMeta, *metaPtr, "accout meta should be equal")
}

func checkAccountReward(
	t *testing.T, ctx sdk.Context, accKey types.AccountKey, reward model.Reward) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	rewardPtr, err := accStorage.GetReward(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, reward, *rewardPtr, "accout reward should be equal")
}

func checkAccountGrantKeyList(
	t *testing.T, ctx sdk.Context, accKey types.AccountKey, grantList model.GrantKeyList) {
	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	grantListPtr, err := accStorage.GetGrantKeyList(ctx, accKey)
	assert.Nil(t, err)
	assert.Equal(t, grantList, *grantListPtr, "accout grantList should be equal")
}

func TestIsAccountExist(t *testing.T) {
	ctx, am := setupTest(t, 1)
	assert.False(t, am.IsAccountExist(ctx, types.AccountKey("user1")))
	createTestAccount(ctx, am, "user1")
	assert.True(t, am.IsAccountExist(ctx, types.AccountKey("user1")))
}

func TestAddCoinToAddress(t *testing.T) {
	ctx, am := setupTest(t, 1)
	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)

	// add coin to non-exist account
	err = am.AddCoinToAddress(ctx, sdk.Address("test"), coin1)
	assert.Nil(t, err)

	bank := model.AccountBank{
		Address: sdk.Address("test"),
		Balance: coin1,
	}
	checkBankKVByAddress(t, ctx, sdk.Address("test"), bank)
	pendingStakeQueue := model.PendingStakeQueue{
		LastUpdatedAt:    ctx.BlockHeader().Time,
		StakeCoinInQueue: sdk.ZeroRat,
		TotalCoin:        coin1,
		PendingStakeList: []model.PendingStake{model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake,
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
			EndTime:   ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake,
			Coin:      coin100,
		})
	pendingStakeQueue.TotalCoin = types.NewCoin(101)
	checkPendingStake(t, ctx, sdk.Address("test"), pendingStakeQueue)

	// add coin to exist bank after previous coin day
	ctx = ctx.WithBlockHeader(
		abci.Header{ChainID: "Lino", Height: 3,
			Time: (ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake + 1)})
	err = am.AddCoinToAddress(ctx, sdk.Address("test"), coin100)
	assert.Nil(t, err)
	bank = model.AccountBank{
		Address: sdk.Address("test"),
		Balance: types.NewCoin(201),
		Stake:   types.NewCoin(101),
	}
	checkBankKVByAddress(t, ctx, sdk.Address("test"), bank)
	pendingStakeQueue.PendingStakeList = []model.PendingStake{model.PendingStake{
		StartTime: ctx.BlockHeader().Time,
		EndTime:   ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake,
		Coin:      coin100,
	}}
	pendingStakeQueue.TotalCoin = coin100
	pendingStakeQueue.LastUpdatedAt = ctx.BlockHeader().Time
	checkPendingStake(t, ctx, sdk.Address("test"), pendingStakeQueue)
}

func TestCreateAccount(t *testing.T) {
	ctx, am := setupTest(t, 1)
	priv := crypto.GenPrivKeyEd25519()
	accKey := types.AccountKey("accKey")
	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)

	// normal test
	assert.False(t, am.IsAccountExist(ctx, accKey))
	err = am.AddCoinToAddress(ctx, priv.PubKey().Address(), coin100)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey,
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), coin0)
	assert.Nil(t, err)

	assert.True(t, am.IsAccountExist(ctx, accKey))
	bank := model.AccountBank{
		Address:  priv.PubKey().Address(),
		Balance:  coin100,
		Username: accKey,
	}
	checkBankKVByAddress(t, ctx, priv.PubKey().Address(), bank)
	pendingStakeQueue := model.PendingStakeQueue{
		LastUpdatedAt:    ctx.BlockHeader().Time,
		StakeCoinInQueue: sdk.ZeroRat,
		TotalCoin:        coin100,
		PendingStakeList: []model.PendingStake{model.PendingStake{
			StartTime: ctx.BlockHeader().Time,
			EndTime:   ctx.BlockHeader().Time + coinDayParams.SecondsToRecoverCoinDayStake,
			Coin:      coin100,
		}}}
	checkPendingStake(t, ctx, priv.PubKey().Address(), pendingStakeQueue)
	accInfo := model.AccountInfo{
		Username:       accKey,
		CreatedAt:      ctx.BlockHeader().Time,
		MasterKey:      priv.PubKey(),
		TransactionKey: priv.Generate(1).PubKey(),
		PostKey:        priv.Generate(2).PubKey(),
		Address:        priv.PubKey().Address(),
	}
	checkAccountInfo(t, ctx, accKey, accInfo)
	accMeta := model.AccountMeta{
		LastActivityAt: ctx.BlockHeader().Time,
	}
	checkAccountMeta(t, ctx, accKey, accMeta)

	reward := model.Reward{coin0, coin0, coin0, coin0}
	checkAccountReward(t, ctx, accKey, reward)

	var grantPubKeyList []model.GrantPubKey
	grantList := model.GrantKeyList{GrantPubKeyList: grantPubKeyList}
	checkAccountGrantKeyList(t, ctx, accKey, grantList)

	// username already took
	err = am.CreateAccount(ctx, accKey,
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), coin0)
	assert.Equal(t, ErrAccountAlreadyExists(accKey), err)

	// bank already registered
	err = am.CreateAccount(ctx, types.AccountKey("newKey"),
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), coin0)
	assert.Equal(t, ErrBankAlreadyRegistered(), err)

	// bank doesn't exist
	priv2 := crypto.GenPrivKeyEd25519()
	err = am.CreateAccount(ctx, types.AccountKey("newKey"),
		priv2.PubKey(), priv2.Generate(1).PubKey(), priv2.Generate(2).PubKey(), coin0)
	assert.Equal(t,
		"Error{311:create account newKey failed,Error{310:account bank is not found,<nil>,0},1}",
		err.Error())

	// register fee doesn't enough
	err = am.AddCoinToAddress(ctx, priv2.PubKey().Address(), coin100)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, types.AccountKey("newKey"),
		priv2.PubKey(), priv2.Generate(1).PubKey(), priv2.Generate(2).PubKey(), types.NewCoin(101))
	assert.Equal(t, ErrRegisterFeeInsufficient(), err)
}

func TestCoinDayByAddress(t *testing.T) {
	ctx, am := setupTest(t, 1)
	accKey := types.AccountKey("accKey")
	priv := crypto.GenPrivKeyEd25519()

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)
	totalCoinDaysSec := coinDayParams.SecondsToRecoverCoinDayStake

	// create bank and account
	err = am.AddCoinToAddress(ctx, priv.PubKey().Address(), coin100)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey,
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), coin0)
	assert.Nil(t, err)

	baseTime1 := ctx.BlockHeader().Time
	baseTime2 := baseTime1 + totalCoinDaysSec*2
	cases := []struct {
		AddCoin           types.Coin
		AtWhen            int64
		ExpectBalance     types.Coin
		ExpectStake       types.Coin
		ExpectStakeInBank types.Coin
	}{
		{coin0, baseTime1 + 3024, coin100, coin0, coin0},
		{coin0, baseTime1 + 3025, coin100, coin1, coin0},
		{coin0, baseTime1 + totalCoinDaysSec/2, coin100, coin50, coin0},
		{coin100, baseTime1 + totalCoinDaysSec/2, coin200, coin50, coin0},
		{coin0, baseTime1 + totalCoinDaysSec + 1, coin200, types.NewCoin(150), coin100},
		{coin0, baseTime1 + totalCoinDaysSec*2 + 1, coin200, coin200, coin200},
		{coin1, baseTime2, types.NewCoin(201), coin200, coin200},
		{coin0, baseTime2 + totalCoinDaysSec/2, types.NewCoin(201), coin200, coin200},
		{coin0, baseTime2 + totalCoinDaysSec/2 + 1,
			types.NewCoin(201), types.NewCoin(201), coin200},
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
	accKey := types.AccountKey("accKey")
	priv := crypto.GenPrivKeyEd25519()
	// create bank and account
	err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), coin100)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey,
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), coin0)
	assert.Nil(t, err)

	coinDayParams, err := am.paramHolder.GetCoinDayParam(ctx)
	assert.Nil(t, err)
	totalCoinDaysSec := coinDayParams.SecondsToRecoverCoinDayStake

	baseTime := ctx.BlockHeader().Time
	baseTime2 := baseTime + totalCoinDaysSec + 1000
	baseTime3 := baseTime2 + totalCoinDaysSec + 1000
	baseTime4 := baseTime3 + totalCoinDaysSec*3/2 + 3

	cases := []struct {
		IsAdd             bool
		Coin              types.Coin
		AtWhen            int64
		ExpectBalance     types.Coin
		ExpectStake       types.Coin
		ExpectStakeInBank types.Coin
	}{
		{true, coin0, baseTime + 3024, coin100, coin0, coin0},
		{true, coin0, baseTime + 3025, coin100, coin1, coin0},
		{false, coin100, baseTime + 3457, coin0, coin0, coin0},
		{true, coin0, baseTime + totalCoinDaysSec + 1, coin0, coin0, coin0},

		{true, coin100, baseTime2, coin100, coin0, coin0},
		{false, coin50, baseTime2 + totalCoinDaysSec/2 + 1, coin50, types.NewCoin(25), coin0},
		{true, coin0, baseTime2 + totalCoinDaysSec + 1, coin50, coin50, coin50},

		{true, coin100, baseTime3, types.NewCoin(150), coin50, coin50},
		{true, coin100, baseTime3 + totalCoinDaysSec/2 + 1, types.NewCoin(250), coin100, coin50},
		{false, coin50, baseTime3 + totalCoinDaysSec*3/4 + 2,
			coin200, types.NewCoin(138), types.NewCoin(50)},
		{true, coin0, baseTime3 + totalCoinDaysSec + 2,
			coin200, types.NewCoin(175), types.NewCoin(150)},
		{true, coin0, baseTime3 + totalCoinDaysSec*3/2 + 2, coin200, coin200, coin200},

		{true, coin1, baseTime4, types.NewCoin(201), coin200, coin200},
		{true, coin0, baseTime4 + totalCoinDaysSec/2 + 1,
			types.NewCoin(201), types.NewCoin(201), coin200},
		{false, coin1, baseTime4 + totalCoinDaysSec/2 + 1, coin200, coin200, coin200},
		{true, coin0, baseTime4 + totalCoinDaysSec + 1, coin200, coin200, coin200},
		{true, coin0, baseTime4 + totalCoinDaysSec*100 + 1, coin200, coin200, coin200},
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
	accKey := types.AccountKey("accKey")
	priv := crypto.GenPrivKeyEd25519()

	err := am.AddCoinToAddress(ctx, priv.PubKey().Address(), c100)
	assert.Nil(t, err)
	err = am.CreateAccount(ctx, accKey,
		priv.PubKey(), priv.Generate(1).PubKey(), priv.Generate(2).PubKey(), coin0)
	assert.Nil(t, err)

	err = am.AddIncomeAndReward(ctx, accKey, c500, c200, c300)
	assert.Nil(t, err)
	reward := model.Reward{c500, c200, c300, c300}
	checkAccountReward(t, ctx, accKey, reward)
	err = am.AddIncomeAndReward(ctx, accKey, c500, c300, c200)
	assert.Nil(t, err)
	reward = model.Reward{c1000, c500, c500, c500}
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
	reward = model.Reward{c1000, c500, c500, c0}
	checkAccountReward(t, ctx, accKey, reward)
}

func TestCheckUserTPSCapacity(t *testing.T) {
	ctx, am := setupTest(t, 1)
	accKey := types.AccountKey("accKey")

	bandwidthParams, err := am.paramHolder.GetBandwidthParam(ctx)
	assert.Nil(t, err)
	secondsToRecoverBandwidth := bandwidthParams.SecondsToRecoverBandwidth

	baseTime := ctx.BlockHeader().Time

	priv := createTestAccount(ctx, am, string(accKey))
	err = am.AddCoinToAddress(ctx, priv.PubKey().Address(), c100)
	assert.Nil(t, err)

	accStorage := model.NewAccountStorage(TestAccountKVStoreKey)
	err = accStorage.SetPendingStakeQueue(
		ctx, priv.PubKey().Address(), &model.PendingStakeQueue{})
	assert.Nil(t, err)

	cases := []struct {
		TPSCapacityRatio     sdk.Rat
		UserStake            types.Coin
		LastActivity         int64
		LastCapacity         types.Coin
		CurrentTime          int64
		ExpectResult         sdk.Error
		ExpectRemainCapacity types.Coin
	}{
		{sdk.NewRat(1, 10), types.NewCoin(10 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime, ErrAccountTPSCapacityNotEnough(accKey), types.NewCoin(0)},
		{sdk.NewRat(1, 10), types.NewCoin(10 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime + secondsToRecoverBandwidth, nil, types.NewCoin(990000)},
		{sdk.NewRat(1, 2), types.NewCoin(10 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime + secondsToRecoverBandwidth, nil, types.NewCoin(950000)},
		{sdk.NewRat(1), types.NewCoin(10 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime + secondsToRecoverBandwidth, nil, types.NewCoin(9 * types.Decimals)},
		{sdk.NewRat(1), types.NewCoin(1 * types.Decimals), baseTime,
			types.NewCoin(10 * types.Decimals), baseTime, nil, types.NewCoin(0)},
		{sdk.NewRat(1), types.NewCoin(10), baseTime, types.NewCoin(1 * types.Decimals),
			baseTime, ErrAccountTPSCapacityNotEnough(accKey), types.NewCoin(1 * types.Decimals)},
		{sdk.NewRat(1), types.NewCoin(1 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime + secondsToRecoverBandwidth/2,
			ErrAccountTPSCapacityNotEnough(accKey), types.NewCoin(0)},
		{sdk.NewRat(1, 2), types.NewCoin(1 * types.Decimals), baseTime, types.NewCoin(0),
			baseTime + secondsToRecoverBandwidth/2, nil, types.NewCoin(0)},
		{sdk.OneRat, types.NewCoin(1 * types.Decimals), 0, types.NewCoin(0),
			baseTime, nil, types.NewCoin(0)},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Time: cs.CurrentTime})
		bank := &model.AccountBank{
			Address: priv.PubKey().Address(),
			Balance: cs.UserStake,
			Stake:   cs.UserStake,
		}
		err = accStorage.SetBankFromAddress(ctx, priv.PubKey().Address(), bank)
		assert.Nil(t, err)
		meta := &model.AccountMeta{
			LastActivityAt:      cs.LastActivity,
			TransactionCapacity: cs.LastCapacity,
		}
		err = accStorage.SetMeta(ctx, accKey, meta)
		assert.Nil(t, err)

		err = am.CheckUserTPSCapacity(ctx, accKey, cs.TPSCapacityRatio)
		assert.Equal(t, cs.ExpectResult, err)

		accMeta := model.AccountMeta{
			LastActivityAt:      ctx.BlockHeader().Time,
			TransactionCapacity: cs.ExpectRemainCapacity,
		}
		if cs.ExpectResult != nil {
			accMeta.LastActivityAt = cs.LastActivity
		}
		checkAccountMeta(t, ctx, accKey, accMeta)
	}
}

func TestGrantPubkey(t *testing.T) {
	ctx, am := setupTest(t, 1)
	user1 := types.AccountKey("user1")
	user2 := types.AccountKey("user2")
	user3 := types.AccountKey("user3")

	createTestAccount(ctx, am, string(user1))
	priv2 := createTestAccount(ctx, am, string(user2))
	priv3 := createTestAccount(ctx, am, string(user3))

	baseTime := ctx.BlockHeader().Time

	cases := []struct {
		user             types.AccountKey
		grantTo          types.AccountKey
		expireTime       int64
		checkTime        int64
		checkGrantUser   types.AccountKey
		checkGrantPubKey crypto.PubKey
		expectResult     sdk.Error
	}{
		{user1, user2, 100, baseTime + 99, user2, priv2.Generate(2).PubKey(), nil},
		{user1, user3, 100, baseTime + 99, user3, priv3.Generate(2).PubKey(), nil},
		{user1, user2, 100, baseTime + 101, user2, priv2.Generate(2).PubKey(),
			ErrCheckAuthenticatePubKeyOwner(user1)},
		{user1, user2, 100, baseTime + 99, user2, priv2.Generate(2).PubKey(), nil},
		{user1, user2, 500, baseTime + 101, user2, priv2.Generate(2).PubKey(), nil},
		{user1, user2, 300, baseTime + 301, user2, priv2.Generate(2).PubKey(),
			ErrCheckAuthenticatePubKeyOwner(user1)},
	}

	for _, cs := range cases {
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 1, Time: baseTime})
		err := am.AuthorizePermission(ctx, cs.user, cs.grantTo, cs.expireTime, 0)
		assert.Nil(t, err)
		ctx = ctx.WithBlockHeader(abci.Header{ChainID: "Lino", Height: 2, Time: cs.checkTime})
		grantUser, err := am.CheckAuthenticatePubKeyOwner(ctx, cs.user, cs.checkGrantPubKey, 0)
		assert.Equal(t, err, cs.expectResult)
		if cs.expectResult == nil {
			assert.Equal(t, grantUser, cs.checkGrantUser)
		}
	}
}

func TestDonationRelationship(t *testing.T) {
	ctx, am := setupTest(t, 1)
	user1 := types.AccountKey("user1")
	user2 := types.AccountKey("user2")
	user3 := types.AccountKey("user3")

	createTestAccount(ctx, am, string(user1))
	createTestAccount(ctx, am, string(user2))
	createTestAccount(ctx, am, string(user3))

	cases := []struct {
		user             types.AccountKey
		donateTo         types.AccountKey
		expectDonateTime int64
	}{
		{user1, user2, 1},
		{user1, user2, 2},
		{user1, user3, 1},
		{user3, user1, 1},
		{user2, user1, 1},
	}

	for _, cs := range cases {
		err := am.UpdateDonationRelationship(ctx, cs.user, cs.donateTo)
		assert.Nil(t, err)
		donateTime, err := am.GetDonationRelationship(ctx, cs.user, cs.donateTo)
		assert.Nil(t, err)
		assert.Equal(t, donateTime, cs.expectDonateTime)
	}
}

func TestAccountRecover(t *testing.T) {
	ctx, am := setupTest(t, 1)
	user1 := types.AccountKey("user1")

	priv := createTestAccount(ctx, am, string(user1))

	cases := []struct {
		user              types.AccountKey
		newPostKey        crypto.PubKey
		newTransactionKey crypto.PubKey
	}{
		{user1, crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey()},
	}

	for _, cs := range cases {
		err := am.RecoverAccount(ctx, cs.user, cs.newPostKey, cs.newTransactionKey)
		assert.Nil(t, err)
		accInfo := model.AccountInfo{
			Username:       cs.user,
			CreatedAt:      ctx.BlockHeader().Time,
			MasterKey:      priv.PubKey(),
			TransactionKey: cs.newTransactionKey,
			PostKey:        cs.newPostKey,
			Address:        priv.PubKey().Address(),
		}
		checkAccountInfo(t, ctx, cs.user, accInfo)
	}
}

func TestIncreaseSequenceByOne(t *testing.T) {
	ctx, am := setupTest(t, 1)
	user1 := types.AccountKey("user1")

	createTestAccount(ctx, am, string(user1))

	cases := []struct {
		user           types.AccountKey
		increaseTimes  int
		expectSequence int64
	}{
		{user1, 1, 1},
		{user1, 100, 101},
	}

	for _, cs := range cases {

		for i := 0; i < cs.increaseTimes; i++ {
			am.IncreaseSequenceByOne(ctx, user1)
		}
		seq, err := am.GetSequence(ctx, user1)
		assert.Nil(t, err)
		assert.Equal(t, cs.expectSequence, seq)
	}
}
