package commands

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	acc "github.com/lino-network/lino/x/account"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/go-crypto"
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

		masterPriv := crypto.GenPrivKeyEd25519()
		transactionPriv := crypto.GenPrivKeyEd25519()
		micropaymentPriv := crypto.GenPrivKeyEd25519()
		postPriv := crypto.GenPrivKeyEd25519()
		fmt.Println("new master private key is:", strings.ToUpper(hex.EncodeToString(masterPriv.Bytes())))
		fmt.Println("new transaction private key is:", strings.ToUpper(hex.EncodeToString(transactionPriv.Bytes())))
		fmt.Println("new micropayment private key is:", strings.ToUpper(hex.EncodeToString(micropaymentPriv.Bytes())))
		fmt.Println("new post private key is:", strings.ToUpper(hex.EncodeToString(postPriv.Bytes())))

		// create the message
		msg := acc.NewRecoverMsg(name, masterPriv.PubKey(), transactionPriv.PubKey(), micropaymentPriv.PubKey(), postPriv.PubKey())

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast(msg, cdc)
		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
