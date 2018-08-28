package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"

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
	cmd.Flags().String(client.FlagUser, "", "user of this transaction")
	cmd.Flags().String(client.FlagPubKey, "", "public key to revoke")
	cmd.Flags().Int64(client.FlagSeconds, 3600, "seconds till expire")
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
		pubKey, err := cryptoAmino.PubKeyFromBytes(pubKeyBytes)
		if err != nil {
			return err
		}
		msg := dev.NewRevokePermissionMsg(username, pubKey)

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
