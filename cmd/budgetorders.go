package cmd

import (
	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
	"github.com/spf13/cobra"
)

var budgetOrdersCmd = &cobra.Command{
	Use:   "budgetorders",
	Short: "Manage budget orders",
}

var boCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a budget order",
	RunE: func(cmd *cobra.Command, args []string) error {
		fromJSON, _ := cmd.Flags().GetString("from-json")
		var req types.BudgetOrderCreate
		if err := parseJSONInput(fromJSON, &req); err != nil {
			return err
		}
		result, err := apiClient.BudgetOrders().Create(&req)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var boGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a budget order by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")
		result, err := apiClient.BudgetOrders().Get(id)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var boListCmd = &cobra.Command{
	Use:   "list",
	Short: "List budget orders",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			if limit == 0 && !cmd.Flags().Changed("limit") {
				limit = defaultPageSize
			}
			result, err := collectAllOffsetPaginated(limit, offset, apiClient.BudgetOrders().List)
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.BudgetOrders().List(limit, offset)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var boUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a budget order",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")
		fromJSON, _ := cmd.Flags().GetString("from-json")
		var req types.BudgetOrderUpdate
		if err := parseJSONInput(fromJSON, &req); err != nil {
			return err
		}
		result, err := apiClient.BudgetOrders().Update(id, &req)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

func init() {
	rootCmd.AddCommand(budgetOrdersCmd)

	boCreateCmd.Flags().String("from-json", "", "JSON input (inline, @file, or @- for stdin)")
	boCreateCmd.MarkFlagRequired("from-json")
	budgetOrdersCmd.AddCommand(boCreateCmd)

	boGetCmd.Flags().Int64("id", 0, "Budget order ID")
	boGetCmd.MarkFlagRequired("id")
	budgetOrdersCmd.AddCommand(boGetCmd)

	boListCmd.Flags().Int("limit", 0, "Max results")
	boListCmd.Flags().Int("offset", 0, "Start offset")
	boListCmd.Flags().Bool("all", false, "Fetch all pages")
	budgetOrdersCmd.AddCommand(boListCmd)

	boUpdateCmd.Flags().Int64("id", 0, "Budget order ID")
	boUpdateCmd.MarkFlagRequired("id")
	boUpdateCmd.Flags().String("from-json", "", "JSON input")
	boUpdateCmd.MarkFlagRequired("from-json")
	budgetOrdersCmd.AddCommand(boUpdateCmd)
}
