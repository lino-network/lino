package vote

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/proposal/model"
)

// GetProposalCmd returns a specific ongoing proposal
func GetOngoingProposalCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "query-ongoing-proposal",
		Short: "Query a specific ongoing proposal",
		RunE:  cmdr.getOngoingProposalCmd,
	}
}

// GetProposalCmd returns a specific expired proposal
func GetExpiredProposalCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "query-expired-proposal",
		Short: "Query a specific expired proposal",
		RunE:  cmdr.getExpiredProposalCmd,
	}
}

type commander struct {
	storeName string
	cdc       *wire.Codec
}

func (c commander) getOngoingProposalCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	if len(args) != 1 {
		return errors.New("You must provide proposal ID")
	}

	proposalID := types.ProposalKey(args[0])

	res, err := ctx.Query(model.GetOngoingProposalKey(proposalID), c.storeName)
	if err != nil {
		return err
	}
	proposal := new(model.Proposal)
	if err := c.cdc.UnmarshalJSON(res, proposal); err != nil {
		return err
	}

	// print out proposal
	output, err := json.MarshalIndent(proposal, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func (c commander) getExpiredProposalCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	if len(args) != 1 {
		return errors.New("You must provide proposal ID")
	}

	proposalID := types.ProposalKey(args[0])

	res, err := ctx.Query(model.GetExpiredProposalKey(proposalID), c.storeName)
	if err != nil {
		return err
	}
	proposal := new(model.Proposal)
	if err := c.cdc.UnmarshalJSON(res, proposal); err != nil {
		return err
	}

	// print out proposal
	output, err := json.MarshalIndent(proposal, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
