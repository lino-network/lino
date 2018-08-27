package vote

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/vote/model"
)

// GetVoterCmd returns target voter information
func GetVoterCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "voter",
		Short: "Query voter",
		RunE:  cmdr.getVoterCmd,
	}
}

// GetVoteCmd returns a voter's vote on a proposal
func GetVoteCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "query-vote",
		Short: "Query a specific vote",
		RunE:  cmdr.getVoteCmd,
	}
}

type commander struct {
	storeName string
	cdc       *wire.Codec
}

func (c commander) getVoterCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide a voter name")
	}

	// find the key to look up the account
	accKey := types.AccountKey(args[0])

	res, err := ctx.Query(model.GetVoterKey(accKey), c.storeName)
	if err != nil {
		return err
	}
	voter := new(model.Voter)
	if err := c.cdc.UnmarshalJSON(res, voter); err != nil {
		return err
	}

	// print out whole voter
	output, err := json.MarshalIndent(voter, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func (c commander) getVoteCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	if len(args) != 2 {
		return errors.New("You must provide proposal ID and voter name")
	}

	proposalID := types.ProposalKey(args[0])
	voter := types.AccountKey(args[1])

	res, err := ctx.Query(model.GetVoteKey(proposalID, voter), c.storeName)
	if err != nil {
		return err
	}
	vote := new(model.Vote)
	if err := c.cdc.UnmarshalJSON(res, vote); err != nil {
		return err
	}

	// print out whole vote
	output, err := json.MarshalIndent(vote, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
