package model

import (
	"github.com/lino-network/lino/types"
)

// PostRow - pk: permlink
type PostRow struct {
	Permlink types.Permlink `json:"permlink"`
	Info     PostInfo       `json:"info"`
	Meta     PostMeta       `json:"meta"`
}

// ToIR -
func (p PostRow) ToIR() PostRowIR {
	return PostRowIR{
		Permlink: p.Permlink,
		Info:     p.Info,
		Meta:     p.Meta.ToIR(),
	}
}

// PostUserRow - pk: (permlink, user)
type PostUserRow struct {
	Permlink       types.Permlink   `json:"permlink"`
	User           types.AccountKey `json:"user"`
	ReportOrUpvote ReportOrUpvote   `json:"report_or_upvote"`
	// XXX(yumin): not exported for upgrade-1
	// View           View             `json:"view"`
	// Donations      Donations        `json:"donations"`
}

// XXX(yumin): not exported for upgrade-1
// PostCommentRow - pk: (permlink, commentPermlink)
// type PostCommentRow struct {
// 	Permlink        types.Permlink `json:"permlink"`
// 	CommentPermlink types.Permlink `json:"comment_permlink"`
// 	Comment         Comment        `json:"comment"`
// }

// PostTables - state of post store.
type PostTables struct {
	Posts     []PostRow     `json:"posts"`
	PostUsers []PostUserRow `json:"post_users"`
	// not exported for upgrade-1
	// PostComments []PostCommentRow `json:"post_comments"`
}

// ToIR -
func (p PostTables) ToIR() *PostTablesIR {
	rst := &PostTablesIR{}
	for _, v := range p.Posts {
		rst.Posts = append(rst.Posts, v.ToIR())
	}
	rst.PostUsers = p.PostUsers
	return rst
}
