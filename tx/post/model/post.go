package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

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
var _ PostInterface = Donations{}

func (_ PostInfo) AssertPostInterface()  {}
func (_ PostMeta) AssertPostInterface()  {}
func (_ Like) AssertPostInterface()      {}
func (_ Report) AssertPostInterface()    {}
func (_ Comment) AssertPostInterface()   {}
func (_ View) AssertPostInterface()      {}
func (_ Donations) AssertPostInterface() {}

// PostInfo can also use to present comment(with parent) or repost(with source)
type PostInfo struct {
	PostID       string                 `json:"post_id"`
	Title        string                 `json:"title"`
	Content      string                 `json:"content"`
	Author       types.AccountKey       `json:"author"`
	ParentAuthor types.AccountKey       `json:"parent_author"`
	ParentPostID string                 `json:"parent_postID"`
	SourceAuthor types.AccountKey       `json:"source_author"`
	SourcePostID string                 `json:"source_postID"`
	Links        []types.IDToURLMapping `json:"links"`
}

// PostMeta stores tiny and frequently updated fields.
type PostMeta struct {
	Created                 int64      `json:"created"`
	LastUpdate              int64      `json:"last_update"`
	LastActivity            int64      `json:"last_activity"`
	AllowReplies            bool       `json:"allow_replies"`
	TotalLikeCount          int64      `json:"total_like_count"`
	TotalDonateCount        int64      `json:"total_donate_count"`
	TotalLikeWeight         int64      `json:"total_like_weight"`
	TotalDislikeWeight      int64      `json:"total_dislike_weight"`
	TotalReward             types.Coin `json:"reward"`
	PenaltyScore            sdk.Rat    `json:"penalty_score"`
	RedistributionSplitRate sdk.Rat    `json:"redistribution_split_rate"`
}

// Like struct, only used in PostLikes
type Like struct {
	Username types.AccountKey `json:"username"`
	Weight   int64            `json:"weight"`
	Created  int64            `json:"created"`
}
type Likes []Like

// Like struct, only used in PostLikes
type Report struct {
	Username types.AccountKey `json:"username"`
	Stake    int64            `json:"stake"`
	Created  int64            `json:"created"`
}
type Reports []Report

// View struct, only used in View
type Comment struct {
	Author  types.AccountKey `json:"author"`
	PostID  string           `json:"post_key"`
	Created int64            `json:"created"`
}
type Comments []Comment

// View struct, only used in View
type View struct {
	Username types.AccountKey `json:"username"`
	Created  int64            `json:"created"`
	Times    int64            `jons:"times"`
}
type Views []View

// Donation struct, only used in Donation
type Donation struct {
	Username types.AccountKey `json:"username"`
	Amount   types.Coin       `json:"amount"`
	Created  int64            `json:"created"`
}
type Donations struct {
	Username     types.AccountKey `json:"username"`
	DonationList []Donation       `json:"donation_list"`
}
