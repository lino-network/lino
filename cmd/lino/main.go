package main

import (
	"encoding/json"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/types"
)

// generate Lino application
func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	app := app.NewLinoBlockchain(
		logger,
		db,
		traceStore,
		// XXX(yumin): previously we use
		// "syncable": PruneSyncable = NewPruningOptions(100, 10000)
		// which means every (10000 * block_time) seconds, a state is kept.
		// If block_time is around 3 seconds, then every ~8.33 hours, a state
		// is kept. When height is 4M, there are about 400 copies in db. Even if
		// it's an immutable iavl, an early state may be just a full copy of the state, as
		// it may be totally different from current and following kepted states.
		// Plus state-sync is not supported for now, we set it to 400000 here.
		baseapp.SetPruning(storetypes.NewPruningOptions(100, 400000)),
	)
	return app
}

func main() {
	cobra.EnableCommandSorting = false

	types.ConfigAndSealCosmosSDKAddress()

	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "lino",
		Short:             "Lino Blockchain (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(app.VersionCmd())

	rootCmd.AddCommand(app.InitCmd(ctx, cdc))

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	executor := cli.PrepareBaseCmd(rootCmd, "BC", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, traceStore io.Writer,
	_ int64, _ bool, _ []string) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	lb := app.NewLinoBlockchain(logger, db, traceStore)
	return lb.ExportAppStateAndValidators()
}
