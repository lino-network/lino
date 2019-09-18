package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/developer/types"
)

const (
	FlagWebsite     = "website"
	FlagDescription = "description"
	FlagAppMeta     = "appmeta"
	FlagIdaPrice    = "ida-price"
	FlagFrom        = "from"
	FlagTo          = "to"
	FlagActive      = "active"
	FlagAmount      = "amount"
	FlagApp         = "app"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Developer tx subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(client.PostCommands(
		GetCmdRegister(cdc),
		GetCmdUpdate(cdc),
		GetCmdIDAIssue(cdc),
		GetCmdIDAMint(cdc),
		GetCmdIDATransfer(cdc),
		GetCmdIDAAuthorize(cdc),
		GetCmdUpdateAffiliated(cdc),
	)...)

	return cmd
}

// GetCmdRegister - register as developer.
func GetCmdRegister(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "register <username> --website <web> --description <desc> --appmeta <meta>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			app := args[0]
			msg := types.NewDeveloperRegisterMsg(
				app,
				viper.GetString(FlagWebsite),
				viper.GetString(FlagDescription),
				viper.GetString(FlagAppMeta))
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagWebsite, "", "website of the app")
	cmd.Flags().String(FlagDescription, "", "description of the app")
	cmd.Flags().String(FlagAppMeta, "", "meta-data of the app")
	return cmd
}

// GetCmdUpdate - update App info
func GetCmdUpdate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update <app> --website <web> --description <desc> --appmeta <meta>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			app := linotypes.AccountKey(args[0])
			msg := types.DeveloperUpdateMsg{
				Username:    app,
				Website:     viper.GetString(FlagWebsite),
				Description: viper.GetString(FlagDescription),
				AppMetaData: viper.GetString(FlagAppMeta),
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagWebsite, "", "website of the app")
	cmd.Flags().String(FlagDescription, "", "description of the app")
	cmd.Flags().String(FlagAppMeta, "", "meta-data of the app")
	return cmd
}

// GetCmdIDAIssue - issue an IDA for the app.
// ida-issue dlivetv --ida-price 1
func GetCmdIDAIssue(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ida-issue",
		Short: "ida-issue <app> --ida-price <amount>, amount is an integer in [1, 1000], 1 = 0.001 USD",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			app := linotypes.AccountKey(args[0])
			msg := types.IDAIssueMsg{
				Username: app,
				IDAPrice: viper.GetInt64(FlagIdaPrice),
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().Int64(FlagIdaPrice, 0,
		"The price of IDA in unit of 0.001 USD, valid range: [1, 1000]")
	return cmd
}

// GetCmdIDAMint - mint IDA for the app
func GetCmdIDAMint(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ida-mint",
		Short: "ida-mint <app> --amount <lino-amount>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			app := linotypes.AccountKey(args[0])
			amount := viper.GetString(FlagAmount)
			msg := types.IDAMintMsg{
				Username: app,
				Amount:   amount,
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().Int64(FlagAmount, 0,
		"Amount of LINO to convert to IDA")
	_ = cmd.MarkFlagRequired(FlagAmount)
	return cmd
}

// GetCmdIDATransfer - transfer ida from or to some account.
func GetCmdIDATransfer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ida-transfer",
		Short: "ida-transfer <signer> <app> --from <foo> --to <bar> --amount <amount>",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			signer := linotypes.AccountKey(args[0])
			app := linotypes.AccountKey(args[1])
			amount := linotypes.IDAStr(viper.GetString(FlagAmount))
			from := linotypes.AccountKey(viper.GetString(FlagFrom))
			to := linotypes.AccountKey(viper.GetString(FlagTo))
			msg := types.IDATransferMsg{
				App:    app,
				Amount: amount,
				From:   from,
				To:     to,
				Signer: signer,
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagAmount, "", "amount of IDA")
	cmd.Flags().String(FlagTo, "", "receipient of this transfer")
	cmd.Flags().String(FlagFrom, "", "sender of this transfer")
	return cmd
}

// GetCmdIDAAuthorize -
func GetCmdIDAAuthorize(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ida-auth",
		Short: "ida-auth <username> <app> --active=true/false",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			user := linotypes.AccountKey(args[0])
			app := linotypes.AccountKey(args[1])
			active := viper.GetBool(FlagActive)
			msg := types.IDAAuthorizeMsg{
				Username: user,
				App:      app,
				Activate: active,
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().Bool(FlagActive, false, "true = active IDA account")
	_ = cmd.MarkFlagRequired(FlagActive)
	return cmd
}

// GetCmdUpdateAffiliated -
func GetCmdUpdateAffiliated(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "affiliated",
		Short: "affiliated <app> <username> --active=true/false",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			app := linotypes.AccountKey(args[0])
			user := linotypes.AccountKey(args[1])
			active := viper.GetBool(FlagActive)
			msg := types.UpdateAffiliatedMsg{
				App:      app,
				Username: user,
				Activate: active,
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().Bool(FlagActive, false, "true = add <username> as an affiliated of <app>")
	_ = cmd.MarkFlagRequired(FlagActive)
	return cmd
}
