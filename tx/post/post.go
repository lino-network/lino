package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

// PostKey key format in KVStore
type PostKey string

// Identifier is used to map the url in post.
type Identifier string

// URL used to link resources like vedio, text or photo.
type URL string

// PostIntreface needs to be implemented by all post related struct.
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
	PostID  string           `json:"post_id"`
	Title   string           `json:"title"`
	Content string           `json:"content"`
	Author  acc.AccountKey   `json:"author"`
	Parent  PostKey          `json:"Parent"`
	Source  PostKey          `json:"source"`
	Created types.Height     `json:"created"`
	Links   []IDToURLMapping `json:"links"`
}

// Donation struct, only used in PostDonations
type IDToURLMapping struct {
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}

// PostMeta stores tiny and frequently updated fields.
type PostMeta struct {
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

// NewLinoAccount return the account pointer
func NewPost(author acc.AccountKey, postID string, postManager *PostManager) *post {
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

func (p *post) checkPostInfo(ctx sdk.Context) (err sdk.Error) {
	if p.postInfo == nil {
		p.postInfo, err = p.postManager.GetPostInfo(ctx, p.GetPostKey())
	}
	return err
}

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
	return nil
}
