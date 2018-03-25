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
	oldwire "github.com/tendermint/go-wire"
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
func defaultOptions(args []string) (json.RawMessage, error) {
	pubKey, secret, err := generateCoinKey()
	if err != nil {
		return nil, err
	}
	fmt.Println("Secret phrase to access coins:")
	fmt.Println(secret)

	b, _ := oldwire.MarshalJSON(*pubKey)
	opts := fmt.Sprintf(`{
	      "accounts": [{
	        "coins": [
	          {
	            "denom": "lino",
	            "amount": 10000000000
	          }
	        ],
	        "name": "Lino",
	        "pub_key": %s
	      }]
	    }`, b)
	fmt.Println("default address:", pubKey.Address())
	return json.RawMessage(opts), nil
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
