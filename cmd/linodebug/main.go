package main

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(hackCmd)
}

var rootCmd = &cobra.Command{
	Use:          "gaiadebug",
	Short:        "Gaia debug tool",
	SilenceUsage: true,
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
