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
	FlagTo          = "to"
	FlagAmount      = "amount"
	FlagMemo        = "memo"
	FlagAddr        = "addr"
	FlagAddrSeq     = "addr-seq"
	FlagAddrPrivKey = "addr-priv-key"
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
		getCmdTransferV2(cdc),
	)...)

	return cmd
}

// GetCmdRegister - register a new user with random generated private keys.
func GetCmdRegister(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register <type:referrer> --amount <amount> --name <name>",
		Short: "register <type:referrer> --amount <amount> --name <name>",
		Args:  cobra.ExactArgs(1),
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
			referrer, err := parseAccOrAddr(referrerArg)
			if err != nil {
				return err
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

	// always 0, in this cmd.
	// cmd.Flags().Uint64(FlagAddrSeq, 0, "sequence# of the new transaction key")
	return cmd
}

// GetCmdBin - register as developer.
func GetCmdBind(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind <type:referrer>",
		Short: "bind <type:referrer> --addr <addr> --addr-priv-key <hex> --addr-seq <seq> --name <name> --amount <amount>",
		Args:  cobra.ExactArgs(1),
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
			referrer, err := parseAccOrAddr(referrerArg)
			if err != nil {
				return err
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

	// always 0, in this cmd.
	// cmd.Flags().Uint64(FlagAddrSeq, 0, "sequence# of the new transaction key")
	return cmd
}

// GetCmdTransferV2 -
func getCmdTransferV2(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "transfer <type:from> --to <type:to> --amount <amount> --memo memo, See help for type",
		Long:  "type is either 'addr' or 'user', e.g. transfer addr:lino158de3... --to user:yxia --amount 10 --memo 'demo'",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			from, err := parseAccOrAddr(args[0])
			if err != nil {
				return err
			}
			to, err := parseAccOrAddr(viper.GetString(FlagTo))
			if err != nil {
				return err
			}
			amount := viper.GetString(FlagAmount)
			memo := viper.GetString(FlagMemo)
			msg := types.TransferV2Msg{
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

func parseAccOrAddr(s string) (rst linotypes.AccOrAddr, err error) {
	comps := strings.Split(s, ":")
	if len(comps) != 2 || !(comps[0] == "addr" || comps[0] == "user") {
		return rst, fmt.Errorf("invalid param: %s", s)
	}
	if comps[0] == "addr" {
		addr, err := sdk.AccAddressFromBech32(comps[1])
		if err != nil {
			return rst, err
		}
		return linotypes.NewAccOrAddrFromAddr(addr), nil
	} else {
		return linotypes.NewAccOrAddrFromAcc(linotypes.AccountKey(comps[1])), nil
	}
}
