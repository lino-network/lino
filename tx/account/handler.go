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
		case ClaimMsg:
			return handleClaimMsg(ctx, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleFollowMsg(ctx sdk.Context, am AccountManager, msg FollowMsg) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Followee) || !am.IsAccountExist(ctx, msg.Follower) {
		return ErrUsernameNotFound().Result()
	}
	// add the "msg.Follower" to the "msg.Followee" 's follower list.
	// add "msg.Followee/msg.Follower" key under "follower" prefix.
	if err := am.SetFollower(ctx, msg.Followee, msg.Follower); err != nil {
		return err.Result()
	}

	// add the "msg.Followee" to the "msg.Follower" 's following list.
	// add "msg.Follower/msg.Followee" key under "following" prefix
	if err := am.SetFollowing(ctx, msg.Follower, msg.Followee); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleUnfollowMsg(ctx sdk.Context, am AccountManager, msg UnfollowMsg) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Followee) || !am.IsAccountExist(ctx, msg.Follower) {
		return ErrUsernameNotFound().Result()
	}

	// add the "msg.Follower" to the "msg.Followee" 's follower list.
	// add "msg.Followee/msg.Follower" key under "follower" prefix.
	if err := am.RemoveFollower(ctx, msg.Followee, msg.Follower); err != nil {
		return err.Result()
	}

	// add the "msg.Followee" to the "msg.Follower" 's following list.
	// add "msg.Follower/msg.Followee" key under "following" prefix
	if err := am.RemoveFollowing(ctx, msg.Follower, msg.Followee); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleTransferMsg(ctx sdk.Context, am AccountManager, msg TransferMsg) sdk.Result {
	// withdraw money from sender's bank
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	if err := am.MinusCoin(ctx, msg.Sender, coin); err != nil {
		return err.Result()
	}

	// both username and address provided
	if len(msg.ReceiverName) != 0 && len(msg.ReceiverAddr) != 0 {
		// check if username and address match
		associatedAddr, err := am.GetBankAddress(ctx, msg.ReceiverName)
		if !bytes.Equal(associatedAddr, msg.ReceiverAddr) || err != nil {
			return ErrTransferHandler(msg.Sender).Result()
		}
	}

	// send coins using username
	if len(msg.ReceiverName) != 0 {
		if err := am.AddCoin(ctx, msg.ReceiverName, coin); err != nil {
			return ErrTransferHandler(msg.Sender).TraceCause(err, "").Result()
		}
		return sdk.Result{}
	}

	if setErr := am.AddCoinToAddress(ctx, msg.ReceiverAddr, coin); setErr != nil {
		return ErrTransferHandler(msg.Sender).TraceCause(setErr, "").Result()
	}
	return sdk.Result{}
}

func handleClaimMsg(ctx sdk.Context, am AccountManager, msg ClaimMsg) sdk.Result {
	// claim reward
	if err := am.ClaimReward(ctx, msg.Username); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
