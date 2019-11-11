package client

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/lino-network/lino/client/encrypt"
)

func GetNowCmd(cdc *amino.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "now",
		Short: "now",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			str, _ := cdc.MarshalJSON(time.Now())
			fmt.Println(string(str))
			return nil
		},
	}
}

func GetGenAddrCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gen-addr [keyfile]",
		Short: "gen-addr [keyfile] prints a lino bech32 address. The private key will be printed if keyfile not specified",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			priv := secp256k1.GenPrivKey()
			pub := priv.PubKey()
			addr := sdk.AccAddress(pub.Address())
			fmt.Printf("addr: %s\n", addr)
			if len(args) > 0 {
				keyfile := args[0]
				return encryptSave(keyfile, []byte(hex.EncodeToString(priv.Bytes())))
			} else {
				fmt.Printf("priv-key: %s\n", strings.ToUpper(hex.EncodeToString(priv.Bytes())))
			}
			return nil
		},
	}
}

func GetAddrOfCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "addr-of <@keyfile>",
		Short: "addr-of <@keyfile> prints the lino bech32 address, @ included, e.g. @foo.key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			file := args[0]
			priv, err := ParsePrivKey(file)
			if err != nil {
				return err
			}
			addr := sdk.AccAddress(priv.PubKey().Address())
			fmt.Printf("addr: %s\n", addr)
			return nil
		},
	}
}

func GetEncryptPrivKey() *cobra.Command {
	return &cobra.Command{
		Use:   "encrypt-key <file>",
		Short: "encrypt-key <file>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Hex-encoded Private key: ")
			privKey, err := terminal.ReadPassword(0)
			if err != nil {
				return err
			}
			fmt.Println("")

			// validate key
			_, err = ParsePrivKey(string(privKey))
			if err != nil {
				return fmt.Errorf("invalid privKey: %+v", err)
			}

			return encryptSave(args[0], privKey)
		},
	}
}

func encryptSave(filepath string, privKey []byte) error {
	fmt.Printf("Password: ")
	pw1, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}
	fmt.Println("")

	fmt.Printf("Password again: ")
	pw2, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}
	fmt.Printf("\n\n")

	if string(pw1) != string(pw2) {
		return fmt.Errorf("password mismatch")
	}

	err = encrypt.EncryptToFile(filepath, privKey, string(pw1))
	if err != nil {
		return err
	}
	fmt.Printf("encrypted key have been wrote to %s.\n", filepath)
	return nil
}
