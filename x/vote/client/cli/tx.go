package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	linotypes "github.com/lino-network/lino/types"
	types "github.com/lino-network/lino/x/vote"
)

const (
	FlagAmount = "amount"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "vote tx subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(client.PostCommands(
		GetCmdStakein(cdc),
		GetCmdStakeout(cdc),
	)...)

	return cmd
}

// GetCmdStakein - stakein
func GetCmdStakein(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake-in",
		Short: "stake-in <username> --amount <lino>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			user := linotypes.AccountKey(args[0])
			amount := viper.GetString(FlagAmount)
			msg := types.StakeInMsg{
				Username: user,
				Deposit:  amount,
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagAmount, "", "amount of stake in")
	_ = cmd.MarkFlagRequired(FlagAmount)
	return cmd
}

// GetCmdStakeout -
func GetCmdStakeout(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake-out",
		Short: "stake-out <username> --amount <lino>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			user := linotypes.AccountKey(args[0])
			amount := viper.GetString(FlagAmount)
			msg := types.StakeOutMsg{
				Username: user,
				Amount:   amount,
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagAmount, "", "amount of stake in")
	_ = cmd.MarkFlagRequired(FlagAmount)
	return cmd
}

// GetCmdClaimInterest -
func GetCmdClaimInterest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-interest",
		Short: "claim-interest username",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			username := linotypes.AccountKey(args[0])
			msg := types.ClaimInterestMsg{
				Username: username,
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	return cmd
}
