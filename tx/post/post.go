package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

// PostKey key format in KVStore
type PostKey string

// Identifier is used to map the url in post
type Identifier string

// URL used to link resources like vedio, text or photo
type URL string

// PostIntreface needs to be implemented by all post related struct
// this is needed in post manager
type PostInterface interface {
	AssertPostInterface()
}

var _ PostInterface = PostInfo{}
var _ PostInterface = PostMeta{}
var _ PostInterface = Like{}
var _ PostInterface = Report{}
var _ PostInterface = Comment{}
var _ PostInterface = View{}
var _ PostInterface = Donation{}

func (_ PostInfo) AssertPostInterface() {}
func (_ PostMeta) AssertPostInterface() {}
func (_ Like) AssertPostInterface()     {}
func (_ Report) AssertPostInterface()   {}
func (_ Comment) AssertPostInterface()  {}
func (_ View) AssertPostInterface()     {}
func (_ Donation) AssertPostInterface() {}

// PostInfo can also use to present comment(with parent) or repost(with source)
type PostInfo struct {
	PostID       string           `json:"post_id"`
	Title        string           `json:"title"`
	Content      string           `json:"content"`
	Author       acc.AccountKey   `json:"author"`
	ParentAuthor acc.AccountKey   `json:"parent_author"`
	ParentPostID string           `json:"parent_postID"`
	SourceAuthor acc.AccountKey   `json:"source_author"`
	SourcePostID string           `json:"source_postID"`
	Links        []IDToURLMapping `json:"links"`
}

// Donation struct, only used in Donation
type IDToURLMapping struct {
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}

// PostMeta stores tiny and frequently updated fields.
type PostMeta struct {
	Created           types.Height `json:"created"`
	LastUpdate        types.Height `json:"last_update"`
	LastActivity      types.Height `json:"last_activity"`
	AllowReplies      bool         `json:"allow_replies"`
	TotalLikeCount    int64        `json:"total_like_count"`
	TotalDonateCount  int64        `json:"total_donate_count"`
	TotalLikeWeight   int64        `json:"total_like_weight"`
	TotalDislikeStake int64        `json:"total_dislike_stake"`
	TotalReportStake  int64        `json:"total_report_stake"`
	TotalReward       sdk.Coins    `json:"reward"`
}

// Like struct, only used in PostLikes
type Like struct {
	Username acc.AccountKey `json:"username"`
	Weight   int64          `json:"weight"`
	Created  types.Height   `json:"created"`
}
type Likes []Like

// Like struct, only used in PostLikes
type Report struct {
	Username acc.AccountKey `json:"username"`
	Stake    int64          `json:"stake"`
	Created  types.Height   `json:"created"`
}
type Reports []Report

// View struct, only used in View
type Comment struct {
	Author  acc.AccountKey `json:"author"`
	PostID  string         `json:"post_key"`
	Created types.Height   `json:"created"`
}
type Comments []Comment

// View struct, only used in View
type View struct {
	Username acc.AccountKey `json:"username"`
	Created  types.Height   `json:"created"`
}
type Views []View

// Donation struct, only used in Donation
type Donation struct {
	Username acc.AccountKey `json:"username"`
	Amount   sdk.Coins      `json:"amount"`
	Created  types.Height   `json:"created"`
}
type Donations []Donation

// GetPostKey try to generate PostKey from acc.AccountKey and PostID
func GetPostKey(author acc.AccountKey, postID string) PostKey {
	return PostKey(string(author) + "#" + postID)
}

// post is the proxy for all storage structs defined above
type post struct {
	author        acc.AccountKey `json:"author"`
	postID        string         `json:"post_ID"`
	postKey       PostKey        `json:"post_key"`
	postManager   *PostManager   `json:"postManager"`
	writePostInfo bool           `json:"write_post_info"`
	writePostMeta bool           `json:"write_post_meta"`
	postInfo      *PostInfo      `json:"post_info"`
	postMeta      *PostMeta      `json:"post_meta"`
}

// create NewProxyPost
func NewProxyPost(author acc.AccountKey, postID string, postManager *PostManager) *post {
	return &post{
		author:      author,
		postID:      postID,
		postManager: postManager,
		postKey:     GetPostKey(author, postID),
	}
}

func (p *post) GetAuthor() acc.AccountKey {
	return p.author
}

func (p *post) GetPostID() string {
	return p.postID
}

func (p *post) GetPostKey() PostKey {
	return p.postKey
}

// check if post exist
func (p *post) IsPostExist(ctx sdk.Context) bool {
	if err := p.checkPostInfo(ctx); err != nil {
		return false
	}
	return true
}

// create the post
func (p *post) CreatePost(ctx sdk.Context, postInfo *PostInfo) sdk.Error {
	if p.IsPostExist(ctx) {
		return ErrPostExist()
	}
	p.writePostInfo = true
	p.postInfo = postInfo
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
func (p *post) AddOrUpdateLikeToPost(ctx sdk.Context, likeToAddOrUpdate Like) sdk.Error {
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
func (p *post) AddComment(ctx sdk.Context, comment Comment) sdk.Error {
	if err := p.UpdateLastActivity(ctx); err != nil {
		return err
	}
	return p.postManager.SetPostComment(ctx, p.GetPostKey(), &comment)
}

// add donation to post donation list
func (p *post) AddDonation(ctx sdk.Context, donation Donation) sdk.Error {
	if err := p.UpdateLastActivity(ctx); err != nil {
		return err
	}
	if err := p.postManager.SetPostDonation(ctx, p.GetPostKey(), &donation); err != nil {
		return err
	}
	p.postMeta.TotalReward = p.postMeta.TotalReward.Plus(donation.Amount)
	p.postMeta.TotalDonateCount = p.postMeta.TotalDonateCount + 1
	return nil
}

// add view to post view list
func (p *post) AddView(ctx sdk.Context, view View) sdk.Error {
	if err := p.UpdateLastActivity(ctx); err != nil {
		return err
	}
	return p.postManager.SetPostView(ctx, p.GetPostKey(), &view)
}

// add report to post report list
// func (p *post) AddReport(ctx sdk.Context, report Report) sdk.Error {
// 	return p.postManager.SetPostReport(ctx, p.GetPostKey(), &report)
// }

// update comment last activity
func (p *post) UpdateLastActivity(ctx sdk.Context) sdk.Error {
	if err := p.checkPostMeta(ctx); err != nil {
		return err
	}
	p.postMeta.LastActivity = types.Height(ctx.BlockHeight())
	p.writePostMeta = true
	return nil
}

// check if PostInfo exists
func (p *post) checkPostInfo(ctx sdk.Context) (err sdk.Error) {
	if p.postInfo == nil {
		p.postInfo, err = p.postManager.GetPostInfo(ctx, p.GetPostKey())
	}
	return err
}

// check if PostMeta exists
func (p *post) checkPostMeta(ctx sdk.Context) (err sdk.Error) {
	if p.postMeta == nil {
		p.postMeta, err = p.postManager.GetPostMeta(ctx, p.GetPostKey())
	}
	return err
}

// apply all changes to storage
func (p *post) Apply(ctx sdk.Context) sdk.Error {
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
func (p *post) clear() {
	p.writePostInfo = false
	p.writePostMeta = false
	p.postInfo = nil
	p.postMeta = nil
}
