package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/validator"

	"github.com/cosmos/cosmos-sdk/wire"
)

// RevokeTxCmd will create a revoke tx and sign it with the given key
func RevokeTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-revoke",
		Short: "revoke a validator",
		RunE:  sendRevokeTx(cdc),
	}
	cmd.Flags().String(client.FlagUser, "", "user of this transaction")
	return cmd
}

// send revoke transaction to the blockchain
func sendRevokeTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		name := viper.GetString(client.FlagUser)

		// // create the message
		msg := validator.NewValidatorRevokeMsg(name)

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast(msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
