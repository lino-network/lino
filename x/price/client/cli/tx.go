package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/lino-network/lino/client"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/price/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Price tx subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(client.PostCommands(
		GetCmdFeedPrice(cdc),
	)...)

	return cmd
}

// GetCmdFeedPrice - feed price
func GetCmdFeedPrice(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feed <username> <amount>",
		Short: "feed <username> <amount>",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			user := linotypes.AccountKey(args[0])
			amount := args[1]
			amt, ok := sdk.NewIntFromString(amount)
			if !ok {
				panic("Invalid price")
			}

			msg := types.FeedPriceMsg{
				Username: user,
				Price:    linotypes.NewMiniDollarFromInt(amt),
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	return cmd
}
