package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lino-network/lino/client"
	linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/x/post/types"
)

const (
	FlagAuthor  = "author"
	FlagPostID  = "post-id"
	FlagTitle   = "title"
	FlagContent = "content"
	FlagPreauth = "preauth"

	FlagDonator = "donator"
	FlagAmount  = "amount"
	FlagMemo    = "memo"
	FlagApp     = "app"
	FlagSigner  = "signer"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Post tx subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(client.PostCommands(
		GetCmdCreatePost(cdc),
		GetCmdDeletePost(cdc),
		GetCmdUpdatePost(cdc),
		GetCmdDonate(cdc),
		GetCmdIDADonate(cdc),
	)...)

	return cmd
}

// GetCmdCreatePost - create post.
func GetCmdCreatePost(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create ",
		Short: "create <created-by> --author <author> --post-id <id> --title <title> --content <content> --preauth=true/false",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			createdBy := linotypes.AccountKey(args[0])
			msg := types.CreatePostMsg{
				Author:    linotypes.AccountKey(viper.GetString(FlagAuthor)),
				PostID:    viper.GetString(FlagPostID),
				Title:     viper.GetString(FlagTitle),
				Content:   viper.GetString(FlagContent),
				CreatedBy: createdBy,
				Preauth:   viper.GetBool(FlagPreauth),
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagAuthor, "", "author of this post")
	cmd.Flags().String(FlagPostID, "", "post id to identify this post for the author")
	cmd.Flags().String(FlagTitle, "", "title for the post")
	cmd.Flags().String(FlagContent, "", "content for the post")
	cmd.Flags().Bool(FlagPreauth, false, "application(developer) that creates the post")
	for _, v := range []string{FlagAuthor, FlagPostID, FlagPreauth} {
		_ = cmd.MarkFlagRequired(v)
	}
	return cmd
}

// GetCmdUpdatePost - update post info.
func GetCmdUpdatePost(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update <author> <postid> --title <title> --content <content>",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			author := linotypes.AccountKey(args[0])
			postid := args[1]

			msg := types.UpdatePostMsg{
				Author:  author,
				PostID:  postid,
				Title:   viper.GetString(FlagTitle),
				Content: viper.GetString(FlagContent),
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagTitle, "", "title for the post")
	cmd.Flags().String(FlagContent, "", "content for the post")
	return cmd
}

// GetCmdDeletePost -
func GetCmdDeletePost(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete <author> <postid>",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			author := linotypes.AccountKey(args[0])
			postid := args[1]
			msg := types.DeletePostMsg{
				Author: author,
				PostID: postid,
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	return cmd
}

// GetCmdDonate -
func GetCmdDonate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "donate",
		Short: "donate <donator> --amount <amount> --author <author> --post-id <id> --app <app> --memo <memo>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			from := linotypes.AccountKey(args[0])
			msg := types.DonateMsg{
				Username: from,
				Amount:   viper.GetString(FlagAmount),
				Author:   linotypes.AccountKey(viper.GetString(FlagAuthor)),
				PostID:   viper.GetString(FlagPostID),
				FromApp:  linotypes.AccountKey(viper.GetString(FlagApp)),
				Memo:     viper.GetString(FlagMemo),
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagAuthor, "", "author of this post")
	cmd.Flags().String(FlagPostID, "", "post id to identify this post for the author")
	cmd.Flags().String(FlagAmount, "", "amount of the donation")
	cmd.Flags().String(FlagMemo, "", "memo of this donation")
	cmd.Flags().String(FlagApp, "", "donation comes from app")
	for _, v := range []string{FlagAuthor, FlagPostID, FlagAmount} {
		_ = cmd.MarkFlagRequired(v)
	}
	return cmd
}

// GetCmdIDADonate -
func GetCmdIDADonate(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ida-donate",
		Short: "ida-donate <signer> --amount <amount> --author <author> --post-id <id> --app <app> --memo <memo> --donator <donator>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.NewCoreContextFromViper().WithTxEncoder(linotypes.TxEncoder(cdc))
			signer := linotypes.AccountKey(args[0])
			msg := types.IDADonateMsg{
				Username: linotypes.AccountKey(viper.GetString(FlagDonator)),
				App:      linotypes.AccountKey(viper.GetString(FlagApp)),
				Amount:   linotypes.IDAStr(viper.GetString(FlagAmount)),
				Author:   linotypes.AccountKey(viper.GetString(FlagAuthor)),
				PostID:   viper.GetString(FlagPostID),
				Memo:     viper.GetString(FlagMemo),
				Signer:   signer,
			}
			return ctx.DoTxPrintResponse(msg)
		},
	}
	cmd.Flags().String(FlagDonator, "", "donator of this transaction")
	cmd.Flags().String(FlagAuthor, "", "author of the target post")
	cmd.Flags().String(FlagPostID, "", "post id of the target post")
	cmd.Flags().String(FlagAmount, "", "amount of the donation")
	cmd.Flags().String(FlagMemo, "", "memo of this donation")
	cmd.Flags().String(FlagApp, "", "App's IDA")
	for _, v := range []string{FlagDonator, FlagAuthor, FlagPostID, FlagAmount, FlagApp} {
		_ = cmd.MarkFlagRequired(v)
	}
	return cmd
}
