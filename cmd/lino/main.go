package main

import (
	"encoding/json"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/lino-network/lino/app"
)

// generate Lino application
func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	app := app.NewLinoBlockchain(logger, db, traceStore,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))
	// after upgrade-1, lino needs to starts
	app.SetImportRequired(true)
	return app
}

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "lino",
		Short:             "Lino Blockchain (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(app.InitCmd(ctx, cdc))

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	executor := cli.PrepareBaseCmd(rootCmd, "BC", app.DefaultNodeHome)
	executor.Execute()
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, traceStore io.Writer,
	_ int64, _ bool, _ []string) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	lb := app.NewLinoBlockchain(logger, db, traceStore)
	return lb.ExportAppStateAndValidators()
}
