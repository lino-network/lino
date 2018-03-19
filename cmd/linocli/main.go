package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/tendermint/tmlibs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/commands"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/commands"
	acccmd "github.com/lino-network/lino/tx/account/commands"

	"github.com/cosmos/cosmos-sdk/examples/basecoin/app"
	"github.com/cosmos/cosmos-sdk/examples/basecoin/types"
)

// linocliCmd is the entry point for this binary
var (
	linocliCmd = &cobra.Command{
		Use:   "linocli",
		Short: "Lino Blockchain light-client",
	}
)

func todoNotImplemented(_ *cobra.Command, _ []string) error {
	return errors.New("TODO: Command not yet implemented")
}

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()

	// add standard rpc, and tx commands
	rpc.AddCommands(linocliCmd)
	linocliCmd.AddCommand(client.LineBreak)
	tx.AddCommands(linocliCmd, cdc)
	linocliCmd.AddCommand(client.LineBreak)

	// TODO(Lino): Customize our own command
	// add query/post commands (custom to binary)
	linocliCmd.AddCommand(
	 	client.GetCommands(
			authcmd.GetAccountCmd("main", cdc, types.GetParseAccount(cdc)),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			acccmd.RegisterTxCmd(cdc),
		)...)
	// linocliCmd.AddCommand(
	// 	client.PostCommands(
	// 		coolcmd.SetTrendTxCmd(cdc),
	// 	)...)

	// add proxy, version and key info
	linocliCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(linocliCmd, "BC", os.ExpandEnv("$HOME/.linocli"))
	executor.Execute()
}
