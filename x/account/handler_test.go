package account

import (
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var (
	user1 = types.AccountKey("user1")
	user2 = types.AccountKey("user2")

	memo = "This is a memo!"
)

func TestFollow(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, gm)

	// create two test users
	createTestAccount(ctx, am, "user1")
	createTestAccount(ctx, am, "user2")

	// let user1 follows user2
	msg := NewFollowMsg("user1", "user2")
	result := handler(ctx, msg)
	assert.Equal(t, sdk.Result{}, result)

	// check user1 in the user2's follower list
	assert.True(t, am.IsMyFollowing(ctx, types.AccountKey("user1"), types.AccountKey("user2")))

	// check user2 in the user1's following list
	assert.True(t, am.IsMyFollower(ctx, types.AccountKey("user2"), types.AccountKey("user1")))
}

func TestFollowUserNotExist(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, gm)

	// create test user
	createTestAccount(ctx, am, "user1")

	// let user2(not exists) follows user1
	msg := NewFollowMsg("user2", "user1")
	result := handler(ctx, msg)

	assert.Equal(t, result, ErrFollowerNotFound("user2").Result())
	assert.False(t, am.IsMyFollower(ctx, types.AccountKey("user1"), types.AccountKey("user2")))

	// let user1 follows user3(not exists)
	msg = NewFollowMsg("user1", "user3")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrFolloweeNotFound("user3").Result())
	assert.False(t, am.IsMyFollowing(ctx, types.AccountKey("user1"), types.AccountKey("user3")))
}

func TestFollowAgain(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, gm)

	// create two test users
	createTestAccount(ctx, am, "user1")
	createTestAccount(ctx, am, "user2")

	// let user1 follows user2 twice
	msg := NewFollowMsg("user1", "user2")
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	msg = NewFollowMsg("user1", "user2")
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// check user1 is user2's only follower
	assert.True(t, am.IsMyFollower(ctx, types.AccountKey("user2"), types.AccountKey("user1")))

	// check user2 is the only one in the user1's following list
	assert.True(t, am.IsMyFollowing(ctx, types.AccountKey("user1"), types.AccountKey("user2")))
}

func TestUnfollow(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, gm)

	// create two test users
	createTestAccount(ctx, am, "user1")
	createTestAccount(ctx, am, "user2")

	// let user1 follows user2
	msg := NewFollowMsg("user1", "user2")
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// let user1 unfollows user2
	msg2 := NewUnfollowMsg("user1", "user2")
	result = handler(ctx, msg2)
	assert.Equal(t, result, sdk.Result{})

	// check user1 is not in the user2's follower list
	assert.False(t, am.IsMyFollower(ctx, types.AccountKey("user2"), types.AccountKey("user1")))

	// check user2 is not in the user1's following list
	assert.False(t, am.IsMyFollowing(ctx, types.AccountKey("user1"), types.AccountKey("user2")))
}

func TestUnfollowUserNotExist(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, gm)
	// create test user
	createTestAccount(ctx, am, "user1")

	// let user2(not exists) unfollows user1
	msg := NewUnfollowMsg("user2", "user1")
	result := handler(ctx, msg)
	assert.Equal(t, result, ErrFollowerNotFound("user2").Result())

	// let user1 unfollows user3(not exists)
	msg = NewUnfollowMsg("user1", "user3")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrFolloweeNotFound("user3").Result())
}

func TestInvalidUnfollow(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, gm)
	// create test user
	createTestAccount(ctx, am, "user1")
	createTestAccount(ctx, am, "user2")
	createTestAccount(ctx, am, "user3")

	// let user1 follows user2
	msg := NewFollowMsg("user1", "user2")
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// let user3 unfollows user1 and user2 unfollows user3 (invalid)
	//this won't make any changes
	msg2 := NewUnfollowMsg("user3", "user1")
	result = handler(ctx, msg2)
	assert.Equal(t, result, sdk.Result{})

	msg3 := NewUnfollowMsg("user2", "user3")
	result = handler(ctx, msg3)
	assert.Equal(t, result, sdk.Result{})

	// check user1 in the user2's follower list
	assert.True(t, am.IsMyFollower(ctx, types.AccountKey("user2"), types.AccountKey("user1")))

	// check user2 in the user1's following list
	assert.True(t, am.IsMyFollowing(ctx, types.AccountKey("user1"), types.AccountKey("user2")))

}

