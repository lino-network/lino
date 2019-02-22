package commands

import (
	"fmt"

	wire "github.com/cosmos/cosmos-sdk/codec"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	post "github.com/lino-network/lino/x/post"
)

// PostTxCmd will create a post tx and sign it with the given key
func UpdatePostTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post",
		Short: "public a post to blockchain",
		RunE:  sendUpdatePostTx(cdc),
	}
	cmd.Flags().String(client.FlagAuthor, "", "author of this post")
	cmd.Flags().String(client.FlagPostID, "", "post id to identify this post for the author")
	cmd.Flags().String(client.FlagTitle, "", "title for the post")
	cmd.Flags().String(client.FlagContent, "", "content for the post")
	return cmd
}

// send update post transaction to the blockchain
func sendUpdatePostTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()

		msg := post.NewUpdatePostMsg(
			viper.GetString(client.FlagAuthor), viper.GetString(client.FlagPostID),
			viper.GetString(client.FlagTitle), viper.GetString(client.FlagContent),
			[]types.IDToURLMapping(nil))

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast([]sdk.Msg{msg}, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
