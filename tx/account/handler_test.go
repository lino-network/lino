package account

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
	"github.com/stretchr/testify/assert"
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
	followerList, _ := lam.GetFollower(ctx, types.AccountKey("user2"))
	idx := findAccountInList(types.AccountKey("user1"), followerList.Follower)
	assert.Equal(t, true, idx >= 0)

	// check user2 in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, types.AccountKey("user1"))
	idx = findAccountInList(types.AccountKey("user2"), followingList.Following)
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

	followerList, _ := lam.GetFollower(ctx, types.AccountKey("user1"))
	assert.Equal(t, 0, len(followerList.Follower))

	// let user1 follows user3(not exists)
	msg = NewFollowMsg("user1", "user3")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrAccountManagerFail("Get follower list failed").Result())

	followingList, _ := lam.GetFollowing(ctx, types.AccountKey("user1"))
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
	followerList, _ := lam.GetFollower(ctx, types.AccountKey("user2"))
	idx := findAccountInList(types.AccountKey("user1"), followerList.Follower)
	assert.Equal(t, 0, idx)
	assert.Equal(t, 1, len(followerList.Follower))

	// check user2 is the only one in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, types.AccountKey("user1"))
	idx = findAccountInList(types.AccountKey("user2"), followingList.Following)
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
	followerList, _ := lam.GetFollower(ctx, types.AccountKey("user2"))
	idx := findAccountInList(types.AccountKey("user1"), followerList.Follower)
	assert.Equal(t, -1, idx)

	// check user2 is not in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, types.AccountKey("user1"))
	idx = findAccountInList(types.AccountKey("user2"), followingList.Following)
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
	followerList, _ := lam.GetFollower(ctx, types.AccountKey("user2"))
	idx := findAccountInList(types.AccountKey("user1"), followerList.Follower)
	assert.Equal(t, true, idx >= 0)
	assert.Equal(t, 1, len(followerList.Follower))

	// check user2 in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, types.AccountKey("user1"))
	idx = findAccountInList(types.AccountKey("user2"), followingList.Following)
	assert.Equal(t, true, idx >= 0)
	assert.Equal(t, 1, len(followingList.Following))

}
