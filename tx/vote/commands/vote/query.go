package vote

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/tx/vote/model"
	"github.com/lino-network/lino/types"
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

// GetProposalCmd returns a specific proposal
func GetProposalCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "query-proposal",
		Short: "Query a specific proposal",
		RunE:  cmdr.getProposalCmd,
	}
}

// GetProposalListCmd returns the proposal list
func GetProposalListCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "proposal-list",
		Short: "Query ongoing and past proposal",
		RunE:  cmdr.getProposalListCmd,
	}
}

type commander struct {
	storeName string
	cdc       *wire.Codec
}

func (c commander) getVoterCmd(cmd *cobra.Command, args []string) error {
	ctx := context.NewCoreContextFromViper()
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
	ctx := context.NewCoreContextFromViper()
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

func (c commander) getProposalCmd(cmd *cobra.Command, args []string) error {
	ctx := context.NewCoreContextFromViper()
	if len(args) != 1 {
		return errors.New("You must provide proposal ID")
	}

	proposalID := types.ProposalKey(args[0])

	res, err := ctx.Query(model.GetProposalKey(proposalID), c.storeName)
	if err != nil {
		return err
	}
	proposal := new(model.Proposal)
	if err := c.cdc.UnmarshalJSON(res, proposal); err != nil {
		return err
	}

	// print out whole vote
	output, err := json.MarshalIndent(proposal, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func (c commander) getProposalListCmd(cmd *cobra.Command, args []string) error {
	ctx := context.NewCoreContextFromViper()
	res, err := ctx.Query(model.GetProposalListKey(), c.storeName)
	if err != nil {
		return err
	}

	proposalList := new(model.ProposalList)
	if err := c.cdc.UnmarshalJSON(res, proposalList); err != nil {
		return err
	}

	// print out whole bank
	output, err := json.MarshalIndent(proposalList, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))

	return nil
}
