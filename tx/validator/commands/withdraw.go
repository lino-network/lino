package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/validator"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

const (
	FlagName = "name"
)

// SendTxCommand will create a send tx and sign it with the given key
func WithdrawTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "withdraw a validator",
		RunE:  withDrawTx(cdc),
	}
	return cmd
}

func withDrawTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCoreContextFromViper()
		name := viper.GetString(FlagName)

		amount, err := sdk.NewRatFromDecimal(viper.GetString(FlagAmount))
		if err != nil {
			return err
		}
		// // create the message
		msg := validator.NewValidatorWithdrawMsg(name, amount)

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast(name, msg, cdc)

		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
