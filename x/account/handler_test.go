package account

import (
	"testing"

	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"
	"github.com/lino-network/lino/x/global"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var (
	user1 = types.AccountKey("user1")
	user2 = types.AccountKey("user2")

	memo = "This is a memo!"
)

func TestTransferNormal(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, &gm)

	accParam, _ := am.paramHolder.GetAccountParam(ctx)
	// create two test users with initial deposit of 100 LNO.
	createTestAccount(ctx, am, "user1")
	createTestAccount(ctx, am, "user2")

	am.AddCoinToUsername(ctx, types.AccountKey("user1"), c2000)

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

		senderSaving, _ := am.GetSavingFromUsername(ctx, tc.msg.Sender)
		receiverSaving, _ := am.GetSavingFromUsername(ctx, tc.msg.Receiver)

		if !senderSaving.IsEqual(tc.wantSenderBalance) {
			t.Errorf("%s: diff sender saving, got %v, want %v", tc.testName, senderSaving, tc.wantSenderBalance)
		}
		if !receiverSaving.IsEqual(tc.wantReceiverBalance) {
			t.Errorf("%s: diff receiver saving, got %v, want %v", tc.testName, receiverSaving, tc.wantReceiverBalance)
		}
	}
}

func BenchmarkNumTransfer(b *testing.B) {
	ctx := getContext(0)
	ph := param.NewParamHolder(testParamKVStoreKey)
	ph.InitParam(ctx)
	accManager := NewAccountManager(testAccountKVStoreKey, ph)
	globalManager := global.NewGlobalManager(testGlobalKVStoreKey, ph)
	handler := NewHandler(accManager, &globalManager)

	// create two test users with initial deposit of 100 LNO.
	createTestAccount(ctx, accManager, "user1")
	createTestAccount(ctx, accManager, "user2")

	accManager.AddCoinToUsername(
		ctx, types.AccountKey("user1"), types.NewCoinFromInt64(100000*int64(b.N)))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		handler(ctx, NewTransferMsg("user1", "user2", "1", ""))
	}
}

func TestSenderCoinNotEnough(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, &gm)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)

	// create two test users
	createTestAccount(ctx, am, "user1")
	createTestAccount(ctx, am, "user2")

	memo := "This is a memo!"

	// let user1 transfers 2000 to user2
	msg := NewTransferMsg("user1", "user2", l2000, memo)
	result := handler(ctx, msg)
	assert.Equal(t, ErrAccountSavingCoinNotEnough().Result(), result)

	acc1Balance, _ := am.GetSavingFromUsername(ctx, types.AccountKey("user1"))
	assert.Equal(t, acc1Balance, accParam.RegisterFee)
}

func TestReceiverUsernameIncorrect(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, &gm)

	// create two test users
	createTestAccount(ctx, am, "user1")
	err := am.AddCoinToUsername(ctx, types.AccountKey("user1"), c2000)
	if err != nil {
		t.Errorf("TestReceiverUsernameIncorrect: failed to add coin to account, got err %v", err)
	}

	memo := "This is a memo!"

	// let user1 transfers 2000 to a random user
	msg := NewTransferMsg("user1", "dnqwondqowindow", l2000, memo)
	result := handler(ctx, msg)
	// fmt.Println(result)
	assert.Equal(t, model.ErrAccountInfoNotFound().Result().Code, result.Code)
}

func TestHandleAccountRecover(t *testing.T) {
	// ctx, am, gm := setupTest(t, 1)
	// handler := NewHandler(am, &gm)
	// accParam, _ := am.paramHolder.GetAccountParam(ctx)
	// user1 := "user1"

	// createTestAccount(ctx, am, user1)

	// testCases := map[string]struct {
	// 	user              string
	// 	newResetKey       crypto.PubKey
	// 	newTransactionKey crypto.PubKey
	// 	newAppKey         crypto.PubKey
	// }{
	// 	"normal case": {
	// 		user:              user1,
	// 		newResetKey:       secp256k1.GenPrivKey().PubKey(),
	// 		newTransactionKey: secp256k1.GenPrivKey().PubKey(),
	// 		newAppKey:         secp256k1.GenPrivKey().PubKey(),
	// 	},
	// }

	// for testName, tc := range testCases {
	// 	msg := NewRecoverMsg(tc.user, tc.newResetKey, tc.newTransactionKey, tc.newAppKey)
	// 	result := handler(ctx, msg)
	// 	if !assert.Equal(t, sdk.Result{}, result) {
	// 		t.Errorf("%s: diff result, got %v, want %v", testName, result, sdk.Result{})
	// 	}

	// 	accInfo := model.AccountInfo{
	// 		Username:       types.AccountKey(tc.user),
	// 		CreatedAt:      ctx.BlockHeader().Time.Unix(),
	// 		SignningKey:    tc.newResetKey,
	// 		TransactionKey: tc.newTransactionKey,
	// 	}
	// 	checkAccountInfo(t, ctx, testName, types.AccountKey(tc.user), accInfo)

	// 	newBank := model.AccountBank{
	// 		Saving:  accParam.RegisterFee,
	// 		CoinDay: accParam.RegisterFee,
	// 	}
	// 	checkBankKVByUsername(t, ctx, testName, types.AccountKey(tc.user), newBank)
	// }
}

