package model

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	linotypes "github.com/lino-network/lino/types"
)

type ValidatorListV1 struct {
	OncallValidators   []linotypes.AccountKey `json:"oncall_validators"`
	AllValidators      []linotypes.AccountKey `json:"all_validators"`
	PreBlockValidators []linotypes.AccountKey `json:"pre_block_validators"`
	LowestPower        linotypes.Coin         `json:"lowest_power"`
	LowestValidator    linotypes.AccountKey   `json:"lowest_validator"`
}

type ValidatorV1 struct {
	ABCIValidator   abci.Validator
	PubKey          crypto.PubKey        `json:"pubkey"`
	Username        linotypes.AccountKey `json:"username"`
	Deposit         linotypes.Coin       `json:"deposit"`
	AbsentCommit    int64                `json:"absent_commit"`
	ByzantineCommit int64                `json:"byzantine_commit"`
	ProducedBlocks  int64                `json:"produced_blocks"`
	Link            string               `json:"link"`
}

// ValidatorRow - pk: (Username)
type ValidatorRow struct {
	Username linotypes.AccountKey `json:"username"`
	// XXX(yumin): type changed.
	Validator ValidatorV1 `json:"validator"`
}

// ToIR -
func (v ValidatorRow) ToIR() ValidatorRowIR {
	return ValidatorRowIR{
		Username:  v.Username,
		Validator: v.Validator.ToIR(),
	}
}

// ValidatorListRow - pk: none
type ValidatorListRow struct {
	List ValidatorListV1 `json:"list"`
}

// ValidatorTables state of validators
type ValidatorTables struct {
	Validators    []ValidatorRow   `json:"validators"`
	ValidatorList ValidatorListRow `json:"validator_list"`
}

// ToIR -
func (v ValidatorTables) ToIR() *ValidatorTablesIR {
	rst := &ValidatorTablesIR{}
	for _, v := range v.Validators {
		rst.Validators = append(rst.Validators, v.ToIR())
	}
	rst.ValidatorList = v.ValidatorList
	return rst
}
