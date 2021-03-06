package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	app "github.com/lino-network/lino/app"
	linoclient "github.com/lino-network/lino/client"
	blockcli "github.com/lino-network/lino/client/blockchain"
	paramcli "github.com/lino-network/lino/param/client/cli"
	"github.com/lino-network/lino/types"
	acccli "github.com/lino-network/lino/x/account/client/cli"
	bwcli "github.com/lino-network/lino/x/bandwidth/client/cli"
	devcli "github.com/lino-network/lino/x/developer/client/cli"
	globalcli "github.com/lino-network/lino/x/global/client/cli"
	postcli "github.com/lino-network/lino/x/post/client/cli"
	pricecli "github.com/lino-network/lino/x/price/client/cli"

	// proposalcli "github.com/lino-network/lino/x/proposal/client/cli"
	repcli "github.com/lino-network/lino/x/reputation/client/cli"
	validatorcli "github.com/lino-network/lino/x/validator/client/cli"
	votecli "github.com/lino-network/lino/x/vote/client/cli"
)

// linocliCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "linocli",
		Short: "Lino Blockchain CLI",
	}
	DefaultCLIHome = os.ExpandEnv("$HOME/.linocli")
)

func main() {
	cobra.EnableCommandSorting = false

	types.ConfigAndSealCosmosSDKAddress()

	cdc := app.MakeCodec()

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")
	// rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
	// 	return initConfig(rootCmd)
	// }

	// Construct Root Command
	rootCmd.AddCommand(
		app.VersionCmd(),
		rpc.StatusCommand(),
		client.ConfigCmd(DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		client.LineBreak,
		linoclient.GetNowCmd(cdc),
		linoclient.GetGenAddrCmd(),
		linoclient.GetAddrOfCmd(),
		linoclient.GetPubKeyOfCmd(),
		linoclient.GetEncryptPrivKey(),
		client.LineBreak,
	)

	executor := cli.PrepareMainCmd(rootCmd, "NS", DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		blockcli.GetQueryCmd(cdc),
		client.LineBreak,
		devcli.GetQueryCmd(cdc),
		acccli.GetQueryCmd(cdc),
		postcli.GetQueryCmd(cdc),
		// proposalcli.GetQueryCmd(cdc),
		validatorcli.GetQueryCmd(cdc),
		globalcli.GetQueryCmd(cdc),
		bwcli.GetQueryCmd(cdc),
		paramcli.GetQueryCmd(cdc),
		repcli.GetQueryCmd(cdc),
		votecli.GetQueryCmd(cdc),
		pricecli.GetQueryCmd(cdc),
	)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		linoclient.GetCmdBroadcast(cdc),
		client.LineBreak,
		devcli.GetTxCmd(cdc),
		acccli.GetTxCmd(cdc),
		postcli.GetTxCmd(cdc),
		// proposalcli.GetTxCmd(cdc),
		validatorcli.GetTxCmd(cdc),
		votecli.GetTxCmd(cdc),
		pricecli.GetTxCmd(cdc),
		client.LineBreak,
	)

	return txCmd
}
