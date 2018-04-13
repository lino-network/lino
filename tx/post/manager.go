package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/tx/post/model"
	"github.com/lino-network/lino/types"
)

// post is the proxy for all storage structs defined above
type PostManager struct {
	postStorage *model.PostStorage `json:"post_storage"`
}

// create NewPostManager
func NewPostManager(key sdk.StoreKey) *PostManager {
	return &PostManager{
		postStorage: model.NewPostStorage(key),
	}
}

func (pm *PostManager) GetRedistributionSplitRate(ctx sdk.Context, postKey types.PostKey) (sdk.Rat, sdk.Error) {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, postKey)
	if err != nil {
		return sdk.ZeroRat, ErrGetRedistributionSplitRate(postKey).TraceCause(err, "")
	}
	return postMeta.RedistributionSplitRate, nil
}

// check if post exist
func (pm *PostManager) IsPostExist(ctx sdk.Context, postKey types.PostKey) bool {
	if postInfo, _ := pm.postStorage.GetPostInfo(ctx, postKey); postInfo == nil {
		return false
	}
	return true
}

// return root source post
func (pm *PostManager) GetSourcePost(
	ctx sdk.Context, postKey types.PostKey) (types.AccountKey, string, sdk.Error) {
	postInfo, err := pm.postStorage.GetPostInfo(ctx, postKey)
	if err != nil {
		return types.AccountKey(""), "", ErrGetRootSourcePost(postKey).TraceCause(err, "")
	}

	// check source post's source, that's the root
	if postInfo.SourceAuthor == types.AccountKey("") || postInfo.SourcePostID == "" {
		return types.AccountKey(""), "", nil
	} else {
		return postInfo.SourceAuthor, postInfo.SourcePostID, nil
	}
}

func (pm *PostManager) setRootSourcePost(ctx sdk.Context, postInfo *model.PostInfo) sdk.Error {
	if postInfo.SourceAuthor == types.AccountKey("") || postInfo.SourcePostID == "" {
		return nil
	}
	postKey := types.GetPostKey(postInfo.Author, postInfo.PostID)
	rootAuthor, rootPostID, err :=
		pm.GetSourcePost(ctx, types.GetPostKey(postInfo.SourceAuthor, postInfo.SourcePostID))
	if err != nil {
		return ErrSetRootSourcePost(postKey).TraceCause(err, "")
	}
	if rootAuthor != types.AccountKey("") && rootPostID != "" {
		postInfo.SourceAuthor = rootAuthor
		postInfo.SourcePostID = rootPostID
	}
	return nil
}

// create the post
func (pm *PostManager) CreatePost(ctx sdk.Context, postCreateParams *PostCreateParams) sdk.Error {
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
	postKey := types.GetPostKey(postInfo.Author, postInfo.PostID)
	if pm.IsPostExist(ctx, postKey) {
		return ErrPostExist(postKey)
	}
	if err := pm.setRootSourcePost(ctx, postInfo); err != nil {
		return ErrCreatePost(postKey).TraceCause(err, "")
	}
	if err := pm.postStorage.SetPostInfo(ctx, postInfo); err != nil {
		return ErrCreatePost(postKey).TraceCause(err, "")
	}
	postMeta := &model.PostMeta{
		Created:                 ctx.BlockHeader().Time,
		LastUpdate:              ctx.BlockHeader().Time,
		LastActivity:            ctx.BlockHeader().Time,
		AllowReplies:            true, // Default
		RedistributionSplitRate: postCreateParams.RedistributionSplitRate,
	}
	if err := pm.postStorage.SetPostMeta(ctx, postKey, postMeta); err != nil {
		return ErrCreatePost(postKey).TraceCause(err, "")
	}
	return nil
}

// add or update like from the user if like exists
func (pm *PostManager) AddOrUpdateLikeToPost(
	ctx sdk.Context, postKey types.PostKey, user types.AccountKey, weight int64) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, postKey)
	if err != nil {
		return ErrAddOrUpdateLikeToPost(postKey).TraceCause(err, "")
	}
	like, _ := pm.postStorage.GetPostLike(ctx, postKey, user)
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
	if err := pm.postStorage.SetPostLike(ctx, postKey, like); err != nil {
		return ErrAddOrUpdateLikeToPost(postKey).TraceCause(err, "")
	}
	if err := pm.postStorage.SetPostMeta(ctx, postKey, postMeta); err != nil {
		return ErrAddOrUpdateLikeToPost(postKey).TraceCause(err, "")
	}
	return nil
}

