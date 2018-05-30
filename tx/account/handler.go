package account

import (
	"fmt"
	"reflect"

	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		case RecoverMsg:
			return handleRecoverMsg(ctx, am, msg)
		case RegisterMsg:
			return handleRegisterMsg(ctx, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account msg type: %v", reflect.TypeOf(msg).Name())
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
	if !am.IsAccountExist(ctx, msg.Receiver) || !am.IsAccountExist(ctx, msg.Sender) {
		return ErrUsernameNotFound().Result()
	}
	// withdraw money from sender's bank
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	if err := am.MinusSavingCoin(ctx, msg.Sender, coin, types.TransferOut); err != nil {
		return err.Result()
	}

	// send coins using username
	if err := am.AddSavingCoin(ctx, msg.Receiver, coin, types.TransferIn); err != nil {
		return ErrTransferHandler(msg.Sender).TraceCause(err, "").Result()
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

func handleRecoverMsg(ctx sdk.Context, am AccountManager, msg RecoverMsg) sdk.Result {
	// recover
	if !am.IsAccountExist(ctx, msg.Username) {
		return ErrUsernameNotFound().Result()
	}
	if err := am.RecoverAccount(
		ctx, msg.Username, msg.NewMasterPubKey, msg.NewTransactionPubKey,
		msg.NewPostPubKey); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle RegisterMsg
func handleRegisterMsg(ctx sdk.Context, am AccountManager, msg RegisterMsg) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Referrer) {
		return ErrUsernameNotFound().Result()
	}
	coin, err := types.LinoToCoin(msg.RegisterFee)
	if err != nil {
		return err.Result()
	}
	if err := am.MinusSavingCoin(ctx, msg.Referrer, coin, types.TransferOut); err != nil {
		return err.Result()
	}
	if err := am.CreateAccount(
		ctx, msg.NewUser, msg.NewMasterPubKey, msg.NewPostPubKey,
		msg.NewTransactionPubKey, coin); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
