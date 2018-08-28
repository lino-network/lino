package model

import (
	types "github.com/lino-network/lino/types"
)

// Voter - a voter in blockchain is account with voter deposit, who can vote for a proposal
type Voter struct {
	Username       types.AccountKey `json:"username"`
	Deposit        types.Coin       `json:"deposit"`
	DelegatedPower types.Coin       `json:"delegated_power"`
}

// Vote - a vote is created by a voter to a proposal
type Vote struct {
	Voter       types.AccountKey `json:"voter"`
	VotingPower types.Coin       `json:"voting_power"`
	Result      bool             `json:"result"`
}

// Delegation - normal user can delegate money to a voter to increase voter's voting power
type Delegation struct {
	Delegator types.AccountKey `json:"delegator"`
	Amount    types.Coin       `json:"amount"`
}

// ReferenceList - record validator to punish the validator who doesn't vote for proposal
type ReferenceList struct {
	AllValidators []types.AccountKey `json:"all_validators"`
}
