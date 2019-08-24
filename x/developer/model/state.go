package model

import (
	"github.com/lino-network/lino/types"
)

// Developer -
type DeveloperV1 struct {
	Username       types.AccountKey `json:"username"`
	Deposit        types.Coin       `json:"deposit"`
	AppConsumption types.Coin       `json:"app_consumption"`
	Website        string           `json:"web_site"`
	Description    string           `json:"description"`
	AppMetaData    string           `json:"app_meta_data"`
}

// DeveloperRow - pk: Username
type DeveloperRow struct {
	Username  types.AccountKey `json:"username"`
	Developer DeveloperV1      `json:"developer"`
}

// DeveloperList is deprecated since upgrade2.
// // DeveloperListTable all developers, pk: none
// type DeveloperListTable struct {
// 	List DeveloperList `json:"list"`
// }

// // DeveloperTables is the state of developer storage, organized as a table.
type DeveloperTables struct {
	Developers []DeveloperRow `json:"developers"`
	// DeveloperList DeveloperListTable `json:"developer_list"`
}

// ToIR -
func (d DeveloperTables) ToIR() DeveloperTablesIR {
	return d
}
