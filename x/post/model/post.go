package model

import (
	"github.com/lino-network/lino/types"
)

// Post - post is created by the CreatedBy.
type Post struct {
	PostID    string           `json:"post_id"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	Author    types.AccountKey `json:"author"`
	CreatedBy types.AccountKey `json:"created_by"`
	CreatedAt int64            `json:"created_at"`
	UpdatedAt int64            `json:"updated_at"`
}
