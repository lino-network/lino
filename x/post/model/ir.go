package model

import (
	"github.com/lino-network/lino/types"
)

// PostMetaIR RedistributionSplitRate rat -> string
type PostMetaIR struct {
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
	RedistributionSplitRate string     `json:"redistribution_split_rate"`
}

// PostRowIR - Meta changed
type PostRowIR struct {
	Permlink types.Permlink `json:"permlink"`
	Info     PostInfo       `json:"info"`
	Meta     PostMetaIR     `json:"meta"`
}

// PostTablesIR - PostRow changed.
type PostTablesIR struct {
	Posts     []PostRowIR   `json:"posts"`
	PostUsers []PostUserRow `json:"post_users"`
}
