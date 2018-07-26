package commands

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	developer "github.com/lino-network/lino/x/developer"
)

func DeveloperUpdateTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "developer-update",
		Short: "developer update",
		RunE:  sendDeveloperUpdateTx(cdc),
	}
	cmd.Flags().String(client.FlagDeveloper, "", "developer name of this transaction")
	cmd.Flags().String(client.FlagWebsite, "", "website of the app")
	cmd.Flags().String(client.FlagDescription, "", "description of the app")
	cmd.Flags().String(client.FlagAppMeta, "", "meta-data of the app")
	return cmd
}

// send register transaction to the blockchain
func sendDeveloperUpdateTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(client.FlagDeveloper)
		msg := developer.NewDeveloperUpdateMsg(
			username, viper.GetString(client.FlagWebsite),
			viper.GetString(client.FlagDescription), viper.GetString(client.FlagAppMeta))

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
