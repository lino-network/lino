package cli

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/lino-network/lino/client"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/types"
)

const (
	FlagTo      = "to"
	FlagAmount  = "amount"
	FlagMemo    = "memo"
	FlagByAddr  = "by-addr"
	FlagAddrSeq = "addr-seq"
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
		// GetCmdRegister(cdc),
		GetCmdTransfer(cdc),
	)...)

	return cmd
}

// GetCmdRegister - register as developer.
func GetCmdRegister(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register <referrer> <amount> <name>",
		Short: "register <referrer> <amount> <name> --by-addr=true/false",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			referrerArg := args[0]
			amount := args[1]
			username := linotypes.AccountKey(args[2])

			txPriv := secp256k1.GenPrivKey()
			signPriv := secp256k1.GenPrivKey()

			fmt.Println(
				"tx private hex-encoded:",
				strings.ToUpper(hex.EncodeToString(txPriv.Bytes())))
			fmt.Println(
				"signing private key hex-encoded:",
				strings.ToUpper(hex.EncodeToString(signPriv.Bytes())))
			isAddr := viper.GetBool(FlagByAddr)
			var referrer linotypes.AccOrAddr
			if isAddr {
				referrer = linotypes.NewAccOrAddrFromAcc(linotypes.AccountKey(referrerArg))
			} else {
				referrer = linotypes.NewAccOrAddrFromAddr(sdk.AccAddress(referrerArg))
			}

			msg := types.RegisterV2Msg{
				Referrer:             referrer,
				NewUser:              username,
				RegisterFee:          amount,
				NewTransactionPubKey: txPriv.PubKey(),
				NewSigningPubKey:     signPriv.PubKey(),
			}
			return ctx.DoTxPrintResponse(msg, client.OptionalSigner{
				PrivKey: txPriv,
				Seq:     0,
			})
		},
	}

	cmd.Flags().Bool(FlagByAddr, false, "register referrer is an address")
	// always 0, in this cmd.
	// cmd.Flags().Uint64(FlagAddrSeq, 0, "sequence# of the new transaction key")
	return cmd
}

// TODO(yumin):
// Add an addition CLI to support register an account for an existing address.

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
			amount := viper.GetString(FlagAmount)
			memo := viper.GetString(FlagMemo)
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
	_ = cmd.MarkFlagRequired(FlagTo)
	_ = cmd.MarkFlagRequired(FlagAmount)
	return cmd
}
