package commands

import (
	"fmt"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	// "github.com/lino-network/lino/types"
	developer "github.com/lino-network/lino/x/developer/types"
)

// DeveloperRegisterTxCmd - register to be developer
func DeveloperRegisterTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "developer-register",
		Short: "developer register",
		RunE:  sendDeveloperRegisterTx(cdc),
	}
	cmd.Flags().String(client.FlagDeveloper, "", "developer name of this transaction")
	cmd.Flags().String(client.FlagWebsite, "", "website of the app")
	cmd.Flags().String(client.FlagDescription, "", "description of the app")
	cmd.Flags().String(client.FlagAppMeta, "", "meta-data of the app")
	return cmd
}

// send register transaction to the blockchain
func sendDeveloperRegisterTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(client.FlagDeveloper)
		msg := developer.NewDeveloperRegisterMsg(
			username,
			viper.GetString(client.FlagWebsite), viper.GetString(client.FlagDescription),
			viper.GetString(client.FlagAppMeta))

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
