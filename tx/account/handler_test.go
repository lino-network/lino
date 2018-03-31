package account

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
)

var (
	l0    = types.LNO(sdk.NewRat(0))
	l100  = types.LNO(sdk.NewRat(100))
	l200  = types.LNO(sdk.NewRat(200))
	l1600 = types.LNO(sdk.NewRat(1600))
	l1800 = types.LNO(sdk.NewRat(1800))
	l1900 = types.LNO(sdk.NewRat(1900))
	l2000 = types.LNO(sdk.NewRat(2000))
)

func TestFollow(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc2 := createTestAccount(ctx, lam, "user2")

	// let user1 follows user2
	msg := NewFollowMsg("user1", "user2")
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// check user1 in the user2's follower list
	assert.Equal(t, true, acc2.IsMyFollower(ctx, "user1"))

	// check user2 in the user1's following list
	assert.Equal(t, true, acc1.IsMyFollowing(ctx, "user2"))
}

func TestFollowUserNotExist(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create test user
	acc1 := createTestAccount(ctx, lam, "user1")

	// let user2(not exists) follows user1
	msg := NewFollowMsg("user2", "user1")
	result := handler(ctx, msg)

	assert.Equal(t, result, ErrUsernameNotFound().Result())
	assert.Equal(t, false, acc1.IsMyFollower(ctx, "user2"))

	// let user1 follows user3(not exists)
	msg = NewFollowMsg("user1", "user3")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrUsernameNotFound().Result())
	assert.Equal(t, false, acc1.IsMyFollowing(ctx, "user3"))
}

func TestFollowAgain(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc2 := createTestAccount(ctx, lam, "user2")

	// let user1 follows user2 twice
	msg := NewFollowMsg("user1", "user2")
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	msg = NewFollowMsg("user1", "user2")
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// check user1 is user2's only follower
	assert.Equal(t, true, acc2.IsMyFollower(ctx, "user1"))

	// check user2 is the only one in the user1's following list
	assert.Equal(t, true, acc1.IsMyFollowing(ctx, "user2"))
}

func TestUnfollow(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc2 := createTestAccount(ctx, lam, "user2")

	// let user1 follows user2
	msg := NewFollowMsg("user1", "user2")
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// let user1 unfollows user2
	msg2 := NewUnfollowMsg("user1", "user2")
	result = handler(ctx, msg2)
	assert.Equal(t, result, sdk.Result{})

	// check user1 is not in the user2's follower list
	assert.Equal(t, false, acc2.IsMyFollower(ctx, "user1"))

	// check user2 is not in the user1's following list
	assert.Equal(t, false, acc1.IsMyFollowing(ctx, "user2"))
}

func TestUnfollowUserNotExist(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)
	// create test user
	createTestAccount(ctx, lam, "user1")

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
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)
	// create test user
	acc1 := createTestAccount(ctx, lam, "user1")
	acc2 := createTestAccount(ctx, lam, "user2")
	createTestAccount(ctx, lam, "user3")

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
	assert.Equal(t, true, acc2.IsMyFollower(ctx, "user1"))

	// check user2 in the user1's following list
	assert.Equal(t, true, acc1.IsMyFollowing(ctx, "user2"))

}

