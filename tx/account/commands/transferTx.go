package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	acc "github.com/lino-network/lino/tx/account"

	sdkcli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// nolint
const (
	FlagReceiverName = "receiver_name"
	FlagReceiverAddr = "receiver_addr"
	FlagAmount       = "amount"
	FlagMemo         = "memo"
)

// SendTxCommand will create a send tx and sign it with the given key
func TransferTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Create and sign a transfer tx",
		RunE:  sendTransferTx(cdc),
	}
	cmd.Flags().String(FlagReceiverName, "", "receiver username")
	cmd.Flags().String(FlagReceiverAddr, "", "receiver address")
	cmd.Flags().String(FlagAmount, "", "amount to transfer")
	cmd.Flags().String(FlagMemo, "", "memo msg")
	return cmd
}

// send register transaction to the blockchain
func sendTransferTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		sender := viper.GetString(sdkcli.FlagName)
		receiverName := viper.GetString(FlagReceiverName)
		receiverAddr, err := hex.DecodeString(viper.GetString(FlagReceiverAddr))
		if err != nil {
			return err
		}
		amount, err := sdk.ParseCoins(viper.GetString(FlagAmount))
		if err != nil {
			return err
		}

		msg := acc.NewTransferMsg(sender, amount, []byte(viper.GetString(FlagMemo)), acc.TransferToUser(receiverName), acc.TransferToAddr(sdk.Address(receiverAddr)))

		// get password
		buf := sdkcli.BufferStdin()
		prompt := fmt.Sprintf("Password to sign with '%s':", sender)
		passphrase, err := sdkcli.GetPassword(prompt, buf)
		if err != nil {
			return err
		}
		// build and sign the transaction, then broadcast to Tendermint
		res, err := builder.SignBuildBroadcast(sender, passphrase, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
