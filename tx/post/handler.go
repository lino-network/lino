package post

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/lino-network/lino/tx/global"
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	dev "github.com/lino-network/lino/tx/developer"
)

func NewHandler(pm PostManager, am acc.AccountManager, gm global.GlobalManager, dm dev.DeveloperManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case CreatePostMsg:
			return handleCreatePostMsg(ctx, msg, pm, am, gm)
		case DonateMsg:
			return handleDonateMsg(ctx, msg, pm, am, gm, dm)
		case LikeMsg:
			return handleLikeMsg(ctx, msg, pm, am, gm)
		case ReportOrUpvoteMsg:
			return handleReportOrUpvoteMsg(ctx, msg, pm, am, gm)
		case ViewMsg:
			return handleViewMsg(ctx, msg, pm, am, gm)
		case UpdatePostMsg:
			return handleUpdatePostMsg(ctx, msg, pm, am)
		case DeletePostMsg:
			return handleDeletePostMsg(ctx, msg, pm, am)
		default:
			errMsg := fmt.Sprintf("Unrecognized post msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle RegisterMsg
func handleCreatePostMsg(ctx sdk.Context, msg CreatePostMsg, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Author) {
		return ErrCreatePostAuthorNotFound(msg.Author).Result()
	}
	permLink := types.GetPermLink(msg.Author, msg.PostID)
	if pm.IsPostExist(ctx, permLink) {
		return ErrCreateExistPost(permLink).Result()
	}
	if len(msg.ParentAuthor) > 0 || len(msg.ParentPostID) > 0 {
		parentPostKey := types.GetPermLink(msg.ParentAuthor, msg.ParentPostID)
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
	permLink := types.GetPermLink(msg.Author, msg.PostID)
	if !pm.IsPostExist(ctx, permLink) {
		return ErrLikeNonExistPost(permLink).Result()
	}
	if err := pm.AddOrUpdateLikeToPost(ctx, permLink, msg.Username, msg.Weight); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// Handle ViewMsg
func handleViewMsg(ctx sdk.Context, msg ViewMsg, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Username) {
		return ErrViewPostUserNotFound(msg.Username).Result()
	}
	permLink := types.GetPermLink(msg.Author, msg.PostID)
	if !pm.IsPostExist(ctx, permLink) {
		return ErrViewNonExistPost(permLink).Result()
	}
	if err := pm.AddOrUpdateViewToPost(ctx, permLink, msg.Username); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// Handle DonateMsg
func handleDonateMsg(
	ctx sdk.Context, msg DonateMsg, pm PostManager, am acc.AccountManager,
	gm global.GlobalManager, dm dev.DeveloperManager) sdk.Result {
	permLink := types.GetPermLink(msg.Author, msg.PostID)
	coin, err := types.LinoToCoin(msg.Amount)
	if err != nil {
		return ErrDonateFailed(permLink).TraceCause(err, "").Result()
	}
	if !am.IsAccountExist(ctx, msg.Username) {
		return ErrDonateUserNotFound(msg.Username).Result()
	}
	if !pm.IsPostExist(ctx, permLink) {
		return ErrDonatePostNotFound(permLink).Result()
	}

	if msg.Username == msg.Author {
		return ErrDonateToSelf(msg.Username).Result()
	}
	if msg.FromApp != "" {
		if !dm.IsDeveloperExist(ctx, msg.FromApp) {
			return ErrDonateFailed(permLink).Result()
		}
	}
	if err := am.MinusSavingCoin(ctx, msg.Username, coin, types.DonationOut); err != nil {
		return ErrAccountSavingCoinNotEnough(permLink).Result()
	}
	sourceAuthor, sourcePostID, err := pm.GetSourcePost(ctx, permLink)
	if err != nil {
		return ErrDonateFailed(permLink).TraceCause(err, "").Result()
	}
	if sourceAuthor != types.AccountKey("") && sourcePostID != "" {
		sourcePermLink := types.GetPermLink(sourceAuthor, sourcePostID)

		redistributionSplitRate, err := pm.GetRedistributionSplitRate(ctx, sourcePermLink)
		if err != nil {
			return ErrDonateFailed(permLink).TraceCause(err, "").Result()
		}
		sourceIncome, err := types.RatToCoin(new(big.Rat).Mul(coin.ToRat(), sdk.OneRat.Sub(redistributionSplitRate).GetRat()))
		if err != nil {
			return err.Result()
		}
		coin = coin.Minus(sourceIncome)
		if err := processDonationFriction(
			ctx, msg.Username, sourceIncome, sourceAuthor, sourcePostID, msg.FromApp, am, pm, gm); err != nil {
			return err.Result()
		}
	}
	if err := processDonationFriction(
		ctx, msg.Username, coin, msg.Author, msg.PostID, msg.FromApp, am, pm, gm); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func processDonationFriction(
	ctx sdk.Context, consumer types.AccountKey, coin types.Coin,
	postAuthor types.AccountKey, postID string, fromApp types.AccountKey,
	am acc.AccountManager, pm PostManager, gm global.GlobalManager) sdk.Error {
	postKey := types.GetPermLink(postAuthor, postID)
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
	frictionCoin, err := types.RatToCoin(new(big.Rat).Mul(coin.ToRat(), consumptionFrictionRate.GetRat()))
	if err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "")
	}
	// evaluate this consumption can get the result, the result is used to get inflation from pool
	evaluateResult, err := evaluateConsumption(ctx, consumer, coin, postAuthor, postID, am, pm, gm)
	if err != nil {
		return err
	}
	rewardEvent := RewardEvent{
		PostAuthor: postAuthor,
		PostID:     postID,
		Consumer:   consumer,
		Evaluate:   evaluateResult,
		Original:   coin,
		Friction:   frictionCoin,
		FromApp:    fromApp,
	}
	if err := gm.AddFrictionAndRegisterContentRewardEvent(
		ctx, rewardEvent, frictionCoin, evaluateResult); err != nil {
		return err
	}

	directDeposit := coin.Minus(frictionCoin)
	if err := pm.AddDonation(ctx, postKey, consumer, directDeposit, types.DirectDeposit); err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "")
	}
	if err := am.AddSavingCoin(ctx, postAuthor, directDeposit, types.DonationIn); err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "")
	}
	if err := gm.AddConsumption(ctx, coin); err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "")
	}
	if err := am.UpdateDonationRelationship(ctx, postAuthor, consumer); err != nil {
		return ErrDonateFailed(postKey).TraceCause(err, "")
	}
	return nil
}

