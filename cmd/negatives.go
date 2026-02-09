package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
	"github.com/spf13/cobra"
)

var negativesCmd = &cobra.Command{
	Use:   "negatives",
	Short: "Manage negative keywords (campaign and ad group level)",
}

// Campaign-level negative keywords

var negCampaignCreateCmd = &cobra.Command{
	Use:   "campaign-create",
	Short: "Create campaign-level negative keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		fromJSON, _ := cmd.Flags().GetString("from-json")
		text, _ := cmd.Flags().GetString("text")
		matchType, _ := cmd.Flags().GetString("match-type")

		var keywords []types.NegativeKeyword
		if fromJSON != "" {
			if err := parseJSONInput(fromJSON, &keywords); err != nil {
				return err
			}
		} else {
			keywords = []types.NegativeKeyword{{Text: text, MatchType: matchType}}
		}

		result, err := apiClient.Negatives().CampaignCreate(campaignID, keywords)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var negCampaignGetCmd = &cobra.Command{
	Use:   "campaign-get",
	Short: "Get a campaign-level negative keyword",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		id, _ := cmd.Flags().GetInt64("id")
		result, err := apiClient.Negatives().CampaignGet(campaignID, id)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var negCampaignListCmd = &cobra.Command{
	Use:   "campaign-list",
	Short: "List campaign-level negative keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			if limit == 0 && !cmd.Flags().Changed("limit") {
				limit = defaultPageSize
			}
			result, err := collectAllOffsetPaginated(limit, offset, func(lim, off int) ([]types.NegativeKeyword, *types.PageDetail, error) {
				return apiClient.Negatives().CampaignList(campaignID, lim, off)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Negatives().CampaignList(campaignID, limit, offset)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var negCampaignFindCmd = &cobra.Command{
	Use:   "campaign-find",
	Short: "Find campaign-level negative keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		selectorJSON, _ := cmd.Flags().GetString("selector-json")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all && !cmd.Flags().Changed("limit") {
			limit = defaultPageSize
		}

		sel := &types.Selector{}
		if selectorJSON != "" {
			parsed, err := parseSelectorJSON(selectorJSON)
			if err != nil {
				return err
			}
			sel = parsed
		}

		if sel.Pagination == nil && (limit > 0 || offset > 0) {
			sel.Pagination = &types.Pagination{}
		}
		if sel.Pagination != nil {
			if (all && sel.Pagination.Limit == 0) || cmd.Flags().Changed("limit") {
				sel.Pagination.Limit = limit
			}
			if cmd.Flags().Changed("offset") && sel.Pagination.Offset == 0 {
				sel.Pagination.Offset = offset
			}
		}

		if all {
			pageSize := 0
			if selectorJSON != "" {
				pageSize = limit
			}
			result, err := collectAllSelectorPaginated(sel, pageSize, func(s *types.Selector) ([]types.NegativeKeyword, *types.PageDetail, error) {
				return apiClient.Negatives().CampaignFind(campaignID, s)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Negatives().CampaignFind(campaignID, sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var negCampaignUpdateCmd = &cobra.Command{
	Use:   "campaign-update",
	Short: "Update campaign-level negative keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		fromJSON, _ := cmd.Flags().GetString("from-json")

		var keywords []types.NegativeKeyword
		if err := parseJSONInput(fromJSON, &keywords); err != nil {
			return err
		}

		result, err := apiClient.Negatives().CampaignUpdate(campaignID, keywords)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var negCampaignDeleteCmd = &cobra.Command{
	Use:   "campaign-delete",
	Short: "Delete campaign-level negative keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		idsStr, _ := cmd.Flags().GetString("ids")
		ids, err := parseIDList(idsStr)
		if err != nil {
			return err
		}
		if err := apiClient.Negatives().CampaignDelete(campaignID, ids); err != nil {
			return err
		}
		fmt.Println("Campaign negative keywords deleted")
		return nil
	},
}

// Ad group-level negative keywords

var negAdGroupCreateCmd = &cobra.Command{
	Use:   "adgroup-create",
	Short: "Create ad group-level negative keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		fromJSON, _ := cmd.Flags().GetString("from-json")
		text, _ := cmd.Flags().GetString("text")
		matchType, _ := cmd.Flags().GetString("match-type")

		var keywords []types.NegativeKeyword
		if fromJSON != "" {
			if err := parseJSONInput(fromJSON, &keywords); err != nil {
				return err
			}
		} else {
			keywords = []types.NegativeKeyword{{Text: text, MatchType: matchType}}
		}

		result, err := apiClient.Negatives().AdGroupCreate(campaignID, adGroupID, keywords)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var negAdGroupGetCmd = &cobra.Command{
	Use:   "adgroup-get",
	Short: "Get an ad group-level negative keyword",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		id, _ := cmd.Flags().GetInt64("id")
		result, err := apiClient.Negatives().AdGroupGet(campaignID, adGroupID, id)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var negAdGroupListCmd = &cobra.Command{
	Use:   "adgroup-list",
	Short: "List ad group-level negative keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			if limit == 0 && !cmd.Flags().Changed("limit") {
				limit = defaultPageSize
			}
			result, err := collectAllOffsetPaginated(limit, offset, func(lim, off int) ([]types.NegativeKeyword, *types.PageDetail, error) {
				return apiClient.Negatives().AdGroupList(campaignID, adGroupID, lim, off)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Negatives().AdGroupList(campaignID, adGroupID, limit, offset)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var negAdGroupFindCmd = &cobra.Command{
	Use:   "adgroup-find",
	Short: "Find ad group-level negative keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		selectorJSON, _ := cmd.Flags().GetString("selector-json")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all && !cmd.Flags().Changed("limit") {
			limit = defaultPageSize
		}

		sel := &types.Selector{}
		if selectorJSON != "" {
			parsed, err := parseSelectorJSON(selectorJSON)
			if err != nil {
				return err
			}
			sel = parsed
		}

		if sel.Pagination == nil && (limit > 0 || offset > 0) {
			sel.Pagination = &types.Pagination{}
		}
		if sel.Pagination != nil {
			if (all && sel.Pagination.Limit == 0) || cmd.Flags().Changed("limit") {
				sel.Pagination.Limit = limit
			}
			if cmd.Flags().Changed("offset") && sel.Pagination.Offset == 0 {
				sel.Pagination.Offset = offset
			}
		}

		if all {
			pageSize := 0
			if selectorJSON != "" {
				pageSize = limit
			}
			result, err := collectAllSelectorPaginated(sel, pageSize, func(s *types.Selector) ([]types.NegativeKeyword, *types.PageDetail, error) {
				return apiClient.Negatives().AdGroupFind(campaignID, adGroupID, s)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Negatives().AdGroupFind(campaignID, adGroupID, sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var negAdGroupUpdateCmd = &cobra.Command{
	Use:   "adgroup-update",
	Short: "Update ad group-level negative keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		fromJSON, _ := cmd.Flags().GetString("from-json")

		var keywords []types.NegativeKeyword
		if err := parseJSONInput(fromJSON, &keywords); err != nil {
			return err
		}

		result, err := apiClient.Negatives().AdGroupUpdate(campaignID, adGroupID, keywords)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var negAdGroupDeleteCmd = &cobra.Command{
	Use:   "adgroup-delete",
	Short: "Delete ad group-level negative keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		idsStr, _ := cmd.Flags().GetString("ids")
		ids, err := parseIDList(idsStr)
		if err != nil {
			return err
		}
		if err := apiClient.Negatives().AdGroupDelete(campaignID, adGroupID, ids); err != nil {
			return err
		}
		fmt.Println("Ad group negative keywords deleted")
		return nil
	},
}

func parseIDList(s string) ([]int64, error) {
	var ids []int64
	for _, part := range strings.Split(s, ",") {
		id, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid ID %q: %w", part, err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func init() {
	rootCmd.AddCommand(negativesCmd)

	// Campaign-level commands
	for _, c := range []*cobra.Command{negCampaignCreateCmd, negCampaignGetCmd, negCampaignListCmd, negCampaignFindCmd, negCampaignUpdateCmd, negCampaignDeleteCmd} {
		c.Flags().Int64("campaign-id", 0, "Campaign ID")
		c.MarkFlagRequired("campaign-id")
		negativesCmd.AddCommand(c)
	}

	negCampaignCreateCmd.Flags().String("text", "", "Negative keyword text")
	negCampaignCreateCmd.Flags().String("match-type", "EXACT", "BROAD or EXACT")
	negCampaignCreateCmd.Flags().String("from-json", "", "JSON input")

	negCampaignGetCmd.Flags().Int64("id", 0, "Keyword ID")
	negCampaignGetCmd.MarkFlagRequired("id")

	negCampaignListCmd.Flags().Int("limit", 0, "Max results")
	negCampaignListCmd.Flags().Int("offset", 0, "Start offset")
	negCampaignListCmd.Flags().Bool("all", false, "Fetch all pages")

	negCampaignFindCmd.Flags().String("selector-json", "", "Selector JSON")
	negCampaignFindCmd.Flags().Int("limit", 20, "Max results")
	negCampaignFindCmd.Flags().Int("offset", 0, "Start offset")
	negCampaignFindCmd.Flags().Bool("all", false, "Fetch all pages")

	negCampaignUpdateCmd.Flags().String("from-json", "", "JSON input")
	negCampaignUpdateCmd.MarkFlagRequired("from-json")

	negCampaignDeleteCmd.Flags().String("ids", "", "Comma-separated keyword IDs")
	negCampaignDeleteCmd.MarkFlagRequired("ids")

	// Ad group-level commands
	for _, c := range []*cobra.Command{negAdGroupCreateCmd, negAdGroupGetCmd, negAdGroupListCmd, negAdGroupFindCmd, negAdGroupUpdateCmd, negAdGroupDeleteCmd} {
		c.Flags().Int64("campaign-id", 0, "Campaign ID")
		c.MarkFlagRequired("campaign-id")
		c.Flags().Int64("adgroup-id", 0, "Ad group ID")
		c.MarkFlagRequired("adgroup-id")
		negativesCmd.AddCommand(c)
	}

	negAdGroupCreateCmd.Flags().String("text", "", "Negative keyword text")
	negAdGroupCreateCmd.Flags().String("match-type", "EXACT", "BROAD or EXACT")
	negAdGroupCreateCmd.Flags().String("from-json", "", "JSON input")

	negAdGroupGetCmd.Flags().Int64("id", 0, "Keyword ID")
	negAdGroupGetCmd.MarkFlagRequired("id")

	negAdGroupListCmd.Flags().Int("limit", 0, "Max results")
	negAdGroupListCmd.Flags().Int("offset", 0, "Start offset")
	negAdGroupListCmd.Flags().Bool("all", false, "Fetch all pages")

	negAdGroupFindCmd.Flags().String("selector-json", "", "Selector JSON")
	negAdGroupFindCmd.Flags().Int("limit", 20, "Max results")
	negAdGroupFindCmd.Flags().Int("offset", 0, "Start offset")
	negAdGroupFindCmd.Flags().Bool("all", false, "Fetch all pages")

	negAdGroupUpdateCmd.Flags().String("from-json", "", "JSON input")
	negAdGroupUpdateCmd.MarkFlagRequired("from-json")

	negAdGroupDeleteCmd.Flags().String("ids", "", "Comma-separated keyword IDs")
	negAdGroupDeleteCmd.MarkFlagRequired("ids")
}
