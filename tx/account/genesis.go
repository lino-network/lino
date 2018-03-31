package account

import (
	crypto "github.com/tendermint/go-crypto"
)

// State to Unmarshal
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Name      string        `json:"name"`
	Lino      int64         `json:"coin"`
	PubKey    crypto.PubKey `json:"pub_key"`
	ValPubKey crypto.PubKey `json:"validator_pub_key"`
}
