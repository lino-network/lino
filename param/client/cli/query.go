package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	// linotypes "github.com/lino-network/lino/types"
	types "github.com/lino-network/lino/param"
	"github.com/lino-network/lino/utils"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the param module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(client.GetCommands(
		getCmdAll(cdc),
	)...)
	return cmd
}

// GetCmdAll -
func getCmdAll(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "all print all blockchain parameters",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			all := []struct {
				path string
				t    interface{}
			}{
				{types.QueryAllocationParam, &types.GlobalAllocationParam{}},
				{types.QueryInfraInternalAllocationParam, &types.InfraInternalAllocationParam{}},
				{types.QueryDeveloperParam, &types.DeveloperParam{}},
				{types.QueryVoteParam, &types.VoteParam{}},
				{types.QueryProposalParam, &types.ProposalParam{}},
				{types.QueryValidatorParam, &types.ValidatorParam{}},
				{types.QueryBandwidthParam, &types.BandwidthParam{}},
				{types.QueryAccountParam, &types.AccountParam{}},
				{types.QueryPostParam, &types.PostParam{}},
				{types.QueryReputationParam, &types.ReputationParam{}},
			}
			for _, v := range all {
				uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, v.path)
				err := utils.CLIQueryJSONPrint(cdc, uri, nil,
					func() interface{} { return v.t })
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}
