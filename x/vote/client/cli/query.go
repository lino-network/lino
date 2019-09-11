package cli

// import (
// 	"fmt"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/codec"
// 	"github.com/spf13/cobra"

// 	// linotypes "github.com/lino-network/lino/types"
// 	"github.com/lino-network/lino/utils"
// 	types "github.com/lino-network/lino/x/vote"
// 	"github.com/lino-network/lino/x/vote/model"
// )

// func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:                        types.ModuleName,
// 		Short:                      "Querying commands for the vote module",
// 		DisableFlagParsing:         true,
// 		SuggestionsMinimumDistance: 2,
// 		RunE:                       client.ValidateCmd,
// 	}
// 	cmd.AddCommand(client.GetCommands(
// 		getCmdInfo(cdc),
// 	)...)
// 	return cmd
// }
