package vote

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/proposal/model"
	"github.com/lino-network/lino/types"
)

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

func (c commander) getProposalCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
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
	ctx := client.NewCoreContextFromViper()
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
