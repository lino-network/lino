package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/client"
	acccmd "github.com/lino-network/lino/tx/account/commands"
	developercmd "github.com/lino-network/lino/tx/developer/commands"
	infracmd "github.com/lino-network/lino/tx/infra/commands"
	postcmd "github.com/lino-network/lino/tx/post/commands"
	proposalcmd "github.com/lino-network/lino/tx/proposal/commands"
	registercmd "github.com/lino-network/lino/tx/register/commands"
	validatorcmd "github.com/lino-network/lino/tx/validator/commands"
	delegatecmd "github.com/lino-network/lino/tx/vote/commands/delegate"
	delegationcmd "github.com/lino-network/lino/tx/vote/commands/delegate"
	votecmd "github.com/lino-network/lino/tx/vote/commands/vote"
	"github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/tendermint/tmlibs/cli"
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
		client.PostCommands(
			validatorcmd.DepositValidatorTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			validatorcmd.WithdrawTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			validatorcmd.RevokeTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			delegationcmd.RevokeDelegateTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			delegationcmd.DelegateTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			delegationcmd.WithdrawDelegateTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			delegatecmd.GetDelegationCmd(types.VoteKVStoreKey, cdc),
		)...)

	linocliCmd.AddCommand(
		client.PostCommands(
			votecmd.DepositVoterTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			votecmd.RevokeVoterTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			votecmd.VoteTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			votecmd.WithdrawVoterTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			votecmd.GetVoterCmd(types.VoteKVStoreKey, cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			proposalcmd.GetProposalCmd(types.VoteKVStoreKey, cdc),
		)...)

	linocliCmd.AddCommand(
		client.GetCommands(
			proposalcmd.GetProposalListCmd(types.VoteKVStoreKey, cdc),
		)...)

	linocliCmd.AddCommand(
		client.GetCommands(
			votecmd.GetVoteCmd(types.VoteKVStoreKey, cdc),
		)...)

	linocliCmd.AddCommand(
		client.PostCommands(
			infracmd.ProviderReportTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			developercmd.DeveloperRegisterTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			developercmd.DeveloperRevokeTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			developercmd.GrantDeveloperTxCmd(cdc),
		)...)

	linocliCmd.AddCommand(
		client.GetCommands(
			acccmd.GetBankCmd(types.AccountKVStoreKey, cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			acccmd.GetAccountCmd(types.AccountKVStoreKey, cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			postcmd.GetPostCmd(types.PostKVStoreKey, cdc),
		)...)

	linocliCmd.AddCommand(
		client.GetCommands(
			infracmd.GetInfraProviderCmd(types.InfraKVStoreKey, cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			infracmd.GetInfraProvidersCmd(types.InfraKVStoreKey, cdc),
		)...)

	linocliCmd.AddCommand(
		client.GetCommands(
			developercmd.GetDeveloperCmd(types.DeveloperKVStoreKey, cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			developercmd.GetDevelopersCmd(types.DeveloperKVStoreKey, cdc),
		)...)

	linocliCmd.AddCommand(
		client.GetCommands(
			validatorcmd.GetValidatorsCmd(types.ValidatorKVStoreKey, cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			validatorcmd.GetValidatorCmd(types.ValidatorKVStoreKey, cdc),
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
