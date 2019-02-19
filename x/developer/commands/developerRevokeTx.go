package commands

import (
	"fmt"

	wire "github.com/cosmos/cosmos-sdk/codec"
	"github.com/lino-network/lino/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	developer "github.com/lino-network/lino/x/developer"
)

// DeveloperRevokeTxCmd - revoke developer
func DeveloperRevokeTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "developer-revoke",
		Short: "developer revoke",
		RunE:  sendDeveloperRevokeTx(cdc),
	}
	cmd.Flags().String(client.FlagDeveloper, "", "developer name of this transaction")
	return cmd
}

// send developer revoke transaction to the blockchain
func sendDeveloperRevokeTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(client.FlagDeveloper)
		msg := developer.NewDeveloperRevokeMsg(username)

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
