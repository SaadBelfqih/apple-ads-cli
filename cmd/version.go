package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.1.0-alpha"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("aads %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
