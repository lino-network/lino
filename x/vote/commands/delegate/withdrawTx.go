package delegate

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/x/vote"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WithdrawDelegateTxCmd will create a send tx and sign it with the given key
func WithdrawDelegateTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delegate-withdraw",
		Short: "withdraw money from a voter",
		RunE:  sendWithdrawDelegateTx(cdc),
	}
	cmd.Flags().String(client.FlagUser, "", "withdraw user")
	cmd.Flags().String(client.FlagVoter, "", "withdraw from")
	cmd.Flags().String(client.FlagAmount, "", "amount to withdraw")
	return cmd
}

func sendWithdrawDelegateTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		user := viper.GetString(client.FlagUser)
		voter := viper.GetString(client.FlagVoter)
		// create the message
		msg := vote.NewDelegatorWithdrawMsg(user, voter, viper.GetString(client.FlagAmount))

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)

		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
