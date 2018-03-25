package account

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var (
	c0    = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 0}}
	c100  = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 100}}
	c200  = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 200}}
	c1600 = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 1600}}
	c1800 = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 1800}}
	c1900 = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 1900}}
	c2000 = sdk.Coins{sdk.Coin{Denom: "lino", Amount: 2000}}
)

func TestFollow(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	createTestAccount(ctx, lam, "user1")
	createTestAccount(ctx, lam, "user2")

	// let user1 follows user2
	msg := NewFollowMsg("user1", "user2")
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// check user1 in the user2's follower list
	followerList, _ := lam.GetFollower(ctx, AccountKey("user2"))
	idx := findAccountInList(AccountKey("user1"), followerList.Follower)
	assert.Equal(t, true, idx >= 0)

	// check user2 in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, AccountKey("user1"))
	idx = findAccountInList(AccountKey("user2"), followingList.Following)
	assert.Equal(t, true, idx >= 0)
}

func TestFollowUserNotExist(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create test user
	createTestAccount(ctx, lam, "user1")

	// let user2(not exists) follows user1
	msg := NewFollowMsg("user2", "user1")
	result := handler(ctx, msg)
	assert.Equal(t, result, ErrAccountManagerFail("Get following list failed").Result())

	followerList, _ := lam.GetFollower(ctx, AccountKey("user1"))
	assert.Equal(t, 0, len(followerList.Follower))

	// let user1 follows user3(not exists)
	msg = NewFollowMsg("user1", "user3")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrAccountManagerFail("Get follower list failed").Result())

	followingList, _ := lam.GetFollowing(ctx, AccountKey("user1"))
	assert.Equal(t, 0, len(followingList.Following))
}

func TestFollowAgain(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	createTestAccount(ctx, lam, "user1")
	createTestAccount(ctx, lam, "user2")

	// let user1 follows user2 twice
	msg := NewFollowMsg("user1", "user2")
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	msg = NewFollowMsg("user1", "user2")
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// check user1 is user2's only follower
	followerList, _ := lam.GetFollower(ctx, AccountKey("user2"))
	idx := findAccountInList(AccountKey("user1"), followerList.Follower)
	assert.Equal(t, 0, idx)
	assert.Equal(t, 1, len(followerList.Follower))

	// check user2 is the only one in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, AccountKey("user1"))
	idx = findAccountInList(AccountKey("user2"), followingList.Following)
	assert.Equal(t, 0, idx)
	assert.Equal(t, 1, len(followingList.Following))
}

func TestUnfollow(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	createTestAccount(ctx, lam, "user1")
	createTestAccount(ctx, lam, "user2")

	// let user1 follows user2
	msg := NewFollowMsg("user1", "user2")
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	// let user1 unfollows user2
	msg2 := NewUnfollowMsg("user1", "user2")
	result = handler(ctx, msg2)
	assert.Equal(t, result, sdk.Result{})

	// check user1 is not in the user2's follower list
	followerList, _ := lam.GetFollower(ctx, AccountKey("user2"))
	idx := findAccountInList(AccountKey("user1"), followerList.Follower)
	assert.Equal(t, -1, idx)

	// check user2 is not in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, AccountKey("user1"))
	idx = findAccountInList(AccountKey("user2"), followingList.Following)
	assert.Equal(t, -1, idx)
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
	assert.Equal(t, result, ErrAccountManagerFail("Get following list failed").Result())

	// let user1 unfollows user3(not exists)
	msg = NewUnfollowMsg("user1", "user3")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrAccountManagerFail("Get follower list failed").Result())
}

func TestInvalidUnfollow(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)
	// create test user
	createTestAccount(ctx, lam, "user1")
	createTestAccount(ctx, lam, "user2")
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
	followerList, _ := lam.GetFollower(ctx, AccountKey("user2"))
	idx := findAccountInList(AccountKey("user1"), followerList.Follower)
	assert.Equal(t, true, idx >= 0)
	assert.Equal(t, 1, len(followerList.Follower))

	// check user2 in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, AccountKey("user1"))
	idx = findAccountInList(AccountKey("user2"), followingList.Following)
	assert.Equal(t, true, idx >= 0)
	assert.Equal(t, 1, len(followingList.Following))

}

