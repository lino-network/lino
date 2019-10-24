package cli

import (
	// "fmt"

	"github.com/cosmos/cosmos-sdk/client"
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
			"pool <poolname>",
			"pool <poolname>",
			types.QuerierRoute, types.QueryPool,
			1, &linotypes.Coin{})(cdc),
		utils.SimpleQueryCmd(
			"supply",
			"supply",
			types.QuerierRoute, types.QuerySupply,
			0, &model.Supply{})(cdc),
	)...)
	return cmd
}
