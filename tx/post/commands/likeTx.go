package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	acc "github.com/lino-network/lino/tx/account"
	post "github.com/lino-network/lino/tx/post"

	"github.com/cosmos/cosmos-sdk/client/builder"
	"github.com/cosmos/cosmos-sdk/wire"
)

// nolint
const (
	FlagLikeUser = "likeUser"
	FlagWeight   = "weight"
	FlagAuthor   = "author"
)

// SendTxCommand will create a send tx and sign it with the given key
func LikeTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "like",
		Short: "like a post or dislike a post",
		RunE:  sendLikeTx(cdc),
	}
	cmd.Flags().String(FlagLikeUser, "", "like user of this transaction")
	cmd.Flags().String(FlagPostID, "", "post id to identify this post for the author")
	cmd.Flags().String(FlagAuthor, "", "title for the post")
	cmd.Flags().String(FlagWeight, "", "content for the post")
	return cmd
}

// send register transaction to the blockchain
func sendLikeTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		username := viper.GetString(FlagLikeUser)
		author := viper.GetString(FlagAuthor)
		postID := viper.GetString(FlagPostID)
		weight, err := strconv.Atoi(viper.GetString(FlagWeight))
		if err != nil {
			return err
		}

		msg := post.NewLikeMsg(acc.AccountKey(username), int64(weight), acc.AccountKey(author), postID)

		// build and sign the transaction, then broadcast to Tendermint
		res, err := builder.SignBuildBroadcast(username, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
