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
		case FollowMsg:
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
		followerList.Follower = append(followerList.Follower, msg.Follower)
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
		followingList.Following = append(followingList.Following, msg.Followee)
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

	for index, user := range followerList.Follower {
		if user == msg.Follower {
			followerList.Follower = append(followerList.Follower[:index], followerList.Follower[index+1:]...)
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

	for index, user := range followingList.Following {
		if user == msg.Followee {
			followingList.Following = append(followingList.Following[:index], followingList.Following[index+1:]...)
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
	if !am.AccountExist(ctx, msg.Sender) {
		return ErrUsernameNotFound("Username not found").Result()
	}

	// check if the sender has enough money
	senderBank, err := am.GetBankFromAccountKey(ctx, msg.Sender)
	if err != nil {
		return ErrAccountManagerFail("Get sender's account bank failed").Result()
	}

	if senderBank.Coins.IsGTE(msg.Amount) == false {
		return ErrAccountManagerFail("Sender's coins are not enough").Result()
	}

	// withdraw money from sender's bank
	senderBank.Coins.Minus(msg.Amount)
	if err := am.SetBankFromAccountKey(ctx, msg.Sender, senderBank); err != nil {
		return ErrAccountManagerFail("Set sender's bank failed").Result()
	}

	// send coins using username
	if am.AccountExist(ctx, msg.ReceiverName) {
		if receiverBank, err := am.GetBankFromAccountKey(ctx, msg.ReceiverName); err == nil {
			receiverBank.Coins.Plus(msg.Amount)
			if setErr := am.SetBankFromAccountKey(ctx, msg.ReceiverName, receiverBank); setErr != nil {
				return ErrAccountManagerFail("Set receiver's bank failed").Result()
			}
			return sdk.Result{}
		}
	}

	// send coins using address
	receiverBank, err := am.GetBankFromAddress(ctx, msg.ReceiverAddr)
	if err == nil {
		// account bank exists
		receiverBank.Coins.Plus(msg.Amount)
	} else {
		// account bank not found, create a new one for this address
		receiverBank = &types.AccountBank{
			Address:  msg.ReceiverAddr,
			Coins:    msg.Amount,
			Username: "",
		}
	}

	if setErr := am.SetBankFromAddress(ctx, msg.ReceiverAddr, receiverBank); setErr != nil {
		return ErrAccountManagerFail("Set receiver's bank failed").Result()
	}
	return sdk.Result{}
}

// helper function
func isInFollowerList(me types.AccountKey, lst *types.Follower) bool {
	for _, user := range lst.Follower {
		if user == me {
			return true
		}
	}
	return false
}

func isInFollowingList(me types.AccountKey, lst *types.Following) bool {
	for _, user := range lst.Following {
		if user == me {
			return true
		}
	}
	return false
}
