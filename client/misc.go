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
		Use:   "gen-addr",
		Short: "gen-addr prints a lino bech32 address and the associated private key hex",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			priv := secp256k1.GenPrivKey()
			pub := priv.PubKey()
			addr := sdk.AccAddress(pub.Address())
			fmt.Printf("addr: %s\n", addr)
			fmt.Printf("priv-key: %s\n", strings.ToUpper(hex.EncodeToString(priv.Bytes())))
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

			filepath := args[0]
			err = encrypt.EncryptToFile(filepath, privKey, string(pw1))
			if err != nil {
				return err
			}
			fmt.Printf("encerypted key have been wrote to %s.\n", filepath)
			return nil
		},
	}
}
