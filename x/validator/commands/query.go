package commands

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/validator/model"
)

// GetValidatorsCmd returns all validators relative information
func GetValidatorsCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "validators",
		Short: "Query validators",
		RunE:  cmdr.getValidatorsCmd,
	}
}

// GetValidatorCmd returns target validator information
func GetValidatorCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "validator",
		Short: "Query validator",
		RunE:  cmdr.getValidatorCmd,
	}
}

type commander struct {
	storeName string
	cdc       *wire.Codec
}

func (c commander) getValidatorsCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	res, err := ctx.Query(model.GetValidatorListKey(), c.storeName)
	if err != nil {
		return err
	}

	validatorList := new(model.ValidatorList)
	if err := c.cdc.UnmarshalJSON(res, validatorList); err != nil {
		return err
	}

	// print out whole bank
	output, err := json.MarshalIndent(validatorList, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))

	return nil
}

func (c commander) getValidatorCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide a username")
	}

	accKey := types.AccountKey(args[0])

	res, err := ctx.Query(model.GetValidatorKey(accKey), c.storeName)
	if err != nil {
		return err
	}
	validator := new(model.Validator)
	if err := c.cdc.UnmarshalJSON(res, validator); err != nil {
		return err
	}

	output, err := json.MarshalIndent(validator, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
