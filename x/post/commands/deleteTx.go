package commands

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	post "github.com/lino-network/lino/x/post"
)

// DeletePostTxCmd deletes a post tx and sign it with the given key
func DeletePostTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post-delete",
		Short: "delete a post to blockchain",
		RunE:  sendDeletePostTx(cdc),
	}
	cmd.Flags().String(client.FlagAuthor, "", "author of this post")
	cmd.Flags().String(client.FlagPostID, "", "post id to identify this post for the author")
	return cmd
}

// send delete post transaction to the blockchain
func sendDeletePostTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		author := viper.GetString(client.FlagAuthor)
		postID := viper.GetString(client.FlagPostID)

		msg := post.NewDeletePostMsg(author, postID)

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast(msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
