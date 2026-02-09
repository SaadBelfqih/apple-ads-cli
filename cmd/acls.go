package cmd

import (
	"github.com/spf13/cobra"
)

var aclsCmd = &cobra.Command{
	Use:   "acls",
	Short: "Manage ACLs and user info",
}

var aclsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List user ACLs",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := apiClient.ACLs().List()
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var aclsMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Get caller details",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := apiClient.ACLs().Me()
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

func init() {
	rootCmd.AddCommand(aclsCmd)
	aclsCmd.AddCommand(aclsListCmd)
	aclsCmd.AddCommand(aclsMeCmd)
}
