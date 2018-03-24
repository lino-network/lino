package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	acc "github.com/lino-network/lino/tx/account"

	sdkcli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// nolint
const (
	FlagIsFollow = "is_follow"
	FlagFollowee = "followee"
)

// SendTxCommand will create a send tx and sign it with the given key
func FollowTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "follow",
		Short: "Create and sign a follow/unfollow tx",
		RunE:  sendFollowTx(cdc),
	}
	cmd.Flags().Bool(FlagIsFollow, true, "false if this is unfollow")
	cmd.Flags().String(FlagFollowee, "", "target to follow or unfollow")
	return cmd
}

// send register transaction to the blockchain
func sendFollowTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		follower := viper.GetString(sdkcli.FlagName)
		followee := viper.GetString(FlagFollowee)
		var msg sdk.Msg
		isFollow := viper.GetBool(FlagIsFollow)
		if isFollow {
			msg = acc.NewFollowMsg(follower, followee)
		} else {
			msg = acc.NewUnfollowMsg(follower, followee)
		}

		// get password
		buf := sdkcli.BufferStdin()
		prompt := fmt.Sprintf("Password to sign with '%s':", follower)
		passphrase, err := sdkcli.GetPassword(prompt, buf)
		if err != nil {
			return err
		}
		// build and sign the transaction, then broadcast to Tendermint
		res, err := builder.SignBuildBroadcast(follower, passphrase, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
