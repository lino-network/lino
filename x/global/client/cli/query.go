package cli

import (
	// "fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	linotypes "github.com/lino-network/lino/types"
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
		utils.SimpleQueryCmd(
			"time",
			"time",
			types.QuerierRoute, types.QueryGlobalTime,
			0, &model.GlobalTime{})(cdc),
		utils.SimpleQueryCmd(
			"event-errors",
			"event-errors",
			types.QuerierRoute, types.QueryGlobalEventErrors,
			0, &([]model.EventError{}))(cdc),
		utils.SimpleQueryCmd(
			"bc-event-errors",
			"bc-event-errors",
			types.QuerierRoute, types.QueryGlobalBCEventErrors,
			0, &([]linotypes.BCEventErr{}))(cdc),
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
