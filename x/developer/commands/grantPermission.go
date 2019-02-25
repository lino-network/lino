package commands

import (
	"fmt"

	wire "github.com/cosmos/cosmos-sdk/codec"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dev "github.com/lino-network/lino/x/developer"
)

// GrantPermissionTxCmd - user grant permission to application
func GrantPermissionTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grant-permission",
		Short: "grant permission to developer",
		RunE:  sendGrantDeveloperTx(cdc),
	}
	cmd.Flags().String(client.FlagUser, "", "user of this transaction")
	cmd.Flags().String(client.FlagDeveloper, "", "developer name to grant")
	cmd.Flags().Int64(client.FlagSeconds, 3600, "seconds till expire")
	cmd.Flags().String(client.FlagPermission, "app", "grant permission")
	return cmd
}

// send grant developer transaction to the blockchain
func sendGrantDeveloperTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(client.FlagUser)
		developer := viper.GetString(client.FlagDeveloper)
		seconds := viper.GetInt64(client.FlagSeconds)
		permissionStr := viper.GetString(client.FlagPermission)
		var permission types.Permission
		switch permissionStr {
		case "app":
			permission = types.AppPermission
		default:
			return errors.New("only app permission are allowed")
		}

		msg := dev.NewGrantPermissionMsg(username, developer, seconds, permission, "0")

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
