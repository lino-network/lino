package genesis

import (
	"encoding/json"

	types "github.com/lino-network/lino/types"
	crypto "github.com/tendermint/go-crypto"
)

// genesis state for blockchain
type GenesisState struct {
	Accounts   []GenesisAccount       `json:"accounts"`
	Developers []GenesisAppDeveloper  `json:"developers"`
	Infra      []GenesisInfraProvider `json:"infra"`
	TotalLino  types.LNO              `json:"total_lino"`
}

// genesis account will get coin to the address and register user
// if genesis account is validator, it will be added to validator list automatically
type GenesisAccount struct {
	Name        string        `json:"name"`
	Lino        types.LNO     `json:"lino"`
	PubKey      crypto.PubKey `json:"pub_key"`
	IsValidator bool          `json:"is_validator"`
	ValPubKey   crypto.PubKey `json:"validator_pub_key"`
}

// register developer in genesis phase
type GenesisAppDeveloper struct {
	Name    string    `json:"name"`
	Deposit types.LNO `json:"deposit"`
}

// register infra provider in genesis phase
type GenesisInfraProvider struct {
	Name string `json:"name"`
}

// generate json format config based on genesis state
func GetGenesisJson(genesisState GenesisState) (string, error) {
	output, err := json.MarshalIndent(genesisState, "", "\t")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// default genesis file, only have one genesis account
func GetDefaultGenesis(pubkey crypto.PubKey, validatorPubKey crypto.PubKey) (string, error) {
	totalLino := "10000000000"
	genesisAcc := GenesisAccount{
		Name:        "Lino",
		Lino:        totalLino,
		PubKey:      pubkey,
		IsValidator: true,
		ValPubKey:   validatorPubKey,
	}
	genesisAppDeveloper := GenesisAppDeveloper{
		Name:    "Lino",
		Deposit: "1000000",
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
