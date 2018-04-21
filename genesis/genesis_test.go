package genesis

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tendermint/go-crypto"
)

func TestGetGenesisJson(t *testing.T) {
	genesisAccPriv := crypto.GenPrivKeyEd25519()
	validatorPriv := crypto.GenPrivKeyEd25519()
	totalLino := int64(10000000000)
	genesisAcc := GenesisAccount{
		Name:        "Lino",
		Lino:        totalLino,
		PubKey:      genesisAccPriv.PubKey(),
		IsValidator: true,
		ValPubKey:   validatorPriv.PubKey(),
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
	assert.Nil(t, err)
	//err := oldwire.UnmarshalJSON(stateJSON, genesisState)
	appGenesisState := new(GenesisState)
	err = json.Unmarshal([]byte(result), appGenesisState)
	assert.Nil(t, err)

	assert.Equal(t, genesisState, *appGenesisState)
}
