package model

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	linotypes "github.com/lino-network/lino/types"
)

// Validator is basic structure records all validator information
type Validator struct {
	ABCIValidator   abci.Validator       `json:"abci_validator"`
	PubKey          crypto.PubKey        `json:"pubkey"`
	Username        linotypes.AccountKey `json:"username"`
	ReceivedVotes   linotypes.Coin       `json:"received_votes"`
	HasRevoked      bool                 `json:"has_revoked"`
	AbsentCommit    int64                `json:"absent_commit"`
	ByzantineCommit int64                `json:"byzantine_commit"`
	ProducedBlocks  int64                `json:"produced_blocks"`
	Link            string               `json:"link"`
}

type ElectionVote struct {
	ValidatorName linotypes.AccountKey `json:"validator_name"`
	Vote          linotypes.Coin       `json:"votes"`
}

type ReceivedVotesStatus struct {
	ValidatorName linotypes.AccountKey `json:"validator_name"`
	ReceivedVotes linotypes.Coin       `json:"received_votes"`
}

type ElectionVoteList struct {
	ElectionVotes []ElectionVote `json:"election_votes"`
}

// ValidatorList
type ValidatorList struct {
	Oncall             []linotypes.AccountKey `json:"oncall"`
	Standby            []linotypes.AccountKey `json:"standby"`
	Candidates         []linotypes.AccountKey `json:"candidates"`
	Jail               []linotypes.AccountKey `json:"jail"`
	PreBlockValidators []linotypes.AccountKey `json:"pre_block_validators"`
	LowestOncallVotes  linotypes.Coin         `json:"lowest_oncall_votes"`
	LowestOncall       linotypes.AccountKey   `json:"lowest_oncall"`
	LowestStandbyVotes linotypes.Coin         `json:"lowest_standby_votes"`
	LowestStandby      linotypes.AccountKey   `json:"lowest_standby"`
}
