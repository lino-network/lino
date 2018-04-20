package vote

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/vote"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
)

// RevokeVoterTxCmd will create a send tx and sign it with the given key
func RevokeVoterTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "voter-revoke",
		Short: "revoke voter",
		RunE:  sendRevokeVoterTx(cdc),
	}
	cmd.Flags().String(FlagUsername, "", "revoke voter")
	return cmd
}

func sendRevokeVoterTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCoreContextFromViper()
		user := viper.GetString(FlagUsername)

		// create the message
		msg := vote.NewVoterRevokeMsg(user)

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast(user, msg, cdc)

		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
