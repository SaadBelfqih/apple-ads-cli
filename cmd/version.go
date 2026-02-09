package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// These are overridden at build time via -ldflags.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown" // UTC ISO-8601 recommended
)

func versionLine() string {
	v := strings.TrimSpace(Version)
	if v == "" {
		v = "dev"
	}

	parts := []string{fmt.Sprintf("aads %s", v)}
	if c := strings.TrimSpace(Commit); c != "" && c != "none" {
		parts = append(parts, "("+c+")")
	}
	if d := strings.TrimSpace(Date); d != "" && d != "unknown" {
		parts = append(parts, d)
	}
	parts = append(parts, runtime.Version())
	return strings.Join(parts, " ")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versionLine())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
