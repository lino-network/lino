package commands

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/builder"
	"github.com/cosmos/cosmos-sdk/wire"
	acc "github.com/lino-network/lino/tx/account"
	"github.com/lino-network/lino/tx/validator"
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
	res, err := builder.Query(validator.ValidatorListKey(), c.storeName)
	if err != nil {
		return err
	}

	validatorList := new(validator.ValidatorList)
	if err := c.cdc.UnmarshalBinary(res, validatorList); err != nil {
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
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide a username")
	}

	// find the key to look up the account
	accKey := acc.AccountKey(args[0])

	res, err := builder.Query(validator.ValidatorKey(accKey), c.storeName)
	if err != nil {
		return err
	}
	validator := new(validator.Validator)
	if err := c.cdc.UnmarshalBinary(res, validator); err != nil {
		return err
	}

	// print out whole bank
	output, err := json.MarshalIndent(validator, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
