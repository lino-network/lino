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
var _ PostInterface = PostLikes{}
var _ PostInterface = PostComments{}
var _ PostInterface = PostViews{}
var _ PostInterface = PostDonations{}

func (_ PostInfo) AssertPostInterface()      {}
func (_ PostMeta) AssertPostInterface()      {}
func (_ PostLikes) AssertPostInterface()     {}
func (_ PostComments) AssertPostInterface()  {}
func (_ PostViews) AssertPostInterface()     {}
func (_ PostDonations) AssertPostInterface() {}

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

// Donation struct, only used in PostDonations
type IDToURLMapping struct {
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}

// PostMeta stores tiny and frequently updated fields.
type PostMeta struct {
	Created      types.Height `json:"created"`
	LastUpdate   types.Height `json:"last_update"`
	LastActivity types.Height `json:"last_activity"`
	AllowReplies bool         `json:"allow_replies"`
}

// PostLikes stores all likes of the post
type PostLikes struct {
	Likes       []Like `json:"likes"`
	TotalWeight int64  `json:"total_weight"`
}

// Like struct, only used in PostLikes
type Like struct {
	Username acc.AccountKey `json:"username"`
	Weight   int64          `json:"weight"`
}

// PostComments stores all comments of the post
type PostComments struct {
	Comments []PostKey `json:"comments"`
}

// PostViews stores all views of the post
type PostViews struct {
	Views []View `json:"views"`
}

// View struct, only used in PostViews
type View struct {
	Username acc.AccountKey `json:"username"`
	Created  types.Height   `json:"created"`
}

// PostDonations stores all donation of the post
type PostDonations struct {
	Donations []Donation `json:"donates"`
	// TODO: Using sdk.Coins for now
	Reward sdk.Coins `json:"reward"`
}

// Donation struct, only used in PostDonations
type Donation struct {
	Username acc.AccountKey `json:"username"`
	Amount   sdk.Coins      `json:"amount"`
	Created  types.Height   `json:"created"`
}

// GetPostKey try to generate PostKey from acc.AccountKey and PostID
func GetPostKey(author acc.AccountKey, postID string) PostKey {
	return PostKey(string(author) + "#" + postID)
}

// post is the proxy for all storage structs defined above
type post struct {
	author             acc.AccountKey `json:"author"`
	postID             string         `json:"post_ID"`
	postKey            PostKey        `json:"post_key"`
	postManager        *PostManager   `json:"postManager"`
	writePostInfo      bool           `json:"write_post_info"`
	writePostMeta      bool           `json:"write_post_meta"`
	writePostLikes     bool           `json:"write_post_likes"`
	writePostComments  bool           `json:"write_post_comments"`
	writePostViews     bool           `json:"write_post_views"`
	writePostDonations bool           `json:"write_post_donations"`
	postInfo           *PostInfo      `json:"post_info"`
	postMeta           *PostMeta      `json:"post_meta"`
	postLikes          *PostLikes     `json:"post_likes"`
	postComments       *PostComments  `json:"post_comments"`
	postViews          *PostViews     `json:"post_views"`
	postDonations      *PostDonations `json:"post_donations"`
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
	p.writePostLikes = true
	p.postLikes = &PostLikes{Likes: []Like{}}
	p.writePostComments = true
	p.postComments = &PostComments{Comments: []PostKey{}}
	p.writePostViews = true
	p.postViews = &PostViews{Views: []View{}}
	p.writePostDonations = true
	p.postDonations = &PostDonations{Donations: []Donation{}, Reward: sdk.Coins{}}
	return nil
}

// add comment to post comment list
func (p *post) AddComment(ctx sdk.Context, comment PostKey) sdk.Error {
	if err := p.checkPostComments(ctx); err != nil {
		return err
	}
	p.writePostComments = true
	p.postComments.Comments = append(p.postComments.Comments, comment)
	if err := p.UpdateLastActivity(ctx); err != nil {
		return err
	}
	return nil
}

// update comment last activity
func (p *post) UpdateLastActivity(ctx sdk.Context) sdk.Error {
	if err := p.checkPostMeta(ctx); err != nil {
		return err
	}
	p.writePostMeta = true
	p.postMeta.LastActivity = types.Height(ctx.BlockHeight())
	return nil
}

// check if PostInfo exists
func (p *post) checkPostInfo(ctx sdk.Context) (err sdk.Error) {
	if p.postInfo == nil {
		p.postInfo, err = p.postManager.GetPostInfo(ctx, p.GetPostKey())
	}
	return err
}

// check if PostComments exists
func (p *post) checkPostComments(ctx sdk.Context) (err sdk.Error) {
	if p.postComments == nil {
		p.postComments, err = p.postManager.GetPostComments(ctx, p.GetPostKey())
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
	if p.writePostLikes {
		if err := p.postManager.SetPostLikes(ctx, p.GetPostKey(), p.postLikes); err != nil {
			return err
		}
	}
	if p.writePostComments {
		if err := p.postManager.SetPostComments(ctx, p.GetPostKey(), p.postComments); err != nil {
			return err
		}
	}
	if p.writePostViews {
		if err := p.postManager.SetPostViews(ctx, p.GetPostKey(), p.postViews); err != nil {
			return err
		}
	}
	if p.writePostDonations {
		if err := p.postManager.SetPostDonations(ctx, p.GetPostKey(), p.postDonations); err != nil {
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
	p.writePostLikes = false
	p.writePostComments = false
	p.writePostViews = false
	p.writePostDonations = false
	p.postInfo = nil
	p.postMeta = nil
	p.postLikes = nil
	p.postComments = nil
	p.postViews = nil
	p.postDonations = nil
}
