package account

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func NewHandler(am types.AccountManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case RegisterMsg:
			return handleFollowMsg(ctx, am, msg)
		case UnfollowMsg:
			return handleUnfollowMsg(ctx, am, msg)
		case TransferMsg:
			return handleTransferMsg(ctx, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle FollowMsg
func handleFollowMsg(ctx sdk.Context, am types.AccountManager, msg FollowMsg) sdk.Result {
	if !am.AccountExist(ctx, msg.Follower) || !am.AccountExist(ctx, msg.Followee) {
		return ErrUsernameNotFound("Username not found").Result()
	}

	// add the "msg.Follower" to the "msg.Followee" 's follower list.
	followerList, err := am.GetFollower(ctx, msg.Followee)
	if err != nil {
		return ErrAccountManagerFail("Get follower list failed").Result()
	}

	if isInFollowerList(msg.Follower, followerList) == false {
		*followerList = append(*followerList, msg.Follower)
		if err := am.SetFollower(ctx, msg.Followee, followerList); err != nil {
			return ErrAccountManagerFail("Set follower failed").Result()
		}
	}

	// add the "msg.Followee" to the "msg.Follower" 's following list.
	followingList, err := am.GetFollowing(ctx, msg.Follower)
	if err != nil {
		return ErrAccountManagerFail("Get following list failed").Result()
	}

	if isInFollowingList(msg.Followee, followingList) == false {
		*followingList = append(*followingList, msg.Followee)
		if err := am.SetFollowing(ctx, msg.Followee, followingList); err != nil {
			return ErrAccountManagerFail("Set following failed").Result()
		}
	}

	return sdk.Result{}
}

// Handle UnfollowMsg
func handleUnfollowMsg(ctx sdk.Context, am types.AccountManager, msg UnfollowMsg) sdk.Result {
	if !am.AccountExist(ctx, msg.Follower) || !am.AccountExist(ctx, msg.Followee) {
		return ErrUsernameNotFound("Username not found").Result()
	}

	// remove the "msg.Follower" from the "msg.Followee" 's follower list.
	followerList, err := am.GetFollower(ctx, msg.Followee)
	if err != nil {
		return ErrAccountManagerFail("Get follower list failed").Result()
	}

	for index, user := range *followerList {
		if user == msg.Follower {
			*followerList = append(*followerList[:index], *followerList[index+1]...)
			if err := am.SetFollower(ctx, msg.Followee, followerList); err != nil {
				return ErrAccountManagerFail("Set follower failed").Result()
			}
			break
		}
	}

	// remove the "msg.Followee" from the "msg.Follower" 's following list.
	followingList, err := am.GetFollowing(ctx, msg.Follower)
	if err != nil {
		return ErrAccountManagerFail("Get following list failed").Result()
	}

	for index, user := range *followingList {
		if user == msg.Followee {
			*followingList = append(*followingList[:index], *followingList[index+1]...)
			if err := am.SetFollowing(ctx, msg.Follower, followingList); err != nil {
				return ErrAccountManagerFail("Set following failed").Result()
			}
			break
		}
	}

	return sdk.Result{}
}

// Handle TransferMsg
func handleTransferMsg(ctx sdk.Context, am types.AccountManager, msg TransferMsg) sdk.Result {
	if !am.AccountExist(ctx, msg.Follower) || !am.AccountExist(ctx, msg.Followee) {
		return ErrUsernameNotFound("Username not found").Result()
	}
}

// helper function
func isInFollowerList(me types.AccountKey, lst *Follower) bool {
	for _, user := range *lst {
		if user == me {
			return true
		}
	}
	return false
}

func isInFollowingList(me types.AccountKey, lst *Following) bool {
	for _, user := range *lst {
		if user == me {
			return true
		}
	}
	return false
}
