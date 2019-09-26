package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	types "github.com/lino-network/lino/x/reputation"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the reputation module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(client.GetCommands(
		getCmdShow(cdc),
		getCmdDetail(cdc),
	)...)
	return cmd
}

// GetCmdShow -
func getCmdShow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "show <username>",
		Short: "show <username>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryReputation, username)
			rst := linotypes.MiniDollar{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdDetail -
func getCmdDetail(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "detail <username>",
		Short: "detail <username>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryDetails, username)
			return utils.CLIQueryStrPrint(cdc, uri, nil)
		},
	}
}
