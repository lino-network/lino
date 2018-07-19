package main

import (
	"encoding/json"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/lino-network/lino/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"
)

// // defaultOptions sets up the app_options for the
// // default genesis file
// func defaultAppState(args []string, addr sdk.Address, coinDenom string) (json.RawMessage, error) {
// 	pubKey, secret, err := generateCoinKey()
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("Secret phrase to access coins:")
// 	fmt.Println(secret)
// 	fmt.Println("Init address:")
// 	fmt.Println(pubKey.Address())

// 	fmt.Println(hex.EncodeToString(privValidator.PrivKey.Bytes()))
// 	result, err := genesis.GetDefaultGenesis(pubKey, privValidator.PubKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return json.RawMessage(result), nil
// }

// generate Lino application
func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewLinoBlockchain(logger, db, traceStore, baseapp.SetPruning(viper.GetString("pruning")))
}

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "lino",
		Short:             "Lino Blockchain (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, app.LinoBlockchainInit(),
		server.ConstructAppCreator(newApp, "lino"),
		server.ConstructAppExporter(exportAppStateAndTMValidators, "lino"))

	executor := cli.PrepareBaseCmd(rootCmd, "BC", app.DefaultNodeHome)
	executor.Execute()
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	lb := app.NewLinoBlockchain(logger, db, traceStore)
	return lb.ExportAppStateAndValidators()
}
