package post

import "time"

type Post struct {
	Author           string    `json:"author"`
	PostID           string    `json:"postID"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	ParentAuthor     string    `json:"parentAuthor"`
	ParentPostID     string    `json:"parentPostID"`
	SourceAuthor     string    `json:"sourceAuthor"`
	SourcePostID     string    `json:"sourcePostID"`
	Links            string    `json:"links"`
	CreatedAt        time.Time `json:"createdAt"`
	TotalDonateCount int64     `json:"totalDonateCount"`
	TotalReward      string    `json:"totalReward"`
	IsDeleted        bool      `json:"is_deleted"`
}
