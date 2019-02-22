package model

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	types "github.com/lino-network/lino/types"
)

// Validator is basic structure records all validator information
type Validator struct {
	ABCIValidator   abci.Validator
	PubKey          crypto.PubKey    `json:"pubkey"`
	Username        types.AccountKey `json:"username"`
	Deposit         types.Coin       `json:"deposit"`
	AbsentCommit    int64            `json:"absent_commit"`
	ByzantineCommit int64            `json:"byzantine_commit"`
	ProducedBlocks  int64            `json:"produced_blocks"`
	Link            string           `json:"link"`
}

// Validator list
type ValidatorList struct {
	OncallValidators   []types.AccountKey `json:"oncall_validators"`
	AllValidators      []types.AccountKey `json:"all_validators"`
	PreBlockValidators []types.AccountKey `json:"pre_block_validators"`
	LowestPower        types.Coin         `json:"lowest_power"`
	LowestValidator    types.AccountKey   `json:"lowest_validator"`
}
