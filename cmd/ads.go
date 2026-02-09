package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
	"github.com/spf13/cobra"
)

var adsCmd = &cobra.Command{
	Use:   "ads",
	Short: "Manage ads",
}

var adsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an ad",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		creativeID, _ := cmd.Flags().GetInt64("creative-id")
		name, _ := cmd.Flags().GetString("name")
		status, _ := cmd.Flags().GetString("status")

		req := &types.AdCreate{
			CreativeID: creativeID,
			Name:       name,
			Status:     status,
		}

		result, err := apiClient.Ads().Create(campaignID, adGroupID, req)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an ad by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		id, _ := cmd.Flags().GetInt64("id")
		result, err := apiClient.Ads().Get(campaignID, adGroupID, id)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ads in an ad group",
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
			result, err := collectAllOffsetPaginated(limit, offset, func(lim, off int) ([]types.Ad, *types.PageDetail, error) {
				return apiClient.Ads().List(campaignID, adGroupID, lim, off)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Ads().List(campaignID, adGroupID, limit, offset)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adsFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find ads in an ad group",
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
			result, err := collectAllSelectorPaginated(sel, pageSize, func(s *types.Selector) ([]types.Ad, *types.PageDetail, error) {
				return apiClient.Ads().Find(campaignID, adGroupID, s)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Ads().Find(campaignID, adGroupID, sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adsFindAllCmd = &cobra.Command{
	Use:   "find-all",
	Short: "Find ads across all campaigns (org-level)",
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
			result, err := collectAllSelectorPaginated(sel, pageSize, apiClient.Ads().FindAll)
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Ads().FindAll(sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an ad",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		id, _ := cmd.Flags().GetInt64("id")
		name, _ := cmd.Flags().GetString("name")
		status, _ := cmd.Flags().GetString("status")

		req := &types.AdUpdate{Name: name, Status: status}
		result, err := apiClient.Ads().Update(campaignID, adGroupID, id, req)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an ad",
	RunE: func(cmd *cobra.Command, args []string) error {
		campaignID, _ := cmd.Flags().GetInt64("campaign-id")
		adGroupID, _ := cmd.Flags().GetInt64("adgroup-id")
		id, _ := cmd.Flags().GetInt64("id")
		if err := apiClient.Ads().Delete(campaignID, adGroupID, id); err != nil {
			return err
		}
		fmt.Println("Ad " + strconv.FormatInt(id, 10) + " deleted")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(adsCmd)

	// Common flags helper
	addAdFlags := func(c *cobra.Command) {
		c.Flags().Int64("campaign-id", 0, "Campaign ID")
		c.MarkFlagRequired("campaign-id")
		c.Flags().Int64("adgroup-id", 0, "Ad group ID")
		c.MarkFlagRequired("adgroup-id")
	}

	// create
	addAdFlags(adsCreateCmd)
	adsCreateCmd.Flags().Int64("creative-id", 0, "Creative ID")
	adsCreateCmd.MarkFlagRequired("creative-id")
	adsCreateCmd.Flags().String("name", "", "Ad name")
	adsCreateCmd.Flags().String("status", "", "ENABLED or PAUSED")
	adsCmd.AddCommand(adsCreateCmd)

	// get
	addAdFlags(adsGetCmd)
	adsGetCmd.Flags().Int64("id", 0, "Ad ID")
	adsGetCmd.MarkFlagRequired("id")
	adsCmd.AddCommand(adsGetCmd)

	// list
	addAdFlags(adsListCmd)
	adsListCmd.Flags().Int("limit", 0, "Max results")
	adsListCmd.Flags().Int("offset", 0, "Start offset")
	adsListCmd.Flags().Bool("all", false, "Fetch all pages")
	adsCmd.AddCommand(adsListCmd)

	// find
	addAdFlags(adsFindCmd)
	adsFindCmd.Flags().String("selector-json", "", "Selector JSON")
	adsFindCmd.Flags().String("field", "", "Filter field")
	adsFindCmd.Flags().String("op", "", "Filter operator")
	adsFindCmd.Flags().String("values", "", "Filter values")
	adsFindCmd.Flags().Int("limit", 20, "Max results")
	adsFindCmd.Flags().Int("offset", 0, "Start offset")
	adsFindCmd.Flags().Bool("all", false, "Fetch all pages")
	adsCmd.AddCommand(adsFindCmd)

	// find-all
	adsFindAllCmd.Flags().String("selector-json", "", "Selector JSON")
	adsFindAllCmd.Flags().String("field", "", "Filter field")
	adsFindAllCmd.Flags().String("op", "", "Filter operator")
	adsFindAllCmd.Flags().String("values", "", "Filter values")
	adsFindAllCmd.Flags().Int("limit", 20, "Max results")
	adsFindAllCmd.Flags().Int("offset", 0, "Start offset")
	adsFindAllCmd.Flags().Bool("all", false, "Fetch all pages")
	adsCmd.AddCommand(adsFindAllCmd)

	// update
	addAdFlags(adsUpdateCmd)
	adsUpdateCmd.Flags().Int64("id", 0, "Ad ID")
	adsUpdateCmd.MarkFlagRequired("id")
	adsUpdateCmd.Flags().String("name", "", "New name")
	adsUpdateCmd.Flags().String("status", "", "ENABLED or PAUSED")
	adsCmd.AddCommand(adsUpdateCmd)

	// delete
	addAdFlags(adsDeleteCmd)
	adsDeleteCmd.Flags().Int64("id", 0, "Ad ID")
	adsDeleteCmd.MarkFlagRequired("id")
	adsCmd.AddCommand(adsDeleteCmd)
}
