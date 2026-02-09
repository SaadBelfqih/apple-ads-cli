package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
	"github.com/spf13/cobra"
)

var keywordsCmd = &cobra.Command{
	Use:   "keywords",
	Short: "Manage targeting keywords",
}

var keywordsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create targeting keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		fromJSON, _ := cmd.Flags().GetString("from-json")
		text, _ := cmd.Flags().GetString("text")
		matchType, _ := cmd.Flags().GetString("match-type")
		bid, _ := cmd.Flags().GetString("bid")

		var keywords []types.Keyword
		if fromJSON != "" {
			if err := parseJSONInput(fromJSON, &keywords); err != nil {
				return err
			}
		} else {
			kw := types.Keyword{Text: text, MatchType: matchType}
			if bid != "" {
				m, err := moneyFromAmount(bid)
				if err != nil {
					return err
				}
				kw.BidAmount = m
			}
			keywords = []types.Keyword{kw}
		}

		result, err := apiClient.Keywords().Create(campaignID, adGroupID, keywords)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var keywordsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a targeting keyword",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		id, _ := cmd.Flags().GetInt64("id")
		result, err := apiClient.Keywords().Get(campaignID, adGroupID, id)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var keywordsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List targeting keywords in an ad group",
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
			result, err := collectAllOffsetPaginated(limit, offset, func(lim, off int) ([]types.Keyword, *types.PageDetail, error) {
				return apiClient.Keywords().List(campaignID, adGroupID, lim, off)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Keywords().List(campaignID, adGroupID, limit, offset)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var keywordsFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find targeting keywords in an ad group",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		selectorJSON, _ := cmd.Flags().GetString("selector-json")
		field, _ := cmd.Flags().GetString("field")
		op, _ := cmd.Flags().GetString("op")
		valuesStr, _ := cmd.Flags().GetString("values")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all && !cmd.Flags().Changed("limit") {
			limit = defaultPageSize
		}

		var values []string
		if valuesStr != "" {
			values = strings.Split(valuesStr, ",")
		}
		sel, err := parseSelector(selectorJSON, field, op, values, limit, offset)
		if err != nil {
			return err
		}
		if all {
			pageSize := 0
			if selectorJSON != "" {
				pageSize = limit
				if offset > 0 {
					if sel.Pagination == nil {
						sel.Pagination = &types.Pagination{}
					}
					sel.Pagination.Offset = offset
				}
			}
			result, err := collectAllSelectorPaginated(sel, pageSize, func(s *types.Selector) ([]types.Keyword, *types.PageDetail, error) {
				return apiClient.Keywords().Find(campaignID, adGroupID, s)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Keywords().Find(campaignID, adGroupID, sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var keywordsFindCampaignCmd = &cobra.Command{
	Use:   "find-campaign",
	Short: "Find targeting keywords across all ad groups in a campaign",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		selectorJSON, _ := cmd.Flags().GetString("selector-json")
		field, _ := cmd.Flags().GetString("field")
		op, _ := cmd.Flags().GetString("op")
		valuesStr, _ := cmd.Flags().GetString("values")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all && !cmd.Flags().Changed("limit") {
			limit = defaultPageSize
		}

		var values []string
		if valuesStr != "" {
			values = strings.Split(valuesStr, ",")
		}
		sel, err := parseSelector(selectorJSON, field, op, values, limit, offset)
		if err != nil {
			return err
		}

		if all {
			pageSize := 0
			if selectorJSON != "" {
				pageSize = limit
				if offset > 0 {
					if sel.Pagination == nil {
						sel.Pagination = &types.Pagination{}
					}
					sel.Pagination.Offset = offset
				}
			}
			result, err := collectAllSelectorPaginated(sel, pageSize, func(s *types.Selector) ([]types.Keyword, *types.PageDetail, error) {
				return apiClient.Keywords().FindCampaign(campaignID, s)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Keywords().FindCampaign(campaignID, sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var keywordsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update targeting keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		fromJSON, _ := cmd.Flags().GetString("from-json")

		var keywords []types.Keyword
		if err := parseJSONInput(fromJSON, &keywords); err != nil {
			return err
		}

		result, err := apiClient.Keywords().Update(campaignID, adGroupID, keywords)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var keywordsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Bulk delete targeting keywords",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		idsStr, _ := cmd.Flags().GetString("ids")

		var ids []int64
		for _, s := range strings.Split(idsStr, ",") {
			id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
			if err != nil {
				return fmt.Errorf("invalid keyword ID %q: %w", s, err)
			}
			ids = append(ids, id)
		}

		if err := apiClient.Keywords().Delete(campaignID, adGroupID, ids); err != nil {
			return err
		}
		fmt.Println("Keywords deleted")
		return nil
	},
}

var keywordsDeleteOneCmd = &cobra.Command{
	Use:   "delete-one",
	Short: "Delete a single targeting keyword",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		id, _ := cmd.Flags().GetInt64("id")
		if err := apiClient.Keywords().DeleteOne(campaignID, adGroupID, id); err != nil {
			return err
		}
		fmt.Println("Keyword " + strconv.FormatInt(id, 10) + " deleted")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(keywordsCmd)

	// create
	keywordsCreateCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	keywordsCreateCmd.MarkFlagRequired("campaign-id")
	keywordsCreateCmd.Flags().Int64("adgroup-id", 0, "Ad group ID")
	keywordsCreateCmd.MarkFlagRequired("adgroup-id")
	keywordsCreateCmd.Flags().String("text", "", "Keyword text")
	keywordsCreateCmd.Flags().String("match-type", "BROAD", "BROAD or EXACT")
	keywordsCreateCmd.Flags().String("bid", "", "Bid amount")
	keywordsCreateCmd.Flags().String("from-json", "", "JSON input (inline, @file, or @- for stdin)")
	keywordsCmd.AddCommand(keywordsCreateCmd)

	// get
	keywordsGetCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	keywordsGetCmd.MarkFlagRequired("campaign-id")
	keywordsGetCmd.Flags().Int64("adgroup-id", 0, "Ad group ID")
	keywordsGetCmd.MarkFlagRequired("adgroup-id")
	keywordsGetCmd.Flags().Int64("id", 0, "Keyword ID")
	keywordsGetCmd.MarkFlagRequired("id")
	keywordsCmd.AddCommand(keywordsGetCmd)

	// list
	keywordsListCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	keywordsListCmd.MarkFlagRequired("campaign-id")
	keywordsListCmd.Flags().Int64("adgroup-id", 0, "Ad group ID")
	keywordsListCmd.MarkFlagRequired("adgroup-id")
	keywordsListCmd.Flags().Int("limit", 0, "Max results")
	keywordsListCmd.Flags().Int("offset", 0, "Start offset")
	keywordsListCmd.Flags().Bool("all", false, "Fetch all pages")
	keywordsCmd.AddCommand(keywordsListCmd)

	// find
	keywordsFindCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	keywordsFindCmd.MarkFlagRequired("campaign-id")
	keywordsFindCmd.Flags().Int64("adgroup-id", 0, "Ad group ID")
	keywordsFindCmd.MarkFlagRequired("adgroup-id")
	keywordsFindCmd.Flags().String("selector-json", "", "Selector JSON")
	keywordsFindCmd.Flags().String("field", "", "Filter field")
	keywordsFindCmd.Flags().String("op", "", "Filter operator")
	keywordsFindCmd.Flags().String("values", "", "Filter values")
	keywordsFindCmd.Flags().Int("limit", 20, "Max results")
	keywordsFindCmd.Flags().Int("offset", 0, "Start offset")
	keywordsFindCmd.Flags().Bool("all", false, "Fetch all pages")
	keywordsCmd.AddCommand(keywordsFindCmd)

	// find-campaign
	keywordsFindCampaignCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	keywordsFindCampaignCmd.MarkFlagRequired("campaign-id")
	keywordsFindCampaignCmd.Flags().String("selector-json", "", "Selector JSON")
	keywordsFindCampaignCmd.Flags().String("field", "", "Filter field")
	keywordsFindCampaignCmd.Flags().String("op", "", "Filter operator")
	keywordsFindCampaignCmd.Flags().String("values", "", "Filter values")
	keywordsFindCampaignCmd.Flags().Int("limit", 20, "Max results")
	keywordsFindCampaignCmd.Flags().Int("offset", 0, "Start offset")
	keywordsFindCampaignCmd.Flags().Bool("all", false, "Fetch all pages")
	keywordsCmd.AddCommand(keywordsFindCampaignCmd)

	// update
	keywordsUpdateCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	keywordsUpdateCmd.MarkFlagRequired("campaign-id")
	keywordsUpdateCmd.Flags().Int64("adgroup-id", 0, "Ad group ID")
	keywordsUpdateCmd.MarkFlagRequired("adgroup-id")
	keywordsUpdateCmd.Flags().String("from-json", "", "JSON input")
	keywordsUpdateCmd.MarkFlagRequired("from-json")
	keywordsCmd.AddCommand(keywordsUpdateCmd)

	// delete (bulk)
	keywordsDeleteCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	keywordsDeleteCmd.MarkFlagRequired("campaign-id")
	keywordsDeleteCmd.Flags().Int64("adgroup-id", 0, "Ad group ID")
	keywordsDeleteCmd.MarkFlagRequired("adgroup-id")
	keywordsDeleteCmd.Flags().String("ids", "", "Comma-separated keyword IDs")
	keywordsDeleteCmd.MarkFlagRequired("ids")
	keywordsCmd.AddCommand(keywordsDeleteCmd)

	// delete-one
	keywordsDeleteOneCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	keywordsDeleteOneCmd.MarkFlagRequired("campaign-id")
	keywordsDeleteOneCmd.Flags().Int64("adgroup-id", 0, "Ad group ID")
	keywordsDeleteOneCmd.MarkFlagRequired("adgroup-id")
	keywordsDeleteOneCmd.Flags().Int64("id", 0, "Keyword ID")
	keywordsDeleteOneCmd.MarkFlagRequired("id")
	keywordsCmd.AddCommand(keywordsDeleteOneCmd)
}
