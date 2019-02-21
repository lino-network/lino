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
	return app.NewLinoBlockchain(logger, db, traceStore,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))
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

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	executor := cli.PrepareBaseCmd(rootCmd, "BC", app.DefaultNodeHome)
	executor.Execute()
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, traceStore io.Writer,
	_ int64, _ bool, _ []string) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	lb := app.NewLinoBlockchain(logger, db, traceStore)
	return lb.ExportAppStateAndValidators()
}
