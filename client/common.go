package client

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

type CommandTxCallback func(cmd *cobra.Command, args []string) error

func PrintIndent(inputs ...interface{}) error {
	for _, input := range inputs {
		output, err := json.MarshalIndent(input, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
	}
	return nil
}
