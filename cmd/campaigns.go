package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
	"github.com/spf13/cobra"
)

var campaignsCmd = &cobra.Command{
	Use:   "campaigns",
	Short: "Manage campaigns",
}

var campaignsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a campaign",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		adamID, _ := cmd.Flags().GetInt64("adam-id")
		budget, _ := cmd.Flags().GetString("budget")
		dailyBudget, _ := cmd.Flags().GetString("daily-budget")
		countries, _ := cmd.Flags().GetString("countries")
		status, _ := cmd.Flags().GetString("status")
		fromJSON, _ := cmd.Flags().GetString("from-json")

		var req types.CampaignCreate
		if fromJSON != "" {
			if err := parseJSONInput(fromJSON, &req); err != nil {
				return err
			}
		} else {
			req.Name = name
			req.AdamID = adamID
			if countries != "" {
				req.CountriesOrRegions = strings.Split(countries, ",")
			}
			if budget != "" {
				m, err := moneyFromAmount(budget)
				if err != nil {
					return err
				}
				req.BudgetAmount = m
			}
			if dailyBudget != "" {
				m, err := moneyFromAmount(dailyBudget)
				if err != nil {
					return err
				}
				req.DailyBudgetAmount = m
			}
			if status != "" {
				req.Status = status
			}
		}

		result, err := apiClient.Campaigns().Create(&req)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var campaignsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a campaign by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")
		result, err := apiClient.Campaigns().Get(id, fieldsFlag)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var campaignsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all campaigns",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			if limit == 0 && !cmd.Flags().Changed("limit") {
				limit = defaultPageSize
			}
			result, err := collectAllOffsetPaginated(limit, offset, func(lim, off int) ([]types.Campaign, *types.PageDetail, error) {
				return apiClient.Campaigns().List(lim, off, fieldsFlag)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, pagination, err := apiClient.Campaigns().List(limit, offset, fieldsFlag)
		if err != nil {
			return err
		}

		if verbose && pagination != nil {
			fmt.Printf("Total: %d, showing %d-%d\n", pagination.TotalResults, pagination.StartIndex, pagination.StartIndex+pagination.ItemsPerPage)
		}
		return printOutput(result)
	},
}

var campaignsFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find campaigns with selector",
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
			result, err := collectAllSelectorPaginated(sel, pageSize, apiClient.Campaigns().Find)
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Campaigns().Find(sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var campaignsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a campaign",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")
		name, _ := cmd.Flags().GetString("name")
		budget, _ := cmd.Flags().GetString("budget")
		dailyBudget, _ := cmd.Flags().GetString("daily-budget")
		status, _ := cmd.Flags().GetString("status")
		countries, _ := cmd.Flags().GetString("countries")

		req := &types.CampaignUpdate{}
		if name != "" {
			req.Name = name
		}
		if budget != "" {
			m, err := moneyFromAmount(budget)
			if err != nil {
				return err
			}
			req.BudgetAmount = m
		}
		if dailyBudget != "" {
			m, err := moneyFromAmount(dailyBudget)
			if err != nil {
				return err
			}
			req.DailyBudgetAmount = m
		}
		if status != "" {
			req.Status = status
		}
		if countries != "" {
			req.CountriesOrRegions = strings.Split(countries, ",")
		}

		result, err := apiClient.Campaigns().Update(id, req)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var campaignsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a campaign",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")
		if err := apiClient.Campaigns().Delete(id); err != nil {
			return err
		}
		fmt.Println("Campaign " + strconv.FormatInt(id, 10) + " deleted")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(campaignsCmd)

	// create
	campaignsCreateCmd.Flags().String("name", "", "Campaign name")
	campaignsCreateCmd.Flags().Int64("adam-id", 0, "App Adam ID")
	campaignsCreateCmd.Flags().String("budget", "", "Total budget amount")
	campaignsCreateCmd.Flags().String("daily-budget", "", "Daily budget amount")
	campaignsCreateCmd.Flags().String("countries", "", "Comma-separated country codes")
	campaignsCreateCmd.Flags().String("status", "", "ENABLED or PAUSED")
	campaignsCreateCmd.Flags().String("from-json", "", "JSON input (inline, @file, or @- for stdin)")
	campaignsCmd.AddCommand(campaignsCreateCmd)

	// get
	campaignsGetCmd.Flags().Int64("id", 0, "Campaign ID")
	campaignsGetCmd.MarkFlagRequired("id")
	campaignsCmd.AddCommand(campaignsGetCmd)

	// list
	campaignsListCmd.Flags().Int("limit", 0, "Max results")
	campaignsListCmd.Flags().Int("offset", 0, "Start offset")
	campaignsListCmd.Flags().Bool("all", false, "Fetch all pages")
	campaignsCmd.AddCommand(campaignsListCmd)

	// find
	campaignsFindCmd.Flags().String("selector-json", "", "Selector JSON (inline or @file)")
	campaignsFindCmd.Flags().String("field", "", "Filter field name")
	campaignsFindCmd.Flags().String("op", "", "Filter operator (EQUALS, CONTAINS, etc.)")
	campaignsFindCmd.Flags().String("values", "", "Comma-separated filter values")
	campaignsFindCmd.Flags().Int("limit", 20, "Max results")
	campaignsFindCmd.Flags().Int("offset", 0, "Start offset")
	campaignsFindCmd.Flags().Bool("all", false, "Fetch all pages")
	campaignsCmd.AddCommand(campaignsFindCmd)

	// update
	campaignsUpdateCmd.Flags().Int64("id", 0, "Campaign ID")
	campaignsUpdateCmd.MarkFlagRequired("id")
	campaignsUpdateCmd.Flags().String("name", "", "New name")
	campaignsUpdateCmd.Flags().String("budget", "", "New total budget")
	campaignsUpdateCmd.Flags().String("daily-budget", "", "New daily budget")
	campaignsUpdateCmd.Flags().String("status", "", "ENABLED or PAUSED")
	campaignsUpdateCmd.Flags().String("countries", "", "Comma-separated country codes")
	campaignsCmd.AddCommand(campaignsUpdateCmd)

	// delete
	campaignsDeleteCmd.Flags().Int64("id", 0, "Campaign ID")
	campaignsDeleteCmd.MarkFlagRequired("id")
	campaignsCmd.AddCommand(campaignsDeleteCmd)
}
