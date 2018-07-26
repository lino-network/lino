package app

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/wire"
)

func TestGetGenesisJson(t *testing.T) {
	resetPriv := secp256k1.GenPrivKey()
	transactionPriv := secp256k1.GenPrivKey()
	appPriv := secp256k1.GenPrivKey()
	validatorPriv := secp256k1.GenPrivKey()

	totalLino := "10000000000"
	genesisAcc := GenesisAccount{
		Name:           "Lino",
		Lino:           totalLino,
		ResetKey:       resetPriv.PubKey(),
		TransactionKey: transactionPriv.PubKey(),
		AppKey:         appPriv.PubKey(),
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

func TestLinoBlockchainGenTx(t *testing.T) {
	cdc := MakeCodec()
	pk := secp256k1.GenPrivKey().PubKey()
	var genTxConfig config.GenTx
	appGenTx, _, validator, err := LinoBlockchainGenTx(cdc, pk, genTxConfig)
	assert.Nil(t, err)
	var genesisAcc GenesisAccount
	err = cdc.UnmarshalJSON(appGenTx, &genesisAcc)
	assert.Nil(t, err)
	assert.Equal(t, genesisAcc.Name, "lino")
	assert.Equal(t, genesisAcc.Lino, "10000000000")
	assert.Equal(t, genesisAcc.IsValidator, true)
	assert.Equal(t, genesisAcc.ValPubKey, pk)
	assert.Equal(t, validator.PubKey, pk)
}

func TestLinoBlockchainGenState(t *testing.T) {
	cdc := MakeCodec()
	appGenTxs := []json.RawMessage{}
	for i := 1; i < 21; i++ {
		genesisAcc := GenesisAccount{
			Name:           "validator" + strconv.Itoa(i),
			Lino:           LNOPerValidator,
			ResetKey:       secp256k1.GenPrivKey().PubKey(),
			TransactionKey: secp256k1.GenPrivKey().PubKey(),
			AppKey:         secp256k1.GenPrivKey().PubKey(),
			IsValidator:    true,
			ValPubKey:      secp256k1.GenPrivKey().PubKey(),
		}
		marshalJson, err := wire.MarshalJSONIndent(cdc, genesisAcc)
		assert.Nil(t, err)
		appGenTxs = append(appGenTxs, json.RawMessage(marshalJson))
	}
	appState, err := LinoBlockchainGenState(cdc, appGenTxs)
	assert.Nil(t, err)

	genesisState := new(GenesisState)
	if err := cdc.UnmarshalJSON(appState, genesisState); err != nil {
		panic(err)
	}
	for i, gacc := range genesisState.Accounts {
		assert.Equal(t, gacc.Name, "validator"+strconv.Itoa(i+1))
		assert.Equal(t, gacc.Lino, LNOPerValidator)
	}
	assert.Equal(t, 1, len(genesisState.Developers))
	assert.Equal(t, 1, len(genesisState.Infra))
}
