package post

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

func NewHandler(pm PostManager, am acc.AccountManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreatePostMsg:
			return handleCreatePostMsg(ctx, pm, am, msg)
		case DonateMsg:
			return handleDonateMsg(ctx, pm, am, msg)
		case LikeMsg:
			return handleLikeMsg(ctx, pm, am, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized account Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}

	}
}

// Handle RegisterMsg
func handleCreatePostMsg(ctx sdk.Context, pm PostManager, am acc.AccountManager, msg CreatePostMsg) sdk.Result {
	account := acc.NewProxyAccount(msg.Author, &am)
	if !account.IsAccountExist(ctx) {
		return acc.ErrUsernameNotFound(string(msg.Author)).Result()
	}
	post := NewProxyPost(msg.Author, msg.PostID, &pm)
	if post.IsPostExist(ctx) {
		return ErrPostExist().Result()
	}
	if err := post.CreatePost(ctx, &msg.PostInfo); err != nil {
		return err.Result()
	}
	if len(msg.ParentAuthor) > 0 || len(msg.ParentPostID) > 0 {
		parentPost := NewProxyPost(msg.ParentAuthor, msg.ParentPostID, &pm)
		if err := parentPost.AddComment(ctx, post.GetPostKey()); err != nil {
			return err.Result()
		}
		if err := parentPost.Apply(ctx); err != nil {
			return err.Result()
		}
	}
	if err := post.Apply(ctx); err != nil {
		return err.Result()
	}
	if err := account.UpdateLastActivity(ctx); err != nil {
		return err.Result()
	}
	if err := account.Apply(ctx); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle LikeMsg
func handleLikeMsg(ctx sdk.Context, pm PostManager, am acc.AccountManager, msg LikeMsg) sdk.Result {
	account := acc.NewProxyAccount(msg.Username, &am)
	if !account.IsAccountExist(ctx) {
		return acc.ErrUsernameNotFound(string(msg.Username)).Result()
	}
	post := NewProxyPost(msg.Author, msg.PostID, &pm)
	if !post.IsPostExist(ctx) {
		return ErrLikePostDoesntExist().Result()
	}
	// TODO: check acitivity burden
	like := Like{Username: msg.Username, Weight: msg.Weight}
	if err := post.AddOrUpdateLikeToPost(ctx, like); err != nil {
		return err.Result()
	}
	if err := account.UpdateLastActivity(ctx); err != nil {
		return err.Result()
	}

	// apply change to storage
	if err := post.Apply(ctx); err != nil {
		return err.Result()
	}
	if err := account.Apply(ctx); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle DonateMsg
func handleDonateMsg(ctx sdk.Context, pm PostManager, am acc.AccountManager, msg DonateMsg) sdk.Result {
	account := acc.NewProxyAccount(msg.Username, &am)
	if !account.IsAccountExist(ctx) {
		return acc.ErrUsernameNotFound(string(msg.Username)).Result()
	}
	post := NewProxyPost(msg.Author, msg.PostID, &pm)
	if !post.IsPostExist(ctx) {
		return ErrDonatePostDoesntExist().Result()
	}
	// TODO: check acitivity burden
	if err := account.MinusCoins(ctx, msg.Amount); err != nil {
		return err.Result()
	}
	donation := Donation{
		Username: msg.Username,
		Amount:   msg.Amount,
		Created:  types.Height(ctx.BlockHeight()),
	}
	if err := post.AddDonation(ctx, donation); err != nil {
		return err.Result()
	}
	if err := account.UpdateLastActivity(ctx); err != nil {
		return err.Result()
	}

	// apply change to storage
	if err := post.Apply(ctx); err != nil {
		return err.Result()
	}
	if err := account.Apply(ctx); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
