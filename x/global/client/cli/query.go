package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	// linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/global/model"
	types "github.com/lino-network/lino/x/global/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the global module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(client.GetCommands(
		// 	getCmdListEvents(cdc),
		// 	getCmdMeta(cdc),
		// 	getCmdInflationPool(cdc),
		// 	getCmdConsumption(cdc),
		// 	getCmdTPS(cdc),
		getCmdTime(cdc),
		// getCmdStakeStats(cdc),
	)...)
	return cmd
}

// // GetCmdListEvents -
// func getCmdListEvents(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "list-events",
// 		Short: "list-events <unix_time>",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			time := args[0]
// 			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryTimeEventList, time)
// 			rst := linotypes.TimeEventList{}
// 			return utils.CLIQueryJSONPrint(cdc, uri, nil,
// 				func() interface{} { return &rst })
// 		},
// 	}
// }

// // GetCmdMeta -
// func getCmdMeta(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "meta",
// 		Short: "meta",
// 		Args:  cobra.ExactArgs(0),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGlobalMeta)
// 			rst := model.GlobalMeta{}
// 			return utils.CLIQueryJSONPrint(cdc, uri, nil,
// 				func() interface{} { return &rst })
// 		},
// 	}
// }

// // GetCmdInflationPool -
// func getCmdInflationPool(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "inflation",
// 		Short: "inflation",
// 		Args:  cobra.ExactArgs(0),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryInflationPool)
// 			rst := model.InflationPool{}
// 			return utils.CLIQueryJSONPrint(cdc, uri, nil,
// 				func() interface{} { return &rst })
// 		},
// 	}
// }

// // GetCmdConsumption -
// func getCmdConsumption(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "consumption",
// 		Short: "consumption",
// 		Args:  cobra.ExactArgs(0),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryConsumptionMeta)
// 			rst := model.ConsumptionMeta{}
// 			return utils.CLIQueryJSONPrint(cdc, uri, nil,
// 				func() interface{} { return &rst })
// 		},
// 	}
// }

// // GetCmdTPS -
// func getCmdTPS(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "tps",
// 		Short: "tps",
// 		Args:  cobra.ExactArgs(0),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTPS)
// 			rst := model.TPS{}
// 			return utils.CLIQueryJSONPrint(cdc, uri, nil,
// 				func() interface{} { return &rst })
// 		},
// 	}
// }

// GetCmdTime -
func getCmdTime(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "time",
		Short: "time",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGlobalTime)
			rst := model.GlobalTime{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// // GetCmdStakeStats -
// func getCmdStakeStats(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "stake-stats",
// 		Short: "stake-stats <day>",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			day := args[0]
// 			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryLinoStakeStat, day)
// 			rst := model.LinoStakeStat{}
// 			return utils.CLIQueryJSONPrint(cdc, uri, nil,
// 				func() interface{} { return &rst })
// 		},
// 	}
// }
