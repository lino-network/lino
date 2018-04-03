package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

// post is the proxy for all storage structs defined above
type PostProxy struct {
	author        acc.AccountKey `json:"author"`
	postID        string         `json:"post_ID"`
	postKey       PostKey        `json:"post_key"`
	postManager   *PostManager   `json:"postManager"`
	writePostInfo bool           `json:"write_post_info"`
	writePostMeta bool           `json:"write_post_meta"`
	postInfo      *PostInfo      `json:"post_info"`
	postMeta      *PostMeta      `json:"post_meta"`
}

// create NewPostProxy
func NewPostProxy(author acc.AccountKey, postID string, postManager *PostManager) *PostProxy {
	return &PostProxy{
		author:      author,
		postID:      postID,
		postManager: postManager,
		postKey:     GetPostKey(author, postID),
	}
}

func (p *PostProxy) GetAuthor() acc.AccountKey {
	return p.author
}

func (p *PostProxy) GetPostID() string {
	return p.postID
}

func (p *PostProxy) GetPostKey() PostKey {
	return p.postKey
}

func (p *PostProxy) GetRedistributionSplitRate(ctx sdk.Context) (sdk.Rat, sdk.Error) {
	if err := p.checkPostInfo(ctx); err != nil {
		return sdk.Rat{}, err
	}
	return p.postInfo.RedistributionSplitRate, nil
}

// check if post exist
func (p *PostProxy) IsPostExist(ctx sdk.Context) bool {
	if err := p.checkPostInfo(ctx); err != nil {
		return false
	}
	return true
}

// return source post proxy
func (p *PostProxy) GetRootSourcePost(ctx sdk.Context) (*PostProxy, sdk.Error) {
	if err := p.checkPostInfo(ctx); err != nil {
		return nil, err
	}
	if len(p.postInfo.SourceAuthor) == 0 && len(p.postInfo.SourcePostID) == 0 {
		return nil, nil
	}
	sourcePost := NewPostProxy(p.postInfo.SourceAuthor, p.postInfo.SourcePostID, p.postManager)
	rootPost, err := sourcePost.GetRootSourcePost(ctx)
	if err != nil {
		return nil, err
	}
	if rootPost == nil {
		return sourcePost, nil
	} else {
		return rootPost, nil
	}
}

func (p *PostProxy) setRootSourcePost(ctx sdk.Context) sdk.Error {
	source, err := p.GetRootSourcePost(ctx)
	if err != nil {
		return err
	}
	if source != nil {
		p.writePostInfo = true
		p.postInfo.SourceAuthor = source.GetAuthor()
		p.postInfo.SourcePostID = source.GetPostID()
	}
	return nil
}

// create the post
func (p *PostProxy) CreatePost(ctx sdk.Context, postInfo *PostInfo) sdk.Error {
	if p.IsPostExist(ctx) {
		return ErrPostExist()
	}
	p.writePostInfo = true
	p.postInfo = postInfo
	if err := p.setRootSourcePost(ctx); err != nil {
		return err
	}
	p.writePostMeta = true
	p.postMeta = &PostMeta{
		Created:      types.Height(ctx.BlockHeight()),
		LastUpdate:   types.Height(ctx.BlockHeight()),
		LastActivity: types.Height(ctx.BlockHeight()),
		AllowReplies: true, // Default
	}
	return nil
}

// add or update like from the user if like exists
func (p *PostProxy) AddOrUpdateLikeToPost(ctx sdk.Context, likeToAddOrUpdate Like) sdk.Error {
	if err := p.checkPostMeta(ctx); err != nil {
		return err
	}
	like, _ := p.postManager.GetPostLike(ctx, p.GetPostKey(), likeToAddOrUpdate.Username)
	if like != nil {
		p.postMeta.TotalLikeWeight -= like.Weight
	} else {
		p.postMeta.TotalLikeCount += 1
	}
	p.writePostMeta = true
	like = &likeToAddOrUpdate
	p.postMeta.TotalLikeWeight += like.Weight
	return p.postManager.SetPostLike(ctx, p.GetPostKey(), like)
}

// add comment to post comment list
func (p *PostProxy) AddComment(ctx sdk.Context, comment Comment) sdk.Error {
	if err := p.UpdateLastActivity(ctx); err != nil {
		return err
	}
	return p.postManager.SetPostComment(ctx, p.GetPostKey(), &comment)
}

// add donation to post donation list
func (p *PostProxy) AddDonation(ctx sdk.Context, donator acc.AccountKey, amount types.Coin) sdk.Error {
	if err := p.UpdateLastActivity(ctx); err != nil {
		return err
	}
	donation := Donation{
		Amount:  amount,
		Created: types.Height(ctx.BlockHeight()),
	}
	donations, _ := p.postManager.GetPostDonations(ctx, p.GetPostKey(), donator)
	if donations == nil {
		donations = &Donations{Username: donator, DonationList: []Donation{}}
	}
	donations.DonationList = append(donations.DonationList, donation)
	if err := p.postManager.SetPostDonations(ctx, p.GetPostKey(), donations); err != nil {
		return err
	}
	p.writePostMeta = true
	p.postMeta.TotalReward = p.postMeta.TotalReward.Plus(donation.Amount)
	p.postMeta.TotalDonateCount = p.postMeta.TotalDonateCount + 1
	return nil
}

// add view to post view list
func (p *PostProxy) AddView(ctx sdk.Context, view View) sdk.Error {
	if err := p.UpdateLastActivity(ctx); err != nil {
		return err
	}
	return p.postManager.SetPostView(ctx, p.GetPostKey(), &view)
}

// update last activity
func (p *PostProxy) UpdateLastActivity(ctx sdk.Context) sdk.Error {
	if err := p.checkPostMeta(ctx); err != nil {
		return err
	}
	p.postMeta.LastActivity = types.Height(ctx.BlockHeight())
	p.writePostMeta = true
	return nil
}

// check if PostInfo exists
func (p *PostProxy) checkPostInfo(ctx sdk.Context) (err sdk.Error) {
	if p.postInfo == nil {
		p.postInfo, err = p.postManager.GetPostInfo(ctx, p.GetPostKey())
	}
	return err
}

// check if PostMeta exists
func (p *PostProxy) checkPostMeta(ctx sdk.Context) (err sdk.Error) {
	if p.postMeta == nil {
		p.postMeta, err = p.postManager.GetPostMeta(ctx, p.GetPostKey())
	}
	return err
}

// apply all changes to storage
func (p *PostProxy) Apply(ctx sdk.Context) sdk.Error {
	if p.writePostInfo {
		if err := p.postManager.SetPostInfo(ctx, p.postInfo); err != nil {
			return err
		}
	}
	if p.writePostMeta {
		if err := p.postManager.SetPostMeta(ctx, p.GetPostKey(), p.postMeta); err != nil {
			return err
		}
	}
	p.clear()

	return nil
}

// clear current post proxy
func (p *PostProxy) clear() {
	p.writePostInfo = false
	p.writePostMeta = false
	p.postInfo = nil
	p.postMeta = nil
}
