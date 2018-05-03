package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/validator"

	"github.com/cosmos/cosmos-sdk/wire"
)

// WithdrawTxCmd will create a withdraw tx and sign it with the given key
func WithdrawTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-withdraw",
		Short: "withdraw a validator",
		RunE:  withDrawTx(cdc),
	}
	cmd.Flags().String(client.FlagUser, "", "user of this transaction")
	return cmd
}

func withDrawTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		name := viper.GetString(client.FlagName)
		// // create the message
		msg := validator.NewValidatorWithdrawMsg(name, viper.GetString(client.FlagAmount))

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast(msg, cdc)

		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
