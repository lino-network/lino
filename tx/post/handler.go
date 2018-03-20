package post

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

func NewHandler(pm types.PostManager, am types.AccountManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreatePostMsg:
			return handleCreatePostMsg(ctx, pm, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle RegisterMsg
func handleCreatePostMsg(ctx sdk.Context, pm types.PostManager, am types.AccountManager, msg CreatePostMsg) sdk.Result {
	_, err := am.GetMeta(ctx, msg.Author)
	if err != nil {
		return err.Result()
	}
	// TODO: check activity burden
	if err := pm.CreatePost(ctx, &msg.Post); err != nil {
		return err.Result()
	}
	// TODO: update user activity
	return sdk.Result{}
}
