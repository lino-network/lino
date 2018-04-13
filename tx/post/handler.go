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
		case ReportOrUpvoteMsg:
			return handleReportOrUpvoteMsg(ctx, msg, pm, am, gm)
		default:
			errMsg := fmt.Sprintf("Unrecognized account Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle RegisterMsg
func handleCreatePostMsg(ctx sdk.Context, msg CreatePostMsg, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Author) {
		return ErrCreatePostAuthorNotFound(msg.Author).Result()
	}
	postKey := types.GetPostKey(msg.Author, msg.PostID)
	if pm.IsPostExist(ctx, postKey) {
		return ErrCreateExistPost(postKey).Result()
	}
	if len(msg.ParentAuthor) > 0 || len(msg.ParentPostID) > 0 {
		parentPostKey := types.GetPostKey(msg.ParentAuthor, msg.ParentPostID)
		if !pm.IsPostExist(ctx, parentPostKey) {
			return ErrCommentInvalidParent(parentPostKey).Result()
		}
		if err := pm.AddComment(ctx, parentPostKey, msg.Author, msg.PostID); err != nil {
			return err.Result()
		}
	}
	if err := pm.CreatePost(ctx, &msg.PostCreateParams); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// Handle LikeMsg
func handleLikeMsg(ctx sdk.Context, msg LikeMsg, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Username) {
		return ErrLikePostUserNotFound(msg.Username).Result()
	}
	postKey := types.GetPostKey(msg.Author, msg.PostID)
	if !pm.IsPostExist(ctx, postKey) {
		return ErrLikeNonExistPost(postKey).Result()
	}
	if err := pm.AddOrUpdateLikeToPost(ctx, postKey, msg.Username, msg.Weight); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// Handle DonateMsg
func handleDonateMsg(ctx sdk.Context, msg DonateMsg, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Result {
	postKey := types.GetPostKey(msg.Author, msg.PostID)
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "").Result()
	}
	if !am.IsAccountExist(ctx, msg.Username) {
		return ErrDonateUserNotFound(msg.Username).Result()
	}
	if !pm.IsPostExist(ctx, postKey) {
		return ErrDonatePostDoesntExist(postKey).Result()
	}
	// TODO: check acitivity burden
	if err := am.MinusCoin(ctx, msg.Username, coin); err != nil {
		return ErrDonateFailed(postKey).Result()
	}
	sourceAuthor, sourcePostID, err := pm.GetSourcePost(ctx, postKey)
	if err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "").Result()
	}
	if sourceAuthor != types.AccountKey("") && sourcePostID != "" {
		sourcePostKey := types.GetPostKey(sourceAuthor, sourcePostID)
		redistributionSplitRate, err := pm.GetRedistributionSplitRate(ctx, sourcePostKey)
		if err != nil {
			return ErrDonateFailed(postKey).TraceCause(err, "").Result()
		}
		sourceIncome := types.RatToCoin(coin.ToRat().Mul(sdk.OneRat.Sub(redistributionSplitRate)))
		coin = coin.Minus(sourceIncome)
		if err := processDonationFriction(ctx, msg.Username, sourceIncome, sourceAuthor, sourcePostID, am, pm, gm); err != nil {
			return err.Result()
		}
	}
	if err := processDonationFriction(ctx, msg.Username, coin, msg.Author, msg.PostID, am, pm, gm); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func processDonationFriction(
	ctx sdk.Context, consumer types.AccountKey, coin types.Coin, postAuthor types.AccountKey,
	postID string, am acc.AccountManager, pm PostManager, gm global.GlobalManager) sdk.Error {
	postKey := types.GetPostKey(postAuthor, postID)
	if coin.IsZero() {
		return nil
	}
	if !am.IsAccountExist(ctx, postAuthor) {
		return ErrDonateAuthorNotFound(postKey, postAuthor)
	}
	consumptionFrictionRate, err := gm.GetConsumptionFrictionRate(ctx)
	if err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "")
	}
	redistribute := types.RatToCoin(coin.ToRat().Mul(consumptionFrictionRate))
	directDeposit := coin.Minus(redistribute)
	if err := pm.AddDonation(ctx, postKey, consumer, directDeposit); err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "")
	}
	if err := am.AddCoin(ctx, postAuthor, directDeposit); err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "")
	}
	if err := gm.AddConsumption(ctx, coin); err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "")
	}
	if err := gm.AddConsumptionFrictionToRewardPool(ctx, redistribute); err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "")
	}
	rewardEvent := RewardEvent{
		PostAuthor: postAuthor,
		PostID:     postID,
		Consumer:   consumer,
		Amount:     coin,
	}
	if err := gm.RegisterContentRewardEvent(ctx, rewardEvent); err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "")
	}
	return nil
}

// Handle ReportMsgOrUpvoteMsg
func handleReportOrUpvoteMsg(
	ctx sdk.Context, msg ReportOrUpvoteMsg, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Username) {
		return ErrReportUserNotFound(msg.Username).Result()
	}
	postKey := types.GetPostKey(msg.Author, msg.PostID)
	stake, err := am.GetStake(ctx, msg.Username)
	if err != nil {
		return ErrReportFailed(postKey).TraceCause(err, "").Result()
	}
	if !pm.IsPostExist(ctx, postKey) {
		return ErrReportPostDoesntExist(postKey).Result()
	}
	sourceAuthor, sourcePostID, err := pm.GetSourcePost(ctx, postKey)
	if err != nil {
		return ErrReportFailed(postKey).TraceCause(err, "").Result()
	}
	if sourceAuthor != types.AccountKey("") && sourcePostID != "" {
		sourcePostKey := types.GetPostKey(sourceAuthor, sourcePostID)
		if err := pm.AddOrUpdateReportOrUpvoteToPost(
			ctx, sourcePostKey, msg.Username, stake, msg.IsReport); err != nil {
			return err.Result()
		}
	}
	if err := pm.AddOrUpdateReportOrUpvoteToPost(
		ctx, postKey, msg.Username, stake, msg.IsReport); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