func TestTransferNormal(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc2 := createTestAccount(ctx, lam, "user2")

	acc1.AddCoin(ctx, types.LinoToCoin(l2000))

	acc1.Apply(ctx)
	acc2.Apply(ctx)

	memo := []byte("This is a memo!")

	// let user1 transfers 200 to user2 (by username)
	msg := NewTransferMsg("user1", l200, memo, TransferToUser("user2"))
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	acc1Balance, _ := acc1.GetBankBalance(ctx)
	acc2Balance, _ := acc2.GetBankBalance(ctx)

	assert.Equal(t, true, acc1Balance.IsEqual(types.LinoToCoin(l1800)))
	assert.Equal(t, true, acc2Balance.IsEqual(types.LinoToCoin(l200)))

	//let user1 transfers 1600 to user2 (by both username and address)
	acc1.clear()
	acc2.clear()

	acc2Addr, _ := acc2.GetBankAddress(ctx)
	msg = NewTransferMsg("user1", l1600, memo, TransferToUser("user2"), TransferToAddr(acc2Addr))
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	acc1Balance, _ = acc1.GetBankBalance(ctx)
	acc2Balance, _ = acc2.GetBankBalance(ctx)

	assert.Equal(t, true, acc1Balance.IsEqual(types.LinoToCoin(l200)))
	assert.Equal(t, true, acc2Balance.IsEqual(types.LinoToCoin(l1800)))

	//let user1 transfers 100 to user2 (by  address)
	acc1.clear()
	acc2.clear()

	msg = NewTransferMsg("user1", l100, memo, TransferToAddr(acc2Addr))
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	acc1Balance, _ = acc1.GetBankBalance(ctx)
	acc2Balance, _ = acc2.GetBankBalance(ctx)

	assert.Equal(t, true, acc1Balance.IsEqual(types.LinoToCoin(l100)))
	assert.Equal(t, true, acc2Balance.IsEqual(types.LinoToCoin(l1900)))

	//let user1 transfers 100 to a random address
	acc1.clear()
	acc2.clear()

	randomAddr := sdk.Address("sdajsdbiqwbdiub")
	msg = NewTransferMsg("user1", l100, memo, TransferToAddr(randomAddr))
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	acc1Balance, _ = acc1.GetBankBalance(ctx)
	generatedBank, _ := lam.GetBankFromAddress(ctx, randomAddr)

	assert.Equal(t, true, acc1Balance.IsEqual(types.LinoToCoin(l0)))
	assert.Equal(t, true, generatedBank.Balance.IsEqual(types.LinoToCoin(l100)))

}

func TestSenderCoinNotEnough(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc2 := createTestAccount(ctx, lam, "user2")

	acc1.AddCoin(ctx, types.LinoToCoin(l1600))

	acc1.Apply(ctx)
	acc2.Apply(ctx)

	memo := []byte("This is a memo!")

	// let user1 transfers 2000 to user2
	msg := NewTransferMsg("user1", l2000, memo, TransferToUser("user2"))
	result := handler(ctx, msg)
	assert.Equal(t, ErrAccountCoinNotEnough().Result(), result)

	acc1Balance, _ := acc1.GetBankBalance(ctx)
	assert.Equal(t, true, acc1Balance.IsEqual(types.LinoToCoin(l1600)))
}

func TestUsernameAddressMismatch(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc2 := createTestAccount(ctx, lam, "user2")

	acc1.AddCoin(ctx, types.LinoToCoin(l2000))
	acc2.AddCoin(ctx, types.LinoToCoin(l2000))

	acc1.Apply(ctx)
	acc2.Apply(ctx)

	memo := []byte("This is a memo!")
	randomAddr := sdk.Address("dqwdnqwdbnqwkjd")

	// let user1 transfers 2000 to user2 (provide both name and address)
	msg := NewTransferMsg("user1", l2000, memo, TransferToUser("user2"), TransferToAddr(randomAddr))
	result := handler(ctx, msg)
	assert.Equal(t, ErrUsernameAddressMismatch().Result(), result)

	acc1Balance, _ := acc1.GetBankBalance(ctx)
	acc2Balance, _ := acc2.GetBankBalance(ctx)

	assert.Equal(t, true, acc1Balance.IsEqual(types.LinoToCoin(l2000)))
	assert.Equal(t, true, acc2Balance.IsEqual(types.LinoToCoin(l2000)))
}

func TestReceiverUsernameIncorrect(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc1.AddCoin(ctx, types.LinoToCoin(l2000))
	acc1.Apply(ctx)

	memo := []byte("This is a memo!")

	// let user1 transfers 2000 to a random user
	msg := NewTransferMsg("user1", l2000, memo, TransferToUser("dnqwondqowindow"))
	result := handler(ctx, msg)
	assert.Equal(t, ErrAddMoneyFailed().Result(), result)

	acc1Balance, _ := acc1.GetBankBalance(ctx)
	assert.Equal(t, true, acc1Balance.IsEqual(types.LinoToCoin(l2000)))
}
