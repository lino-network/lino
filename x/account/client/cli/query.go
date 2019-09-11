package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	// linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/account/model"
	"github.com/lino-network/lino/x/account/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the account module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(client.GetCommands(
		getCmdInfo(cdc),
		getCmdBank(cdc),
		getCmdMeta(cdc),
		getCmdListGrants(cdc),
	)...)
	return cmd
}

// GetCmdInfo - get account info
func getCmdInfo(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "info USERNAME",
		Short: "info USERNAME",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			user := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryAccountInfo, user)
			rst := model.AccountInfo{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdBank -
func getCmdBank(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "bank",
		Short: "bank USERNAME",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			user := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryAccountBank, user)
			rst := model.AccountBank{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// GetCmdMeta -
func getCmdMeta(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "meta",
		Short: "meta USERNAME",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			user := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryAccountMeta, user)
			rst := model.AccountMeta{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// getCmdListGrants -
func getCmdListGrants(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list-grants USERNAME",
		Short: "list-grants USERNAME",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			user := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s",
				types.QuerierRoute, types.QueryAccountAllGrantPubKeys, user)
			rst := make([]model.GrantPermission, 0)
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}

// func (c commander) getAccountsCmd(cmd *cobra.Command, args []string) error {
// 	ctx := client.NewCoreContextFromViper()

// 	resKVs, err := ctx.QuerySubspace(c.cdc, model.GetAccountInfoPrefix(), c.storeName)
// 	if err != nil {
// 		return err
// 	}
// 	var accounts []model.AccountInfo
// 	for _, KV := range resKVs {
// 		var info model.AccountInfo
// 		if err := c.cdc.UnmarshalBinaryLengthPrefixed(KV.Value, &info); err != nil {
// 			return err
// 		}
// 		accounts = append(accounts, info)
// 	}

// 	if err := client.PrintIndent(accounts); err != nil {
// 		return err
// 	}
// 	return nil
// }
