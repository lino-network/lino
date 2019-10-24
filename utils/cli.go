package utils

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func CLIQueryJSONPrint(cdc *codec.Codec, uri string, data []byte, rstTypeFactory func() interface{}) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	res, _, err := cliCtx.QueryWithData(uri, data)
	if err != nil {
		fmt.Printf("Failed to Query and Print: %s, because %s", uri, err)
		return nil
	}

	rst := rstTypeFactory()
	cdc.MustUnmarshalJSON(res, rst)
	out, err := cdc.MarshalJSONIndent(rst, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func SimpleQueryCmd(use, short, route, store string, nargs int, rstPointer interface{}) func(*codec.Codec) *cobra.Command {
	return func(cdc *codec.Codec) *cobra.Command {
		return &cobra.Command{
			Use:   use,
			Short: short,
			Args:  cobra.ExactArgs(nargs),
			RunE: func(cmd *cobra.Command, args []string) error {
				uri := fmt.Sprintf("custom/%s/%s", route, store)
				for i := 0; i < nargs; i++ {
					uri += "/" + args[i]
				}
				return CLIQueryJSONPrint(cdc, uri, nil,
					func() interface{} { return rstPointer })
			},
		}
	}
}
