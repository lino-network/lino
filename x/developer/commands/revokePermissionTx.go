package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	crypto "github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dev "github.com/lino-network/lino/x/developer"
)

func RevokePermissionTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-permission",
		Short: "revoke permission",
		RunE:  sendRevokePermissionTx(cdc),
	}
	cmd.Flags().String(client.FlagUser, "", "user of this transaction")
	cmd.Flags().String(client.FlagPubKey, "", "public key to revoke")
	cmd.Flags().Int64(client.FlagSeconds, 3600, "seconds till expire")
	cmd.Flags().String(client.FlagPermission, "app", "grant permission")
	return cmd
}

// send grant developer transaction to the blockchain
func sendRevokePermissionTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(client.FlagUser)
		pubKeyBytes, err := hex.DecodeString(viper.GetString(client.FlagPubKey))
		if err != nil {
			return err
		}
		pubKey, err := crypto.PubKeyFromBytes(pubKeyBytes)
		if err != nil {
			return err
		}
		fmt.Println(pubKey)
		permissionStr := viper.GetString(client.FlagPermission)

		var permission types.Permission
		switch permissionStr {
		case "app":
			permission = types.AppPermission
		default:
			return errors.New("only app permission are allowed")
		}

		msg := dev.NewRevokePermissionMsg(username, pubKey, permission)

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
