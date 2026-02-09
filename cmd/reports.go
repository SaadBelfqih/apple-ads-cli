package cmd

import (
	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
	"github.com/spf13/cobra"
)

var reportsCmd = &cobra.Command{
	Use:   "reports",
	Short: "Generate reports",
}

func buildReportRequest(cmd *cobra.Command) (*types.ReportingRequest, error) {
	startTime, _ := cmd.Flags().GetString("start-time")
	endTime, _ := cmd.Flags().GetString("end-time")
	granularity, _ := cmd.Flags().GetString("granularity")
	groupBy, _ := cmd.Flags().GetString("group-by")
	selectorJSON, _ := cmd.Flags().GetString("selector-json")

	req := &types.ReportingRequest{
		StartTime:         startTime,
		EndTime:           endTime,
		Granularity:       granularity,
		ReturnRowTotals:   true,
		ReturnGrandTotals: true,
	}

	if groupBy != "" {
		req.GroupBy = []string{groupBy}
	}

	if selectorJSON != "" {
		sel, err := parseSelectorJSON(selectorJSON)
		if err != nil {
			return nil, err
		}
		req.Selector = sel
	} else {
		req.Selector = &types.Selector{
			OrderBy: []*types.Sorting{
				{Field: "localSpend", SortOrder: "DESCENDING"},
			},
		}
	}

	return req, nil
}

func addReportFlags(cmd *cobra.Command) {
	cmd.Flags().String("start-time", "", "Start time (YYYY-MM-DD)")
	cmd.MarkFlagRequired("start-time")
	cmd.Flags().String("end-time", "", "End time (YYYY-MM-DD)")
	cmd.MarkFlagRequired("end-time")
	cmd.Flags().String("granularity", "DAILY", "HOURLY, DAILY, WEEKLY, or MONTHLY")
	cmd.Flags().String("group-by", "", "Group by dimension (e.g., countryOrRegion)")
	cmd.Flags().String("selector-json", "", "Selector JSON for filtering")
}

var reportsCampaignsCmd = &cobra.Command{
	Use:   "campaigns",
	Short: "Campaign-level reports",
	RunE: func(cmd *cobra.Command, args []string) error {
		req, err := buildReportRequest(cmd)
		if err != nil {
			return err
		}
		result, err := apiClient.Reports().Campaigns(req)
		if err != nil {
			return err
		}
		return printRawJSON(result)
	},
}

var reportsAdGroupsCmd = &cobra.Command{
	Use:   "adgroups",
	Short: "Ad group-level reports",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		req, err := buildReportRequest(cmd)
		if err != nil {
			return err
		}
		result, err := apiClient.Reports().AdGroups(campaignID, req)
		if err != nil {
			return err
		}
		return printRawJSON(result)
	},
}

var reportsKeywordsCmd = &cobra.Command{
	Use:   "keywords",
	Short: "Keyword-level reports",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		agID, _ := cmd.Flags().GetInt64("adgroup-id")
		req, err := buildReportRequest(cmd)
		if err != nil {
			return err
		}
		var adGroupID *int64
		if agID > 0 {
			adGroupID = &agID
		}
		result, err := apiClient.Reports().Keywords(campaignID, adGroupID, req)
		if err != nil {
			return err
		}
		return printRawJSON(result)
	},
}

var reportsSearchTermsCmd = &cobra.Command{
	Use:   "searchterms",
	Short: "Search term-level reports",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		agID, _ := cmd.Flags().GetInt64("adgroup-id")
		req, err := buildReportRequest(cmd)
		if err != nil {
			return err
		}
		var adGroupID *int64
		if agID > 0 {
			adGroupID = &agID
		}
		result, err := apiClient.Reports().SearchTerms(campaignID, adGroupID, req)
		if err != nil {
			return err
		}
		return printRawJSON(result)
	},
}

var reportsAdsCmd = &cobra.Command{
	Use:   "ads",
	Short: "Ad-level reports",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		req, err := buildReportRequest(cmd)
		if err != nil {
			return err
		}
		result, err := apiClient.Reports().Ads(campaignID, req)
		if err != nil {
			return err
		}
		return printRawJSON(result)
	},
}

func init() {
	rootCmd.AddCommand(reportsCmd)

	addReportFlags(reportsCampaignsCmd)
	reportsCmd.AddCommand(reportsCampaignsCmd)

	addReportFlags(reportsAdGroupsCmd)
	reportsAdGroupsCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	reportsAdGroupsCmd.MarkFlagRequired("campaign-id")
	reportsCmd.AddCommand(reportsAdGroupsCmd)

	addReportFlags(reportsKeywordsCmd)
	reportsKeywordsCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	reportsKeywordsCmd.MarkFlagRequired("campaign-id")
	reportsKeywordsCmd.Flags().Int64("adgroup-id", 0, "Ad group ID (optional, scopes to ad group)")
	reportsCmd.AddCommand(reportsKeywordsCmd)

	addReportFlags(reportsSearchTermsCmd)
	reportsSearchTermsCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	reportsSearchTermsCmd.MarkFlagRequired("campaign-id")
	reportsSearchTermsCmd.Flags().Int64("adgroup-id", 0, "Ad group ID (optional, scopes to ad group)")
	reportsCmd.AddCommand(reportsSearchTermsCmd)

	addReportFlags(reportsAdsCmd)
	reportsAdsCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	reportsAdsCmd.MarkFlagRequired("campaign-id")
	reportsCmd.AddCommand(reportsAdsCmd)
}
