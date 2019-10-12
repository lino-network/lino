package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/proposal/types"
)

const (
	FlagProposalID = "proposal-id"
	FlagResult     = "result"
	FlagLink       = "link"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "vote tx subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(client.PostCommands(
		GetCmdVote(cdc),
	)...)

	return cmd
}

// GetCmdVote -
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "vote <voter> --proposal-id <id> --result=true/false",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			voter := args[0]
			id := viper.GetInt64(FlagProposalID)
			result := viper.GetBool(FlagResult)
			msg := types.NewVoteProposalMsg(voter, id, result)
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().Int64(FlagProposalID, -1, "proposal id")
	cmd.Flags().Bool(FlagResult, true, "vote result")
	_ = cmd.MarkFlagRequired(FlagProposalID)
	_ = cmd.MarkFlagRequired(FlagResult)
	return cmd
}
