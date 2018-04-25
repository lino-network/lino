package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// nolint
const (
	FlagSender       = "sender"
	FlagReceiverName = "receiver_name"
	FlagReceiverAddr = "receiver_addr"
	FlagAmount       = "amount"
	FlagMemo         = "memo"
)

// TransferTxCmd will create a transfer tx and sign it with the given key
func TransferTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Create and sign a transfer tx",
		RunE:  sendTransferTx(cdc),
	}
	cmd.Flags().String(FlagSender, "", "money sender")
	cmd.Flags().String(FlagReceiverName, "", "receiver username")
	cmd.Flags().String(FlagReceiverAddr, "", "receiver address")
	cmd.Flags().String(FlagAmount, "", "amount to transfer")
	cmd.Flags().String(FlagMemo, "", "memo msg")
	return cmd
}

// send transfer transaction to the blockchain
func sendTransferTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCoreContextFromViper()
		sender := viper.GetString(FlagSender)
		receiverName := viper.GetString(FlagReceiverName)
		receiverAddr, err := hex.DecodeString(viper.GetString(FlagReceiverAddr))
		if err != nil {
			return err
		}
		msg := acc.NewTransferMsg(sender, types.LNO(viper.GetString(FlagAmount)), []byte(viper.GetString(FlagMemo)),
			acc.TransferToUser(receiverName), acc.TransferToAddr(sdk.Address(receiverAddr)))

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast(sender, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
