package commands

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	post "github.com/lino-network/lino/x/post"
)

// PostTxCmd will create a post tx and sign it with the given key
func PostTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post",
		Short: "public a post to blockchain",
		RunE:  sendPostTx(cdc),
	}
	cmd.Flags().String(client.FlagAuthor, "", "author of this post")
	cmd.Flags().String(client.FlagPostID, "", "post id to identify this post for the author")
	cmd.Flags().String(client.FlagTitle, "", "title for the post")
	cmd.Flags().String(client.FlagContent, "", "content for the post")
	cmd.Flags().String(client.FlagParentAuthor, "", "parent post author name")
	cmd.Flags().String(client.FlagParentPostID, "", "parent post id")
	cmd.Flags().String(client.FlagSourceAuthor, "", "source post author name")
	cmd.Flags().String(client.FlagSourcePostID, "", "source post id")
	cmd.Flags().String(client.FlagRedistributionSplitRate, "0", "redistribution split rate")
	return cmd
}

// send post transaction to the blockchain
func sendPostTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		author := viper.GetString(client.FlagAuthor)
		msg := post.CreatePostMsg{
			Author:                  types.AccountKey(author),
			PostID:                  viper.GetString(client.FlagPostID),
			Title:                   viper.GetString(client.FlagTitle),
			Content:                 viper.GetString(client.FlagContent),
			ParentAuthor:            types.AccountKey(viper.GetString(client.FlagParentAuthor)),
			ParentPostID:            viper.GetString(client.FlagParentPostID),
			SourceAuthor:            types.AccountKey(viper.GetString(client.FlagSourceAuthor)),
			SourcePostID:            viper.GetString(client.FlagSourcePostID),
			RedistributionSplitRate: viper.GetString(client.FlagRedistributionSplitRate),
		}

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast(msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
