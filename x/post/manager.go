package post

import (
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/post/model"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

func (pm PostManager) GetRedistributionSplitRate(ctx sdk.Context, permlink types.Permlink) (sdk.Rat, sdk.Error) {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return sdk.ZeroRat(), err
	}
	return postMeta.RedistributionSplitRate, nil
}

func (pm PostManager) GetCreatedTimeAndReward(ctx sdk.Context, permlink types.Permlink) (int64, types.Coin, sdk.Error) {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return 0, types.NewCoinFromInt64(0), err
	}
	return postMeta.CreatedAt, postMeta.TotalReward, nil
}

// check if post exist
func (pm PostManager) DoesPostExist(ctx sdk.Context, permlink types.Permlink) bool {
	return pm.postStorage.DoesPostExist(ctx, permlink)
}

// return root source post
func (pm PostManager) GetSourcePost(
	ctx sdk.Context, permlink types.Permlink) (types.AccountKey, string, sdk.Error) {
	postInfo, err := pm.postStorage.GetPostInfo(ctx, permlink)
	if err != nil {
		return types.AccountKey(""), "", err
	}

	// check source post's source, that's the root
	if postInfo.SourceAuthor == types.AccountKey("") || postInfo.SourcePostID == "" {
		return types.AccountKey(""), "", nil
	}

	return postInfo.SourceAuthor, postInfo.SourcePostID, nil
}

func (pm PostManager) setRootSourcePost(ctx sdk.Context, postInfo *model.PostInfo) sdk.Error {
	if postInfo.SourceAuthor == types.AccountKey("") || postInfo.SourcePostID == "" {
		return nil
	}
	permlink := types.GetPermlink(postInfo.Author, postInfo.PostID)
	rootAuthor, rootPostID, err :=
		pm.GetSourcePost(ctx, types.GetPermlink(postInfo.SourceAuthor, postInfo.SourcePostID))
	if err != nil {
		return ErrGetSourcePost(permlink)
	}
	if rootAuthor != types.AccountKey("") && rootPostID != "" {
		postInfo.SourceAuthor = rootAuthor
		postInfo.SourcePostID = rootPostID
	}
	return nil
}

// create the post
func (pm PostManager) CreatePost(
	ctx sdk.Context, author types.AccountKey, postID string,
	sourceAuthor types.AccountKey, sourcePostID string,
	parentAuthor types.AccountKey, parentPostID string,
	content string, title string, redistributionSplitRate sdk.Rat,
	links []types.IDToURLMapping) sdk.Error {
	postInfo := &model.PostInfo{
		PostID:       postID,
		Title:        title,
		Content:      content,
		Author:       author,
		ParentAuthor: parentAuthor,
		ParentPostID: parentPostID,
		SourceAuthor: sourceAuthor,
		SourcePostID: sourcePostID,
		Links:        links,
	}
	permlink := types.GetPermlink(postInfo.Author, postInfo.PostID)
	if pm.DoesPostExist(ctx, permlink) {
		return ErrPostAlreadyExist(permlink)
	}
	if err := pm.setRootSourcePost(ctx, postInfo); err != nil {
		return ErrCreatePostSourceInvalid(permlink)
	}
	if err := pm.postStorage.SetPostInfo(ctx, postInfo); err != nil {
		return err
	}
	postMeta := &model.PostMeta{
		CreatedAt:               ctx.BlockHeader().Time,
		LastUpdatedAt:           ctx.BlockHeader().Time,
		LastActivityAt:          ctx.BlockHeader().Time,
		AllowReplies:            true, // Default
		IsDeleted:               false,
		RedistributionSplitRate: redistributionSplitRate.Round(types.PrecisionFactor),
	}
	if err := pm.postStorage.SetPostMeta(ctx, permlink, postMeta); err != nil {
		return err
	}
	return nil
}

func (pm PostManager) UpdatePost(
	ctx sdk.Context, author types.AccountKey, postID, title, content string,
	links []types.IDToURLMapping, redistributionSplitRate sdk.Rat) sdk.Error {
	permlink := types.GetPermlink(author, postID)
	postInfo, err := pm.postStorage.GetPostInfo(ctx, permlink)
	if err != nil {
		return err
	}
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return err
	}

	postInfo.Title = title
	postInfo.Content = content
	postInfo.Links = links
	postMeta.RedistributionSplitRate = redistributionSplitRate

	if err := pm.postStorage.SetPostInfo(ctx, postInfo); err != nil {
		return err
	}
	if err := pm.postStorage.SetPostMeta(ctx, permlink, postMeta); err != nil {
		return err
	}
	return nil
}

// add or update view from the user if view exists
func (pm PostManager) AddOrUpdateViewToPost(
	ctx sdk.Context, permlink types.Permlink, user types.AccountKey) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return err
	}
	view, _ := pm.postStorage.GetPostView(ctx, permlink, user)
	// Revoke previous
	if view == nil {
		view = &model.View{Username: user}
	}
	postMeta.TotalViewCount += 1
	view.Times += 1
	view.LastViewAt = ctx.BlockHeader().Time
	if err := pm.postStorage.SetPostView(ctx, permlink, view); err != nil {
		return err
	}
	if err := pm.postStorage.SetPostMeta(ctx, permlink, postMeta); err != nil {
		return err
	}
	return nil
}

