package commands

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	crypto "github.com/tendermint/tendermint/crypto"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/validator"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// DepositValidatorTxCmd will create a send tx and sign it with the given key
func DepositValidatorTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-deposit",
		Short: "register a validator",
		RunE:  sendDepositValidatorTx(cdc),
	}
	cmd.Flags().String(client.FlagUser, "", "user of this transaction")
	cmd.Flags().String(client.FlagAmount, "", "amount of the donation")
	cmd.Flags().String(client.FlagPubKey, "", "validator pub key")
	return cmd
}

// send register transaction to the blockchain
func sendDepositValidatorTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		name := viper.GetString(client.FlagUser)
		pubKeyHex := viper.GetString(client.FlagPubKey)
		keyBytes, err := hex.DecodeString(pubKeyHex)
		if err != nil {
			return err
		}

		pubKey, err := crypto.PubKeyFromBytes(keyBytes)
		if err != nil {
			return err
		}

		// // create the message
		msg := validator.NewValidatorDepositMsg(
			name, types.LNO(viper.GetString(client.FlagAmount)), pubKey, viper.GetString(client.FlagLink))

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
