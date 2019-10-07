package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/price/model"
	"github.com/lino-network/lino/x/price/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the price module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(client.GetCommands(
		getCmdCurrent(cdc),
		getCmdHistory(cdc),
	)...)
	return cmd
}

// GetCmdCurrent -
func getCmdCurrent(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "current",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPriceCurrent)
			rst := linotypes.NewMiniDollar(0)
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdHistory -
func getCmdHistory(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "history",
		Short: "history",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryPriceHistory)
			rst := make([]model.FeedHistory, 0)
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}
