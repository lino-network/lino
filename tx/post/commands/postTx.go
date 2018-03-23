package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	acc "github.com/lino-network/lino/tx/account"
	post "github.com/lino-network/lino/tx/post"

	sdkcli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
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

// SendTxCommand will create a send tx and sign it with the given key
func PostTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post",
		Short: "public a post to blockchain",
		RunE:  sendPostTx(cdc),
	}
	cmd.Flags().String(FlagPostID, "", "post id to identify this post for the author")
	cmd.Flags().String(FlagTitle, "", "title for the post")
	cmd.Flags().String(FlagContent, "", "content for the post")
	cmd.Flags().String(FlagParentAuthor, "", "parent author name")
	cmd.Flags().String(FlagParentPostID, "", "parent post id")
	return cmd
}

// send register transaction to the blockchain
func sendPostTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		author := viper.GetString(sdkcli.FlagName)
		postInfo := post.PostInfo{
			Author:       acc.AccountKey(author),
			PostID:       viper.GetString(FlagPostID),
			Title:        viper.GetString(FlagTitle),
			Content:      viper.GetString(FlagContent),
			ParentAuthor: acc.AccountKey(viper.GetString(FlagParentAuthor)),
			ParentPostID: viper.GetString(FlagParentPostID),
		}

		msg := post.NewCreatePostMsg(postInfo)

		// get password
		buf := sdkcli.BufferStdin()
		prompt := fmt.Sprintf("Password to sign with '%s':", author)
		passphrase, err := sdkcli.GetPassword(prompt, buf)
		if err != nil {
			return err
		}
		// build and sign the transaction, then broadcast to Tendermint
		res, err := builder.SignBuildBroadcast(author, passphrase, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
