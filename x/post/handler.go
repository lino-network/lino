package post

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"

	linotypes "github.com/lino-network/lino/types"
	types "github.com/lino-network/lino/x/post/types"
)

type CreatePostMsg = types.CreatePostMsg
type UpdatePostMsg = types.UpdatePostMsg
type DeletePostMsg = types.DeletePostMsg
type DonateMsg = types.DonateMsg
type IDADonateMsg = types.IDADonateMsg

// NewHandler - Handle all "post" type messages.
func NewHandler(pm PostKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreatePostMsg:
			return handleCreatePostMsg(ctx, msg, pm)
		case UpdatePostMsg:
			return handleUpdatePostMsg(ctx, msg, pm)
		case DeletePostMsg:
			return handleDeletePostMsg(ctx, msg, pm)
		case DonateMsg:
			return handleDonateMsg(ctx, msg, pm)
		case IDADonateMsg:
			return handleIDADonateMsg(ctx, msg, pm)
		default:
			errMsg := fmt.Sprintf("Unrecognized post msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle createPostMsg
func handleCreatePostMsg(ctx sdk.Context, msg CreatePostMsg, pm PostKeeper) sdk.Result {
	err := pm.CreatePost(ctx, msg.Author, msg.PostID, msg.CreatedBy, msg.Content, msg.Title)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleUpdatePostMsg(ctx sdk.Context, msg UpdatePostMsg, pm PostKeeper) sdk.Result {
	err := pm.UpdatePost(ctx, msg.Author, msg.PostID, msg.Title, msg.Content)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleDeletePostMsg(ctx sdk.Context, msg DeletePostMsg, pm PostKeeper) sdk.Result {
	err := pm.DeletePost(ctx, linotypes.GetPermlink(msg.Author, msg.PostID))
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle DonateMsg
func handleDonateMsg(ctx sdk.Context, msg DonateMsg, pm PostKeeper) sdk.Result {
	amount, err := linotypes.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	err = pm.LinoDonate(ctx, msg.Username, amount, msg.Author, msg.PostID, msg.FromApp)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleIDADonateMsg(ctx sdk.Context, msg IDADonateMsg, pm PostKeeper) sdk.Result {
	// amount must be an positive integer.
	amount, err := msg.Amount.ToIDA()
	if err != nil {
		return err.Result()
	}
	err = pm.IDADonate(ctx, msg.Username, amount, msg.Author, msg.PostID, msg.App)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
