package model

import (
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	CreatedAt               int64      `json:"created_at"`
	LastUpdatedAt           int64      `json:"last_updated_at"`
	LastActivityAt          int64      `json:"last_activity_at"`
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
	Username  types.AccountKey `json:"username"`
	Weight    int64            `json:"weight"`
	CreatedAt int64            `json:"created_at"`
}
type Likes []Like

// ReportOrUpvote struct, only used in ReportOrUpvotes
type ReportOrUpvote struct {
	Username  types.AccountKey `json:"username"`
	Stake     types.Coin       `json:"stake"`
	CreatedAt int64            `json:"created_at"`
	IsReport  bool             `json:"is_report"`
}
type ReportOrUpvotes []ReportOrUpvote

type Comment struct {
	Author    types.AccountKey `json:"author"`
	PostID    string           `json:"post_key"`
	CreatedAt int64            `json:"created_at"`
}
type Comments []Comment

// View struct
type View struct {
	Username   types.AccountKey `json:"username"`
	LastViewAt int64            `json:"last_view_at"`
	Times      int64            `jons:"times"`
}

// Donation struct, only used in Donation
type Donation struct {
	Amount    types.Coin `json:"amount"`
	CreatedAt int64      `json:"created_at"`
}
type Donations struct {
	Username     types.AccountKey `json:"username"`
	DonationList []Donation       `json:"donation_list"`
}
