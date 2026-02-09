package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
	"github.com/spf13/cobra"
)

var adgroupsCmd = &cobra.Command{
	Use:   "adgroups",
	Short: "Manage ad groups",
}

var adgroupsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an ad group",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		fromJSON, _ := cmd.Flags().GetString("from-json")

		var req types.AdGroupCreate
		if fromJSON != "" {
			if err := parseJSONInput(fromJSON, &req); err != nil {
				return err
			}
		} else {
			name, _ := cmd.Flags().GetString("name")
			bid, _ := cmd.Flags().GetString("default-bid")
			searchMatch, _ := cmd.Flags().GetBool("search-match")
			status, _ := cmd.Flags().GetString("status")

			req.Name = name
			req.AutomatedKeywordsOptIn = searchMatch
			if bid != "" {
				m, err := moneyFromAmount(bid)
				if err != nil {
					return err
				}
				req.DefaultBidAmount = m
			}
			if status != "" {
				req.Status = status
			}
		}

		result, err := apiClient.AdGroups().Create(campaignID, &req)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adgroupsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an ad group by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		id, _ := cmd.Flags().GetInt64("id")
		result, err := apiClient.AdGroups().Get(campaignID, id, fieldsFlag)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adgroupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ad groups in a campaign",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			if limit == 0 && !cmd.Flags().Changed("limit") {
				limit = defaultPageSize
			}
			result, err := collectAllOffsetPaginated(limit, offset, func(lim, off int) ([]types.AdGroup, *types.PageDetail, error) {
				return apiClient.AdGroups().List(campaignID, lim, off, fieldsFlag)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.AdGroups().List(campaignID, limit, offset, fieldsFlag)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adgroupsFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find ad groups in a campaign with selector",
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
			result, err := collectAllSelectorPaginated(sel, pageSize, func(s *types.Selector) ([]types.AdGroup, *types.PageDetail, error) {
				return apiClient.AdGroups().Find(campaignID, s)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.AdGroups().Find(campaignID, sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adgroupsFindAllCmd = &cobra.Command{
	Use:   "find-all",
	Short: "Find ad groups across all campaigns (org-level)",
	RunE: func(cmd *cobra.Command, args []string) error {
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
			result, err := collectAllSelectorPaginated(sel, pageSize, apiClient.AdGroups().FindAll)
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.AdGroups().FindAll(sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adgroupsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an ad group",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		id, _ := cmd.Flags().GetInt64("id")
		name, _ := cmd.Flags().GetString("name")
		bid, _ := cmd.Flags().GetString("default-bid")
		status, _ := cmd.Flags().GetString("status")
		searchMatch, _ := cmd.Flags().GetString("search-match")

		req := &types.AdGroupUpdate{}
		if name != "" {
			req.Name = name
		}
		if bid != "" {
			m, err := moneyFromAmount(bid)
			if err != nil {
				return err
			}
			req.DefaultBidAmount = m
		}
		if status != "" {
			req.Status = status
		}
		if searchMatch != "" {
			v := searchMatch == "true"
			req.AutomatedKeywordsOptIn = &v
		}

		result, err := apiClient.AdGroups().Update(campaignID, id, req)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adgroupsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an ad group",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		id, _ := cmd.Flags().GetInt64("id")
		if err := apiClient.AdGroups().Delete(campaignID, id); err != nil {
			return err
		}
		fmt.Println("Ad group " + strconv.FormatInt(id, 10) + " deleted")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(adgroupsCmd)

	// create
	adgroupsCreateCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	adgroupsCreateCmd.MarkFlagRequired("campaign-id")
	adgroupsCreateCmd.Flags().String("name", "", "Ad group name")
	adgroupsCreateCmd.Flags().String("default-bid", "", "Default bid amount")
	adgroupsCreateCmd.Flags().Bool("search-match", false, "Enable automated keywords (Search Match)")
	adgroupsCreateCmd.Flags().String("status", "", "ENABLED or PAUSED")
	adgroupsCreateCmd.Flags().String("from-json", "", "JSON input (inline, @file, or @- for stdin)")
	adgroupsCmd.AddCommand(adgroupsCreateCmd)

	// get
	adgroupsGetCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	adgroupsGetCmd.MarkFlagRequired("campaign-id")
	adgroupsGetCmd.Flags().Int64("id", 0, "Ad group ID")
	adgroupsGetCmd.MarkFlagRequired("id")
	adgroupsCmd.AddCommand(adgroupsGetCmd)

	// list
	adgroupsListCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	adgroupsListCmd.MarkFlagRequired("campaign-id")
	adgroupsListCmd.Flags().Int("limit", 0, "Max results")
	adgroupsListCmd.Flags().Int("offset", 0, "Start offset")
	adgroupsListCmd.Flags().Bool("all", false, "Fetch all pages")
	adgroupsCmd.AddCommand(adgroupsListCmd)

	// find
	adgroupsFindCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	adgroupsFindCmd.MarkFlagRequired("campaign-id")
	adgroupsFindCmd.Flags().String("selector-json", "", "Selector JSON")
	adgroupsFindCmd.Flags().String("field", "", "Filter field")
	adgroupsFindCmd.Flags().String("op", "", "Filter operator")
	adgroupsFindCmd.Flags().String("values", "", "Filter values (comma-separated)")
	adgroupsFindCmd.Flags().Int("limit", 20, "Max results")
	adgroupsFindCmd.Flags().Int("offset", 0, "Start offset")
	adgroupsFindCmd.Flags().Bool("all", false, "Fetch all pages")
	adgroupsCmd.AddCommand(adgroupsFindCmd)

	// find-all
	adgroupsFindAllCmd.Flags().String("selector-json", "", "Selector JSON")
	adgroupsFindAllCmd.Flags().String("field", "", "Filter field")
	adgroupsFindAllCmd.Flags().String("op", "", "Filter operator")
	adgroupsFindAllCmd.Flags().String("values", "", "Filter values (comma-separated)")
	adgroupsFindAllCmd.Flags().Int("limit", 20, "Max results")
	adgroupsFindAllCmd.Flags().Int("offset", 0, "Start offset")
	adgroupsFindAllCmd.Flags().Bool("all", false, "Fetch all pages")
	adgroupsCmd.AddCommand(adgroupsFindAllCmd)

	// update
	adgroupsUpdateCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	adgroupsUpdateCmd.MarkFlagRequired("campaign-id")
	adgroupsUpdateCmd.Flags().Int64("id", 0, "Ad group ID")
	adgroupsUpdateCmd.MarkFlagRequired("id")
	adgroupsUpdateCmd.Flags().String("name", "", "New name")
	adgroupsUpdateCmd.Flags().String("default-bid", "", "New default bid")
	adgroupsUpdateCmd.Flags().String("status", "", "ENABLED or PAUSED")
	adgroupsUpdateCmd.Flags().String("search-match", "", "true or false")
	adgroupsCmd.AddCommand(adgroupsUpdateCmd)

	// delete
	adgroupsDeleteCmd.Flags().Int64("campaign-id", 0, "Campaign ID")
	adgroupsDeleteCmd.MarkFlagRequired("campaign-id")
	adgroupsDeleteCmd.Flags().Int64("id", 0, "Ad group ID")
	adgroupsDeleteCmd.MarkFlagRequired("id")
	adgroupsCmd.AddCommand(adgroupsDeleteCmd)
}
