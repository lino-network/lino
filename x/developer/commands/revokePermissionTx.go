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

// RevokePermissionTxCmd - user revoke granted public key permission
func RevokePermissionTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-permission",
		Short: "revoke permission",
		RunE:  sendRevokePermissionTx(cdc),
	}
	cmd.Flags().String(client.FlagRevokeFrom, "", "revoke from app")
	cmd.Flags().String(client.FlagUser, "", "user of this transaction")
	cmd.Flags().Int64(client.FlagSeconds, 3600, "seconds till expire")
	return cmd
}

// send grant developer transaction to the blockchain
func sendRevokePermissionTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		revokeFrom := viper.GetString(client.FlagRevokeFrom)
		permissionStr := viper.GetString(client.FlagPermission)
		var permission types.Permission
		switch permissionStr {
		case "app":
			permission = types.AppPermission
		case "preauth":
			permission = types.PreAuthorizationPermission
		default:
			return errors.New("only app permission are allowed")
		}
		username := viper.GetString(client.FlagUser)
		msg := dev.NewRevokePermissionMsg(username, revokeFrom, int(permission))

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
