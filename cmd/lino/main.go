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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/genesis"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/words"

	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tmtypes "github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

// linoCmd is the entry point for this binary
var (
	context = server.NewDefaultContext()
	linoCmd = &cobra.Command{
		Use:               "lino",
		Short:             "Lino Blockchain (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(context),
	}
)

// defaultOptions sets up the app_options for the
// default genesis file
func defaultAppState(args []string, addr sdk.Address, coinDenom string) (json.RawMessage, error) {
	pubKey, secret, err := generateCoinKey()
	if err != nil {
		return nil, err
	}
	fmt.Println("Secret phrase to access coins:")
	fmt.Println(secret)

	config, err := tcmd.ParseConfig()
	if err != nil {
		return nil, err
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

	result, err := genesis.GetDefaultGenesis(*pubKey, privValidator.PubKey)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(result), nil
}

// generate Lino application
func generateApp(rootDir string, logger log.Logger) (abci.Application, error) {
	dbMain, err := dbm.NewGoLevelDB("LinoBlockchain-main", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbAcc, err := dbm.NewGoLevelDB("LinoBlockchain-acc", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbPost, err := dbm.NewGoLevelDB("LinoBlockchain-post", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbVal, err := dbm.NewGoLevelDB("LinoBlockchain-val", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbVote, err := dbm.NewGoLevelDB("LinoBlockchain-vote", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbInfra, err := dbm.NewGoLevelDB("LinoBlockchain-infra", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbDeveloper, err := dbm.NewGoLevelDB("LinoBlockchain-developer", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbGlobal, err := dbm.NewGoLevelDB("LinoBlockchain-global", filepath.Join(rootDir, "data"))
	if err != nil {
		return nil, err
	}
	dbs := map[string]dbm.DB{
		"main":      dbMain,
		"acc":       dbAcc,
		"post":      dbPost,
		"val":       dbVal,
		"vote":      dbVote,
		"infra":     dbInfra,
		"developer": dbDeveloper,
		"global":    dbGlobal,
	}
	lb := app.NewLinoBlockchain(logger, dbs)
	return lb, nil
}

func main() {
	server.AddCommands(linoCmd, defaultAppState, generateApp, context)
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
