package commands

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/infra/model"

	"github.com/cosmos/cosmos-sdk/wire"
)

// GetInfraProviderCmd returns target voter information
func GetInfraProviderCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "infra-provider",
		Short: "Query infra provider",
		RunE:  cmdr.getInfraProviderCmd,
	}
}

// GetInfraProvidersCmd returns all validators relative information
func GetInfraProvidersCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "infra-providers",
		Short: "Query infra providers",
		RunE:  cmdr.getInfraProvidersCmd,
	}
}

type commander struct {
	storeName string
	cdc       *wire.Codec
}

func (c commander) getInfraProviderCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide a infra provider name")
	}

	// find the key to look up the account
	accKey := types.AccountKey(args[0])

	res, err := ctx.Query(model.GetInfraProviderKey(accKey), c.storeName)
	if err != nil {
		return err
	}
	provider := new(model.InfraProvider)
	if err := c.cdc.UnmarshalJSON(res, provider); err != nil {
		return err
	}

	// print out whole infra provider
	output, err := json.MarshalIndent(provider, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}

func (c commander) getInfraProvidersCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	res, err := ctx.Query(model.GetInfraProviderListKey(), c.storeName)
	if err != nil {
		return err
	}

	providerList := new(model.InfraProviderList)
	if err := c.cdc.UnmarshalJSON(res, providerList); err != nil {
		return err
	}

	output, err := json.MarshalIndent(providerList, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))

	return nil
}
