package cli

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/lino-network/lino/client"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/types"
)

const (
	FlagTo     = "to"
	FlagAmount = "amount"
	FlagMemo   = "memo"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Account tx subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(client.PostCommands(
		GetCmdRegister(cdc),
		GetCmdTransfer(cdc),
	)...)

	return cmd
}

// GetCmdRegister - register as developer.
func GetCmdRegister(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register <referrer> <amount> <name>",
		Short: "register <referrer> <amount> <name>",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			referrer := linotypes.AccountKey(args[0])
			amount := linotypes.LNO(args[1])
			username := linotypes.AccountKey(args[2])

			resetPriv := secp256k1.GenPrivKey()
			transactionPriv := secp256k1.GenPrivKey()

			fmt.Println(
				"reset private key is:",
				strings.ToUpper(hex.EncodeToString(resetPriv.Bytes())))
			fmt.Println(
				"transaction private key is:",
				strings.ToUpper(hex.EncodeToString(transactionPriv.Bytes())))

			msg := types.RegisterMsg{
				Referrer:             referrer,
				NewUser:              username,
				RegisterFee:          amount,
				NewResetPubKey:       resetPriv.PubKey(),
				NewTransactionPubKey: transactionPriv.PubKey(),
				NewAppPubKey:         transactionPriv.PubKey(),
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	return cmd
}

// GetCmdTransfer -
func GetCmdTransfer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "transfer <from> --to <bar> --amount <amount> --memo memo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			from := linotypes.AccountKey(args[0])
			to := linotypes.AccountKey(viper.GetString(FlagTo))
			amount := linotypes.LNO(viper.GetString(FlagAmount))
			memo := linotypes.LNO(viper.GetString(FlagMemo))
			msg := types.TransferMsg{
				Sender:   from,
				Receiver: to,
				Amount:   amount,
				Memo:     memo,
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagTo, "", "receiver username")
	cmd.Flags().String(FlagAmount, "", "amount to transfer")
	cmd.Flags().String(FlagMemo, "", "memo msg")
	cmd.MarkFlagRequired(FlagTo)
	cmd.MarkFlagRequired(FlagAmount)
	return cmd
}
