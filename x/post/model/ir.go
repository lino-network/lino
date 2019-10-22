package model

import (
	"github.com/lino-network/lino/types"
)

// PostIR - is the IR of Post.
type PostIR struct {
	PostID    string           `json:"post_id"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	Author    types.AccountKey `json:"author"`
	CreatedBy types.AccountKey `json:"created_by"`
	CreatedAt int64            `json:"created_at"`
	UpdatedAt int64            `json:"updated_at"`
	IsDeleted bool             `json:"is_deleted"`
}

// PostTablesIR - is the Post State.
type PostTablesIR struct {
	Version           int              `json:"version"`
	Posts             []PostIR         `json:"posts"`
	ConsumptionWindow types.MiniDollar `json:"consumption_window"`
}
