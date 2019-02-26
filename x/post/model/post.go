package model

import (
	"github.com/lino-network/lino/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Identifier is used to map the url in post
type Identifier string

// URL used to link resources such as vedio, text or photo
type URL string

// PostInfo - can also use to present comment(with parent) or repost(with source)
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

// PostMeta - stores tiny and frequently updated fields.
type PostMeta struct {
	CreatedAt               int64      `json:"created_at"`
	LastUpdatedAt           int64      `json:"last_updated_at"`
	LastActivityAt          int64      `json:"last_activity_at"`
	AllowReplies            bool       `json:"allow_replies"`
	IsDeleted               bool       `json:"is_deleted"`
	TotalDonateCount        int64      `json:"total_donate_count"`
	TotalReportCoinDay      types.Coin `json:"total_report_coin_day"`
	TotalUpvoteCoinDay      types.Coin `json:"total_upvote_coin_day"`
	TotalViewCount          int64      `json:"total_view_count"`
	TotalReward             types.Coin `json:"total_reward"`
	RedistributionSplitRate sdk.Dec    `json:"redistribution_split_rate"`
}

// ToIR -
func (pm PostMeta) ToIR() PostMetaIR {
	return PostMetaIR{
		CreatedAt:               pm.CreatedAt,
		LastUpdatedAt:           pm.LastUpdatedAt,
		LastActivityAt:          pm.LastActivityAt,
		AllowReplies:            pm.AllowReplies,
		IsDeleted:               pm.IsDeleted,
		TotalDonateCount:        pm.TotalDonateCount,
		TotalReportCoinDay:      pm.TotalReportCoinDay,
		TotalUpvoteCoinDay:      pm.TotalUpvoteCoinDay,
		TotalViewCount:          pm.TotalViewCount,
		TotalReward:             pm.TotalReward,
		RedistributionSplitRate: pm.RedistributionSplitRate.String(), // XXX(yumin): rat to dec
	}
}

// ReportOrUpvote - report or upvote from a user to a post
type ReportOrUpvote struct {
	Username  types.AccountKey `json:"username"`
	CoinDay   types.Coin       `json:"coin_day"`
	CreatedAt int64            `json:"created_at"`
	IsReport  bool             `json:"is_report"`
}

// Comment - comment list store dy a post
type Comment struct {
	Author    types.AccountKey `json:"author"`
	PostID    string           `json:"post_id"`
	CreatedAt int64            `json:"created_at"`
}

// View - from a user to a post
type View struct {
	Username   types.AccountKey `json:"username"`
	LastViewAt int64            `json:"last_view_at"`
	Times      int64            `jons:"times"`
}

// Donations - record a user donation behavior to a post
type Donations struct {
	Username types.AccountKey `json:"username"`
	Times    int64            `json:"times"`
	Amount   types.Coin       `json:"amount"`
}
