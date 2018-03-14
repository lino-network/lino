package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-crypto"
)

type PostKey []byte
type Identifier string
type URL string
type LinkMapping map[Identifier]URL
type LikeMapping map[AccountKey]Like

// Lino Post, which can also used to present comment(with parent) or repost(with source)
// key: Username + Sequence PostKey
type Post struct {
	Title   string      `json:"title"`
	Content string      `json:"content"`
	Author  AccountKey  `json:"author"`
	Parent  PostKey     `json:"Parent"`
	Source  PostKey     `json:"source"`
	Created uint64      `json:"created"`
	Links   LinkMapping `json:"links"`
}

// PostMeta stores tiny and frequently updated fields.
// key: Username + Sequence PostKey
type PostMeta struct {
	LastUpdate   uint64 `json:"last_update"`
	LastActivity uint64 `json:"last_activity"`
	AllowReplies uint64 `json:"allow_replies"`
	AllowLike    uint64 `json:"allow_like"`
}

// like list of the post
// key: Username + Sequence PostKey
type LikeList struct {
	Likes       LikeMapping `json:"likes"`
	TotalWeight int64       `json:"total_weight"`
}

type Like struct {
	Username AccountKey `json:"username"`
	Weight   int64      `json:"weight"`
}

// all comments of the post
// key: Username + Sequence PostKey
type CommentList struct {
	Comments []PostKey `json:"comments"`
}

// all views of the post
// key: Username + Sequence PostKey
type ViewList struct {
	Views []View `json:"views"`
}

type View struct {
	Username AccountKey `json:"username"`
	When     uint64     `json:"when"`
}

// all donation of this post
// key: Username + Sequence PostKey
type DonateList struct {
	Donates []Donate `json:"donates"`
	Reward  Coins    `json:"reward"`
}

type Donate struct {
	Username AccountKey `json:"username"`
	Amount   sdk.Coins  `json:"amount"`
	When     uint64     `json:"when"`
}
