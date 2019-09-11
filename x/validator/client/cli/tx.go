package cli

import (
	"os/user"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	pvm "github.com/tendermint/tendermint/privval"

	"github.com/lino-network/lino/client"
	linotypes "github.com/lino-network/lino/types"
	types "github.com/lino-network/lino/x/validator"
)

const (
	FlagUser   = "user"
	FlagAmount = "amount"
	FlagLink   = "link"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "validator tx subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(client.PostCommands(
		GetCmdDeposit(cdc),
		GetCmdRevoke(cdc),
		GetCmdWithdraw(cdc),
	)...)

	return cmd
}

// GetCmdDeposit -
func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit",
		Short: "deposit user --amount <amount> --link <link>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			validator := args[0]
			amount := linotypes.LNO(viper.GetString(FlagAmount))
			link := viper.GetString(FlagLink)
			pubKey, err := getLocalUserPubKey()
			if err != nil {
				return err
			}
			msg := types.NewValidatorDepositMsg(validator, amount, pubKey, link)
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagAmount, "", "amount of the donation")
	cmd.Flags().String(FlagLink, "", "link of the validator")
	for _, v := range []string{FlagUser, FlagAmount, FlagLink} {
		cmd.MarkFlagRequired(v)
	}
	return cmd
}

func getLocalUserPubKey() (crypto.PubKey, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	root := usr.HomeDir + "/.lino/"

	tmConfig := cfg.DefaultConfig()
	tmConfig = tmConfig.SetRoot(root)

	privValFile := tmConfig.PrivValidatorKeyFile()
	privValStateFile := tmConfig.PrivValidatorStateFile()

	var privValidator *pvm.FilePV
	if cmn.FileExists(privValFile) {
		privValidator = pvm.LoadFilePV(privValFile, privValStateFile)
	} else {
		privValidator = pvm.GenFilePV(privValFile, privValStateFile)
		privValidator.Save()
	}
	pubKey := privValidator.GetPubKey()
	return pubKey, nil
}

// GetCmdRevoke -
func GetCmdRevoke(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke",
		Short: "revoke <username>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			user := args[0]
			msg := types.NewValidatorRevokeMsg(user)
			return ctx.DoTxPrintResponse(msg)
		},
	}
	return cmd
}

// GetCmdWithdraw -
func GetCmdWithdraw(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "withdraw <username> --amount <amount>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			user := args[0]
			msg := types.NewValidatorWithdrawMsg(user, viper.GetString(FlagAmount))
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagAmount, "", "amount of the donation")
	cmd.MarkFlagRequired(FlagAmount)
	return cmd
}