func TestTransferNormal(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, gm)

	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	// create two test users with initial deposit of 100 LNO.
	createTestAccount(ctx, am, "user1")
	createTestAccount(ctx, am, "user2")

	am.AddSavingCoin(
		ctx, types.AccountKey("user1"), c2000, "", "", types.TransferIn)

	testCases := []struct {
		testName            string
		msg                 TransferMsg
		wantOK              bool
		wantSenderBalance   types.Coin
		wantReceiverBalance types.Coin
	}{
		{
			testName:            "user1 transfers 200 LNO to user2 (by username)",
			msg:                 NewTransferMsg("user1", "user2", l200, memo),
			wantOK:              true,
			wantSenderBalance:   c1800.Plus(accParam.RegisterFee),
			wantReceiverBalance: c200.Plus(accParam.RegisterFee),
		},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.msg)

		if result.IsOK() != tc.wantOK {
			t.Errorf("%s diff result, got %v, want %v", tc.testName, result.IsOK(), tc.wantOK)
		}

		senderSaving, _ := am.GetSavingFromBank(ctx, tc.msg.Sender)
		receiverSaving, _ := am.GetSavingFromBank(ctx, tc.msg.Receiver)

		if !senderSaving.IsEqual(tc.wantSenderBalance) {
			t.Errorf("%s: diff sender saving, got %v, want %v", tc.testName, senderSaving, tc.wantSenderBalance)
		}
		if !receiverSaving.IsEqual(tc.wantReceiverBalance) {
			t.Errorf("%s: diff receiver saving, got %v, want %v", tc.testName, receiverSaving, tc.wantReceiverBalance)
		}
	}
}

func TestSenderCoinNotEnough(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, gm)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)

	// create two test users
	createTestAccount(ctx, am, "user1")
	createTestAccount(ctx, am, "user2")

	memo := "This is a memo!"

	// let user1 transfers 2000 to user2
	msg := NewTransferMsg("user1", "user2", l2000, memo)
	result := handler(ctx, msg)
	assert.Equal(t, ErrAccountSavingCoinNotEnough().Result(), result)

	acc1Balance, _ := am.GetSavingFromBank(ctx, types.AccountKey("user1"))
	assert.Equal(t, acc1Balance, accParam.RegisterFee)
}

func TestReceiverUsernameIncorrect(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, gm)

	// create two test users
	createTestAccount(ctx, am, "user1")

	memo := "This is a memo!"

	// let user1 transfers 2000 to a random user
	msg := NewTransferMsg("user1", "dnqwondqowindow", l2000, memo)
	result := handler(ctx, msg)
	assert.Equal(t, ErrReceiverNotFound("dnqwondqowindow").Result().Code, result.Code)
}

func TestHandleAccountRecover(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, gm)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	user1 := "user1"

	createTestAccount(ctx, am, user1)

	testCases := map[string]struct {
		user              string
		newResetKey       crypto.PubKey
		newTransactionKey crypto.PubKey
		newAppKey         crypto.PubKey
	}{
		"normal case": {
			user:              user1,
			newResetKey:       secp256k1.GenPrivKey().PubKey(),
			newTransactionKey: secp256k1.GenPrivKey().PubKey(),
			newAppKey:         secp256k1.GenPrivKey().PubKey(),
		},
	}

	for testName, tc := range testCases {
		msg := NewRecoverMsg(tc.user, tc.newResetKey, tc.newTransactionKey, tc.newAppKey)
		result := handler(ctx, msg)
		if !assert.Equal(t, sdk.Result{}, result) {
			t.Errorf("%s: diff result, got %v, want %v", testName, result, sdk.Result{})
		}

		accInfo := model.AccountInfo{
			Username:       types.AccountKey(tc.user),
			CreatedAt:      ctx.BlockHeader().Time.Unix(),
			ResetKey:       tc.newResetKey,
			TransactionKey: tc.newTransactionKey,
			AppKey:         tc.newAppKey,
		}
		checkAccountInfo(t, ctx, testName, types.AccountKey(tc.user), accInfo)

		newBank := model.AccountBank{
			Saving:  accParam.RegisterFee,
			NumOfTx: 1,
			CoinDay: accParam.RegisterFee,
		}
		checkBankKVByUsername(t, ctx, testName, types.AccountKey(tc.user), newBank)
	}
}

