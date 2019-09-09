package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	txutils "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	// server "github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/lino-network/lino/app"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	acccmd "github.com/lino-network/lino/x/account/commands"
	developercmd "github.com/lino-network/lino/x/developer/commands"
	infracmd "github.com/lino-network/lino/x/infra/commands"
	postcmd "github.com/lino-network/lino/x/post/client/cli"
	proposalcmd "github.com/lino-network/lino/x/proposal/commands"
	validatorcmd "github.com/lino-network/lino/x/validator/commands"
	delegatecmd "github.com/lino-network/lino/x/vote/commands/delegate"
	delegationcmd "github.com/lino-network/lino/x/vote/commands/delegate"
	votecmd "github.com/lino-network/lino/x/vote/commands/vote"
)

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()
	// add standard rpc, and tx commands
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint state querying subcommands",
	}
	tendermintCmd.AddCommand(
		rpc.BlockCommand(),
		rpc.ValidatorCommand(cdc),
	)

	// XXX(yumin): before major-update-1, it's tx.AddCommands.
	tendermintCmd.AddCommand(
		// txutils.SearchTxCmd(cdc),
		txutils.QueryTxCmd(cdc),
	)

	advancedCmd := &cobra.Command{
		Use:   "advanced",
		Short: "Advanced subcommands",
	}

	// TODO(yumin): ServeCommand now requires a callback func to register all routes
	// do we use this lcd or not?
	// advancedCmd.AddCommand(
	// 	tendermintCmd,
	// 	lcd.ServeCommand(cdc),
	// )

	linocliCmd.AddCommand(
		advancedCmd,
		client.LineBreak,
	)

	linocliCmd.AddCommand(
		client.PostCommands(
			acccmd.RegisterTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			acccmd.RecoverTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			acccmd.TransferTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			postcmd.PostTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			postcmd.UpdatePostTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			postcmd.DeletePostTxCmd(cdc),
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
			votecmd.WithdrawVoterTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			votecmd.GetVoterCmd(types.VoteKVStoreKey, cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			proposalcmd.GetOngoingProposalCmd(types.VoteKVStoreKey, cdc),
		)...)

	linocliCmd.AddCommand(
		client.GetCommands(
			proposalcmd.GetExpiredProposalCmd(types.VoteKVStoreKey, cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			proposalcmd.VoteProposalTxCmd(cdc),
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
			developercmd.GrantPermissionTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			developercmd.RevokePermissionTxCmd(cdc),
		)...)
	linocliCmd.AddCommand(
		client.PostCommands(
			developercmd.DeveloperUpdateTxCmd(cdc),
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
			acccmd.GetAccountsCmd(types.AccountKVStoreKey, cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			postcmd.GetPostCmd(types.PostKVStoreKey, cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			postcmd.GetPostsCmd(types.PostKVStoreKey, cdc),
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
			validatorcmd.GetValidatorsCmd(types.ValidatorKVStoreKey, cdc),
		)...)
	linocliCmd.AddCommand(
		client.GetCommands(
			validatorcmd.GetValidatorCmd(types.ValidatorKVStoreKey, cdc),
		)...)

	// add proxy, version and key info
	linocliCmd.AddCommand(
		keys.Commands(),
		client.LineBreak,
		// server.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(linocliCmd, "BC", os.ExpandEnv("$HOME/.linocli"))
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}
