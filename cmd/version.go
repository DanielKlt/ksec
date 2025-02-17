package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.0.1"

var VersionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "version command",
	Long:    "Print the version of ksec cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", version)
	},
}
