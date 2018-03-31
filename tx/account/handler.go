package account

import (
	"bytes"
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func NewHandler(am AccountManager) sdk.Handler {
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
func handleFollowMsg(ctx sdk.Context, am AccountManager, msg FollowMsg) sdk.Result {
	proxyFollowee := NewProxyAccount(msg.Followee, &am)
	proxyFollower := NewProxyAccount(msg.Follower, &am)
	if !proxyFollowee.IsAccountExist(ctx) || !proxyFollower.IsAccountExist(ctx) {
		return ErrUsernameNotFound().Result()
	}
	// add the "msg.Follower" to the "msg.Followee" 's follower list.
	// add "msg.Followee/msg.Follower" key under "follower" prefix.
	// if findAccountInList(msg.Follower, followerList.Follower) == -1 {
	// 	followerList.Follower = append(followerList.Follower, msg.Follower)
	if err := proxyFollowee.SetFollower(ctx, msg.Follower); err != nil {
		return err.Result()
	}

	// add the "msg.Followee" to the "msg.Follower" 's following list.
	// add "msg.Follower/msg.Followee" key under "following" prefix
	//if findAccountInList(msg.Followee, followingList.Following) == -1 {
	//followingList.Following = append(followingList.Following, msg.Followee)
	if err := proxyFollower.SetFollowing(ctx, msg.Followee); err != nil {
		return err.Result()
	}
	//}
	proxyFollowee.Apply(ctx)
	proxyFollower.Apply(ctx)
	return sdk.Result{}
}

// Handle UnfollowMsg
func handleUnfollowMsg(ctx sdk.Context, am AccountManager, msg UnfollowMsg) sdk.Result {
	proxyFollowee := NewProxyAccount(msg.Followee, &am)
	proxyFollower := NewProxyAccount(msg.Follower, &am)

	if !proxyFollowee.IsAccountExist(ctx) || !proxyFollower.IsAccountExist(ctx) {
		return ErrUsernameNotFound().Result()
	}

	// add the "msg.Follower" to the "msg.Followee" 's follower list.
	// add "msg.Followee/msg.Follower" key under "follower" prefix.
	// if findAccountInList(msg.Follower, followerList.Follower) == -1 {
	// 	followerList.Follower = append(followerList.Follower, msg.Follower)
	if err := proxyFollowee.RemoveFollower(ctx, msg.Follower); err != nil {
		return err.Result()
	}

	// add the "msg.Followee" to the "msg.Follower" 's following list.
	// add "msg.Follower/msg.Followee" key under "following" prefix
	//if findAccountInList(msg.Followee, followingList.Following) == -1 {
	//followingList.Following = append(followingList.Following, msg.Followee)
	if err := proxyFollower.RemoveFollowing(ctx, msg.Followee); err != nil {
		return err.Result()
	}
	//}
	proxyFollowee.Apply(ctx)
	proxyFollower.Apply(ctx)
	return sdk.Result{}
}

// Handle TransferMsg
func handleTransferMsg(ctx sdk.Context, am AccountManager, msg TransferMsg) sdk.Result {
	// withdraw money from sender's bank
	accSender := NewProxyAccount(msg.Sender, &am)
	if err := accSender.MinusCoin(ctx, types.LinoToCoin(msg.Amount)); err != nil {
		return err.Result()
	}

	// both username and address provided
	if len(msg.ReceiverName) != 0 && len(msg.ReceiverAddr) != 0 {
		// check if username and address match
		associatedAddr, err := NewProxyAccount(msg.ReceiverName, &am).GetBankAddress(ctx)
		if !bytes.Equal(associatedAddr, msg.ReceiverAddr) || err != nil {
			return ErrUsernameAddressMismatch().Result()
		}
	}

	// send coins using username
	if len(msg.ReceiverName) != 0 {
		accReceiver := NewProxyAccount(msg.ReceiverName, &am)
		if err := accReceiver.AddCoin(ctx, types.LinoToCoin(msg.Amount)); err != nil {
			return ErrAddMoneyFailed().Result()
		}
		accSender.Apply(ctx)
		accReceiver.Apply(ctx)
		return sdk.Result{}
	}

	// send coins using address (even no account bank associated with this addr)
	receiverBank, err := am.GetBankFromAddress(ctx, msg.ReceiverAddr)
	if err == nil {
		// account bank exists
		receiverBank.Balance = receiverBank.Balance.Plus(types.LinoToCoin(msg.Amount))
	} else {
		// account bank not found, create a new one for this address
		receiverBank = &AccountBank{
			Address: msg.ReceiverAddr,
			Balance: types.LinoToCoin(msg.Amount),
		}
	}

	if setErr := am.SetBankFromAddress(ctx, msg.ReceiverAddr, receiverBank); setErr != nil {
		return setErr.Result()
	}
	accSender.Apply(ctx)
	return sdk.Result{}
}
