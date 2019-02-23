package model

import (
	"github.com/lino-network/lino/types"
)

// DeveloperRow - pk: Username
type DeveloperRow struct {
	Username  types.AccountKey `json:"username"`
	Developer Developer        `json:"developer"`
}

// DeveloperListTable all developers, pk: none
type DeveloperListTable struct {
	List DeveloperList `json:"list"`
}

// DeveloperTables is the state of developer storage, organized as a table.
type DeveloperTables struct {
	Developers    []DeveloperRow     `json:"developers"`
	DeveloperList DeveloperListTable `json:"developer_list"`
}

// ToIR -
func (d DeveloperTables) ToIR() DeveloperTablesIR {
	return d
}
