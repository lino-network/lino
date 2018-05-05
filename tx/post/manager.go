package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/tx/post/model"
	"github.com/lino-network/lino/types"
)

const rewardEvent = 0x1

// post is the proxy for all storage structs defined above
type PostManager struct {
	postStorage model.PostStorage `json:"post_storage"`
	paramHolder param.ParamHolder `json:"param_holder"`
}

// create NewPostManager
func NewPostManager(key sdk.StoreKey, holder param.ParamHolder) PostManager {
	return PostManager{
		postStorage: model.NewPostStorage(key),
		paramHolder: holder,
	}
}

func (pm PostManager) GetRedistributionSplitRate(ctx sdk.Context, permLink types.PermLink) (sdk.Rat, sdk.Error) {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permLink)
	if err != nil {
		return sdk.ZeroRat, ErrGetRedistributionSplitRate(permLink).TraceCause(err, "")
	}
	return postMeta.RedistributionSplitRate, nil
}

func (pm PostManager) GetCreatedTimeAndReward(ctx sdk.Context, permLink types.PermLink) (int64, types.Coin, sdk.Error) {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permLink)
	if err != nil {
		return 0, types.NewCoin(0), ErrGetCreatedTime(permLink).TraceCause(err, "")
	}
	return postMeta.Created, postMeta.TotalReward, nil
}

// check if post exist
func (pm PostManager) IsPostExist(ctx sdk.Context, permLink types.PermLink) bool {
	if postInfo, _ := pm.postStorage.GetPostInfo(ctx, permLink); postInfo == nil {
		return false
	}
	return true
}

// return root source post
func (pm PostManager) GetSourcePost(
	ctx sdk.Context, permLink types.PermLink) (types.AccountKey, string, sdk.Error) {
	postInfo, err := pm.postStorage.GetPostInfo(ctx, permLink)
	if err != nil {
		return types.AccountKey(""), "", ErrGetRootSourcePost(permLink).TraceCause(err, "")
	}

	// check source post's source, that's the root
	if postInfo.SourceAuthor == types.AccountKey("") || postInfo.SourcePostID == "" {
		return types.AccountKey(""), "", nil
	} else {
		return postInfo.SourceAuthor, postInfo.SourcePostID, nil
	}
}

func (pm PostManager) setRootSourcePost(ctx sdk.Context, postInfo *model.PostInfo) sdk.Error {
	if postInfo.SourceAuthor == types.AccountKey("") || postInfo.SourcePostID == "" {
		return nil
	}
	permLink := types.GetPermLink(postInfo.Author, postInfo.PostID)
	rootAuthor, rootPostID, err :=
		pm.GetSourcePost(ctx, types.GetPermLink(postInfo.SourceAuthor, postInfo.SourcePostID))
	if err != nil {
		return ErrSetRootSourcePost(permLink).TraceCause(err, "")
	}
	if rootAuthor != types.AccountKey("") && rootPostID != "" {
		postInfo.SourceAuthor = rootAuthor
		postInfo.SourcePostID = rootPostID
	}
	return nil
}

// create the post
func (pm PostManager) CreatePost(ctx sdk.Context, postCreateParams *PostCreateParams) sdk.Error {
	postInfo := &model.PostInfo{
		PostID:       postCreateParams.PostID,
		Title:        postCreateParams.Title,
		Content:      postCreateParams.Content,
		Author:       postCreateParams.Author,
		ParentAuthor: postCreateParams.ParentAuthor,
		ParentPostID: postCreateParams.ParentPostID,
		SourceAuthor: postCreateParams.SourceAuthor,
		SourcePostID: postCreateParams.SourcePostID,
		Links:        postCreateParams.Links,
	}
	permLink := types.GetPermLink(postInfo.Author, postInfo.PostID)
	if pm.IsPostExist(ctx, permLink) {
		return ErrPostExist(permLink)
	}
	if err := pm.setRootSourcePost(ctx, postInfo); err != nil {
		return ErrCreatePostSourceInvalid(permLink)
	}
	if err := pm.postStorage.SetPostInfo(ctx, postInfo); err != nil {
		return ErrCreatePost(permLink).TraceCause(err, "")
	}
	splitRate, err := sdk.NewRatFromDecimal(postCreateParams.RedistributionSplitRate)
	if err != nil {
		return ErrCreatePost(permLink).TraceCause(err, "")
	}
	postMeta := &model.PostMeta{
		Created:                 ctx.BlockHeader().Time,
		LastUpdate:              ctx.BlockHeader().Time,
		LastActivity:            ctx.BlockHeader().Time,
		AllowReplies:            true, // Default
		RedistributionSplitRate: splitRate,
	}
	if err := pm.postStorage.SetPostMeta(ctx, permLink, postMeta); err != nil {
		return ErrCreatePost(permLink).TraceCause(err, "")
	}
	return nil
}