// add or update view from the user if view exists
func (pm PostManager) GetReportOrUpvoteInterval(ctx sdk.Context) (int64, sdk.Error) {
	postParam, err := pm.paramHolder.GetPostParam(ctx)
	if err != nil {
		return 0, err
	}
	return postParam.ReportOrUpvoteInterval, nil
}

// add or update report or upvote from the user if exist
func (pm PostManager) ReportOrUpvoteToPost(
	ctx sdk.Context, permlink types.Permlink, user types.AccountKey,
	stake types.Coin, isReport bool) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return err
	}
	postMeta.LastActivityAt = ctx.BlockHeader().Time

	reportOrUpvote, _ := pm.postStorage.GetPostReportOrUpvote(ctx, permlink, user)

	if reportOrUpvote != nil {
		if reportOrUpvote.IsReport {
			postMeta.TotalReportStake = postMeta.TotalReportStake.Minus(reportOrUpvote.Stake)
		} else {
			postMeta.TotalUpvoteStake = postMeta.TotalUpvoteStake.Minus(reportOrUpvote.Stake)
		}
	}
	reportOrUpvote =
		&model.ReportOrUpvote{Username: user, Stake: stake, CreatedAt: ctx.BlockHeader().Time}
	if isReport {
		postMeta.TotalReportStake = postMeta.TotalReportStake.Plus(reportOrUpvote.Stake)
		reportOrUpvote.IsReport = true
	} else {
		postMeta.TotalUpvoteStake = postMeta.TotalUpvoteStake.Plus(reportOrUpvote.Stake)
		reportOrUpvote.IsReport = false
	}
	if err := pm.postStorage.SetPostReportOrUpvote(ctx, permlink, reportOrUpvote); err != nil {
		return err
	}
	if err := pm.postStorage.SetPostMeta(ctx, permlink, postMeta); err != nil {
		return err
	}
	return nil
}

// add comment to post comment list
func (pm PostManager) AddComment(
	ctx sdk.Context, permlink types.Permlink, commentAuthor types.AccountKey, commentPostID string) sdk.Error {
	comment := &model.Comment{
		Author:    commentAuthor,
		PostID:    commentPostID,
		CreatedAt: ctx.BlockHeader().Time,
	}
	if err := pm.postStorage.SetPostComment(ctx, permlink, comment); err != nil {
		return err
	}

	return nil
}

// add donation to post donation list
func (pm PostManager) AddDonation(
	ctx sdk.Context, permlink types.Permlink, donator types.AccountKey,
	amount types.Coin, donationType types.DonationType) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return err
	}
	donations, _ := pm.postStorage.GetPostDonations(ctx, permlink, donator)
	if donations == nil {
		donations = &model.Donations{Username: donator, Amount: types.NewCoinFromInt64(0), Times: 0}
	}
	donations.Amount = donations.Amount.Plus(amount)
	donations.Times = donations.Times + 1
	if err := pm.postStorage.SetPostDonations(ctx, permlink, donations); err != nil {
		return err
	}
	postMeta.TotalReward = postMeta.TotalReward.Plus(amount)
	postMeta.TotalDonateCount = postMeta.TotalDonateCount + 1
	if err := pm.postStorage.SetPostMeta(ctx, permlink, postMeta); err != nil {
		return err
	}
	return nil
}

// DeletePost triggered by censorship proposal
func (pm PostManager) DeletePost(ctx sdk.Context, permlink types.Permlink) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return err
	}
	postMeta.IsDeleted = true
	postMeta.RedistributionSplitRate = sdk.OneRat()
	if err := pm.postStorage.SetPostMeta(ctx, permlink, postMeta); err != nil {
		return err
	}
	postInfo, err := pm.postStorage.GetPostInfo(ctx, permlink)
	if err != nil {
		return err
	}
	postInfo.Title = ""
	postInfo.Content = ""
	postInfo.Links = nil

	if err := pm.postStorage.SetPostInfo(ctx, postInfo); err != nil {
		return err
	}
	return nil
}

func (pm PostManager) IsDeleted(ctx sdk.Context, permlink types.Permlink) (bool, sdk.Error) {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return false, err
	}
	return postMeta.IsDeleted, nil
}

// get penalty score from report and upvote
func (pm PostManager) GetPenaltyScore(ctx sdk.Context, permlink types.Permlink) (sdk.Rat, sdk.Error) {
	sourceAuthor, sourcePostID, err := pm.GetSourcePost(ctx, permlink)
	if err != nil {
		return sdk.ZeroRat(), err
	}
	if sourceAuthor != types.AccountKey("") && sourcePostID != "" {
		paneltyScore, err := pm.GetPenaltyScore(ctx, types.GetPermlink(sourceAuthor, sourcePostID))
		if err != nil {
			return sdk.ZeroRat(), err
		}
		return paneltyScore, nil
	}
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return sdk.ZeroRat(), err
	}
	if postMeta.TotalReportStake.IsZero() {
		return sdk.ZeroRat(), nil
	}
	if postMeta.TotalUpvoteStake.IsZero() {
		return sdk.OneRat(), nil
	}
	penaltyScore := postMeta.TotalReportStake.ToRat().Quo(postMeta.TotalUpvoteStake.ToRat()).Round(types.PrecisionFactor)
	if penaltyScore.GT(sdk.OneRat()) {
		return sdk.OneRat(), nil
	}
	return penaltyScore, nil
}
