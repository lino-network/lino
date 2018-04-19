package genesis

import (
	"encoding/json"
	crypto "github.com/tendermint/go-crypto"
)

// State to Unmarshal
type GenesisState struct {
	Accounts   []GenesisAccount       `json:"accounts"`
	Developers []GenesisAppDeveloper  `json:"developers"`
	Infra      []GenesisInfraProvider `json:"infra"`
	TotalLino  int64                  `json:"total_lino"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Name        string        `json:"name"`
	Lino        int64         `json:"lino"`
	PubKey      crypto.PubKey `json:"pub_key"`
	IsValidator bool          `json:"is_validator"`
	ValPubKey   crypto.PubKey `json:"validator_pub_key"`
}

type GenesisAppDeveloper struct {
	Name    string `json:"name"`
	Deposit int64  `json:"deposit"`
}

type GenesisInfraProvider struct {
	Name string `json:"name"`
}

func GetGenesisJson(genesisState GenesisState) (string, error) {
	output, err := json.MarshalIndent(genesisState, "", "\t")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func GetDefaultGenesis(pubkey crypto.PubKey, validatorPubKey crypto.PubKey) (string, error) {
	totalLino := int64(10000000000)
	genesisAcc := GenesisAccount{
		Name:        "Lino",
		Lino:        totalLino,
		PubKey:      pubkey,
		IsValidator: true,
		ValPubKey:   validatorPubKey,
	}
	genesisAppDeveloper := GenesisAppDeveloper{
		Name:    "Lino",
		Deposit: 1000000,
	}
	genesisInfraProvider := GenesisInfraProvider{
		Name: "Lino",
	}
	genesisState := GenesisState{
		Accounts:   []GenesisAccount{genesisAcc},
		TotalLino:  totalLino,
		Developers: []GenesisAppDeveloper{genesisAppDeveloper},
		Infra:      []GenesisInfraProvider{genesisInfraProvider},
	}

	result, err := GetGenesisJson(genesisState)
	if err != nil {
		return "", err
	}
	return result, err
}
