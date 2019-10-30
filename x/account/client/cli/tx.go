package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"
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
	FlagName        = "name"
	FlagAmount      = "amount"
	FlagRegFee      = "reg-fee"
	FlagMemo        = "memo"
	FlagAddr        = "addr"
	FlagAddrSeq     = "addr-seq"
	FlagAddrPrivKey = "addr-priv-key"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Account tx subcommands, type is either 'addr' or 'user'.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(client.PostCommands(
		getCmdRegister(cdc),
		getCmdTransferV2(cdc),
		getCmdBind(cdc),
	)...)

	return cmd
}

// getCmdRegister - register a new user with random generated private keys.
func getCmdRegister(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register <type:referrer> --reg-fee <amount> --name <name>",
		Short: "register <type:referrer> --reg-fee <amount> --name <name>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			referrer, err := parseAccOrAddr(args[0])
			if err != nil {
				return err
			}
			amount := viper.GetString(FlagRegFee)
			username := linotypes.AccountKey(viper.GetString(FlagName))

			txPriv := secp256k1.GenPrivKey()
			signPriv := secp256k1.GenPrivKey()

			fmt.Println(
				"tx private hex-encoded:",
				strings.ToUpper(hex.EncodeToString(txPriv.Bytes())))
			fmt.Println(
				"signing private key hex-encoded:",
				strings.ToUpper(hex.EncodeToString(signPriv.Bytes())))
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
	cmd.Flags().String(FlagRegFee, "", "amount of register fee")
	cmd.Flags().String(FlagName, "", "name of new user")
	_ = cmd.MarkFlagRequired(FlagRegFee)
	_ = cmd.MarkFlagRequired(FlagName)
	return cmd
}

// getCmdBind - bind username to an address.
func getCmdBind(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind <type:referrer>",
		Short: "bind <type:referrer> --addr <addr> --addr-priv-key <key> --addr-seq <seq> --name <name> --reg-fee <amount>",
		Long:  "bind <type:referrer> --addr <addr> --addr-priv-key <key> --addr-seq <seq> --name <name> --reg-fee <amount> will bind <addr> with <name> as username. The signing key will be the same as transaction key.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			referrer, err := parseAccOrAddr(args[0])
			if err != nil {
				return err
			}
			amount := viper.GetString(FlagRegFee)
			username := linotypes.AccountKey(viper.GetString(FlagName))
			addr, err := sdk.AccAddressFromBech32(viper.GetString(FlagAddr))
			if err != nil {
				return err
			}
			addrPriv, err := client.ParsePrivKey(viper.GetString(FlagAddrPrivKey))
			if err != nil {
				return err
			}
			addrSeq, err := strconv.ParseInt(viper.GetString(FlagAddrSeq), 10, 64)
			if err != nil {
				return err
			}

			if addr.String() != sdk.AccAddress(addrPriv.PubKey().Address()).String() {
				return fmt.Errorf("address and priv-key mismatch, priv-key's addr: %s",
					sdk.AccAddress(addrPriv.PubKey().Address()))
			}

			msg := types.RegisterV2Msg{
				Referrer:             referrer,
				NewUser:              username,
				RegisterFee:          amount,
				NewTransactionPubKey: addrPriv.PubKey(),
				NewSigningPubKey:     addrPriv.PubKey(),
			}
			return ctx.DoTxPrintResponse(msg, client.OptionalSigner{
				PrivKey: addrPriv,
				Seq:     uint64(addrSeq),
			})
		},
	}

	cmd.Flags().String(FlagAddr, "", "address to be binded with the name")
	cmd.Flags().String(FlagAddrPrivKey, "", "private hex of the address")
	cmd.Flags().String(FlagAddrSeq, "", "sequence # of the address")
	cmd.Flags().String(FlagRegFee, "", "amount of register fee")
	cmd.Flags().String(FlagName, "", "name of new user")
	_ = cmd.MarkFlagRequired(FlagAddr)
	_ = cmd.MarkFlagRequired(FlagAddrPrivKey)
	_ = cmd.MarkFlagRequired(FlagAddrSeq)
	_ = cmd.MarkFlagRequired(FlagRegFee)
	_ = cmd.MarkFlagRequired(FlagName)
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
