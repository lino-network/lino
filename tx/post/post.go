package post

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"
)

// PostKey key format in KVStore
type PostKey string

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
	PostID       string           `json:"post_id"`
	Title        string           `json:"title"`
	Content      string           `json:"content"`
	Author       acc.AccountKey   `json:"author"`
	ParentAuthor acc.AccountKey   `json:"parent_author"`
	ParentPostID string           `json:"parent_postID"`
	SourceAuthor acc.AccountKey   `json:"source_author"`
	SourcePostID string           `json:"source_postID"`
	Links        []IDToURLMapping `json:"links"`
}

// Donation struct, only used in Donation
type IDToURLMapping struct {
	Identifier string `json:"identifier"`
	URL        string `json:"url"`
}

// PostMeta stores tiny and frequently updated fields.
type PostMeta struct {
	Created                 types.Height `json:"created"`
	RedistributionSplitRate sdk.Rat      `json:"redistribution_split_rate"`
	LastUpdate              types.Height `json:"last_update"`
	LastActivity            types.Height `json:"last_activity"`
	AllowReplies            bool         `json:"allow_replies"`
	TotalLikeCount          int64        `json:"total_like_count"`
	TotalDonateCount        int64        `json:"total_donate_count"`
	TotalLikeWeight         int64        `json:"total_like_weight"`
	TotalDislikeStake       int64        `json:"total_dislike_stake"`
	TotalReportStake        int64        `json:"total_report_stake"`
	TotalReward             types.Coin   `json:"reward"`
	PenaltyScore            sdk.Rat      `json:"penalty_score"`
}

// Like struct, only used in PostLikes
type Like struct {
	Username acc.AccountKey `json:"username"`
	Weight   int64          `json:"weight"`
	Created  types.Height   `json:"created"`
}
type Likes []Like

// Like struct, only used in PostLikes
type Report struct {
	Username acc.AccountKey `json:"username"`
	Stake    int64          `json:"stake"`
	Created  types.Height   `json:"created"`
}
type Reports []Report

// View struct, only used in View
type Comment struct {
	Author  acc.AccountKey `json:"author"`
	PostID  string         `json:"post_key"`
	Created types.Height   `json:"created"`
}
type Comments []Comment

// View struct, only used in View
type View struct {
	Username acc.AccountKey `json:"username"`
	Created  types.Height   `json:"created"`
	Times    int64          `jons:"times"`
}
type Views []View

// Donation struct, only used in Donation
type Donation struct {
	Username acc.AccountKey `json:"username"`
	Amount   types.Coin     `json:"amount"`
	Created  types.Height   `json:"created"`
}
type Donations struct {
	Username     acc.AccountKey `json:"username"`
	DonationList []Donation     `json:"donation_list"`
}

// GetPostKey try to generate PostKey from acc.AccountKey and PostID
func GetPostKey(author acc.AccountKey, postID string) PostKey {
	return PostKey(string(author) + "#" + postID)
}
