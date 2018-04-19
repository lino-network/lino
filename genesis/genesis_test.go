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
	genesisState := GenesisState{
		Accounts:  []GenesisAccount{genesisAcc},
		TotalLino: totalLino,
	}

	result, err := GetGenesisJson(genesisState)
	assert.Nil(t, err)
	//err := oldwire.UnmarshalJSON(stateJSON, genesisState)
	appGenesisState := new(GenesisState)
	err = json.Unmarshal([]byte(result), appGenesisState)
	assert.Nil(t, err)

	assert.Equal(t, genesisState, *appGenesisState)
}