func TestTransferNormal(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc2 := createTestAccount(ctx, lam, "user2")

	acc1.AddCoins(ctx, c2000)

	acc1.Apply(ctx)
	acc2.Apply(ctx)

	memo := []byte("This is a memo!")

	// let user1 transfers 200 to user2 (by username)
	msg := NewTransferMsg("user1", c200, memo, TransferToUser("user2"))
	result := handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	acc1Balance, _ := acc1.GetBankBalance(ctx)
	acc2Balance, _ := acc2.GetBankBalance(ctx)

	assert.Equal(t, true, acc1Balance.IsEqual(c1800))
	assert.Equal(t, true, acc2Balance.IsEqual(c200))

	//let user1 transfers 1600 to user2 (by both username and address)
	acc1.clear()
	acc2.clear()

	acc2Addr, _ := acc2.GetBankAddress(ctx)
	msg = NewTransferMsg("user1", c1600, memo, TransferToUser("user2"), TransferToAddr(acc2Addr))
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	acc1Balance, _ = acc1.GetBankBalance(ctx)
	acc2Balance, _ = acc2.GetBankBalance(ctx)

	assert.Equal(t, true, acc1Balance.IsEqual(c200))
	assert.Equal(t, true, acc2Balance.IsEqual(c1800))

	//let user1 transfers 100 to user2 (by  address)
	acc1.clear()
	acc2.clear()

	msg = NewTransferMsg("user1", c100, memo, TransferToAddr(acc2Addr))
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	acc1Balance, _ = acc1.GetBankBalance(ctx)
	acc2Balance, _ = acc2.GetBankBalance(ctx)

	assert.Equal(t, true, acc1Balance.IsEqual(c100))
	assert.Equal(t, true, acc2Balance.IsEqual(c1900))

	//let user1 transfers 100 to a random address
	acc1.clear()
	acc2.clear()

	randomAddr := sdk.Address("sdajsdbiqwbdiub")
	msg = NewTransferMsg("user1", c100, memo, TransferToAddr(randomAddr))
	result = handler(ctx, msg)
	assert.Equal(t, result, sdk.Result{})

	acc1Balance, _ = acc1.GetBankBalance(ctx)
	generatedBank, _ := lam.GetBankFromAddress(ctx, randomAddr)

	assert.Equal(t, true, acc1Balance.IsEqual(c0))
	assert.Equal(t, true, generatedBank.Balance.IsEqual(c100))

}

func TestSenderCoinNotEnough(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc2 := createTestAccount(ctx, lam, "user2")

	acc1.AddCoins(ctx, c200)

	acc1.Apply(ctx)
	acc2.Apply(ctx)

	memo := []byte("This is a memo!")

	// let user1 transfers 2000 to user2
	msg := NewTransferMsg("user1", c2000, memo, TransferToUser("user2"))
	result := handler(ctx, msg)
	assert.Equal(t, ErrAccountManagerFail("Account bank's coins are not enough").Result(), result)

	acc1Balance, _ := acc1.GetBankBalance(ctx)
	assert.Equal(t, true, acc1Balance.IsEqual(c200))
}

func TestUsernameAddressMismatch(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc2 := createTestAccount(ctx, lam, "user2")

	acc1.AddCoins(ctx, c2000)
	acc2.AddCoins(ctx, c2000)

	acc1.Apply(ctx)
	acc2.Apply(ctx)

	memo := []byte("This is a memo!")
	randomAddr := sdk.Address("dqwdnqwdbnqwkjd")

	// let user1 transfers 2000 to user2 (provide both name and address)
	msg := NewTransferMsg("user1", c2000, memo, TransferToUser("user2"), TransferToAddr(randomAddr))
	result := handler(ctx, msg)
	assert.Equal(t, ErrAccountManagerFail("Username and address mismatch").Result(), result)

	acc1Balance, _ := acc1.GetBankBalance(ctx)
	acc2Balance, _ := acc2.GetBankBalance(ctx)

	assert.Equal(t, true, acc1Balance.IsEqual(c2000))
	assert.Equal(t, true, acc2Balance.IsEqual(c2000))
}

func TestReceiverUsernameIncorrect(t *testing.T) {
	lam := newLinoAccountManager()
	ctx := getContext()
	handler := NewHandler(lam)

	// create two test users
	acc1 := createTestAccount(ctx, lam, "user1")
	acc1.AddCoins(ctx, c2000)
	acc1.Apply(ctx)

	memo := []byte("This is a memo!")

	// let user1 transfers 2000 to a random user
	msg := NewTransferMsg("user1", c2000, memo, TransferToUser("dnqwondqowindow"))
	result := handler(ctx, msg)
	assert.Equal(t, ErrAccountManagerFail("Add money to receiver's bank failed").Result(), result)

	acc1Balance, _ := acc1.GetBankBalance(ctx)
	assert.Equal(t, true, acc1Balance.IsEqual(c2000))
}
