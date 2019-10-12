package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	// linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/proposal/model"
	"github.com/lino-network/lino/x/proposal/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the proposal module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(client.GetCommands(
		getCmdOngoing(cdc),
		getCmdExpired(cdc),
	)...)
	return cmd
}

// GetCmdListOngoing -
func getCmdOngoing(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ongoing",
		Short: "ongoing <proposal_id> print the ongoing proposals",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pid := args[0]
			uri := fmt.Sprintf(
				"custom/%s/%s/%s", types.QuerierRoute, types.QueryOngoingProposal, pid)
			var rst model.Proposal
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdExpired -
func getCmdExpired(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "expired",
		Short: "expired <proposal_id> print the expired proposal",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			pid := args[0]
			uri := fmt.Sprintf(
				"custom/%s/%s/%s", types.QuerierRoute, types.QueryExpiredProposal, pid)
			var rst model.Proposal
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}
