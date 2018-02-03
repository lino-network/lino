package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/light-client/commands"
	txcmd "github.com/tendermint/light-client/commands/txs"

	btypes "github.com/lino-network/lino/types"
)

//-------------------------
// SendTx

// SendTxCmd is CLI command to send tokens between basecoin accounts
var PostTxCmd = &cobra.Command{
	Use:   "post",
	Short: "send a short post",
	RunE:  commands.RequireInit(doPostTx),
}

//nolint
const (
	FlagTitle    = "title"
	FlagContent  = "content"
)

func init() {
	flags := PostTxCmd.Flags()
	flags.String(FlagTitle, "", "Post Title")
	flags.String(FlagAddress, "", "Post author")
	flags.String(FlagContent, "", "Post content")
	flags.Int(FlagSequence, -1, "Sequence number for this post")
}

// runDemo is an example of how to make a tx
func doPostTx(cmd *cobra.Command, args []string) error {
	// load data from json or flags
	tx := new(btypes.PostTx)
	err := readPostTxFlags(tx)
	if err != nil {
		return err
	}

	// Wrap and add signer
	post := &PostTx{
		chainID: commands.GetChainID(),
		Tx:      tx,
	}
	post.AddSigner(txcmd.GetSigner())
	// Sign if needed and post.  This it the work-horse
	bres, err := txcmd.SignAndPostTx(post)
	if err != nil {
		return err
	}

	// Output result
	return txcmd.OutputTx(bres)
}

func readPostTxFlags(tx *btypes.PostTx) error {
	//parse the fee and amounts into coin types
	poster, err := ParseChainAddress(viper.GetString(FlagAddress))
	if err != nil {
		return err
	}
	tx.Address = poster
	tx.Title = viper.GetString(FlagTitle)
	tx.Content = viper.GetString(FlagContent)
	tx.Sequence = viper.GetInt(FlagSequence)
	return nil
}
