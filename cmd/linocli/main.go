package main

import (
	"os"

	"github.com/spf13/cobra"

	keycmd "github.com/tendermint/go-crypto/cmd"
	"github.com/tendermint/light-client/commands"
	"github.com/tendermint/light-client/commands/proofs"
	"github.com/tendermint/light-client/commands/proxy"
	rpccmd "github.com/tendermint/light-client/commands/rpc"
	"github.com/tendermint/light-client/commands/seeds"
	"github.com/tendermint/light-client/commands/txs"
	"github.com/tendermint/tmlibs/cli"

	bcmd "github.com/lino-network/lino/cmd/linocli/commands"
	coincmd "github.com/lino-network/lino/cmd/linocoin/commands"
)

// BaseCli represents the base command when called without any subcommands
var BaseCli = &cobra.Command{
	Use:   "linocli",
	Short: "Light client for tendermint",
	Long: `Linocli is an version of tmcli including custom logic to
present a nice (not raw hex) interface to the linocoin blockchain structure.

This is a useful tool, but also serves to demonstrate how one can configure
tmcli to work for any custom abci app.
`,
}

func main() {
	commands.AddBasicFlags(BaseCli)

	// Prepare queries
	pr := proofs.RootCmd
	// These are default parsers, but optional in your app (you can remove key)
	pr.AddCommand(proofs.TxCmd)
	pr.AddCommand(proofs.KeyCmd)
	pr.AddCommand(bcmd.AccountQueryCmd)
	pr.AddCommand(bcmd.PostQueryCmd)

	// you will always want this for the base send command
	proofs.TxPresenters.Register("base", bcmd.BaseTxPresenter{})
	tr := txs.RootCmd
	tr.AddCommand(bcmd.SendTxCmd)
	tr.AddCommand(bcmd.PostTxCmd)

	// Set up the various commands to use
	BaseCli.AddCommand(
		commands.InitCmd,
		commands.ResetCmd,
		keycmd.RootCmd,
		seeds.RootCmd,
		rpccmd.RootCmd,
		pr,
		tr,
		proxy.RootCmd,
		coincmd.VersionCmd,
		bcmd.AutoCompleteCmd,
	)

	cmd := cli.PrepareMainCmd(BaseCli, "BC", os.ExpandEnv("$HOME/.basecli"))
	cmd.Execute()
}
