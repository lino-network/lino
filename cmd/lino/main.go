package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/genesis"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/keys"
	"github.com/tendermint/go-crypto/keys/words"
	"github.com/tendermint/tmlibs/cli"
	"github.com/tendermint/tmlibs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	pvm "github.com/tendermint/tendermint/types/priv_validator"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
)

// linoCmd is the entry point for this binary
var (
	context = server.NewDefaultContext()
	linoCmd = &cobra.Command{
		Use:               "lino",
		Short:             "Lino Blockchain (node)",
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
	fmt.Println("Init address:")
	fmt.Println(pubKey.Address())

	// private validator
	privValFile := context.Config.PrivValidatorFile()
	var privValidator *pvm.FilePV
	if cmn.FileExists(privValFile) {
		privValidator = pvm.LoadFilePV(privValFile)
	} else {
		privValidator = pvm.GenFilePV(privValFile)
		privValidator.Save()
	}

	fmt.Println(hex.EncodeToString(privValidator.PrivKey.Bytes()))
	result, err := genesis.GetDefaultGenesis(pubKey, privValidator.PubKey)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(result), nil
}

// generate Lino application
func generateApp(rootDir string, logger log.Logger) (abci.Application, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := dbm.NewGoLevelDB("linoblockchain", dataDir)
	if err != nil {
		return nil, err
	}
	lb := app.NewLinoBlockchain(logger, db)
	return lb, nil
}

func main() {
	server.AddCommands(linoCmd, defaultAppState, generateApp, context)
	// prepare and add flags
	rootDir := os.ExpandEnv("$HOME/.lino")
	executor := cli.PrepareBaseCmd(linoCmd, "BC", rootDir)
	executor.Execute()
}

func generateCoinKey() (crypto.PubKey, string, error) {
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

	return info.PubKey, secret, nil
}
