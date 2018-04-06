package genesis

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-crypto"
)

func TestGetGenesisJson(t *testing.T) {
	genesisAccPriv := crypto.GenPrivKeyEd25519()
	validatorPriv := crypto.GenPrivKeyEd25519()
	totalLino := int64(10000000000)
	genesisAcc := GenesisAccount{
		Name:      "Lino",
		Lino:      totalLino,
		PubKey:    genesisAccPriv.PubKey(),
		ValPubKey: validatorPriv.PubKey(),
	}
	globalState := GlobalState{
		TotalLino:                totalLino,
		GrowthRate:               sdk.NewRat(98, 1000),
		InfraAllocation:          sdk.NewRat(20, 100),
		ContentCreatorAllocation: sdk.NewRat(50, 100),
		DeveloperAllocation:      sdk.NewRat(20, 100),
		ValidatorAllocation:      sdk.NewRat(10, 100),
		ConsumptionFrictionRate:  sdk.NewRat(1, 100),
		FreezingPeriodHr:         24 * 7,
	}
	genesisState := GenesisState{
		Accounts:    []GenesisAccount{genesisAcc},
		GlobalState: globalState,
	}

	result, err := GetGenesisJson(genesisState)
	assert.Nil(t, err)
	//err := oldwire.UnmarshalJSON(stateJSON, genesisState)
	appGenesisState := new(GenesisState)
	err = json.Unmarshal([]byte(result), appGenesisState)
	assert.Nil(t, err)

	assert.Equal(t, genesisState, *appGenesisState)
}
