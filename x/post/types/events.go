package types

import (
	linotypes "github.com/lino-network/lino/types"
)

// RewardEventV1 - used only for importing.
type RewardEventV1 struct {
	PostAuthor linotypes.AccountKey `json:"post_author"`
	PostID     string               `json:"post_id"`
	Consumer   linotypes.AccountKey `json:"consumer"`
	Evaluate   linotypes.Coin       `json:"evaluate"`
	Original   linotypes.Coin       `json:"original"`
	Friction   linotypes.Coin       `json:"friction"`
	FromApp    linotypes.AccountKey `json:"from_app"`
}

// RewardEvent - when donation occurred, a reward event will be register
// at 7 days later. After 7 days reward event will be executed and send
// inflation to author.
type RewardEvent struct {
	PostAuthor linotypes.AccountKey `json:"post_author"`
	PostID     string               `json:"post_id"`
	Consumer   linotypes.AccountKey `json:"consumer"`
	Evaluate   linotypes.MiniDollar `json:"evaluate"`
	FromApp    linotypes.AccountKey `json:"from_app"`
}
