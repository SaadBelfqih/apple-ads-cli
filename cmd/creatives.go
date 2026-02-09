package cmd

import (
	"strings"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
	"github.com/spf13/cobra"
)

var creativesCmd = &cobra.Command{
	Use:   "creatives",
	Short: "Manage creatives (org-level)",
}

var creativesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a creative",
	RunE: func(cmd *cobra.Command, args []string) error {
		adamID, _ := cmd.Flags().GetInt64("adam-id")
		productPageID, _ := cmd.Flags().GetString("product-page-id")
		name, _ := cmd.Flags().GetString("name")

		req := &types.CreativeCreate{
			AdamID:        adamID,
			ProductPageID: productPageID,
			Name:          name,
		}

		result, err := apiClient.Creatives().Create(req)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var creativesGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a creative by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")
		result, err := apiClient.Creatives().Get(id)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var creativesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all creatives",
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		all, _ := cmd.Flags().GetBool("all")

		if all {
			if limit == 0 && !cmd.Flags().Changed("limit") {
				limit = defaultPageSize
			}
			result, err := collectAllOffsetPaginated(limit, offset, apiClient.Creatives().List)
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Creatives().List(limit, offset)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var creativesFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find creatives with selector",
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
			result, err := collectAllSelectorPaginated(sel, pageSize, apiClient.Creatives().Find)
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.Creatives().Find(sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

func init() {
	rootCmd.AddCommand(creativesCmd)

	// create
	creativesCreateCmd.Flags().Int64("adam-id", 0, "App Adam ID")
	creativesCreateCmd.MarkFlagRequired("adam-id")
	creativesCreateCmd.Flags().String("product-page-id", "", "Product page ID")
	creativesCreateCmd.Flags().String("name", "", "Creative name")
	creativesCmd.AddCommand(creativesCreateCmd)

	// get
	creativesGetCmd.Flags().Int64("id", 0, "Creative ID")
	creativesGetCmd.MarkFlagRequired("id")
	creativesCmd.AddCommand(creativesGetCmd)

	// list
	creativesListCmd.Flags().Int("limit", 0, "Max results")
	creativesListCmd.Flags().Int("offset", 0, "Start offset")
	creativesListCmd.Flags().Bool("all", false, "Fetch all pages")
	creativesCmd.AddCommand(creativesListCmd)

	// find
	creativesFindCmd.Flags().String("selector-json", "", "Selector JSON")
	creativesFindCmd.Flags().String("field", "", "Filter field")
	creativesFindCmd.Flags().String("op", "", "Filter operator")
	creativesFindCmd.Flags().String("values", "", "Filter values")
	creativesFindCmd.Flags().Int("limit", 20, "Max results")
	creativesFindCmd.Flags().Int("offset", 0, "Start offset")
	creativesFindCmd.Flags().Bool("all", false, "Fetch all pages")
	creativesCmd.AddCommand(creativesFindCmd)
}