func TestHandleRegister(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)

	handler := NewHandler(am, gm)
	referrer := "referrer"

	createTestAccount(ctx, am, referrer)
	am.AddSavingCoin(
		ctx, types.AccountKey(referrer), types.NewCoinFromInt64(100*types.Decimals),
		"", "", types.TransferIn)

	testCases := []struct {
		testName                string
		registerMsg             RegisterMsg
		expectResult            sdk.Result
		expectReferrerSaving    types.Coin
		expectNewAccountSaving  types.Coin
		expectNewAccountCoinDay types.Coin
	}{
		{
			testName: "normal case",
			registerMsg: NewRegisterMsg(
				"referrer", "user1", "1",
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
			),
			expectResult:            sdk.Result{},
			expectReferrerSaving:    c100,
			expectNewAccountSaving:  c0,
			expectNewAccountCoinDay: c0,
		},
		{
			testName: "account already exist",
			registerMsg: NewRegisterMsg(
				"referrer", "user1", "1",
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
			),
			expectResult:            ErrAccountAlreadyExists("user1").Result(),
			expectReferrerSaving:    types.NewCoinFromInt64(99 * types.Decimals),
			expectNewAccountSaving:  c0,
			expectNewAccountCoinDay: c0,
		},
		{
			testName: "account register fee insufficient",
			registerMsg: NewRegisterMsg(
				"referrer", "user2", "0.1",
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
			),
			expectResult:            ErrRegisterFeeInsufficient().Result(),
			expectReferrerSaving:    types.NewCoinFromInt64(99 * types.Decimals),
			expectNewAccountSaving:  c0,
			expectNewAccountCoinDay: c0,
		},
		{
			testName: "referrer deposit insufficient",
			registerMsg: NewRegisterMsg(
				"referrer", "user2", "1000",
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
			),
			expectResult:            ErrAccountSavingCoinNotEnough().Result(),
			expectReferrerSaving:    types.NewCoinFromInt64(99 * types.Decimals),
			expectNewAccountSaving:  c0,
			expectNewAccountCoinDay: c0,
		},
		{
			testName: "register deposit larger than register fee but less than full coin day limit",
			registerMsg: NewRegisterMsg(
				"referrer", "user3", "1.5",
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
			),
			expectResult:            sdk.Result{},
			expectReferrerSaving:    types.NewCoinFromInt64(9750000),
			expectNewAccountSaving:  types.NewCoinFromInt64(50000),
			expectNewAccountCoinDay: types.NewCoinFromInt64(50000),
		},
		{
			testName: "register deposit larger than register fee and full coin day limit",
			registerMsg: NewRegisterMsg(
				"referrer", "user4", "2.5",
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
			),
			expectResult:            sdk.Result{},
			expectReferrerSaving:    types.NewCoinFromInt64(95 * types.Decimals),
			expectNewAccountSaving:  types.NewCoinFromInt64(150000),
			expectNewAccountCoinDay: types.NewCoinFromInt64(1 * types.Decimals),
		},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.registerMsg)
		if !assert.Equal(t, tc.expectResult, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectResult)
		}

		if result.Code == sdk.ABCICodeOK {
			if !am.DoesAccountExist(ctx, tc.registerMsg.NewUser) {
				t.Errorf("%s: account %s doesn't exist", tc.testName, tc.registerMsg.NewUser)
			}

			resetKey, err := am.GetResetKey(ctx, tc.registerMsg.NewUser)
			if err != nil {
				t.Errorf("%s: failed to get reset key, got err %v", tc.testName, err)
			}
			if !resetKey.Equals(tc.registerMsg.NewResetPubKey) {
				t.Errorf("%s: diff reset key, got %v, want %v", tc.testName, resetKey, tc.registerMsg.NewResetPubKey)
			}

			txKey, err := am.GetTransactionKey(ctx, tc.registerMsg.NewUser)
			if err != nil {
				t.Errorf("%s: failed to get transaction key, got err %v", tc.testName, err)
			}
			if !txKey.Equals(tc.registerMsg.NewTransactionPubKey) {
				t.Errorf("%s: diff transaction key, got %v, want %v", tc.testName, txKey, tc.registerMsg.NewTransactionPubKey)
			}

			appKey, err := am.GetAppKey(ctx, tc.registerMsg.NewUser)
			if err != nil {
				t.Errorf("%s: failed to get app key, got err %v", tc.testName, err)
			}
			if !appKey.Equals(tc.registerMsg.NewAppPubKey) {
				t.Errorf("%s: diff app key, got %v, want %v", tc.testName, appKey, tc.registerMsg.NewAppPubKey)
			}

			bank, err := am.storage.GetBankFromAccountKey(ctx, tc.registerMsg.NewUser)
			if err != nil {
				t.Errorf("%s: failed to get bank, got err %v", tc.testName, err)
			}
			if !bank.Saving.IsEqual(tc.expectNewAccountSaving) {
				t.Errorf("%s: diff saving, got %v, want %v", tc.testName, bank.Saving, tc.expectNewAccountSaving)
			}
			if !bank.CoinDay.IsEqual(tc.expectNewAccountCoinDay) {
				t.Errorf("%s: diff coin day, got %v, want %v", tc.testName, bank.Saving, tc.expectNewAccountSaving)
			}

			accMeta, _ := am.storage.GetMeta(ctx, tc.registerMsg.NewUser)
			if !accMeta.TransactionCapacity.IsEqual(tc.expectNewAccountCoinDay) {
				t.Errorf("%s: diff transaction capacity, got %v, want %v", tc.testName, accMeta.TransactionCapacity, tc.expectNewAccountCoinDay)
			}
			pool, err := gm.GetDeveloperMonthlyInflation(ctx)
			if err != nil {
				t.Errorf("%s: failed to get inflation, got err %v", tc.testName, err)
			}
			if !pool.IsEqual(accParam.RegisterFee) {
				t.Errorf("%s: diff developer inflation, got %v, want %v", tc.testName, pool, accParam.RegisterFee)
			}
		}

		saving, err := am.GetSavingFromBank(ctx, tc.registerMsg.Referrer)
		if err != nil {
			t.Errorf("%s: failed to get saving from bank, got err %v", tc.testName, err)
		}
		if !saving.IsEqual(tc.expectReferrerSaving) {
			t.Errorf("%s: diff saving, got %v, want %v", tc.testName, saving, tc.expectReferrerSaving)
		}
		gm.GetDeveloperMonthlyInflation(ctx)
	}
}

func TestHandleUpdateAccountMsg(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, gm)

	createTestAccount(ctx, am, "accKey")

	testCases := []struct {
		testName         string
		updateAccountMsg UpdateAccountMsg
		expectResult     sdk.Result
	}{
		{
			testName:         "normal update",
			updateAccountMsg: NewUpdateAccountMsg("accKey", "{'link':'https://lino.network'}"),
			expectResult:     sdk.Result{},
		},
		{
			testName:         "invalid username",
			updateAccountMsg: NewUpdateAccountMsg("invalid", "{'link':'https://lino.network'}"),
			expectResult:     ErrAccountNotFound("invalid").Result(),
		},
	}
	for _, tc := range testCases {
		result := handler(ctx, tc.updateAccountMsg)
		if !assert.Equal(t, result, tc.expectResult) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectResult)
		}
	}
}
