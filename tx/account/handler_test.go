package account

import (
	"fmt"
	"testing"

	"github.com/lino-network/lino/tx/account/model"
	"github.com/lino-network/lino/types"

	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

var (
	user1 = types.AccountKey("user1")
	user2 = types.AccountKey("user2")

	memo = "This is a memo!"
)

func TestFollow(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)

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
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)

	// create test user
	createTestAccount(ctx, am, "user1")

	// let user2(not exists) follows user1
	msg := NewFollowMsg("user2", "user1")
	result := handler(ctx, msg)

	assert.Equal(t, result, ErrUsernameNotFound().Result())
	assert.False(t, am.IsMyFollower(ctx, types.AccountKey("user1"), types.AccountKey("user2")))

	// let user1 follows user3(not exists)
	msg = NewFollowMsg("user1", "user3")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrUsernameNotFound().Result())
	assert.False(t, am.IsMyFollowing(ctx, types.AccountKey("user1"), types.AccountKey("user3")))
}

func TestFollowAgain(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)

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
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)

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
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)
	// create test user
	createTestAccount(ctx, am, "user1")

	// let user2(not exists) unfollows user1
	msg := NewUnfollowMsg("user2", "user1")
	result := handler(ctx, msg)
	assert.Equal(t, result, ErrUsernameNotFound().Result())

	// let user1 unfollows user3(not exists)
	msg = NewUnfollowMsg("user1", "user3")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrUsernameNotFound().Result())
}

func TestInvalidUnfollow(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)
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
	ctx, am, accParam := setupTest(t, 1)
	handler := NewHandler(am)

	// create two test users with initial deposit of 100 LNO.
	createTestAccount(ctx, am, "user1")
	createTestAccount(ctx, am, "user2")

	am.AddSavingCoin(ctx, types.AccountKey("user1"), c2000)

	receiverAddr, _ := am.GetBankAddress(ctx, user2)

	testCases := []struct {
		testName            string
		msg                 TransferMsg
		wantOK              bool
		wantSenderBalance   types.Coin
		wantReceiverBalance types.Coin
	}{
		{testName: "user1 transfers 200 LNO to user2 (by username)",
			msg: TransferMsg{
				Sender:       user1,
				ReceiverName: user2,
				Amount:       l200,
				Memo:         memo,
			},
			wantOK:              true,
			wantSenderBalance:   c1800.Plus(accParam.RegisterFee),
			wantReceiverBalance: c200.Plus(accParam.RegisterFee),
		},
		{testName: "user1 transfers 1600 LNO to user2 (by both username and address)",
			msg: TransferMsg{
				Sender:       user1,
				ReceiverName: user2,
				ReceiverAddr: receiverAddr,
				Amount:       l1600,
				Memo:         memo,
			},
			wantOK:              true,
			wantSenderBalance:   c200.Plus(accParam.RegisterFee),
			wantReceiverBalance: c1800.Plus(accParam.RegisterFee),
		},
		{testName: "user1 transfers 100 LNO to user2 (by address)",
			msg: TransferMsg{
				Sender:       user1,
				ReceiverAddr: receiverAddr,
				Amount:       l100,
				Memo:         memo,
			},
			wantOK:              true,
			wantSenderBalance:   c100.Plus(accParam.RegisterFee),
			wantReceiverBalance: c1900.Plus(accParam.RegisterFee),
		},
		{testName: "user2 transfers 100 LNO to a random address",
			msg: TransferMsg{
				Sender:       user2,
				ReceiverAddr: sdk.Address("sdajsdbiqwbdiub"),
				Amount:       l100,
				Memo:         memo,
			},
			wantOK:              true,
			wantSenderBalance:   c1800.Plus(accParam.RegisterFee),
			wantReceiverBalance: c100,
		},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.msg)

		if result.IsOK() != tc.wantOK {
			t.Errorf("%s handler(%v): got %v, want %v, err:%v", tc.testName, tc.msg, result.IsOK(), tc.wantOK, result)
		}

		senderSaving, _ := am.GetSavingFromBank(ctx, tc.msg.Sender)
		var receiverSaving types.Coin
		if tc.msg.ReceiverName != "" {
			receiverSaving, _ = am.GetSavingFromBank(ctx, tc.msg.ReceiverName)
		} else {
			bank, _ := am.storage.GetBankFromAddress(ctx, tc.msg.ReceiverAddr)
			receiverSaving = bank.Saving
		}

		if !senderSaving.IsEqual(tc.wantSenderBalance) {
			t.Errorf("%s get sender bank Saving(%v): got %v, want %v", tc.testName, tc.msg.Sender, senderSaving, tc.wantSenderBalance)
		}
		if !receiverSaving.IsEqual(tc.wantReceiverBalance) {
			t.Errorf("%s: get receiver bank Saving(%v): got %v, want %v", tc.testName, tc.msg.ReceiverName, receiverSaving, tc.wantReceiverBalance)
		}
	}
}

