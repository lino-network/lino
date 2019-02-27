package commands

import (
	"encoding/json"
	"fmt"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/account/model"

	wire "github.com/cosmos/cosmos-sdk/codec"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// GetBankCmd returns a query bank that will display the
// state of the bank at a given address
func GetBankCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "username <username>",
		Short: "Query bank balance",
		RunE:  cmdr.getBankCmd,
	}
}

// GetAccountCmd returns a query account that will display the
// state of the account at a given username
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

// GetAccountCmd returns a query account that will display the
// state of the account at a given username
func GetAccountsCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmdr := commander{
		storeName,
		cdc,
	}
	return &cobra.Command{
		Use:   "accounts",
		Short: "Query all accounts",
		RunE:  cmdr.getAccountsCmd,
	}
}

type commander struct {
	storeName string
	cdc       *wire.Codec
}

func (c commander) getBankCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide an address")
	}

	username := types.AccountKey(args[0])

	res, err := ctx.Query(model.GetAccountBankKey(username), c.storeName)
	if err != nil {
		return err
	}

	bank := new(model.AccountBank)
	if err := c.cdc.UnmarshalBinaryLengthPrefixed(res, bank); err != nil {
		return err
	}

	output, err := json.MarshalIndent(bank, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(output))

	return nil
}

func (c commander) getAccountCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()
	if len(args) != 1 || len(args[0]) == 0 {
		return errors.New("You must provide aa username")
	}

	// find the key to look up the account
	accKey := types.AccountKey(args[0])

	res, err := ctx.Query(model.GetAccountInfoKey(accKey), c.storeName)
	if err != nil {
		return err
	}
	info := new(model.AccountInfo)
	if err := c.cdc.UnmarshalBinaryLengthPrefixed(res, info); err != nil {
		return err
	}

	res, err = ctx.Query(model.GetAccountBankKey(accKey), c.storeName)
	if err != nil {
		return err
	}
	bank := new(model.AccountBank)
	if err := c.cdc.UnmarshalBinaryLengthPrefixed(res, bank); err != nil {
		return err
	}

	res, err = ctx.Query(model.GetAccountMetaKey(accKey), c.storeName)
	if err != nil {
		return err
	}
	meta := new(model.AccountMeta)
	if err := c.cdc.UnmarshalBinaryLengthPrefixed(res, meta); err != nil {
		return err
	}

	if err := client.PrintIndent(info, bank, meta); err != nil {
		return err
	}
	return nil
}

func (c commander) getAccountsCmd(cmd *cobra.Command, args []string) error {
	ctx := client.NewCoreContextFromViper()

	resKVs, err := ctx.QuerySubspace(c.cdc, model.GetAccountInfoPrefix(), c.storeName)
	if err != nil {
		return err
	}
	var accounts []model.AccountInfo
	for _, KV := range resKVs {
		var info model.AccountInfo
		if err := c.cdc.UnmarshalBinaryLengthPrefixed(KV.Value, &info); err != nil {
			return err
		}
		accounts = append(accounts, info)
	}

	if err := client.PrintIndent(accounts); err != nil {
		return err
	}
	return nil
}
