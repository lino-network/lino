package app

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/go-crypto"
)

func TestGetGenesisJson(t *testing.T) {
	masterPriv := crypto.GenPrivKeyEd25519()
	transactionPriv := crypto.GenPrivKeyEd25519()
	postPriv := crypto.GenPrivKeyEd25519()
	validatorPriv := crypto.GenPrivKeyEd25519()
	totalLino := "10000000000"
	genesisAcc := GenesisAccount{
		Name:           "Lino",
		Lino:           totalLino,
		MasterKey:      masterPriv.PubKey(),
		TransactionKey: transactionPriv.PubKey(),
		PostKey:        postPriv.PubKey(),
		IsValidator:    true,
		ValPubKey:      validatorPriv.PubKey(),
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

	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	appState, err := wire.MarshalJSONIndent(cdc, genesisState)
	assert.Nil(t, err)
	//err := oldwire.UnmarshalJSON(stateJSON, genesisState)
	appGenesisState := new(GenesisState)
	err = cdc.UnmarshalJSON([]byte(appState), appGenesisState)
	assert.Nil(t, err)

	assert.Equal(t, genesisState, *appGenesisState)
}
