package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/lino-network/lino/client"
	developer "github.com/lino-network/lino/tx/developer"
	"github.com/lino-network/lino/types"
)

const (
	FlagDeveloper = "developer"
	FlagDeposit   = "deposit"
)

func DeveloperRegisterTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "developer-register",
		Short: "developer register",
		RunE:  sendDeveloperRegisterTx(cdc),
	}
	cmd.Flags().String(FlagDeveloper, "", "developer name of this transaction")
	cmd.Flags().String(FlagDeposit, "", "deposit of the registration")
	return cmd
}

// send register transaction to the blockchain
func sendDeveloperRegisterTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.NewCoreContextFromViper()
		username := viper.GetString(FlagDeveloper)

		deposit, err := sdk.NewRatFromDecimal(viper.GetString(FlagDeposit))
		if err != nil {
			return err
		}
		msg := developer.NewDeveloperRegisterMsg(username, types.LNO(deposit))

		// build and sign the transaction, then broadcast to Tendermint
		res, signErr := ctx.SignBuildBroadcast(username, msg, cdc)
		if signErr != nil {
			return signErr
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
