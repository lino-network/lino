package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	post "github.com/lino-network/lino/tx/post"
	"github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

const (
	FlagDonator = "donator"
	FlagAmount  = "amount"
)

// DonateTxCmd will create a donate tx and sign it with the given key
func DonateTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donate",
		Short: "donate to a post",
		RunE:  sendDonateTx(cdc),
	}
	cmd.Flags().String(FlagDonator, "", "donator of this transaction")
	cmd.Flags().String(FlagAuthor, "", "author of the target post")
	cmd.Flags().String(FlagPostID, "", "post id of the target post")
	cmd.Flags().String(FlagAmount, "", "amount of the donation")
	return cmd
}

// send donate transaction to the blockchain
func sendDonateTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCoreContextFromViper()
		username := viper.GetString(FlagDonator)
		author := viper.GetString(FlagAuthor)
		postID := viper.GetString(FlagPostID)

		amount, err := sdk.NewRatFromDecimal(viper.GetString(FlagAmount))
		if err != nil {
			return err
		}
		msg := post.NewDonateMsg(types.AccountKey(username), types.LNO(amount), types.AccountKey(author), postID, "")

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast(username, msg, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
