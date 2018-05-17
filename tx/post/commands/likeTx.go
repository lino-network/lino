package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	post "github.com/lino-network/lino/tx/post"

	"github.com/cosmos/cosmos-sdk/wire"
)

// LikeTxCmd will create a like tx and sign it with the given key
func LikeTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "like",
		Short: "like a post or dislike a post",
		RunE:  sendLikeTx(cdc),
	}
	cmd.Flags().String(client.FlagLikeUser, "", "like user of this transaction")
	cmd.Flags().String(client.FlagPostID, "", "post id to identify this post for the author")
	cmd.Flags().String(client.FlagAuthor, "", "title for the post")
	cmd.Flags().String(client.FlagWeight, "", "content for the post")
	return cmd
}

// send like transaction to the blockchain
func sendLikeTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		username := viper.GetString(client.FlagLikeUser)
		author := viper.GetString(client.FlagAuthor)
		postID := viper.GetString(client.FlagPostID)
		weight, err := strconv.Atoi(viper.GetString(client.FlagWeight))
		if err != nil {
			return err
		}

		msg := post.NewLikeMsg(username, int64(weight), author, postID)

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast(msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
