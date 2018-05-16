package model

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/types"
)

// Identifier is used to map the url in post
type Identifier string

// URL used to link resources like vedio, text or photo
type URL string

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
	IsDeleted               bool       `json:"is_deleted"`
	TotalLikeCount          int64      `json:"total_like_count"`
	TotalDonateCount        int64      `json:"total_donate_count"`
	TotalLikeWeight         int64      `json:"total_like_weight"`
	TotalDislikeWeight      int64      `json:"total_dislike_weight"`
	TotalReportStake        types.Coin `json:"total_report_stake"`
	TotalUpvoteStake        types.Coin `json:"total_upvote_stake"`
	TotalViewCount          int64      `json:"total_view_count"`
	TotalReward             types.Coin `json:"reward"`
	PenaltyScore            sdk.Rat    `json:"penalty_score"`
	RedistributionSplitRate sdk.Rat    `json:"redistribution_split_rate"`
}

// Like struct, only used in Likes
type Like struct {
	Username types.AccountKey `json:"username"`
	Weight   int64            `json:"weight"`
	Created  int64            `json:"created"`
}
type Likes []Like

// ReportOrUpvote struct, only used in ReportOrUpvotes
type ReportOrUpvote struct {
	Username types.AccountKey `json:"username"`
	Stake    types.Coin       `json:"stake"`
	Created  int64            `json:"created"`
	IsReport bool             `json:"is_report"`
}
type ReportOrUpvotes []ReportOrUpvote

type Comment struct {
	Author  types.AccountKey `json:"author"`
	PostID  string           `json:"post_key"`
	Created int64            `json:"created"`
}
type Comments []Comment

// View struct
type View struct {
	Username types.AccountKey `json:"username"`
	LastView int64            `json:"last_view"`
	Times    int64            `jons:"times"`
}

// Donation struct, only used in Donation
type Donation struct {
	Amount  types.Coin `json:"amount"`
	Created int64      `json:"created"`
}
type Donations struct {
	Username     types.AccountKey `json:"username"`
	DonationList []Donation       `json:"donation_list"`
}
