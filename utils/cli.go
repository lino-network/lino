package utils

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
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
