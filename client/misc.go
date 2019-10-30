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

func GetGenAddrCmd(cdc *amino.Codec) *cobra.Command {
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
