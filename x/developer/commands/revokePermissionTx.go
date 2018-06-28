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
	crypto "github.com/tendermint/go-crypto"

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
	cmd.Flags().Int64(client.FlagTimes, 10, "grant signing times")
	cmd.Flags().String(client.FlagPermission, "post", "grant permission")
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
		permissionStr := viper.GetString(client.FlagPermission)

		var permission types.Permission
		switch permissionStr {
		case "post":
			permission = types.PostPermission
		case "micropayment":
			permission = types.MicropaymentPermission
		default:
			return errors.New("only post and micropayment permission are allowed")
		}

		msg := dev.NewRevokePermissionMsg(username, pubKey, permission)

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast(msg, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
