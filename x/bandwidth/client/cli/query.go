package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	// linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/bandwidth/model"
	"github.com/lino-network/lino/x/bandwidth/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the bandwidth module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(client.GetCommands(
		getCmdInfo(cdc),
		getCmdBlock(cdc),
		getCmdApp(cdc),
	)...)
	return cmd
}

// GetCmdInfo -
func getCmdInfo(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "info",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryBandwidthInfo)
			rst := model.BandwidthInfo{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdBlock -
func getCmdBlock(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "block",
		Short: "block",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryBlockInfo)
			rst := model.BlockInfo{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdApp -
func getCmdApp(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "app",
		Short: "app <app>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s",
				types.QuerierRoute, types.QueryAppBandwidthInfo, app)
			rst := model.AppBandwidthInfo{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}
