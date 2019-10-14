package cli

import (
	"os/user"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	pvm "github.com/tendermint/tendermint/privval"

	"github.com/lino-network/lino/client"
	linotypes "github.com/lino-network/lino/types"
	types "github.com/lino-network/lino/x/validator/types"
)

const (
	FlagUser       = "user"
	FlagAmount     = "amount"
	FlagLink       = "link"
	FlagValidators = "validators"
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
		GetCmdRegister(cdc),
		GetCmdRevoke(cdc),
		GetCmdVote(cdc),
	)...)

	return cmd

}

// GetCmdRegister -
func GetCmdRegister(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "register user --link <link>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			validator := args[0]
			link := viper.GetString(FlagLink)
			pubKey, err := getLocalUserPubKey()
			if err != nil {
				return err
			}
			msg := types.NewValidatorRegisterMsg(validator, pubKey, link)
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagLink, "", "link of the validator")
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

// GetCmdVote -
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote",
		Short: "vote voter --validators val1,val,...",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			voter := args[0]
			validators := strings.Split(viper.GetString(FlagValidators), ",")
			msg := types.NewVoteValidatorMsg(voter, validators)
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagValidators, "", "a comma-separated string, the list of validators")
	_ = cmd.MarkFlagRequired(FlagValidators)
	return cmd
}
