package model

import (
	"github.com/lino-network/lino/types"
)

// ValidatorRow - pk: (Username)
type ValidatorRow struct {
	Username types.AccountKey `json:"username"`
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
	List ValidatorList `json:"list"`
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
