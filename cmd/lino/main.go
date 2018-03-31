package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/lino-network/lino/app"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/words"

	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tmtypes "github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

// linoCmd is the entry point for this binary
var (
	linoCmd = &cobra.Command{
		Use:   "lino",
		Short: "Lino Blockchain (server)",
	}
)

// defaultOptions sets up the app_options for the
// default genesis file
func defaultOptions(args []string) (json.RawMessage, string, cmn.HexBytes, error) {
	pubKey, secret, err := generateCoinKey()
	if err != nil {
		return nil, "", nil, err
	}
	fmt.Println("Secret phrase to access coins:")
	fmt.Println(secret)

	config, err := tcmd.ParseConfig()
	if err != nil {
		return nil, "", nil, err
	}
	// private validator
	privValFile := config.PrivValidatorFile()
	var privValidator *tmtypes.PrivValidatorFS
	if cmn.FileExists(privValFile) {
		privValidator = tmtypes.LoadPrivValidatorFS(privValFile)
	} else {
		privValidator = tmtypes.GenPrivValidatorFS(privValFile)
		privValidator.Save()
	}

	pubKeyBytes, err := json.Marshal(*pubKey)
	if err != nil {
		return nil, "", nil, err
	}
	valPubKeyBytes, err := json.Marshal(privValidator.PubKey)
	if err != nil {
		return nil, "", nil, err
	}

	opts := fmt.Sprintf(`{
	      "accounts": [{
	        "coin": [
	          {
	            "amount": 10000000000
	          }
	        ],
	        "name": "Lino",
	        "pub_key": %s,
	        "validator_pub_key": %s
	      }],
	      "global_state": {
	      	"total_lino":
	          {
	            "amount": 10000000000
	          },
	      	"growth_rate": {
	      		"num": 98,
	      		"denum": 1000
	      	},
	      	"infra_allocation": {
	      		"num": 20,
	      		"denum": 100
	      	},
	      	"content_creator_allocation": {
	      		"num": 55,
	      		"denum": 100
	      	},
	      	"developer_allocation": {
	      		"num": 20,
	      		"denum": 100
	      	},
	      	"validator_allocation": {
	      		"num": 5,
	      		"denum": 100
	      	},
	      	"consumption_friction_rate": {
	      		"num": 1,
	      		"denum": 100
	      	},
	      	"freezing_period_hr": 168
	      }
	    }`, pubKeyBytes, valPubKeyBytes)
	fmt.Println("default address:", pubKey.Address())
	return json.RawMessage(opts), secret, pubKey.Address(), nil
}

// generate Lino application
func generateApp(rootDir string, logger log.Logger) (abci.Application, error) {
	db, err := dbm.NewGoLevelDB("lino", rootDir)
	if err != nil {
		return nil, err
	}
	lb := app.NewLinoBlockchain(logger, db)
	return lb, nil
}

func main() {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).
		With("module", "main")

	linoCmd.AddCommand(
		server.InitCmd(defaultOptions, logger),
		server.StartCmd(generateApp, logger),
		server.UnsafeResetAllCmd(logger),
		version.VersionCmd,
	)

	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.lino")
	executor := cli.PrepareBaseCmd(linoCmd, "BC", rootDir)
	executor.Execute()
}

func generateCoinKey() (*crypto.PubKey, string, error) {
	// construct an in-memory key store
	codec, err := words.LoadCodec("english")
	if err != nil {
		return nil, "", err
	}
	keybase := keys.New(
		dbm.NewMemDB(),
		codec,
	)

	// generate a private key, with recovery phrase
	info, secret, err := keybase.Create("name", "pass", keys.AlgoEd25519)
	if err != nil {
		return nil, "", err
	}

	return &info.PubKey, secret, nil
}
