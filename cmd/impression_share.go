package cmd

import (
	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
	"github.com/spf13/cobra"
)

var impressionShareCmd = &cobra.Command{
	Use:   "impression-share",
	Short: "Manage impression share reports",
}

var isCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an impression share report",
	RunE: func(cmd *cobra.Command, args []string) error {
		fromJSON, _ := cmd.Flags().GetString("from-json")

		var req types.CustomReportRequest
		if fromJSON != "" {
			if err := parseJSONInput(fromJSON, &req); err != nil {
				return err
			}
		} else {
			startTime, _ := cmd.Flags().GetString("start-time")
			endTime, _ := cmd.Flags().GetString("end-time")
			granularity, _ := cmd.Flags().GetString("granularity")

			req = types.CustomReportRequest{
				StartTime:   startTime,
				EndTime:     endTime,
				Granularity: granularity,
			}
		}

		result, err := apiClient.ImpressionShare().Create(&req)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var isGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an impression share report by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")
		result, err := apiClient.ImpressionShare().Get(id)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var isListCmd = &cobra.Command{
	Use:   "list",
	Short: "List impression share reports",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			if limit == 0 && !cmd.Flags().Changed("limit") {
				limit = defaultPageSize
			}
			result, err := collectAllOffsetPaginated(limit, offset, apiClient.ImpressionShare().List)
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.ImpressionShare().List(limit, offset)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

func init() {
	rootCmd.AddCommand(impressionShareCmd)

	// create
	isCreateCmd.Flags().String("from-json", "", "JSON input (inline, @file, or @- for stdin)")
	isCreateCmd.Flags().String("start-time", "", "Start time (YYYY-MM-DD)")
	isCreateCmd.Flags().String("end-time", "", "End time (YYYY-MM-DD)")
	isCreateCmd.Flags().String("granularity", "DAILY", "DAILY, WEEKLY, or MONTHLY")
	impressionShareCmd.AddCommand(isCreateCmd)

	// get
	isGetCmd.Flags().Int64("id", 0, "Report ID")
	isGetCmd.MarkFlagRequired("id")
	impressionShareCmd.AddCommand(isGetCmd)

	// list
	isListCmd.Flags().Int("limit", 0, "Max results")
	isListCmd.Flags().Int("offset", 0, "Start offset")
	isListCmd.Flags().Bool("all", false, "Fetch all pages")
	impressionShareCmd.AddCommand(isListCmd)
}
