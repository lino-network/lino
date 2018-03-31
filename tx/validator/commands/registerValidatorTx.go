package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/validator"
	"github.com/lino-network/lino/types"

	sdkcli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/builder"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	tmtypes "github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
)

const (
	FlagAmount = "amount"
)

// SendTxCommand will create a send tx and sign it with the given key
func RegisterValidatorTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "regval",
		Short: "register a validator",
		RunE:  sendRegisterValidatorTx(cdc),
	}
	cmd.Flags().String(FlagAmount, "", "amount of the donation")
	return cmd
}

// send register transaction to the blockchain
func sendRegisterValidatorTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		name := viper.GetString(sdkcli.FlagName)

		config, err := tcmd.ParseConfig()
		if err != nil {
			return err
		}
		// private validator
		privValFile := config.PrivValidatorFile()
		var privValidator *tmtypes.PrivValidatorFS
		if cmn.FileExists(privValFile) {
			privValidator = tmtypes.LoadPrivValidatorFS(privValFile)
		} else {
			privValidator = tmtypes.GenPrivValidatorFS(privValFile)
			privValidator.Save()
		}

		if err != nil {
			return err
		}

		amount, err := sdk.NewRatFromDecimal(viper.GetString(FlagAmount))
		if err != nil {
			return err
		}
		// // create the message
		msg := validator.NewValidatorDepositMsg(name, types.LNO(amount), privValidator.PubKey)

		// build and sign the transaction, then broadcast to Tendermint
		res, err := builder.SignBuildBroadcast(name, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