// add or update like from the user if like exists
func (pm PostManager) AddOrUpdateLikeToPost(
	ctx sdk.Context, permLink types.PermLink, user types.AccountKey, weight int64) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permLink)
	if err != nil {
		return ErrAddOrUpdateLikeToPost(permLink).TraceCause(err, "")
	}
	like, _ := pm.postStorage.GetPostLike(ctx, permLink, user)
	// Revoke privous
	if like != nil {
		if like.Weight > 0 {
			postMeta.TotalLikeWeight -= like.Weight
		}
		if like.Weight < 0 {
			postMeta.TotalDislikeWeight += like.Weight
		}
		like.Weight = weight
	} else {
		postMeta.TotalLikeCount += 1
		like = &model.Like{Username: user, Weight: weight, Created: ctx.BlockHeader().Time}
	}
	if like.Weight > 0 {
		postMeta.TotalLikeWeight += like.Weight
	}
	if like.Weight < 0 {
		postMeta.TotalDislikeWeight -= like.Weight
	}
	postMeta.LastActivity = ctx.BlockHeader().Time
	if err := pm.postStorage.SetPostLike(ctx, permLink, like); err != nil {
		return ErrAddOrUpdateLikeToPost(permLink).TraceCause(err, "")
	}
	if err := pm.postStorage.SetPostMeta(ctx, permLink, postMeta); err != nil {
		return ErrAddOrUpdateLikeToPost(permLink).TraceCause(err, "")
	}
	return nil
}

// add or update report or upvote from the user if exist
func (pm PostManager) ReportOrUpvoteToPost(
	ctx sdk.Context, permLink types.PermLink, user types.AccountKey, stake types.Coin, isReport bool, isRevoke bool) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permLink)
	if err != nil {
		return ErrAddOrUpdateReportOrUpvoteToPost(permLink).TraceCause(err, "")
	}
	postMeta.LastActivity = ctx.BlockHeader().Time

	reportOrUpvote, _ := pm.postStorage.GetPostReportOrUpvote(ctx, permLink, user)
	// Revoke privous
	if reportOrUpvote != nil {
		if reportOrUpvote.IsReport {
			postMeta.TotalReportStake = postMeta.TotalReportStake.Minus(reportOrUpvote.Stake)
		} else {
			postMeta.TotalUpvoteStake = postMeta.TotalUpvoteStake.Minus(reportOrUpvote.Stake)
		}
		reportOrUpvote.Stake = stake
		if isRevoke {
			if err := pm.postStorage.SetPostMeta(ctx, permLink, postMeta); err != nil {
				return ErrAddOrUpdateReportOrUpvoteToPost(permLink).TraceCause(err, "")
			}
			return pm.postStorage.RemovePostReportOrUpvote(ctx, permLink, user)
		}
	} else {
		if isRevoke {
			return ErrRevokeReportOrUpvoteToPost(permLink)
		}
		reportOrUpvote =
			&model.ReportOrUpvote{Username: user, Stake: stake, Created: ctx.BlockHeader().Time}
	}
	if isReport {
		postMeta.TotalReportStake = postMeta.TotalReportStake.Plus(reportOrUpvote.Stake)
		reportOrUpvote.IsReport = true
	} else {
		postMeta.TotalUpvoteStake = postMeta.TotalUpvoteStake.Plus(reportOrUpvote.Stake)
		reportOrUpvote.IsReport = false
	}
	if err := pm.postStorage.SetPostReportOrUpvote(ctx, permLink, reportOrUpvote); err != nil {
		return ErrAddOrUpdateReportOrUpvoteToPost(permLink).TraceCause(err, "")
	}
	if err := pm.postStorage.SetPostMeta(ctx, permLink, postMeta); err != nil {
		return ErrAddOrUpdateReportOrUpvoteToPost(permLink).TraceCause(err, "")
	}
	return nil
}

