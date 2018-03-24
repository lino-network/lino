package client

import (
	"github.com/spf13/cobra"
)

type CommandTxCallback func(cmd *cobra.Command, args []string) error
