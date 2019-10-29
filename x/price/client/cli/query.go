package cli

import (
	// "fmt"

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
		utils.SimpleQueryCmd(
			"current",
			"current",
			types.QuerierRoute, types.QueryPriceCurrent,
			0, &linotypes.MiniDollar{})(cdc),
		utils.SimpleQueryCmd(
			"history",
			"history",
			types.QuerierRoute, types.QueryPriceHistory,
			0, &([]model.FeedHistory{}))(cdc),
		utils.SimpleQueryCmd(
			"last-feed <username>",
			"last-feed <username>",
			types.QuerierRoute, types.QueryLastFeed,
			1, &model.FedPrice{})(cdc),
	)...)
	return cmd
}
