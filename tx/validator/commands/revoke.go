package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/validator"

	sdkcli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	"github.com/cosmos/cosmos-sdk/wire"
)

// SendTxCommand will create a send tx and sign it with the given key
func RevokeTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "revoke a validator",
		RunE:  sendRevokeTx(cdc),
	}
	return cmd
}

// send register transaction to the blockchain
func sendRevokeTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		name := viper.GetString(sdkcli.FlagName)

		// // create the message
		msg := validator.NewValidatorRevokeMsg(name)

		// build and sign the transaction, then broadcast to Tendermint
		res, err := builder.SignBuildBroadcast(name, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
