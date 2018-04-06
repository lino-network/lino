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

type Vote struct {
	Voter  acc.AccountKey `json:"voter"`
	Result bool           `json:"result"`
}

type Delegation struct {
	Delegator acc.AccountKey `json:"delegator"`
	Amount    types.Coin     `json:"amount"`
}

var voterRegisterFee = types.Coin{Amount: 1000 * types.Decimals}
var proposalRegisterFee = types.Coin{Amount: 2000 * types.Decimals}
