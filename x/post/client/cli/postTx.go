package cli

import (
	"fmt"

	wire "github.com/cosmos/cosmos-sdk/codec"
	"github.com/lino-network/lino/client"
	linotypes "github.com/lino-network/lino/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	post "github.com/lino-network/lino/x/post"
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
	cmd.Flags().String(FlagCreatedBy, "", "application(developer) that creates the post")
	return cmd
}

// send post transaction to the blockchain
func sendPostTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		msg := post.CreatePostMsg{
			Author:    linotypes.AccountKey(viper.GetString(FlagAuthor)),
			PostID:    viper.GetString(FlagPostID),
			Title:     viper.GetString(FlagTitle),
			Content:   viper.GetString(FlagContent),
			CreatedBy: linotypes.AccountKey(viper.GetString(FlagCreatedBy)),
		}

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)
		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
