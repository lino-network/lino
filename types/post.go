package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Identifier is used to map the url in post.
type Identifier string

// URL used to link resources like vedio, text or photo.
type URL string

// PostIntreface needs to be implemented by all post related struct.
type PostInterface interface {
	AssertPostInterface()
}

var _ PostInterface = Post{}
var _ PostInterface = PostMeta{}
var _ PostInterface = PostLikes{}
var _ PostInterface = PostComments{}
var _ PostInterface = PostViews{}
var _ PostInterface = PostDonations{}

// implement PostInterface
func (_ Post) AssertPostInterface()          {}
func (_ PostMeta) AssertPostInterface()      {}
func (_ PostLikes) AssertPostInterface()     {}
func (_ PostComments) AssertPostInterface()  {}
func (_ PostViews) AssertPostInterface()     {}
func (_ PostDonations) AssertPostInterface() {}

// Post can also use to present comment(with parent) or repost(with source)
type Post struct {
	PostID       string         `json:"post_id"`
	Title        string         `json:"title"`
	Content      string         `json:"content"`
	Author       AccountKey     `json:"author"`
	ParentAuthor AccountKey     `json:"parent_author"`
	ParentPostID string         `json:"parent_postID"`
	SourceAuthor AccountKey     `json:"source_author"`
	SourcePostID string         `json:"source_postID"`
	Links        IDToURLMapping `json:"links"`
}

// Donation struct, only used in PostDonations
type IDToURLMapping struct {
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}

// PostMeta stores tiny and frequently updated fields.
type PostMeta struct {
	Created      Height `json:"created_time"`
	LastUpdate   Height `json:"last_update"`
	LastActivity Height `json:"last_activity"`
	AllowReplies bool   `json:"allow_replies"`
}

// PostLikes stores all likes of the post
type PostLikes struct {
	Likes       []Like `json:"likes"`
	TotalWeight int64  `json:"total_weight"`
}

// Like struct, only used in PostLikes
type Like struct {
	Username AccountKey `json:"username"`
	Weight   int64      `json:"weight"`
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
	Username AccountKey `json:"username"`
	Created  Height     `json:"created"`
}

// PostDonations stores all donation of the post
type PostDonations struct {
	Donations []Donation `json:"donations"`
	// TODO: Using sdk.Coins for now
	Reward sdk.Coins `json:"reward"`
}

// Donation struct, only used in PostDonations
type Donation struct {
	Username AccountKey `json:"username"`
	Amount   sdk.Coins  `json:"amount"`
	Created  Height     `json:"created"`
}

// PostManager is the bridge to get post from store
type PostManager interface {
	CreatePost(ctx sdk.Context, post *Post) sdk.Error
	GetPost(ctx sdk.Context, postKey PostKey) (*Post, sdk.Error)
	SetPost(ctx sdk.Context, post *Post) sdk.Error

	GetPostMeta(ctx sdk.Context, postKey PostKey) (*PostMeta, sdk.Error)
	SetPostMeta(ctx sdk.Context, postKey PostKey, postMeta *PostMeta) sdk.Error

	GetPostLikes(ctx sdk.Context, postKey PostKey) (*PostLikes, sdk.Error)
	SetPostLikes(ctx sdk.Context, postKey PostKey, PostLikes *PostLikes) sdk.Error

	GetPostComments(ctx sdk.Context, postKey PostKey) (*PostComments, sdk.Error)
	SetPostComments(ctx sdk.Context, postKey PostKey, PostComments *PostComments) sdk.Error

	GetPostViews(ctx sdk.Context, postKey PostKey) (*PostViews, sdk.Error)
	SetPostViews(ctx sdk.Context, postKey PostKey, PostViews *PostViews) sdk.Error

	GetPostDonations(ctx sdk.Context, postKey PostKey) (*PostDonations, sdk.Error)
	SetPostDonations(ctx sdk.Context, postKey PostKey, PostDonations *PostDonations) sdk.Error
}

// GetPostKey try to generate PostKey from AccountKey and PostID
func GetPostKey(author AccountKey, postID string) PostKey {
	return PostKey(string(author) + "#" + postID)
}
