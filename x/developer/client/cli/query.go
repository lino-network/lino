package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/developer/model"
	"github.com/lino-network/lino/x/developer/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the developer module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(client.GetCommands(
		getCmdShow(cdc),
		getCmdList(cdc),
		getCmdListAffiliated(cdc),
		getCmdIDAShow(cdc),
		getCmdIDABalance(cdc),
		getCmdReservePool(cdc),
		getCmdIDAStats(cdc),
	)...)
	return cmd
}

// GetCmdShow queries information about a name
func getCmdShow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "show [username]",
		Short: "show username's developer detail",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryDeveloper, name)
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &model.Developer{} })
		},
	}
}

// GetCmdList lists all developers
func getCmdList(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list all developers",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDeveloperList)
			rst := make(map[string]model.Developer)
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdListAffiliated - list all affiliated accounts.
func getCmdListAffiliated(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list-affiliated [app]",
		Short: "list all affiliated accounts of app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryAffiliated, app)
			rst := make([]linotypes.AccountKey, 0)
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdIDAShow -
func getCmdIDAShow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ida-show [app]",
		Short: "show ida details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryIDA, app)
			rst := model.AppIDA{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdIDABalance -
func getCmdIDABalance(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ida-balance app user",
		Short: "return ida balance of a user of the app",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := args[0]
			user := args[1]
			uri := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, types.QueryIDABalance, app, user)
			rst := types.QueryResultIDABalance{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdReservePool -
func getCmdReservePool(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reserve-pool",
		Short: "reserve-pool",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			uri := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryReservePool)
			rst := model.ReservePool{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdIDAStats -
func getCmdIDAStats(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ida-stats",
		Short: "ida-stats",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryIDAStats, app)
			rst := model.AppIDAStats{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}
