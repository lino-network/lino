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

// GetPostCmd returns a query post that will display the
// info and meta of the post at a given author and postID
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

	res, err := builder.Query(post.GetPostInfoKey(postKey), c.storeName)
	if err != nil {
		return err
	}
	postInfo := new(post.PostInfo)
	if err := c.cdc.UnmarshalBinary(res, postInfo); err != nil {
		return err
	}

	res, err = builder.Query(post.GetPostMetaKey(postKey), c.storeName)
	if err != nil {
		return err
	}
	postMeta := new(post.PostMeta)
	if err := c.cdc.UnmarshalBinary(res, postMeta); err != nil {
		return err
	}

	if err := client.PrintIndent(postInfo, postMeta); err != nil {
		return err
	}

	return nil
}
