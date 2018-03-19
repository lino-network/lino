package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Identifier is used to map the url in post
type Identifier string

// URL used to link resources like vedio, text or photo.
type URL string

// LinkMapping records mapping between Identifier and URL
type LinkMapping map[Identifier]URL

// LikeMapping records mapping between AccountKey and Like
type LikeMapping map[AccountKey]Like

// Post can also use to present comment(with parent) or repost(with source)
type Post struct {
	Key     PostKey     `json:"key"`
	Title   string      `json:"title"`
	Content string      `json:"content"`
	Author  AccountKey  `json:"author"`
	Parent  PostKey     `json:"Parent"`
	Source  PostKey     `json:"source"`
	Created Height      `json:"created"`
	Links   LinkMapping `json:"links"`
}

// PostMeta stores tiny and frequently updated fields.
type PostMeta struct {
	LastUpdate   Height `json:"last_update"`
	LastActivity Height `json:"last_activity"`
	AllowReplies bool   `json:"allow_replies"`
	AllowLike    bool   `json:"allow_like"`
}

// PostLikes stores all likes of the post
type PostLikes struct {
	Likes       LikeMapping `json:"likes"`
	TotalWeight int64       `json:"total_weight"`
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
	Donations []Donation `json:"donates"`
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
