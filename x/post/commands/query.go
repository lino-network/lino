package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/post/model"
)

// GetPostCmd returns a query post that will display the
// info and meta of the post at a given author and postID
func GetPostCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "getpost <author> <postID>",
		Short: "Query a post",
		RunE:  cmdr.getPostCmd,
	}
}

type commander struct {
	storeName string
	cdc       *wire.Codec
}

func (c commander) getPostCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	if len(args) != 2 || len(args[0]) == 0 || len(args[1]) == 0 {
		return errors.New("You must provide an valid author and post id")
	}

	// find the key to look up the account
	author := args[0]
	postID := args[1]
	postKey := types.GetPermlink(types.AccountKey(author), postID)

	res, err := ctx.Query(model.GetPostInfoKey(postKey), c.storeName)
	if err != nil {
		return err
	}
	postInfo := new(model.PostInfo)
	if err := c.cdc.UnmarshalJSON(res, postInfo); err != nil {
		return err
	}

	res, err = ctx.Query(model.GetPostMetaKey(postKey), c.storeName)
	if err != nil {
		return err
	}
	postMeta := new(model.PostMeta)
	if err := c.cdc.UnmarshalJSON(res, postMeta); err != nil {
		return err
	}

	if err := client.PrintIndent(postInfo, postMeta); err != nil {
		return err
	}

	return nil
}

// GetPostsCmd returns a query post that will display the
// info and meta of the post at a given author and postID
func GetPostsCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "posts <author>",
		Short: "Query posts of an author",
		RunE:  cmdr.getPostsCmd,
	}
}

func (c commander) getPostsCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide an valid author")
	}

	// find the key to look up the account
	author := types.AccountKey(args[0])

	resKVs, err := ctx.QuerySubspace(
		c.cdc, append(model.GetPostInfoPrefix(author), types.PermlinkSeparator...), c.storeName)
	if err != nil {
		return err
	}
	var posts []model.PostInfo
	for _, KV := range resKVs {
		var info model.PostInfo
		if err := c.cdc.UnmarshalJSON(KV.Value, &info); err != nil {
			return err
		}
		posts = append(posts, info)
	}

	if err := client.PrintIndent(posts); err != nil {
		return err
	}
	return nil
}
