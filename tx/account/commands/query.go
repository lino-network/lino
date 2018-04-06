package commands

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/account/model"
	"github.com/lino-network/lino/types"
)

// GetBankCmd returns a query bank that will display the
// state of the bank at a given address
func GetBankCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "address <address>",
		Short: "Query bank balance",
		RunE:  cmdr.getBankCmd,
	}
}

// GetBankCmd returns a query bank that will display the
// state of the bank at a given address
func GetAccountCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "username <username>",
		Short: "Query account",
		RunE:  cmdr.getAccountCmd,
	}
}

type commander struct {
	storeName string
	cdc       *wire.Codec
}

func (c commander) getBankCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide an address")
	}

	// find the key to look up the account
	addr := args[0]
	bz, err := hex.DecodeString(addr)
	if err != nil {
		return err
	}
	key := sdk.Address(bz)

	res, err := builder.Query(model.GetAccountBankKey(key), c.storeName)
	if err != nil {
		return err
	}

	bank := new(model.AccountBank)
	if err := c.cdc.UnmarshalBinary(res, bank); err != nil {
		return err
	}

	// print out whole bank
	output, err := json.MarshalIndent(bank, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))

	return nil
}

func (c commander) getAccountCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide aa username")
	}

	// find the key to look up the account
	accKey := types.AccountKey(args[0])

	res, err := builder.Query(model.GetAccountInfoKey(accKey), c.storeName)
	if err != nil {
		return err
	}
	info := new(model.AccountInfo)
	if err := c.cdc.UnmarshalBinary(res, info); err != nil {
		return err
	}

	res, err = builder.Query(model.GetAccountBankKey(info.Address), c.storeName)
	if err != nil {
		return err
	}
	bank := new(model.AccountBank)
	if err := c.cdc.UnmarshalBinary(res, bank); err != nil {
		return err
	}

	res, err = builder.Query(model.GetAccountMetaKey(accKey), c.storeName)
	if err != nil {
		return err
	}
	meta := new(model.AccountMeta)
	if err := c.cdc.UnmarshalBinary(res, meta); err != nil {
		return err
	}

	if err := client.PrintIndent(info, bank, meta); err != nil {
		return err
	}
	return nil
}
