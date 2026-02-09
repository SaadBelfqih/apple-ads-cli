package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/SaadBelfqih/apple-ads-cli/internal/updatecheck"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for updates (no auto-install)",
	Long:  "Checks GitHub Releases for a newer version and prints the result. This command does not download or install updates.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := updatecheck.CheckLatest(context.Background(), os.Stdout, Version); err != nil {
			if strings.Contains(err.Error(), "404") {
				fmt.Fprintln(os.Stderr, "Update check unavailable (repo may be private). If you're a contributor, set AADS_GITHUB_TOKEN or GH_TOKEN and try again.")
				return nil
			}
			if strings.Contains(err.Error(), "no releases found") {
				fmt.Fprintln(os.Stdout, "No releases found yet.")
				fmt.Fprintln(os.Stdout, "Create a tag like v0.1.0-alpha.1 and push it to publish your first GitHub Release.")
				return nil
			}
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
