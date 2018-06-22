package account

import (
	"fmt"
	"testing"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"

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

	am.AddSavingCoin(
		ctx, types.AccountKey("user1"), c2000, "", "", types.TransferIn)

	testCases := []struct {
		testName            string
		msg                 TransferMsg
		wantOK              bool
		wantSenderBalance   types.Coin
		wantReceiverBalance types.Coin
	}{
		{testName: "user1 transfers 200 LNO to user2 (by username)",
			msg:                 NewTransferMsg("user1", "user2", l200, memo),
			wantOK:              true,
			wantSenderBalance:   c1800.Plus(accParam.RegisterFee),
			wantReceiverBalance: c200.Plus(accParam.RegisterFee),
		},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.msg)

		if result.IsOK() != tc.wantOK {
			t.Errorf(
				"%s handler(%v): got %v, want %v, err:%v",
				tc.testName, tc.msg, result.IsOK(), tc.wantOK, result)
		}

		senderSaving, _ := am.GetSavingFromBank(ctx, tc.msg.Sender)
		receiverSaving, _ := am.GetSavingFromBank(ctx, tc.msg.Receiver)

		if !senderSaving.IsEqual(tc.wantSenderBalance) {
			t.Errorf(
				"%s get sender bank Saving(%v): got %v, want %v",
				tc.testName, tc.msg.Sender, senderSaving, tc.wantSenderBalance)
		}
		if !receiverSaving.IsEqual(tc.wantReceiverBalance) {
			t.Errorf(
				"%s: get receiver bank Saving(%v): got %v, want %v",
				tc.testName, tc.msg.Receiver, receiverSaving, tc.wantReceiverBalance)
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
	msg := NewTransferMsg("user1", "user2", l2000, memo)
	result := handler(ctx, msg)
	assert.Equal(t, ErrAccountSavingCoinNotEnough().Result(), result)

	acc1Balance, _ := am.GetSavingFromBank(ctx, types.AccountKey("user1"))
	assert.Equal(t, acc1Balance, accParam.RegisterFee)
}

func TestReceiverUsernameIncorrect(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)

	// create two test users
	createTestAccount(ctx, am, "user1")

	memo := "This is a memo!"

	// let user1 transfers 2000 to a random user
	msg := NewTransferMsg("user1", "dnqwondqowindow", l2000, memo)
	result := handler(ctx, msg)
	assert.Equal(t, ErrUsernameNotFound().Result().Code, result.Code)
}

func TestHandleAccountRecover(t *testing.T) {
	ctx, am, accParam := setupTest(t, 1)
	handler := NewHandler(am)
	user1 := "user1"

	createTestAccount(ctx, am, user1)

	testCases := map[string]struct {
		user              string
		newMasterKey      crypto.PubKey
		newPostKey        crypto.PubKey
		newTransactionKey crypto.PubKey
	}{
		"normal case": {
			user1, crypto.GenPrivKeyEd25519().PubKey(),
			crypto.GenPrivKeyEd25519().PubKey(), crypto.GenPrivKeyEd25519().PubKey(),
		},
	}

	for testName, tc := range testCases {
		msg := NewRecoverMsg(tc.user, tc.newMasterKey, tc.newTransactionKey, tc.newPostKey)
		result := handler(ctx, msg)
		assert.Equal(
			t, sdk.Result{}, result, fmt.Sprintf("%s: got %v, want %v", testName, result, sdk.Result{}))
		accInfo := model.AccountInfo{
			Username:       types.AccountKey(tc.user),
			CreatedAt:      ctx.BlockHeader().Time,
			MasterKey:      tc.newMasterKey,
			TransactionKey: tc.newTransactionKey,
			PostKey:        tc.newPostKey,
		}
		checkAccountInfo(t, ctx, types.AccountKey(tc.user), accInfo)
		newBank := model.AccountBank{
			Saving:  accParam.RegisterFee,
			NumOfTx: 1,
		}
		checkBankKVByUsername(t, ctx, types.AccountKey(tc.user), newBank)
	}
}

