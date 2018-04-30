package delegate

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/vote"

	"github.com/cosmos/cosmos-sdk/wire"
)

// WithdrawDelegateTxCmd will create a send tx and sign it with the given key
func WithdrawDelegateTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-withdraw",
		Short: "withdraw money from a voter",
		RunE:  sendWithdrawDelegateTx(cdc),
	}
	cmd.Flags().String(FlagUsername, "", "withdraw user")
	cmd.Flags().String(FlagVoter, "", "withdraw from")
	cmd.Flags().String(FlagAmount, "", "amount to withdraw")
	return cmd
}

func sendWithdrawDelegateTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		user := viper.GetString(FlagUsername)
		voter := viper.GetString(FlagVoter)
		// create the message
		msg := vote.NewDelegatorWithdrawMsg(user, voter, viper.GetString(FlagAmount))

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast(user, msg, cdc)

		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
