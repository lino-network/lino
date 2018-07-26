package app

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/types"
	"github.com/spf13/pflag"
	crypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	flagName       = "name"
	flagClientHome = "home-client"
	flagOWK        = "owk"
)

// get app init parameters for server init command
func LinoBlockchainInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)

	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(flagName, "", "validator moniker, required")
	fsAppGenTx.String(flagClientHome, DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(flagOWK, false, "overwrite the accounts created")

	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         LinoBlockchainGenTx,
		AppGenState:      LinoBlockchainGenState,
	}
}

// genesis state for blockchain
type GenesisState struct {
	Accounts   []GenesisAccount       `json:"accounts"`
	Developers []GenesisAppDeveloper  `json:"developers"`
	Infra      []GenesisInfraProvider `json:"infra"`
}

// genesis account will get coin to the address and register user
// if genesis account is validator, it will be added to validator list automatically
type GenesisAccount struct {
	Name           string        `json:"name"`
	Lino           types.LNO     `json:"lino"`
	ResetKey       crypto.PubKey `json:"reset_key"`
	TransactionKey crypto.PubKey `json:"transaction_key"`
	AppKey         crypto.PubKey `json:"app_key"`
	IsValidator    bool          `json:"is_validator"`
	ValPubKey      crypto.PubKey `json:"validator_pub_key"`
}

// register developer in genesis phase
type GenesisAppDeveloper struct {
	Name        string    `json:"name"`
	Deposit     types.LNO `json:"deposit"`
	Website     string    `json:"web_site"`
	Description string    `json:"description"`
	AppMetaData string    `json:"app_meta_data"`
}

// register infra provider in genesis phase
type GenesisInfraProvider struct {
	Name string `json:"name"`
}

func LinoBlockchainGenTx(cdc *wire.Codec, pk crypto.PubKey, genTxConfig config.GenTx) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {
	resetPriv := secp256k1.GenPrivKey()
	transactionPriv := secp256k1.GenPrivKey()
	appPriv := secp256k1.GenPrivKey()

	fmt.Println("reset private key is:", strings.ToUpper(hex.EncodeToString(resetPriv.Bytes())))
	fmt.Println("transaction private key is:", strings.ToUpper(hex.EncodeToString(transactionPriv.Bytes())))
	fmt.Println("app private key is:", strings.ToUpper(hex.EncodeToString(appPriv.Bytes())))

	totalLino := "10000000000"
	genesisAcc := GenesisAccount{
		Name:           "lino",
		Lino:           totalLino,
		ResetKey:       resetPriv.PubKey(),
		TransactionKey: transactionPriv.PubKey(),
		AppKey:         appPriv.PubKey(),
		IsValidator:    true,
		ValPubKey:      pk,
	}

	var bz []byte
	bz, err = wire.MarshalJSONIndent(cdc, genesisAcc)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)

	validator = tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  1000,
	}
	return
}

// default genesis file, only have one genesis account
func LinoBlockchainGenState(cdc *wire.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {
	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}

	// totalLino := "10000000000"
	genesisState := GenesisState{
		Accounts:   []GenesisAccount{},
		Developers: []GenesisAppDeveloper{},
		Infra:      []GenesisInfraProvider{},
	}

	for _, genesisAccRaw := range appGenTxs {
		var genesisAcc GenesisAccount
		err = cdc.UnmarshalJSON(genesisAccRaw, &genesisAcc)
		if err != nil {
			return
		}
		genesisState.Accounts = append(genesisState.Accounts, genesisAcc)
	}
	genesisAppDeveloper := GenesisAppDeveloper{
		Name:        "lino",
		Deposit:     "1000000",
		Website:     "https://lino.network/",
		Description: "",
		AppMetaData: "",
	}
	genesisState.Developers = append(genesisState.Developers, genesisAppDeveloper)
	genesisInfraProvider := GenesisInfraProvider{
		Name: "lino",
	}
	genesisState.Infra = append(genesisState.Infra, genesisInfraProvider)

	appState, err = wire.MarshalJSONIndent(cdc, genesisState)
	return
}
