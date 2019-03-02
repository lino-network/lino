package model

import (
	types "github.com/lino-network/lino/types"
)

// VoterRow - pk: username
type VoterRow struct {
	Username types.AccountKey `json:"username"`
	Voter    Voter            `json:"voter"`
}

// DelegationRow - pk: (voter, delegator)
type DelegationRow struct {
	Voter      types.AccountKey `json:"username"`
	Delegator  types.AccountKey `json:"delegator"`
	Delegation Delegation       `json:"delegation"`
}

// ReferenceListTable - no pk
type ReferenceListTable struct {
	List ReferenceList `json:"list"`
}

// VoterTables - state of voter
type VoterTables struct {
	Voters        []VoterRow         `json:"voters"`
	Delegations   []DelegationRow    `json:"delegations"`
	ReferenceList ReferenceListTable `json:"reference_list"`
}

// ToIR - same
func (v VoterTables) ToIR() VoterTablesIR {
	return v
}
