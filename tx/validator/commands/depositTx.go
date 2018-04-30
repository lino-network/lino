package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	"github.com/lino-network/lino/tx/validator"
	"github.com/lino-network/lino/types"

	sdkcli "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/wire"
	cfg "github.com/tendermint/tendermint/config"
	pvm "github.com/tendermint/tendermint/types/priv_validator"
	cmn "github.com/tendermint/tmlibs/common"
)

const (
	FlagAmount = "amount"
)

// DepositValidatorTxCmd will create a send tx and sign it with the given key
func DepositValidatorTxCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-deposit",
		Short: "register a validator",
		RunE:  sendDepositValidatorTx(cdc),
	}
	cmd.Flags().String(FlagAmount, "", "amount of the donation")
	return cmd
}

// send register transaction to the blockchain
func sendDepositValidatorTx(cdc *wire.Codec) client.CommandTxCallback {
	return func(cmd *cobra.Command, args []string) error {
		ctx := client.NewCoreContextFromViper()
		name := viper.GetString(sdkcli.FlagName)

		config := cfg.DefaultConfig()
		// private validator
		privValFile := config.PrivValidatorFile()
		var privValidator *pvm.FilePV
		if cmn.FileExists(privValFile) {
			privValidator = pvm.LoadFilePV(privValFile)
		} else {
			privValidator = pvm.GenFilePV(privValFile)
			privValidator.Save()
		}

		// // create the message
		msg := validator.NewValidatorDepositMsg(name, types.LNO(viper.GetString(FlagAmount)), privValidator.PubKey)

		// build and sign the transaction, then broadcast to Tendermint
		res, err := ctx.SignBuildBroadcast(name, msg, cdc)

		if err != nil {
			return err
		}

		fmt.Printf("Committed at block %d. Hash: %s\n", res.Height, res.Hash.String())
		return nil
	}
}
