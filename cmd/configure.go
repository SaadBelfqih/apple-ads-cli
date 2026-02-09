package cmd

import (
	"github.com/SaadBelfqih/apple-ads-cli/internal/config"
	"github.com/spf13/cobra"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Interactive setup for Apple Ads API credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.RunInteractiveSetup()
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