// add comment to post comment list
func (pm PostManager) AddComment(
	ctx sdk.Context, permLink types.PermLink, commentUser types.AccountKey, commentPostID string) sdk.Error {
	comment := &model.Comment{Author: commentUser, PostID: commentPostID, Created: ctx.BlockHeader().Time}
	return pm.postStorage.SetPostComment(ctx, permLink, comment)
}

// add donation to post donation list
func (pm PostManager) AddDonation(
	ctx sdk.Context, permLink types.PermLink, donator types.AccountKey, amount types.Coin) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permLink)
	if err != nil {
		return ErrAddDonation(permLink).TraceCause(err, "")
	}
	donation := model.Donation{
		Amount:  amount,
		Created: ctx.BlockHeader().Time,
	}
	donations, _ := pm.postStorage.GetPostDonations(ctx, permLink, donator)
	if donations == nil {
		donations = &model.Donations{Username: donator, DonationList: []model.Donation{}}
	}
	donations.DonationList = append(donations.DonationList, donation)
	if err := pm.postStorage.SetPostDonations(ctx, permLink, donations); err != nil {
		return ErrAddDonation(permLink).TraceCause(err, "")
	}
	postMeta.TotalReward = postMeta.TotalReward.Plus(donation.Amount)
	postMeta.TotalDonateCount = postMeta.TotalDonateCount + 1
	if err := pm.postStorage.SetPostMeta(ctx, permLink, postMeta); err != nil {
		return ErrAddDonation(permLink).TraceCause(err, "")
	}
	return nil
}

// add view to post view list
func (pm PostManager) AddView(ctx sdk.Context, permLink types.PermLink, user types.AccountKey) sdk.Error {
	view, _ := pm.postStorage.GetPostView(ctx, permLink, user)
	if view != nil {
		view.Times += 1
	} else {
		view = &model.View{Username: user, Created: ctx.BlockHeader().Time, Times: 1}
	}

	return pm.postStorage.SetPostView(ctx, permLink, view)
}

// get penalty score from report and upvote
func (pm PostManager) GetPenaltyScore(ctx sdk.Context, permLink types.PermLink) (sdk.Rat, sdk.Error) {
	author, postID, err := pm.GetSourcePost(ctx, permLink)
	if err != nil {
		return sdk.ZeroRat, ErrGetPenaltyScore(permLink).TraceCause(err, "")
	}
	if author != types.AccountKey("") && postID != "" {
		paneltyScore, err := pm.GetPenaltyScore(ctx, types.GetPermLink(author, postID))
		if err != nil {
			return sdk.ZeroRat, err
		}
		return paneltyScore, nil
	}
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permLink)
	if err != nil {
		return sdk.ZeroRat, ErrGetPenaltyScore(permLink).TraceCause(err, "")
	}
	if postMeta.TotalReportStake.IsZero() {
		return sdk.ZeroRat, nil
	}
	if postMeta.TotalUpvoteStake.IsZero() {
		return sdk.OneRat, nil
	}
	penaltyScore := postMeta.TotalReportStake.ToRat().Quo(postMeta.TotalUpvoteStake.ToRat())
	if penaltyScore.LT(sdk.ZeroRat) {
		return sdk.ZeroRat, nil
	}
	if penaltyScore.GT(sdk.OneRat) {
		return sdk.OneRat, nil
	}
	return postMeta.TotalReportStake.ToRat().Quo(postMeta.TotalUpvoteStake.ToRat()), nil
}
