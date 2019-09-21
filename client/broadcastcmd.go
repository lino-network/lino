package client

import (
	"encoding/hex"
	"fmt"

	// "strings"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
	"github.com/cosmos/cosmos-sdk/codec"

	linotypes "github.com/lino-network/lino/types"
)

// GetCmdBoradcast -
func GetCmdBroadcast(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broadcast",
		Short: "broadcast <tx-hex>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := NewCoreBroadcastContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			hexstr := args[0]
			hexbytes, err := hex.DecodeString(hexstr)
			if err != nil {
				return err
			}
			res, err := ctx.BroadcastTx(hexbytes)
			if err != nil {
				return err
			}
			fmt.Println(res.String())
			return nil
		},
	}
	return cmd
}
