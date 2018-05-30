package commands

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/tendermint/go-crypto"
)

// SendTxCommand will create a send tx and sign it with the given key
func RegisterTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Create and sign a register tx",
		RunE:  sendRegisterTx(cdc),
	}
	cmd.Flags().String(client.FlagUser, "", "user of this transaction")
	return cmd
}

// send register transaction to the blockchain
func sendRegisterTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		name := viper.GetString(client.FlagUser)
		referrer := viper.GetString(client.FlagReferrer)
		amount := viper.GetString(client.FlagAmount)
		pubKey, err := GetPubKey()

		if err != nil {
			return err
		}
		transactionPriv := crypto.GenPrivKeyEd25519()
		postPriv := crypto.GenPrivKeyEd25519()
		fmt.Println("transaction private key is:", strings.ToUpper(hex.EncodeToString(transactionPriv.Bytes())))
		fmt.Println("post private key is:", strings.ToUpper(hex.EncodeToString(postPriv.Bytes())))

		// // create the message
		msg := acc.NewRegisterMsg(
			referrer, name, types.LNO(amount),
			pubKey, postPriv.PubKey(), transactionPriv.PubKey())

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcastBasedOnKeyBase(name, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}

// Get the public key from the name flag
func GetPubKey() (pubKey crypto.PubKey, err error) {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	name := viper.GetString(client.FlagUser)
	if name == "" {
		return nil, errors.Errorf("must provide a name using --name")
	}

	info, err := keybase.Get(name)
	if err != nil {
		return nil, errors.Errorf("No key for: %s", name)
	}

	return info.PubKey, nil
}
