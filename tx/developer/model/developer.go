package model

import (
	types "github.com/lino-network/lino/types"
)

type Developer struct {
	Username       types.AccountKey `json:"username"`
	AppConsumption int64            `json:"app_consumption"`
}

type DeveloperList struct {
	AllDevelopers []types.AccountKey `json:"all_developers"`
}
