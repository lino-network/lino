package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/post/model"
	"github.com/lino-network/lino/types"
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
	ctx := context.NewCoreContextFromViper()
	if len(args) != 2 || len(args[0]) == 0 || len(args[1]) == 0 {
		return errors.New("You must provide an valid author and post id")
	}

	// find the key to look up the account
	author := args[0]
	postID := args[1]
	postKey := types.GetPostKey(types.AccountKey(author), postID)

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
