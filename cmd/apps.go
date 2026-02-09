package cmd

import (
	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Search and manage app info",
}

var appsSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for iOS apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		returnOwned, _ := cmd.Flags().GetBool("return-owned-apps")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			if limit == 0 && !cmd.Flags().Changed("limit") {
				limit = defaultPageSize
			}
			result, err := collectAllOffsetPaginated(limit, offset, func(lim, off int) ([]types.AppInfo, *types.PageDetail, error) {
				return apiClient.Apps().Search(query, returnOwned, lim, off)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Apps().Search(query, returnOwned, limit, offset)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var appsEligibilityCmd = &cobra.Command{
	Use:   "eligibility",
	Short: "Check app eligibility for Apple Ads",
	RunE: func(cmd *cobra.Command, args []string) error {
		selectorJSON, _ := cmd.Flags().GetString("selector-json")
		adamID, _ := cmd.Flags().GetString("adam-id")

		var sel *types.Selector
		if selectorJSON != "" {
			var err error
			sel, err = parseSelectorJSON(selectorJSON)
			if err != nil {
				return err
			}
		} else if adamID != "" {
			sel = &types.Selector{
				Conditions: []*types.Condition{
					{Field: "adamId", Operator: "EQUALS", Values: []string{adamID}},
				},
			}
		} else {
			sel = &types.Selector{}
		}

		result, err := apiClient.Apps().Eligibility(sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var appsDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Get app details",
	RunE: func(cmd *cobra.Command, args []string) error {
		adamID, _ := cmd.Flags().GetInt64("adam-id")
		result, err := apiClient.Apps().Details(adamID)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var appsLocalizedCmd = &cobra.Command{
	Use:   "localized",
	Short: "Get localized app details",
	RunE: func(cmd *cobra.Command, args []string) error {
		adamID, _ := cmd.Flags().GetInt64("adam-id")
		result, err := apiClient.Apps().LocalizedDetails(adamID)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

func init() {
	rootCmd.AddCommand(appsCmd)

	appsSearchCmd.Flags().String("query", "", "Search query")
	appsSearchCmd.MarkFlagRequired("query")
	appsSearchCmd.Flags().Bool("return-owned-apps", false, "Only return owned apps")
	appsSearchCmd.Flags().Int("limit", 0, "Max results")
	appsSearchCmd.Flags().Int("offset", 0, "Start offset")
	appsSearchCmd.Flags().Bool("all", false, "Fetch all pages")
	appsCmd.AddCommand(appsSearchCmd)

	appsEligibilityCmd.Flags().String("selector-json", "", "Selector JSON")
	appsEligibilityCmd.Flags().String("adam-id", "", "App Adam ID")
	appsCmd.AddCommand(appsEligibilityCmd)

	appsDetailsCmd.Flags().Int64("adam-id", 0, "App Adam ID")
	appsDetailsCmd.MarkFlagRequired("adam-id")
	appsCmd.AddCommand(appsDetailsCmd)

	appsLocalizedCmd.Flags().Int64("adam-id", 0, "App Adam ID")
	appsLocalizedCmd.MarkFlagRequired("adam-id")
	appsCmd.AddCommand(appsLocalizedCmd)
}
