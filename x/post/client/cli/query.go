package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/post/model"
	"github.com/lino-network/lino/x/post/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the post module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(client.GetCommands(
		getCmdInfo(cdc),
		utils.SimpleQueryCmd(
			"cw",
			"cw prints the consumption competition metadata, unit: miniDollar",
			types.QuerierRoute, types.QueryConsumptionWindow, 0, &linotypes.MiniDollar{})(cdc),
	)...)
	return cmd
}

// GetCmdInfo -
func getCmdInfo(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "info <permlink>",
		Short: "info <permlink>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			permlink := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryPostInfo, permlink)
			rst := model.Post{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}
