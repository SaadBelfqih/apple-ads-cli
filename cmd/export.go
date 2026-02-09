package cmd

import "github.com/spf13/cobra"

// RootCommand returns the Cobra root command for this CLI.
// It is intended for documentation generation and integrations.
func RootCommand() *cobra.Command {
	return rootCmd
}
