package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	// linotypes "github.com/lino-network/lino/types"
	"github.com/lino-network/lino/utils"
	"github.com/lino-network/lino/x/vote/model"
	types "github.com/lino-network/lino/x/vote/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the vote module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(client.GetCommands(
		getCmdVoter(cdc),
		utils.SimpleQueryCmd(
			"stake-stats <day>", "stake-stats <day>",
			types.QuerierRoute, types.QueryStakeStats,
			1, &model.LinoStakeStat{})(cdc),
	)...)
	return cmd
}

// GetCmdVoter -
func getCmdVoter(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "voter",
		Short: "voter [username]",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			user := args[0]
			uri := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryVoter, user)
			rst := model.Voter{}
			return utils.CLIQueryJSONPrint(cdc, uri, nil,
				func() interface{} { return &rst })
		},
	}
}
