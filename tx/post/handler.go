package post

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/global"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

func NewHandler(pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreatePostMsg:
			return handleCreatePostMsg(ctx, msg, pm, am, gm)
		case DonateMsg:
			return handleDonateMsg(ctx, msg, pm, am, gm)
		case LikeMsg:
			return handleLikeMsg(ctx, msg, pm, am, gm)
		default:
			errMsg := fmt.Sprintf("Unrecognized account Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle RegisterMsg
func handleCreatePostMsg(ctx sdk.Context, msg CreatePostMsg, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Result {
	account := acc.NewProxyAccount(msg.Author, &am)
	if !account.IsAccountExist(ctx) {
		return acc.ErrUsernameNotFound().Result()
	}
	post := NewPostProxy(msg.Author, msg.PostID, &pm)
	if post.IsPostExist(ctx) {
		return ErrPostExist().Result()
	}
	if err := post.CreatePost(ctx, &msg.PostInfo); err != nil {
		return err.Result()
	}
	if len(msg.ParentAuthor) > 0 || len(msg.ParentPostID) > 0 {
		parentPost := NewPostProxy(msg.ParentAuthor, msg.ParentPostID, &pm)
		comment := Comment{Author: post.GetAuthor(), PostID: post.GetPostID(), Created: types.Height(ctx.BlockHeight())}
		if err := parentPost.AddComment(ctx, comment); err != nil {
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
func handleLikeMsg(ctx sdk.Context, msg LikeMsg, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Result {
	account := acc.NewProxyAccount(msg.Username, &am)
	if !account.IsAccountExist(ctx) {
		return acc.ErrUsernameNotFound().Result()
	}
	post := NewPostProxy(msg.Author, msg.PostID, &pm)
	if !post.IsPostExist(ctx) {
		return ErrLikePostDoesntExist().Result()
	}
	// TODO: check acitivity burden
	like := Like{Username: msg.Username, Weight: msg.Weight, Created: types.Height(ctx.BlockHeight())}
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
func handleDonateMsg(ctx sdk.Context, msg DonateMsg, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Result {
	globalProxy := global.NewGlobalProxy(&gm)

	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return err.Result()
	}
	account := acc.NewProxyAccount(msg.Username, &am)
	if !account.IsAccountExist(ctx) {
		return acc.ErrUsernameNotFound().Result()
	}
	post := NewPostProxy(msg.Author, msg.PostID, &pm)
	if !post.IsPostExist(ctx) {
		return ErrDonatePostDoesntExist().Result()
	}
	// TODO: check acitivity burden
	if err := account.MinusCoin(ctx, coin); err != nil {
		return err.Result()
	}
	sourcePost, err := post.GetRootSourcePost(ctx)
	if err != nil {
		return err.Result()
	}
	if sourcePost != nil {
		redistributionSplitRate, err := sourcePost.GetRedistributionSplitRate(ctx)
		if err != nil {
			return err.Result()
		}
		sourceIncome := types.Coin{sdk.NewRat(coin.Amount).Mul(sdk.OneRat.Sub(redistributionSplitRate)).Evaluate()}
		coin.Amount -= sourceIncome.Amount
		if err := ProcessDonationFriction(ctx, msg.Username, sourceIncome, sourcePost, am, globalProxy); err != nil {
			return err.Result()
		}
		if err := sourcePost.Apply(ctx); err != nil {
			return err.Result()
		}
	}
	if err := ProcessDonationFriction(ctx, msg.Username, coin, post, am, globalProxy); err != nil {
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

func ProcessDonationFriction(
	ctx sdk.Context, consumer acc.AccountKey, coin types.Coin,
	post *PostProxy, am acc.AccountManager, globalProxy *global.GlobalProxy) sdk.Error {
	if coin.IsZero() {
		return nil
	}
	authorAccount := acc.NewProxyAccount(post.GetAuthor(), &am)
	if !authorAccount.IsAccountExist(ctx) {
		return acc.ErrUsernameNotFound()
	}
	consumptionFrictionRate, err := globalProxy.GetConsumptionFrictionRate(ctx)
	if err != nil {
		return err
	}
	redistribute := types.Coin{sdk.NewRat(coin.Amount).Mul(consumptionFrictionRate).Evaluate()}
	directDeposit := coin.Minus(redistribute)
	if err := post.AddDonation(ctx, consumer, directDeposit); err != nil {
		return err
	}
	if err := authorAccount.AddCoin(ctx, directDeposit); err != nil {
		return err
	}
	if err := globalProxy.AddConsumption(ctx, coin); err != nil {
		return err
	}
	if err := globalProxy.AddRedistributeCoin(ctx, redistribute); err != nil {
		return err
	}
	rewardEvent := RewardEvent{
		PostAuthor: post.GetAuthor(),
		PostID:     post.GetPostID(),
		Consumer:   consumer,
		Amount:     coin,
	}
	if err := globalProxy.RegisterRedistributionEvent(ctx, rewardEvent); err != nil {
		return err
	}
	if err := authorAccount.Apply(ctx); err != nil {
		return err
	}
	return nil
}
