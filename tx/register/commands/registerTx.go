package commands

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/register"

	sdkcli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
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
	return cmd
}

// send register transaction to the blockchain
func sendRegisterTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		name := viper.GetString(sdkcli.FlagName)
		pubKey, err := GetPubKey()

		if err != nil {
			return err
		}

		// // create the message
		msg := register.NewRegisterMsg(name, *pubKey)
		// fmt.Println(fmt.Sprintf("pubkey: %v", *pubKey))
		// fmt.Println(fmt.Sprintf("pubkey to bytes: %v", string(pubKey.Bytes())))
		// get password
		buf := sdkcli.BufferStdin()
		prompt := fmt.Sprintf("Password to sign with '%s':", name)
		passphrase, err := sdkcli.GetPassword(prompt, buf)
		if err != nil {
			return err
		}
		// build and sign the transaction, then broadcast to Tendermint
		res, err := builder.SignBuildBroadcast(name, passphrase, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}

// Get the public key from the name flag
func GetPubKey() (pubKey *crypto.PubKey, err error) {
	keybase, err := keys.GetKeyBase()
	if err != nil {
		return nil, err
	}

	name := viper.GetString(sdkcli.FlagName)
	if name == "" {
		return nil, errors.Errorf("must provide a name using --name")
	}

	info, err := keybase.Get(name)
	if err != nil {
		return nil, errors.Errorf("No key for: %s", name)
	}

	return &info.PubKey, nil
}
