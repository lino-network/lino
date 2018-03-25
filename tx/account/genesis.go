package account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// State to Unmarshal
type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Name   string        `json:"name"`
	Coins  sdk.Coins     `json:"coins"`
	PubKey crypto.PubKey `json:"pub_key"`
}