func TestHandleRegister(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	accParam, _ := am.paramHolder.GetAccountParam(ctx)

	handler := NewHandler(am, &gm)
	referrer := "referrer"

	createTestAccount(ctx, am, referrer)
	am.AddCoinToUsername(ctx, types.AccountKey(referrer), types.NewCoinFromInt64(100*types.Decimals))

	testCases := []struct {
		testName               string
		registerMsg            RegisterMsg
		expectResult           sdk.Result
		expectReferrerSaving   types.Coin
		expectNewAccountSaving types.Coin
	}{
		{
			testName: "normal case",
			registerMsg: NewRegisterMsg(
				"referrer", "user1", "1",
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
			),
			expectResult:           sdk.Result{},
			expectReferrerSaving:   c100,
			expectNewAccountSaving: c0,
		},
		{
			testName: "account already exist",
			registerMsg: NewRegisterMsg(
				"referrer", "user1", "1",
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
			),
			expectResult:           ErrAccountAlreadyExists("user1").Result(),
			expectReferrerSaving:   types.NewCoinFromInt64(99 * types.Decimals),
			expectNewAccountSaving: c0,
		},
		{
			testName: "account register fee insufficient",
			registerMsg: NewRegisterMsg(
				"referrer", "user2", "0.1",
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
			),
			expectResult:           ErrRegisterFeeInsufficient().Result(),
			expectReferrerSaving:   types.NewCoinFromInt64(99 * types.Decimals),
			expectNewAccountSaving: c0,
		},
		{
			testName: "referrer deposit insufficient",
			registerMsg: NewRegisterMsg(
				"referrer", "user2", "1000",
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
				secp256k1.GenPrivKey().PubKey(),
			),
			expectResult:           ErrAccountSavingCoinNotEnough().Result(),
			expectReferrerSaving:   types.NewCoinFromInt64(98 * types.Decimals),
			expectNewAccountSaving: c0,
		},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.registerMsg)
		if !assert.Equal(t, tc.expectResult, result) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectResult)
		}

		if result.Code.IsOK() {
			if !am.DoesAccountExist(ctx, tc.registerMsg.NewUser) {
				t.Errorf("%s: account %s doesn't exist", tc.testName, tc.registerMsg.NewUser)
			}

			// resetKey, err := am.GetResetKey(ctx, tc.registerMsg.NewUser)
			// if err != nil {
			// 	t.Errorf("%s: failed to get reset key, got err %v", tc.testName, err)
			// }
			// if !resetKey.Equals(tc.registerMsg.NewResetPubKey) {
			// 	t.Errorf("%s: diff reset key, got %v, want %v", tc.testName, resetKey, tc.registerMsg.NewResetPubKey)
			// }

			txKey, err := am.GetTransactionKey(ctx, tc.registerMsg.NewUser)
			if err != nil {
				t.Errorf("%s: failed to get transaction key, got err %v", tc.testName, err)
			}
			if !txKey.Equals(tc.registerMsg.NewResetPubKey) {
				t.Errorf("%s: diff transaction key, got %v, want %v", tc.testName, txKey, tc.registerMsg.NewResetPubKey)
			}

			signingKey, err := am.GetSigningKey(ctx, tc.registerMsg.NewUser)
			if err != nil {
				t.Errorf("%s: failed to get app key, got err %v", tc.testName, err)
			}
			if !signingKey.Equals(tc.registerMsg.NewTransactionPubKey) {
				t.Errorf("%s: diff app key, got %v, want %v", tc.testName, signingKey, tc.registerMsg.NewTransactionPubKey)
			}

			info, err := am.storage.GetInfo(ctx, tc.registerMsg.NewUser)
			if err != nil {
				t.Errorf("%s: failed to get info, got err %v", tc.testName, err)
			}
			bank, err := am.storage.GetBank(ctx, info.Address)
			if err != nil {
				t.Errorf("%s: failed to get bank, got err %v", tc.testName, err)
			}
			if !bank.Saving.IsEqual(tc.expectNewAccountSaving) {
				t.Errorf("%s: diff saving, got %v, want %v", tc.testName, bank.Saving, tc.expectNewAccountSaving)
			}
			// if !bank.CoinDay.IsEqual(tc.expectNewAccountCoinDay) {
			// 	t.Errorf("%s: diff coin day, got %v, want %v", tc.testName, bank.Saving, tc.expectNewAccountSaving)
			// }

			// accMeta, _ := am.storage.GetMeta(ctx, tc.registerMsg.NewUser)
			pool, err := gm.GetValidatorHourlyInflation(ctx)
			if err != nil {
				t.Errorf("%s: failed to get inflation, got err %v", tc.testName, err)
			}
			if !pool.IsEqual(accParam.RegisterFee) {
				t.Errorf("%s: diff validator inflation, got %v, want %v", tc.testName, pool, accParam.RegisterFee)
			}
		}

		saving, err := am.GetSavingFromUsername(ctx, tc.registerMsg.Referrer)
		if err != nil {
			t.Errorf("%s: failed to get saving from bank, got err %v", tc.testName, err)
		}
		if !saving.IsEqual(tc.expectReferrerSaving) {
			t.Errorf("%s: diff saving, got %v, want %v", tc.testName, saving, tc.expectReferrerSaving)
		}
		// gm.GetDeveloperMonthlyInflation(ctx)
	}
}

func TestHandleUpdateAccountMsg(t *testing.T) {
	ctx, am, gm := setupTest(t, 1)
	handler := NewHandler(am, &gm)

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
			expectResult:     model.ErrAccountMetaNotFound().Result(),
		},
	}
	for _, tc := range testCases {
		result := handler(ctx, tc.updateAccountMsg)
		if !assert.Equal(t, result, tc.expectResult) {
			t.Errorf("%s: diff result, got %v, want %v", tc.testName, result, tc.expectResult)
		}
	}
}
