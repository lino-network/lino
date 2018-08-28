package model

import (
	types "github.com/lino-network/lino/types"
)

// Developer - developer is account with developer deposit, can get developer inflation
type Developer struct {
	Username       types.AccountKey `json:"username"`
	Deposit        types.Coin       `json:"deposit"`
	AppConsumption types.Coin       `json:"app_consumption"`
	Website        string           `json:"web_site"`
	Description    string           `json:"description"`
	AppMetaData    string           `json:"app_meta_data"`
}

// DeveloperList - list of developers
type DeveloperList struct {
	AllDevelopers []types.AccountKey `json:"all_developers"`
}
