package app

import (
	"fmt"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var Version = ""

// VersionCmd - print current version of binary.
func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "version prints the version of this binary",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Println(Version)
			return nil
		},
	}

	return cmd
}
