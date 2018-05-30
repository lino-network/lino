package commands

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	dev "github.com/lino-network/lino/tx/developer"
)

func GrantDeveloperTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant-developer",
		Short: "grant posting permission to developer",
		RunE:  sendGrantDeveloperTx(cdc),
	}
	cmd.Flags().String(client.FlagUser, "", "user of this transaction")
	cmd.Flags().String(client.FlagDeveloper, "", "developer name to grant")
	cmd.Flags().Int64(client.FlagSeconds, 3600, "seconds till expire")
	return cmd
}

// send grant developer transaction to the blockchain
func sendGrantDeveloperTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(client.FlagUser)
		developer := viper.GetString(client.FlagDeveloper)
		seconds := viper.GetInt64(client.FlagSeconds)
		msg := dev.NewGrantDeveloperMsg(username, developer, seconds, 0)

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast(msg, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
