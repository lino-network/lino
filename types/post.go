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

// LikeList stores all likes of the post
type LikeList struct {
	Likes       LikeMapping `json:"likes"`
	TotalWeight int64       `json:"total_weight"`
}

// Like struct, only used in LikeList
type Like struct {
	Username AccountKey `json:"username"`
	Weight   int64      `json:"weight"`
}

// CommentList stores all comments of the post
type CommentList struct {
	Comments []PostKey `json:"comments"`
}

// ViewList stores all views of the post
type ViewList struct {
	Views []View `json:"views"`
}

// View struct, only used in ViewList
type View struct {
	Username AccountKey `json:"username"`
	Created  Height     `json:"created"`
}

// DonateList stores all donation of the post
type DonateList struct {
	Donates []Donate `json:"donates"`
	// TODO: Using sdk.Coins for now
	Reward sdk.Coins `json:"reward"`
}

// Donate struct, only used in DonateList
type Donate struct {
	Username AccountKey `json:"username"`
	Amount   sdk.Coins  `json:"amount"`
	Created  Height     `json:"created"`
}
