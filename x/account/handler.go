package account

import (
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/lino-network/lino/recorder/user"

	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/global"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler - Handle all "account" type messages.
func NewHandler(am AccountManager, gm global.GlobalManager) sdk.Handler {
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
			return handleRegisterMsg(ctx, am, gm, msg)
		case UpdateAccountMsg:
			return handleUpdateAccountMsg(ctx, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleFollowMsg(ctx sdk.Context, am AccountManager, msg FollowMsg) sdk.Result {
	if !am.DoesAccountExist(ctx, msg.Followee) {
		return ErrFolloweeNotFound(msg.Followee).Result()
	}
	if !am.DoesAccountExist(ctx, msg.Follower) {
		return ErrFollowerNotFound(msg.Follower).Result()
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
	if !am.DoesAccountExist(ctx, msg.Followee) {
		return ErrFolloweeNotFound(msg.Followee).Result()
	}
	if !am.DoesAccountExist(ctx, msg.Follower) {
		return ErrFollowerNotFound(msg.Follower).Result()
	}

	// remove the "msg.Follower" from the "msg.Followee" 's follower list.
	// remove "msg.Followee/msg.Follower" key under "follower" prefix.
	if err := am.RemoveFollower(ctx, msg.Followee, msg.Follower); err != nil {
		return err.Result()
	}

	// remove the "msg.Followee" from the "msg.Follower" 's following list.
	// remove "msg.Follower/msg.Followee" key under "following" prefix
	if err := am.RemoveFollowing(ctx, msg.Follower, msg.Followee); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleTransferMsg(ctx sdk.Context, am AccountManager, msg TransferMsg) sdk.Result {
	if !am.DoesAccountExist(ctx, msg.Receiver) {
		return ErrReceiverNotFound(msg.Receiver).Result()
	}

	if !am.DoesAccountExist(ctx, msg.Sender) {
		return ErrSenderNotFound(msg.Sender).Result()
	}
	// withdraw money from sender's bank
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	if err := am.MinusSavingCoin(
		ctx, msg.Sender, coin, msg.Receiver, msg.Memo, types.TransferOut); err != nil {
		return err.Result()
	}

	// send coins using username
	if err := am.AddSavingCoin(
		ctx, msg.Receiver, coin, msg.Sender, msg.Memo, types.TransferIn); err != nil {
		return err.Result()
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
	if !am.DoesAccountExist(ctx, msg.Username) {
		return ErrAccountNotFound(msg.Username).Result()
	}
	if err := am.RecoverAccount(
		ctx, msg.Username, msg.NewResetPubKey, msg.NewTransactionPubKey,
		msg.NewAppPubKey); err != nil {
		return err.Result()
	}
	am.recorder.UserRepository.UpdatePubKey(
		string(msg.Username),
		hex.EncodeToString(msg.NewResetPubKey.Bytes()),
		hex.EncodeToString(msg.NewTransactionPubKey.Bytes()),
		hex.EncodeToString(msg.NewAppPubKey.Bytes()))
	return sdk.Result{}
}

// Handle RegisterMsg
func handleRegisterMsg(ctx sdk.Context, am AccountManager, gm global.GlobalManager, msg RegisterMsg) sdk.Result {
	if !am.DoesAccountExist(ctx, msg.Referrer) {
		return ErrReferrerNotFound(msg.Referrer).Result()
	}
	coin, err := types.LinoToCoin(msg.RegisterFee)
	if err != nil {
		return err.Result()
	}
	accParams, err := am.paramHolder.GetAccountParam(ctx)
	if err != nil {
		return err.Result()
	}
	if accParams.RegisterFee.IsGT(coin) {
		return ErrRegisterFeeInsufficient().Result()
	}
	if err := am.MinusSavingCoin(
		ctx, msg.Referrer, coin, msg.NewUser, "", types.TransferOut); err != nil {
		return err.Result()
	}
	// the open account fee will be added to developer inflation pool
	if err := gm.AddToDeveloperInflationPool(ctx, accParams.RegisterFee); err != nil {
		return err.Result()
	}

	if err := am.CreateAccount(
		ctx, msg.Referrer, msg.NewUser, msg.NewResetPubKey, msg.NewTransactionPubKey,
		msg.NewAppPubKey, coin.Minus(accParams.RegisterFee)); err != nil {
		return err.Result()
	}

	balanceCoin := coin.Minus(accParams.RegisterFee)
	balance, _ := balanceCoin.ToInt64()
	user := &user.User{
		Username:          string(msg.NewUser),
		CreatedAt:         ctx.BlockHeader().Time,
		Referrer:          string(msg.Referrer),
		ResetPubKey:       hex.EncodeToString(msg.NewResetPubKey.Bytes()),
		TransactionPubKey: hex.EncodeToString(msg.NewTransactionPubKey.Bytes()),
		AppPubKey:         hex.EncodeToString(msg.NewAppPubKey.Bytes()),
		Saving:            balance,
	}
	am.recorder.UserRepository.Add(user)
	return sdk.Result{}
}

// Handle RegisterMsg
func handleUpdateAccountMsg(ctx sdk.Context, am AccountManager, msg UpdateAccountMsg) sdk.Result {
	if !am.DoesAccountExist(ctx, msg.Username) {
		return ErrAccountNotFound(msg.Username).Result()
	}
	if err := am.UpdateJSONMeta(ctx, msg.Username, msg.JSONMeta); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
