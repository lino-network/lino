package commands

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	developer "github.com/lino-network/lino/tx/developer"
)

func DeveloperRegisterTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "developer-register",
		Short: "developer register",
		RunE:  sendDeveloperRegisterTx(cdc),
	}
	cmd.Flags().String(client.FlagDeveloper, "", "developer name of this transaction")
	cmd.Flags().String(client.FlagDeposit, "", "deposit of the registration")
	return cmd
}

// send register transaction to the blockchain
func sendDeveloperRegisterTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(client.FlagDeveloper)
		msg := developer.NewDeveloperRegisterMsg(username, types.LNO(viper.GetString(client.FlagDeposit)))

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast(msg, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