func evaluateConsumption(
	ctx sdk.Context, consumer types.AccountKey, coin types.Coin, postAuthor types.AccountKey,
	postID string, am acc.AccountManager, pm PostManager, gm global.GlobalManager) (types.Coin, sdk.Error) {
	numOfConsumptionOnAuthor, err := am.GetDonationRelationship(ctx, consumer, postAuthor)
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	created, totalReward, err := pm.GetCreatedTimeAndReward(ctx, types.GetPermLink(postAuthor, postID))
	if err != nil {
		return types.NewCoinFromInt64(0), err
	}
	return gm.EvaluateConsumption(ctx, coin, numOfConsumptionOnAuthor, created, totalReward)
}

// Handle ReportMsgOrUpvoteMsg
func handleReportOrUpvoteMsg(
	ctx sdk.Context, msg ReportOrUpvoteMsg, pm PostManager, am acc.AccountManager, gm global.GlobalManager) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Username) {
		return ErrReportUserNotFound(msg.Username).Result()
	}

	postKey := types.GetPermLink(msg.Author, msg.PostID)
	if !pm.IsPostExist(ctx, postKey) {
		return ErrReportPostDoesntExist(postKey).Result()
	}

	stake, err := am.GetStake(ctx, msg.Username)
	if err != nil {
		return ErrReportFailed(postKey).TraceCause(err, "").Result()
	}

	sourceAuthor, sourcePostID, err := pm.GetSourcePost(ctx, postKey)
	if err != nil {
		return ErrReportFailed(postKey).TraceCause(err, "").Result()
	}
	if sourceAuthor != types.AccountKey("") && sourcePostID != "" {
		sourcePermLink := types.GetPermLink(sourceAuthor, sourcePostID)
		if err := pm.ReportOrUpvoteToPost(
			ctx, sourcePermLink, msg.Username, stake, msg.IsReport); err != nil {
			return err.Result()
		}
	} else {
		if err := pm.ReportOrUpvoteToPost(
			ctx, postKey, msg.Username, stake, msg.IsReport); err != nil {
			return err.Result()
		}
	}
	return sdk.Result{}
}

func handleUpdatePostMsg(
	ctx sdk.Context, msg UpdatePostMsg, pm PostManager, am acc.AccountManager) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Author) {
		return ErrUpdatePostAuthorNotFound(msg.Author).Result()
	}
	permLink := types.GetPermLink(msg.Author, msg.PostID)
	if !pm.IsPostExist(ctx, permLink) {
		return ErrUpdatePostNotFound(permLink).Result()
	}
	if isDeleted, err := pm.IsDeleted(ctx, permLink); isDeleted || err != nil {
		return ErrUpdatePostIsDeleted(permLink).Result()
	}

	splitRate, err := sdk.NewRatFromDecimal(msg.RedistributionSplitRate)
	if err != nil {
		return err.Result()
	}
	if err := pm.UpdatePost(
		ctx, msg.Author, msg.PostID, msg.Title, msg.Content, msg.Links, splitRate); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleDeletePostMsg(
	ctx sdk.Context, msg DeletePostMsg, pm PostManager, am acc.AccountManager) sdk.Result {
	if !am.IsAccountExist(ctx, msg.Author) {
		return ErrDeletePostAuthorNotFound(msg.Author).Result()
	}
	permLink := types.GetPermLink(msg.Author, msg.PostID)
	if !pm.IsPostExist(ctx, permLink) {
		return ErrDeletePostNotFound(permLink).Result()
	}

	if err := pm.DeletePost(ctx, permLink); err != nil {
		return err.Result()
	}
	return sdk.Result{}
}
