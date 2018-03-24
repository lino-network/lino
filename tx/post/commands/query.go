package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/builder"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	acc "github.com/lino-network/lino/tx/account"
	post "github.com/lino-network/lino/tx/post"
)

// GetBankCmd returns a query bank that will display the
// state of the bank at a given address
func GetPostCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "post <author> <postID>",
		Short: "Query a post",
		RunE:  cmdr.getPostCmd,
	}
}

type commander struct {
	storeName string
	cdc       *wire.Codec
}

func (c commander) getPostCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 2 || len(args[0]) == 0 || len(args[1]) == 0 {
		return errors.New("You must provide an valid author and post id")
	}

	// find the key to look up the account
	author := args[0]
	postID := args[1]
	postKey := post.GetPostKey(acc.AccountKey(author), postID)

	res, err := builder.Query(post.PostInfoKey(postKey), c.storeName)
	if err != nil {
		return err
	}
	postInfo := new(post.PostInfo)
	if err := c.cdc.UnmarshalBinary(res, postInfo); err != nil {
		return err
	}

	res, err = builder.Query(post.PostMetaKey(postKey), c.storeName)
	if err != nil {
		return err
	}
	postMeta := new(post.PostMeta)
	if err := c.cdc.UnmarshalBinary(res, postMeta); err != nil {
		return err
	}

	res, err = builder.Query(post.PostLikesKey(postKey), c.storeName)
	if err != nil {
		return err
	}
	postLikes := new(post.PostLikes)
	if err := c.cdc.UnmarshalBinary(res, postLikes); err != nil {
		return err
	}

	res, err = builder.Query(post.PostCommentsKey(postKey), c.storeName)
	if err != nil {
		return err
	}
	postComments := new(post.PostComments)
	if err := c.cdc.UnmarshalBinary(res, postComments); err != nil {
		return err
	}

	res, err = builder.Query(post.PostViewsKey(postKey), c.storeName)
	if err != nil {
		return err
	}
	postViews := new(post.PostViews)
	if err := c.cdc.UnmarshalBinary(res, postViews); err != nil {
		return err
	}

	res, err = builder.Query(post.PostDonationKey(postKey), c.storeName)
	if err != nil {
		return err
	}
	postDonations := new(post.PostDonations)
	if err := c.cdc.UnmarshalBinary(res, postDonations); err != nil {
		return err
	}

	if err := client.PrintIndent(postInfo, postMeta, postLikes, postComments, postViews, postDonations); err != nil {
		return err
	}

	return nil
}
