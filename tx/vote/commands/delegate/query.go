package delegate

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/vote/model"
	"github.com/lino-network/lino/types"

	"github.com/cosmos/cosmos-sdk/wire"
)

// GetDelegationCmd returns the delegator's delegation
func GetDelegationCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "delegation",
		Short: "Query a specific delegation",
		RunE:  cmdr.getDelegationCmd,
	}
}

type commander struct {
	storeName string
	cdc       *wire.Codec
}

func (c commander) getDelegationCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	if len(args) != 2 {
		return errors.New("You must provide voter and delegator name")
	}

	voter := types.AccountKey(args[0])
	delegator := types.AccountKey(args[1])

	res, err := ctx.Query(model.GetDelegationKey(voter, delegator), c.storeName)
	if err != nil {
		return err
	}
	delegation := new(model.Delegation)
	if err := c.cdc.UnmarshalJSON(res, delegation); err != nil {
		return err
	}

	// print out whole delegation
	output, err := json.MarshalIndent(delegation, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
