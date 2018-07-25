package commands

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dev "github.com/lino-network/lino/x/developer"
)

func PreAuthorizationPermissionTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pre-authorization-permission",
		Short: "grant pre authorization permission to developer",
		RunE:  sendPreAuthorizationTx(cdc),
	}
	cmd.Flags().String(client.FlagUser, "", "user of this transaction")
	cmd.Flags().String(client.FlagDeveloper, "", "developer name to grant")
	cmd.Flags().Int64(client.FlagSeconds, 3600, "seconds till expire")
	cmd.Flags().String(client.FlagGrantAmount, "grant-amount", "granted amount")
	return cmd
}

// send grant pre authorization transaction to the blockchain
func sendPreAuthorizationTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(client.FlagUser)
		developer := viper.GetString(client.FlagDeveloper)
		seconds := viper.GetInt64(client.FlagSeconds)
		amount := viper.GetString(client.FlagGrantAmount)

		msg := dev.NewPreAuthorizationMsg(username, developer, seconds, amount)

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
