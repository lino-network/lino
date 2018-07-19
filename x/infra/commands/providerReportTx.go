package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	infra "github.com/lino-network/lino/x/infra"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

const ()

func ProviderReportTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider-report",
		Short: "provider report usage",
		RunE:  sendProviderReportTx(cdc),
	}
	cmd.Flags().String(client.FlagProvider, "", "reporter of this transaction")
	cmd.Flags().String(client.FlagUsage, "", "usage of the report")
	return cmd
}

// send provider report transaction to the blockchain
func sendProviderReportTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(client.FlagProvider)
		usage, err := strconv.ParseInt(viper.GetString(client.FlagUsage), 10, 64)
		if err != nil {
			return err
		}
		msg := infra.NewProviderReportMsg(username, usage)

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
