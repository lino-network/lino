package commands

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/tx/developer/model"
	"github.com/lino-network/lino/types"
)

// GetDeveloperCmd returns target developer information
func GetDeveloperCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "developer",
		Short: "Query developer",
		RunE:  cmdr.getDeveloperCmd,
	}
}

// GetDevelopersCmd returns all developers relative information
func GetDevelopersCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "developers",
		Short: "Query developers",
		RunE:  cmdr.getDevelopersCmd,
	}
}

type commander struct {
	storeName string
	cdc       *wire.Codec
}

func (c commander) getDeveloperCmd(cmd *cobra.Command, args []string) error {
	ctx := context.NewCoreContextFromViper()
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide a developer name")
	}

	// find the key to look up the account
	accKey := types.AccountKey(args[0])

	res, err := ctx.Query(model.GetDeveloperKey(accKey), c.storeName)
	if err != nil {
		return err
	}
	developer := new(model.Developer)
	if err := c.cdc.UnmarshalJSON(res, developer); err != nil {
		return err
	}

	// print out whole developer
	output, err := json.MarshalIndent(developer, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func (c commander) getDevelopersCmd(cmd *cobra.Command, args []string) error {
	ctx := context.NewCoreContextFromViper()
	res, err := ctx.Query(model.GetDeveloperListKey(), c.storeName)
	if err != nil {
		return err
	}

	developerList := new(model.DeveloperList)
	if err := c.cdc.UnmarshalJSON(res, developerList); err != nil {
		return err
	}

	// print out whole bank
	output, err := json.MarshalIndent(developerList, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))

	return nil
}
