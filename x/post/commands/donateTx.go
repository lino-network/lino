package commands

import (
	"fmt"

	wire "github.com/cosmos/cosmos-sdk/codec"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	post "github.com/lino-network/lino/x/post"
)

// DonateTxCmd will create a donate tx and sign it with the given key
func DonateTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donate",
		Short: "donate to a post",
		RunE:  sendDonateTx(cdc),
	}
	cmd.Flags().String(client.FlagDonator, "", "donator of this transaction")
	cmd.Flags().String(client.FlagAuthor, "", "author of the target post")
	cmd.Flags().String(client.FlagPostID, "", "post id of the target post")
	cmd.Flags().String(client.FlagAmount, "", "amount of the donation")
	cmd.Flags().String(client.FlagMemo, "", "memo of this donation")
	return cmd
}

// send donate transaction to the blockchain
func sendDonateTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(client.FlagDonator)
		author := viper.GetString(client.FlagAuthor)
		postID := viper.GetString(client.FlagPostID)
		msg := post.NewDonateMsg(
			username, types.LNO(viper.GetString(client.FlagAmount)),
			author, postID, "", viper.GetString(client.FlagMemo))

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
