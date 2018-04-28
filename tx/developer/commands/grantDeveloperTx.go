package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	dev "github.com/lino-network/lino/tx/developer"

	"github.com/cosmos/cosmos-sdk/wire"
)

const (
	FlagUser    = "user"
	FlagSeconds = "seconds"
)

func GrantDeveloperTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant-developer",
		Short: "grant posting permission to developer",
		RunE:  sendGrantDeveloperTx(cdc),
	}
	cmd.Flags().String(FlagUser, "", "user of this transaction")
	cmd.Flags().String(FlagDeveloper, "", "developer name to grant")
	cmd.Flags().Int64(FlagSeconds, 3600, "seconds till expire")
	return cmd
}

// send grant developer transaction to the blockchain
func sendGrantDeveloperTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(FlagUser)
		developer := viper.GetString(FlagDeveloper)
		seconds := viper.GetInt64(FlagSeconds)
		msg := dev.NewGrantDeveloperMsg(username, developer, seconds, 0)

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast(username, msg, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
