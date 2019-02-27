package commands

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/developer/model"

	wire "github.com/cosmos/cosmos-sdk/codec"
)

// GetDeveloperCmd - returns target developer information
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

// GetDevelopersCmd - returns all developers relative information
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
	ctx := client.NewCoreContextFromViper()
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide a developer name")
	}

	accKey := types.AccountKey(args[0])
	res, err := ctx.Query(model.GetDeveloperKey(accKey), c.storeName)
	if err != nil {
		return err
	}
	developer := new(model.Developer)
	if err := c.cdc.UnmarshalBinaryLengthPrefixed(res, developer); err != nil {
		return err
	}

	output, err := json.MarshalIndent(developer, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func (c commander) getDevelopersCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	res, err := ctx.Query(model.GetDeveloperListKey(), c.storeName)
	if err != nil {
		return err
	}

	developerList := new(model.DeveloperList)
	if err := c.cdc.UnmarshalBinaryLengthPrefixed(res, developerList); err != nil {
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