func TestHandleRegister(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)
	referrer := "referrer"

	createTestAccount(ctx, am, referrer)
	am.AddSavingCoin(
		ctx, types.AccountKey(referrer), types.NewCoinFromInt64(100*types.Decimals),
		"", "", types.TransferIn)

	testCases := []struct {
		testName             string
		registerMsg          RegisterMsg
		expectResult         sdk.Result
		expectReferrerSaving types.Coin
	}{
		{"normal case",
			NewRegisterMsg(
				"referrer", "user1", "1",
				crypto.GenPrivKeySecp256k1().PubKey(),
				crypto.GenPrivKeySecp256k1().PubKey(),
				crypto.GenPrivKeySecp256k1().PubKey(),
			),
			sdk.Result{}, c100,
		},
		{"account already exist",
			NewRegisterMsg(
				"referrer", "user1", "1",
				crypto.GenPrivKeySecp256k1().PubKey(),
				crypto.GenPrivKeySecp256k1().PubKey(),
				crypto.GenPrivKeySecp256k1().PubKey(),
			),
			ErrAccountAlreadyExists("user1").Result(),
			types.NewCoinFromInt64(99 * types.Decimals),
		},
		{"account register fee insufficient",
			NewRegisterMsg(
				"referrer", "user2", "0.1",
				crypto.GenPrivKeySecp256k1().PubKey(),
				crypto.GenPrivKeySecp256k1().PubKey(),
				crypto.GenPrivKeySecp256k1().PubKey(),
			),
			ErrRegisterFeeInsufficient().Result(),
			types.NewCoinFromInt64(9890000),
		},
		{"referrer deposit insufficient",
			NewRegisterMsg(
				"referrer", "user2", "1000",
				crypto.GenPrivKeySecp256k1().PubKey(),
				crypto.GenPrivKeySecp256k1().PubKey(),
				crypto.GenPrivKeySecp256k1().PubKey(),
			),
			ErrAccountSavingCoinNotEnough().Result(),
			types.NewCoinFromInt64(9890000),
		},
	}

	for _, tc := range testCases {
		result := handler(ctx, tc.registerMsg)
		assert.Equal(t, tc.expectResult, result)
		if result.Code == sdk.ABCICodeOK {
			assert.True(t, am.DoesAccountExist(ctx, tc.registerMsg.NewUser))
			txKey, err := am.GetTransactionKey(ctx, tc.registerMsg.NewUser)
			assert.Nil(t, err)
			assert.Equal(t, txKey, tc.registerMsg.NewTransactionPubKey)
			postKey, err := am.GetPostKey(ctx, tc.registerMsg.NewUser)
			assert.Nil(t, err)
			assert.Equal(t, postKey, tc.registerMsg.NewPostPubKey)
			masterKey, err := am.GetMasterKey(ctx, tc.registerMsg.NewUser)
			assert.Nil(t, err)
			assert.Equal(t, masterKey, tc.registerMsg.NewMasterPubKey)
		}
		saving, err := am.GetSavingFromBank(ctx, tc.registerMsg.Referrer)
		assert.Nil(t, err)
		assert.Equal(t, tc.expectReferrerSaving, saving)
	}
}

func TesthandleUpdateAccountMsg(t *testing.T) {
	ctx, am, _ := setupTest(t, 1)
	handler := NewHandler(am)

	createTestAccount(ctx, am, "accKey")

	cases := []struct {
		testName         string
		updateAccountMsg UpdateAccountMsg
		expectResult     sdk.Result
	}{
		{"normal update",
			NewUpdateAccountMsg("accKey", "{'link':'https://lino.network'}"),
			sdk.Result{},
		},
		{"invalid username",
			NewUpdateAccountMsg("invalid", "{'link':'https://lino.network'}"),
			ErrUsernameNotFound().Result(),
		},
	}
	for _, cs := range cases {
		result := handler(ctx, cs.updateAccountMsg)
		assert.Equal(t, result, cs.expectResult)
	}
}
