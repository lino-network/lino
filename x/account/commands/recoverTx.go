package commands

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	acc "github.com/lino-network/lino/x/account"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// RecoverCommand will create a send tx and sign it with the given key
func RecoverTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recover",
		Short: "Create a recover tx",
		RunE:  sendRecoverTx(cdc),
	}
	cmd.Flags().String(client.FlagUser, "", "user of this transaction")
	return cmd
}

// send recover transaction to the blockchain
func sendRecoverTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		name := viper.GetString(client.FlagUser)

		resetPriv := secp256k1.GenPrivKey()
		transactionPriv := secp256k1.GenPrivKey()
		appPriv := secp256k1.GenPrivKey()
		fmt.Println("new reset private key is:", strings.ToUpper(hex.EncodeToString(resetPriv.Bytes())))
		fmt.Println("new transaction private key is:", strings.ToUpper(hex.EncodeToString(transactionPriv.Bytes())))
		fmt.Println("new app private key is:", strings.ToUpper(hex.EncodeToString(appPriv.Bytes())))

		// create the message
		msg := acc.NewRecoverMsg(name, resetPriv.PubKey(), transactionPriv.PubKey(), appPriv.PubKey())

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
