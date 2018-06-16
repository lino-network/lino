package model

import (
	types "github.com/lino-network/lino/types"
)

type Voter struct {
	Username       types.AccountKey `json:"username"`
	Deposit        types.Coin       `json:"deposit"`
	DelegatedPower types.Coin       `json:"delegated_power"`
}

type Vote struct {
	Voter       types.AccountKey `json:"voter"`
	VotingPower types.Coin       `json:"voting_power"`
	Result      bool             `json:"result"`
}

type Delegation struct {
	Delegator types.AccountKey `json:"delegator"`
	Amount    types.Coin       `json:"amount"`
}

type ReferenceList struct {
	// OngoingProposal []types.ProposalKey `json:"ongoing_proposal"`
	AllValidators []types.AccountKey `json:"all_validators"`
}

type DelegateeList struct {
	DelegateeList []types.AccountKey `json:"delegatee_list"`
}