func TestSenderCoinNotEnough(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	handler := NewHandler(am)

	// create two test users
	createTestAccount(ctx, am, "user1")
	createTestAccount(ctx, am, "user2")

	memo := "This is a memo!"

	// let user1 transfers 2000 to user2
	msg := NewTransferMsg("user1", l2000, memo, TransferToUser("user2"))
	result := handler(ctx, msg)
	assert.Equal(t, ErrAccountSavingCoinNotEnough().Result(), result)

	acc1Balance, _ := am.GetSavingFromBank(ctx, types.AccountKey("user1"))
	assert.Equal(t, acc1Balance, accParam.RegisterFee)
}

func TestUsernameAddressMismatch(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)

	// create two test users
	createTestAccount(ctx, am, "user1")
	createTestAccount(ctx, am, "user2")

	am.AddSavingCoin(ctx, types.AccountKey("user1"), c2000)
	am.AddSavingCoin(ctx, types.AccountKey("user2"), c2000)

	memo := "This is a memo!"
	randomAddr := sdk.Address("dqwdnqwdbnqwkjd")

	// let user1 transfers 2000 Lino to user2 (provide both name and address)
	msg := NewTransferMsg(
		"user1", l1999, memo, TransferToUser("user2"), TransferToAddr(randomAddr))
	result := handler(ctx, msg)
	assert.Equal(t, ErrTransferHandler(msg.Sender).Result(), result)
}

func TestReceiverUsernameIncorrect(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)

	// create two test users
	createTestAccount(ctx, am, "user1")
	am.AddSavingCoin(ctx, types.AccountKey("user1"), c2000)

	memo := "This is a memo!"

	// let user1 transfers 2000 to a random user
	msg := NewTransferMsg("user1", l2000, memo, TransferToUser("dnqwondqowindow"))
	result := handler(ctx, msg)
	assert.Equal(t, ErrTransferHandler(msg.Sender).Result().Code, result.Code)
}

func TestHandleAccountRecover(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)
	user1 := types.AccountKey("user1")

	priv := createTestAccount(ctx, am, string(user1))

	testCases := map[string]struct {
		user              types.AccountKey
		newPostKey        crypto.PubKey
		newTransactionKey crypto.PubKey
	}{
		"normal case": {
			user1, crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
		},
	}

	for testName, tc := range testCases {
		msg := RecoverMsg{user1, tc.newPostKey, tc.newTransactionKey}
		result := handler(ctx, msg)
		assert.Equal(
			t, sdk.Result{}, result, fmt.Sprintf("%s: got %v, want %v", testName, result, sdk.Result{}))
		accInfo := model.AccountInfo{
			Username:       tc.user,
			CreatedAt:      ctx.BlockHeader().Time,
			MasterKey:      priv.PubKey(),
			TransactionKey: tc.newTransactionKey,
			PostKey:        tc.newPostKey,
			Address:        priv.PubKey().Address(),
		}
		checkAccountInfo(t, ctx, tc.user, accInfo)
	}
}

func TestSavingAndChecking(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	handler := NewHandler(am)
	user1 := types.AccountKey("user1")

	createTestAccount(ctx, am, string(user1))
	am.AddSavingCoin(ctx, user1, c2000)

	testCases := []struct {
		testName             string
		user                 string
		fromSavingToChecking bool
		amount               types.LNO
		expectResult         sdk.Result
		expectChecking       types.Coin
		expectSaving         types.Coin
	}{
		{"transfer from saving to checking",
			string(user1), true, types.LNO("200"), sdk.Result{},
			c200, accParam.RegisterFee.Plus(c1800),
		},
		{"transfer from checking to saving",
			string(user1), false, types.LNO("200"), sdk.Result{},
			c0, accParam.RegisterFee.Plus(c2000),
		},
		{"transfer from checking to saving if checking is insufficient",
			string(user1), false, types.LNO("200"),
			ErrAccountCheckingCoinNotEnough().Result(),
			c0, accParam.RegisterFee.Plus(c2000),
		},
		{"transfer from saving to checking if saving is insufficient",
			string(user1), false, types.LNO("2001"),
			ErrAccountSavingCoinNotEnough().Result(),
			c0, accParam.RegisterFee.Plus(c2000),
		},
	}

	for _, tc := range testCases {
		var msg sdk.Msg
		if tc.fromSavingToChecking {
			msg = NewSavingToCheckingMsg(tc.user, tc.amount)
		} else {
			msg = NewCheckingToSavingMsg(tc.user, tc.amount)
		}
		result := handler(ctx, msg)
		assert.Equal(t, tc.expectResult, result,
			fmt.Sprintf("%s: got %v, want %v", tc.testName, tc.expectResult, result))
		accSaving, err := am.GetSavingFromBank(ctx, user1)
		assert.Nil(t, err, fmt.Sprintf("%s: got err %v", tc.testName, err))
		assert.Equal(t, tc.expectSaving, accSaving,
			fmt.Sprintf("%s: expect saving %v, got %v", tc.testName, tc.expectSaving, accSaving))
		accChecking, err := am.GetCheckingFromBank(ctx, user1)
		assert.Nil(t, err, fmt.Sprintf("%s: got err %v", tc.testName, err))
		assert.Equal(t, tc.expectChecking, accChecking,
			fmt.Sprintf("%s: expect saving %v, got %v", tc.testName, tc.expectChecking, accChecking))
	}
}
