package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	post "github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
)

// nolint
const (
	FlagPostID       = "post_ID"
	FlagTitle        = "title"
	FlagContent      = "content"
	FlagParentAuthor = "parent_author"
	FlagParentPostID = "parent_post_ID"
)

// PostTxCmd will create a post tx and sign it with the given key
func PostTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post",
		Short: "public a post to blockchain",
		RunE:  sendPostTx(cdc),
	}
	cmd.Flags().String(FlagAuthor, "", "author of this post")
	cmd.Flags().String(FlagPostID, "", "post id to identify this post for the author")
	cmd.Flags().String(FlagTitle, "", "title for the post")
	cmd.Flags().String(FlagContent, "", "content for the post")
	cmd.Flags().String(FlagParentAuthor, "", "parent author name")
	cmd.Flags().String(FlagParentPostID, "", "parent post id")
	return cmd
}

// send post transaction to the blockchain
func sendPostTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCoreContextFromViper()
		author := viper.GetString(FlagAuthor)
		postCreateParams := post.PostCreateParams{
			Author:                  types.AccountKey(author),
			PostID:                  viper.GetString(FlagPostID),
			Title:                   viper.GetString(FlagTitle),
			Content:                 viper.GetString(FlagContent),
			ParentAuthor:            types.AccountKey(viper.GetString(FlagParentAuthor)),
			ParentPostID:            viper.GetString(FlagParentPostID),
			RedistributionSplitRate: "0",
		}

		msg := post.NewCreatePostMsg(postCreateParams)

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast(author, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
