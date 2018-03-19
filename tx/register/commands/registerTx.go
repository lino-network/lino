package commands

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/wire"
)

type registerTxCallback func(cmd *cobra.Command, args []string) error

// SendTxCommand will create a send tx and sign it with the given key
func RegisterTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Create and sign a send tx",
		RunE:  sendRegisterTx(cdc),
	}
	cmd.Flags().String(client.FlagName, "", "register new username")
	return cmd
}

// send register transaction to the blockchain
func sendRegisterTx(cdc *wire.Codec) registerTxCallback {
	// return func(cmd *cobra.Command, args []string) error {
	// 	name := viper.GetString(client.FlagName)
	// 	// get the address from the name flag
	// 	addr, err := builder.GetFromAddress()
	// 	if err != nil {
	// 		return err
	// 	}

		// create the message
		// msg := register.NewRegisterMsg(name, addr)

		// get password
		// buf := client.BufferStdin()
		// prompt := fmt.Sprintf("Password to sign with '%s':", name)
		// passphrase, err := client.GetPassword(prompt, buf)
		// if err != nil {
		// 	return err
		// }

		// build and sign the transaction, then broadcast to Tendermint
		// res, err := builder.SignBuildBroadcast(msg, cdc)

		// if err != nil {
		// 	return err
		// }

		//fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	// }
}