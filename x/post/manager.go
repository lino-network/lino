package post

import (
	"github.com/lino-network/lino/param"
	"github.com/lino-network/lino/recorder"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/post/model"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PostManager struct {
	postStorage model.PostStorage
	paramHolder param.ParamHolder
	recorder    recorder.Recorder
}

// NewPostManager - create a new post manager
func NewPostManager(key sdk.StoreKey, holder param.ParamHolder, recorder recorder.Recorder) PostManager {
	return PostManager{
		postStorage: model.NewPostStorage(key),
		paramHolder: holder,
		recorder:    recorder,
	}
}

// GetRedistributionSplitRate - get post redistribution split rate
func (pm PostManager) GetRedistributionSplitRate(ctx sdk.Context, permlink types.Permlink) (sdk.Rat, sdk.Error) {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return sdk.ZeroRat(), err
	}
	return postMeta.RedistributionSplitRate, nil
}

// GetCreatedTimeAndReward - get post created time and reward for evaluate of content value
func (pm PostManager) GetCreatedTimeAndReward(ctx sdk.Context, permlink types.Permlink) (int64, types.Coin, sdk.Error) {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return 0, types.NewCoinFromInt64(0), err
	}
	return postMeta.CreatedAt, postMeta.TotalReward, nil
}

// DoesPostExist - check if post exist
func (pm PostManager) DoesPostExist(ctx sdk.Context, permlink types.Permlink) bool {
	return pm.postStorage.DoesPostExist(ctx, permlink)
}

// GetSourcePost - return root source post
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
		CreatedAt:               ctx.BlockHeader().Time.Unix(),
		LastUpdatedAt:           ctx.BlockHeader().Time.Unix(),
		LastActivityAt:          ctx.BlockHeader().Time.Unix(),
		AllowReplies:            true, // Default
		IsDeleted:               false,
		RedistributionSplitRate: redistributionSplitRate.Round(types.PrecisionFactor),
	}
	if err := pm.postStorage.SetPostMeta(ctx, permlink, postMeta); err != nil {
		return err
	}
	return nil
}

// UpdatePost - update post title, content and links. Can't update a deleted post
func (pm PostManager) UpdatePost(
	ctx sdk.Context, author types.AccountKey, postID, title, content string,
	links []types.IDToURLMapping) sdk.Error {
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
	// postMeta.RedistributionSplitRate = redistributionSplitRate
	postMeta.LastUpdatedAt = ctx.BlockHeader().Time.Unix()

	if err := pm.postStorage.SetPostInfo(ctx, postInfo); err != nil {
		return err
	}
	if err := pm.postStorage.SetPostMeta(ctx, permlink, postMeta); err != nil {
		return err
	}
	return nil
}

// AddOrUpdateViewToPost - add or update view from the user if view exists
func (pm PostManager) AddOrUpdateViewToPost(
	ctx sdk.Context, permlink types.Permlink, user types.AccountKey) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return err
	}
	view, _ := pm.postStorage.GetPostView(ctx, permlink, user)
	// override previous
	if view == nil {
		view = &model.View{Username: user}
	}
	postMeta.TotalViewCount++
	view.Times++
	view.LastViewAt = ctx.BlockHeader().Time.Unix()
	if err := pm.postStorage.SetPostView(ctx, permlink, view); err != nil {
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
		CreatedAt: ctx.BlockHeader().Time.Unix(),
	}
	if err := pm.postStorage.SetPostComment(ctx, permlink, comment); err != nil {
		return err
	}
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return err
	}
	postMeta.LastActivityAt = ctx.BlockHeader().Time.Unix()
	if err := pm.postStorage.SetPostMeta(ctx, permlink, postMeta); err != nil {
		return err
	}
	return nil
}

// AddDonation - add donation to post donation list
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
	postMeta.LastActivityAt = ctx.BlockHeader().Time.Unix()
	if err := pm.postStorage.SetPostMeta(ctx, permlink, postMeta); err != nil {
		return err
	}
	return nil
}

// DeletePost - delete post by author or content censorship
func (pm PostManager) DeletePost(ctx sdk.Context, permlink types.Permlink) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return err
	}
	postMeta.IsDeleted = true
	postMeta.RedistributionSplitRate = sdk.OneRat()
	postMeta.LastUpdatedAt = ctx.BlockHeader().Time.Unix()
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

// IsDeleted - check if a post is deleted or not
func (pm PostManager) IsDeleted(ctx sdk.Context, permlink types.Permlink) (bool, sdk.Error) {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return false, err
	}
	return postMeta.IsDeleted, nil
}

// UpdateLastActivityAt - update post last activity at
func (pm PostManager) UpdateLastActivityAt(ctx sdk.Context, permlink types.Permlink) sdk.Error {
	postMeta, err := pm.postStorage.GetPostMeta(ctx, permlink)
	if err != nil {
		return err
	}
	postMeta.LastActivityAt = ctx.BlockHeader().Time.Unix()
	if err := pm.postStorage.SetPostMeta(ctx, permlink, postMeta); err != nil {
		return err
	}
	return nil
}

// GetPenaltyScore - get penalty score from report and upvote
func (pm PostManager) GetPenaltyScore(ctx sdk.Context, reputation types.Coin) (sdk.Rat, sdk.Error) {
	if reputation.IsNotNegative() {
		return sdk.ZeroRat(), nil
	}
	reputation = types.NewCoinFromInt64(0).Minus(reputation)
	postParam, err := pm.paramHolder.GetPostParam(ctx)
	if err != nil {
		return sdk.OneRat(), err
	}
	// if max report reputation is zero, any negative reputation should result in max penalty score
	if postParam.MaxReportReputation.IsZero() {
		return sdk.OneRat(), nil
	}
	if reputation.IsGTE(postParam.MaxReportReputation) {
		return sdk.OneRat(), nil
	}
	penaltyScore := reputation.ToRat().Quo(postParam.MaxReportReputation.ToRat())
	if penaltyScore.GT(sdk.OneRat()) {
		return sdk.OneRat(), nil
	}
	return penaltyScore, nil
}
