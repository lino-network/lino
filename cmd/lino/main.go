package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/genesis"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/words"

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

	pubKeyBytes, err := json.Marshal(*pubKey)
	if err != nil {
		return nil, "", nil, err
	}

	opts := fmt.Sprintf(`{
	      "accounts": [{
	        "name": "Lino",
	        "lino": 10000000000,
	        "pub_key": %s,
	      }]
	    }`, pubKeyBytes)
	fmt.Println("default address:", pubKey.Address())

	genesisState := new(genesis.GenesisState)

	//err := oldwire.UnmarshalJSON(stateJSON, genesisState)
	err = json.Unmarshal(json.RawMessage(opts), genesisState)
	if err != nil {
		panic(err) // TODO(Cosmos) https://github.com/cosmos/cosmos-sdk/issues/468
	}
	fmt.Println(genesisState)
	return json.RawMessage(opts), secret, pubKey.Address(), nil
}

// generate Lino application
func generateApp(rootDir string, logger log.Logger) (abci.Application, error) {
	dbAcc, err := dbm.NewGoLevelDB("LinoBlockchain-acc", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbPost, err := dbm.NewGoLevelDB("LinoBlockchain-post", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbs := map[string]dbm.DB{
		"acc":  dbAcc,
		"post": dbPost,
	}
	lb := app.NewLinoBlockchain(logger, dbs)
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
