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
	acccmd "github.com/lino-network/lino/tx/account/commands"
	postcmd "github.com/lino-network/lino/tx/post/commands"
	registercmd "github.com/lino-network/lino/tx/register/commands"

	"github.com/lino-network/lino/app"
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

	linocliCmd.AddCommand(
		client.PostCommands(
			registercmd.RegisterTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			acccmd.TransferTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			acccmd.FollowTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			postcmd.PostTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			postcmd.LikeTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			postcmd.DonateTxCmd(cdc),
		)...)

	linocliCmd.AddCommand(
		client.GetCommands(
			acccmd.GetBankCmd("account", cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			acccmd.GetAccountCmd("account", cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			postcmd.GetPostCmd("post", cdc),
		)...)

	// add proxy, version and key info
	linocliCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(linocliCmd, "BC", os.ExpandEnv("$HOME/.linocli"))
	executor.Execute()
}
