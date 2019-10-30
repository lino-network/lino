package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	linotypes "github.com/lino-network/lino/types"
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
		utils.SimpleQueryCmd(
			"info <username>",
			"info <username>",
			types.QuerierRoute, types.QueryAccountInfo,
			1, &model.AccountInfo{})(cdc),
		utils.SimpleQueryCmd(
			"bank <username>",
			"bank <username>",
			types.QuerierRoute, types.QueryAccountBank,
			1, &model.AccountBank{})(cdc),
		utils.SimpleQueryCmd(
			"bank-addr <address>",
			"bank-addr <address>",
			types.QuerierRoute, types.QueryAccountBankByAddress,
			1, &model.AccountBank{})(cdc),
		utils.SimpleQueryCmd(
			"meta <username>",
			"meta <username>",
			types.QuerierRoute, types.QueryAccountMeta,
			1, &model.AccountMeta{})(cdc),
		utils.SimpleQueryCmd(
			"supply",
			"supply",
			types.QuerierRoute, types.QuerySupply,
			0, &model.Supply{})(cdc),
		getQueryPoolCmds(cdc),
	)...)
	return cmd
}

// getQueryPoolCmds - return a commands that queries the pool.
func getQueryPoolCmds(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool [poolname]",
		Short: "pool prints balance of all pools, if not specified, unit: LINO",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			poolnames := linotypes.ListPools()
			if len(args) == 1 {
				poolnames = []linotypes.PoolName{linotypes.PoolName(args[0])}
			}

			fmt.Printf("|%-25s|%-10s|\n", "Pool Name", "LINO")
			for _, poolname := range poolnames {
				uri := fmt.Sprintf("custom/%s/%s/%s",
					types.QuerierRoute, types.QueryPool, poolname)
				res, _, err := cliCtx.QueryWithData(uri, nil)
				if err != nil {
					fmt.Printf("Failed to Query and Print: %s, because %s", uri, err)
					return nil
				}
				rst := &linotypes.Coin{}
				cdc.MustUnmarshalJSON(res, rst)
				fmt.Printf("|%-25s|%-10s|\n", poolname, rst.ToLino())
			}
			return nil
		},
	}
}
