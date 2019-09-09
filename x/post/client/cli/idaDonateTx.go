package cli

import (
	"fmt"

	wire "github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	post "github.com/lino-network/lino/x/post/types"
	linotypes "github.com/lino-network/lino/types"
)

// IDADonateTxCmd will create a donate tx and sign it with the given key
func IDADonateTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "idadonate",
		Short: "donate to a post using ida",
		RunE:  sendIDADonateTx(cdc),
	}
	cmd.Flags().String(FlagDonator, "", "donator of this transaction")
	cmd.Flags().String(FlagAuthor, "", "author of the target post")
	cmd.Flags().String(FlagPostID, "", "post id of the target post")
	cmd.Flags().String(FlagAmount, "", "amount of the donation")
	cmd.Flags().String(FlagMemo, "", "memo of this donation")
	cmd.Flags().String(FlagApp, "", "App's IDA")
	cmd.Flags().String(FlagSigner, "", "signer of the msg")

	return cmd
}

// send donate transaction to the blockchain
func sendIDADonateTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := linotypes.AccountKey(viper.GetString(FlagDonator))
		author := linotypes.AccountKey(viper.GetString(FlagAuthor))
		postID := viper.GetString(FlagPostID)
		app := linotypes.AccountKey(viper.GetString(FlagApp))
		amount := linotypes.IDAStr(viper.GetString(FlagAmount))
		memo := viper.GetString(FlagMemo)
		msg := post.IDADonateMsg{
			Username: username,
			App: app,
			Amount: amount,
			Author: author,
			PostID: postID,
			Memo: memo,
			Signer: username,
		}

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
