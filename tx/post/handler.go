package post

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
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
		return ErrPostCreateNonExistAuthor().Result()
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

// Handle DonateMsg
func handleLikeMsg(ctx sdk.Context, pm types.PostManager, am types.AccountManager, msg LikeMsg) sdk.Result {
	_, err := am.GetMeta(ctx, msg.Username)
	if err != nil {
		return err.Result()
	}
	postKey := types.GetPostKey(msg.Author, msg.PostID)
	postMeta, err := pm.GetPostMeta(ctx, postKey)
	if err != nil {
		return err.Result()
	}
	// TODO: check acitivity burden
	postLikes, err := pm.GetPostLikes(ctx, postKey)
	if err != nil {
		return err.Result()
	}
	index := getLikeFromList(postLikes.Likes, username)
	if index == -1 {
		like := types.Like{
			Username: msg.Username,
			Weight:   msg.Weight,
		}
		postLikes.Likes = append(postLikes.Likes, like)
		postLikes.TotalWeight += like.Weight
	} else {
		postLikes.TotalWeight -= postLikes.Likes[index].Weight
		postLikes.Likes[index].Weight = msg.Weight
		postLikes.TotalWeight += postLikes.Likes[index].Weight
	}
	postMeta.LastActicity = types.Height(ctx.BlockHeight())

	if err := pm.SetPostLikes(ctx, postKey, postLikes); err != nil {
		return err.Result()
	}
	if err := pm.SetPostMeta(ctx, postKey, postMeta); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

// Handle DonateMsg
func handleDonateMsg(ctx sdk.Context, pm types.PostManager, am types.AccountManager, msg DonateMsg) sdk.Result {
	_, err := am.GetMeta(ctx, msg.Username)
	if err != nil {
		return err.Result()
	}
	postKey := types.GetPostKey(msg.Author, msg.PostID)
	postMeta, err := pm.GetPostMeta(ctx, postKey)
	if err != nil {
		return err.Result()
	}
	// TODO: check acitivity burden
	postDonations, err := pm.GetPostDonations(ctx, postKey)
	if err != nil {
		return err.Result()
	}
	bank, err := GetBankFromAccountKey(ctx, msg.Username)

	if msg.Amount.IsGTE(bank.Coins) {
		return ErrPostDonateInsufficient().Result()
	}
	donation := types.Donation{
		Username: msg.Username,
		Amount:   msg.Amount,
		Created:  types.Height(ctx.BlockHeight()),
	}
	postDonations.Donations = append(postDonations.Donations, donation)
	postDonations.Reward = postDonations.Reward.Plus(donation.Amount)
	bank.Coins = bank.Coins.Minus(msg.Amount)
	postMeta.LastActicity = types.Height(ctx.BlockHeight())

	if err := am.SetBank(ctx, bank.Address, bank); err != nil {
		return err.Result()
	}
	if err := pm.SetPostDonations(ctx, postKey, postDonations); err != nil {
		return err.Result()
	}
	if err := pm.SetPostMeta(ctx, postKey, postMeta); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func getLikeFromList(likes []types.Like, user types.AccountKey) int64 {
	for i, like := range likes {
		if like.Username == user {
			return i
		}
	}
	return -1
}
