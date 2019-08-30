package commands

import (
	"fmt"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"

	wire "github.com/cosmos/cosmos-sdk/codec"
	acctypes "github.com/lino-network/lino/x/account/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TransferTxCmd will create a transfer tx and sign it with the given key
func TransferTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Create and sign a transfer tx",
		RunE:  sendTransferTx(cdc),
	}
	cmd.Flags().String(client.FlagSender, "", "money sender")
	cmd.Flags().String(client.FlagReceiver, "", "receiver username")
	cmd.Flags().String(client.FlagAmount, "", "amount to transfer")
	cmd.Flags().String(client.FlagMemo, "", "memo msg")
	return cmd
}

// send transfer transaction to the blockchain
func sendTransferTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		sender := viper.GetString(client.FlagSender)
		receiver := viper.GetString(client.FlagReceiver)
		msg := acctypes.NewTransferMsg(
			sender, receiver, types.LNO(viper.GetString(client.FlagAmount)), viper.GetString(client.FlagMemo))

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
