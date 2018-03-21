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
	flag := isInFollowerList(types.AccountKey("user1"), followerList)
	assert.Equal(t, true, flag)

	// check user2 in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, types.AccountKey("user1"))
	flag = isInFollowingList(types.AccountKey("user2"), followingList)
	assert.Equal(t, true, flag)
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
	assert.Equal(t, result, ErrUsernameNotFound("Username not found").Result())

	// let user1 follows user3(not exists)
	msg = NewFollowMsg("user1", "user3")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrUsernameNotFound("Username not found").Result())
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
	flag := isInFollowerList(types.AccountKey("user1"), followerList)
	assert.Equal(t, true, flag)
	assert.Equal(t, 1, len(followerList.Follower))

	// check user2 is the only one in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, types.AccountKey("user1"))
	flag = isInFollowingList(types.AccountKey("user2"), followingList)
	assert.Equal(t, true, flag)
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
	flag := isInFollowerList(types.AccountKey("user1"), followerList)
	assert.Equal(t, false, flag)

	// check user2 is not in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, types.AccountKey("user1"))
	flag = isInFollowingList(types.AccountKey("user2"), followingList)
	assert.Equal(t, false, flag)
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
	assert.Equal(t, result, ErrUsernameNotFound("Username not found").Result())

	// let user1 unfollows user3(not exists)
	msg = NewUnfollowMsg("user1", "user3")
	result = handler(ctx, msg)
	assert.Equal(t, result, ErrUsernameNotFound("Username not found").Result())
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
	flag := isInFollowerList(types.AccountKey("user1"), followerList)
	assert.Equal(t, true, flag)
	assert.Equal(t, 1, len(followerList.Follower))

	// check user2 in the user1's following list
	followingList, _ := lam.GetFollowing(ctx, types.AccountKey("user1"))
	flag = isInFollowingList(types.AccountKey("user2"), followingList)
	assert.Equal(t, true, flag)
	assert.Equal(t, 1, len(followingList.Following))

}
