package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SaadBelfqih/apple-ads-cli/internal/api"
	"github.com/SaadBelfqih/apple-ads-cli/internal/config"
	"github.com/SaadBelfqih/apple-ads-cli/internal/output"
	"github.com/SaadBelfqih/apple-ads-cli/internal/updatecheck"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	verbose      bool
	orgIDFlag    string
	fieldsFlag   string
	currencyFlag string

	apiClient *api.Client

	activeOrgID               string
	defaultCurrencyFromConfig string
)

var rootCmd = &cobra.Command{
	Use:   "aads",
	Short: "Apple Ads CLI (Campaign Management API v5)",
	Long:  "A command-line interface for Apple's Apple Ads Campaign Management API v5, with safe retries, auto-pagination, and multiple output formats.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Lightweight update check (cached). Never installs anything, only prints a notice.
		// Skip for version/help/configure/update to avoid noisy output.
		switch cmd.Name() {
		case "configure", "version", "help", "update":
			// no-op
		default:
			updatecheck.MaybeNotify(cmd.Context(), os.Stderr, Version)
		}

		// Skip client init for commands that don't need it
		if cmd.Name() == "configure" || cmd.Name() == "version" || cmd.Name() == "help" {
			return nil
		}
		// Also skip for parent commands (e.g., "campaigns" without subcommand)
		if !cmd.HasParent() || cmd.HasSubCommands() && len(args) == 0 {
			return nil
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w\nRun 'aads configure' to set up", err)
		}

		if orgIDFlag != "" {
			cfg.OrgID = orgIDFlag
		}
		activeOrgID = cfg.OrgID
		defaultCurrencyFromConfig = cfg.DefaultCurrency

		requiresOrgID := true
		if cmd.Parent() != nil && cmd.Parent().Name() == "acls" {
			// Per Apple docs, X-AP-Context isn't required for Get User ACL and Get Me Details.
			requiresOrgID = false
		}

		if requiresOrgID {
			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("invalid config: %w\nRun 'aads configure' to fix (or set env vars like AADS_CLIENT_ID, AADS_TEAM_ID, AADS_KEY_ID, AADS_ORG_ID, AADS_PRIVATE_KEY_PATH)", err)
			}
		} else {
			if err := cfg.ValidateAuth(); err != nil {
				return fmt.Errorf("invalid config: %w\nRun 'aads configure' to fix (or set env vars like AADS_CLIENT_ID, AADS_TEAM_ID, AADS_KEY_ID, AADS_PRIVATE_KEY_PATH)", err)
			}
		}

		client, err := api.NewClient(cfg)
		if err != nil {
			return fmt.Errorf("init client: %w", err)
		}

		client.SetVerbose(verbose)
		apiClient = client
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Enable `aads --version` (no shorthand to avoid conflict with `-v/--verbose`).
	rootCmd.Version = versionLine()
	rootCmd.SetVersionTemplate("{{.Version}}\n")

	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "Output format: json, table, yaml")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVar(&orgIDFlag, "org-id", "", "Override org ID from config")
	rootCmd.PersistentFlags().StringVar(&fieldsFlag, "fields", "", "Comma-separated fields for partial fetch")
	rootCmd.PersistentFlags().StringVar(&currencyFlag, "currency", "", "Override currency for money fields (e.g., USD)")
}

func getOutputFormat() output.Format {
	switch outputFormat {
	case "table":
		return output.FormatTable
	case "yaml":
		return output.FormatYAML
	default:
		return output.FormatJSON
	}
}

func printOutput(data any) error {
	return output.Print(os.Stdout, getOutputFormat(), data)
}

func printRawJSON(data []byte) error {
	// Parse and re-output through our formatter for consistent formatting
	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		// If we can't parse, just write raw
		os.Stdout.Write(data)
		fmt.Fprintln(os.Stdout)
		return nil
	}
	return output.Print(os.Stdout, getOutputFormat(), parsed)
}
