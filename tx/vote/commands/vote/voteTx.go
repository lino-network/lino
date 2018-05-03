package vote

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/vote"

	"github.com/cosmos/cosmos-sdk/wire"
)

// VoteTxCmd will create a vote tx and sign it with the given key
func VoteTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "vote a voter",
		RunE:  sendVoteTx(cdc),
	}
	cmd.Flags().String(client.FlagVoter, "", "voter for the proposal")
	cmd.Flags().Int64(client.FlagProposalID, -1, "proposal id")
	cmd.Flags().Bool(client.FlagResult, true, "vote result")
	return cmd
}

func sendVoteTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		voter := viper.GetString(client.FlagVoter)
		id := viper.GetInt64(client.FlagProposalID)
		result := viper.GetBool(client.FlagResult)

		// create the message
		msg := vote.NewVoteMsg(voter, id, result)

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast(msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
