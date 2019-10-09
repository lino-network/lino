package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	types "github.com/lino-network/lino/x/validator"
	model "github.com/lino-network/lino/x/validator/model"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the validator module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(client.GetCommands(
		getCmdShow(cdc),
		getCmdList(cdc),
		getCmdVoteInfo(cdc),
	)...)
	return cmd
}

// GetCmdShow -
func getCmdShow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "show <username>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			user := linotypes.AccountKey(args[0])
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryValidator, user)
			rst := model.Validator{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdList -
func getCmdList(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list all validators",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidatorList)
			rst := model.ValidatorList{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdVoteInfo
func getCmdVoteInfo(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "voteinfo",
		Short: "voteinfo <username>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			user := linotypes.AccountKey(args[0])
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryElectionVoteList, user)
			rst := model.ElectionVoteList{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}
