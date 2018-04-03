package vote

import (
	acc "github.com/lino-network/lino/tx/account"
	types "github.com/lino-network/lino/types"
)

type Voter struct {
	Username       acc.AccountKey `json:"username"`
	Deposit        types.Coin     `json:"deposit"`
	DelegatedPower types.Coin     `json:"delegated_power"`
}

type Delegation struct {
	Delegator acc.AccountKey `json:"delegator"`
	Amount    types.Coin     `json:"amount"`
}

var valRegisterFee = types.Coin{Amount: 100 * types.Decimals}