// add or update report or upvote from the user if exist
func (pm *PostManager) AddOrUpdateReportOrUpvoteToPost(
	ctx sdk.Context, postKey types.PostKey, user types.AccountKey, stake types.Coin, isReport bool) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, postKey)
	if err != nil {
		return ErrAddOrUpdateReportOrUpvoteToPost(postKey).TraceCause(err, "")
	}
	reportOrUpvote, _ := pm.postStorage.GetPostReportOrUpvote(ctx, postKey, user)
	// Revoke privous
	if reportOrUpvote != nil {
		if reportOrUpvote.IsReport {
			postMeta.TotalReportStake = postMeta.TotalReportStake.Minus(reportOrUpvote.Stake)
		} else {
			postMeta.TotalUpvoteStake = postMeta.TotalUpvoteStake.Minus(reportOrUpvote.Stake)
		}
		reportOrUpvote.Stake = stake
	} else {
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
	if err := pm.postStorage.SetPostReportOrUpvote(ctx, postKey, reportOrUpvote); err != nil {
		return ErrAddOrUpdateReportOrUpvoteToPost(postKey).TraceCause(err, "")
	}
	if err := pm.postStorage.SetPostMeta(ctx, postKey, postMeta); err != nil {
		return ErrAddOrUpdateReportOrUpvoteToPost(postKey).TraceCause(err, "")
	}
	return nil
}

// add comment to post comment list
func (pm *PostManager) AddComment(
	ctx sdk.Context, postKey types.PostKey, commentUser types.AccountKey, commentPostID string) sdk.Error {
	comment := &model.Comment{Author: commentUser, PostID: commentPostID, Created: ctx.BlockHeader().Time}
	return pm.postStorage.SetPostComment(ctx, postKey, comment)
}

// add donation to post donation list
func (pm *PostManager) AddDonation(
	ctx sdk.Context, postKey types.PostKey, donator types.AccountKey, amount types.Coin) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, postKey)
	if err != nil {
		return ErrAddDonation(postKey).TraceCause(err, "")
	}
	donation := model.Donation{
		Amount:  amount,
		Created: ctx.BlockHeader().Time,
	}
	donations, _ := pm.postStorage.GetPostDonations(ctx, postKey, donator)
	if donations == nil {
		donations = &model.Donations{Username: donator, DonationList: []model.Donation{}}
	}
	donations.DonationList = append(donations.DonationList, donation)
	if err := pm.postStorage.SetPostDonations(ctx, postKey, donations); err != nil {
		return ErrAddDonation(postKey).TraceCause(err, "")
	}
	postMeta.TotalReward = postMeta.TotalReward.Plus(donation.Amount)
	postMeta.TotalDonateCount = postMeta.TotalDonateCount + 1
	if err := pm.postStorage.SetPostMeta(ctx, postKey, postMeta); err != nil {
		return ErrAddDonation(postKey).TraceCause(err, "")
	}
	return nil
}

// add view to post view list
func (pm *PostManager) AddView(ctx sdk.Context, postKey types.PostKey, user types.AccountKey) sdk.Error {
	view, _ := pm.postStorage.GetPostView(ctx, postKey, user)
	if view != nil {
		view.Times += 1
	} else {
		view = &model.View{Username: user, Created: ctx.BlockHeader().Time, Times: 1}
	}

	return pm.postStorage.SetPostView(ctx, postKey, view)
}

// update last activity
func (pm *PostManager) UpdateLastActivity(ctx sdk.Context, postKey types.PostKey) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, postKey)
	if err != nil {
		return ErrUpdateLastActivity(postKey).TraceCause(err, "")
	}
	postMeta.LastActivity = ctx.BlockHeader().Time
	if err := pm.postStorage.SetPostMeta(ctx, postKey, postMeta); err != nil {
		return ErrUpdateLastActivity(postKey).TraceCause(err, "")
	}
	return nil
}
