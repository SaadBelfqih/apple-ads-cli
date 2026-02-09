package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
	"github.com/spf13/cobra"
)

var adRejectionsCmd = &cobra.Command{
	Use:   "ad-rejections",
	Short: "Manage ad rejection reasons",
}

var adRejFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find ad creative rejection reasons",
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
			result, err := collectAllSelectorPaginated(sel, pageSize, apiClient.AdRejections().Find)
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.AdRejections().Find(sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adRejGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an ad creative rejection reason by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetInt64("id")
		if id == 0 {
			legacy, _ := cmd.Flags().GetString("product-page-id")
			if legacy != "" {
				parsed, err := strconv.ParseInt(legacy, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid product-page-id %q (expected integer): %w", legacy, err)
				}
				id = parsed
			}
		}
		if id == 0 {
			return fmt.Errorf("--id is required")
		}

		result, err := apiClient.AdRejections().Get(id)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var adRejFindAssetsCmd = &cobra.Command{
	Use:   "find-assets",
	Short: "Find app assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		selectorJSON, _ := cmd.Flags().GetString("selector-json")
		adamID, _ := cmd.Flags().GetInt64("adam-id")
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

		// Back-compat: if caller didn't pass --adam-id, try to infer it from the selector.
		if adamID == 0 && sel != nil {
			for _, c := range sel.Conditions {
				if c == nil || c.Field != "adamId" || c.Operator != "EQUALS" || len(c.Values) != 1 {
					continue
				}
				parsed, err := strconv.ParseInt(c.Values[0], 10, 64)
				if err != nil {
					continue
				}
				adamID = parsed
				break
			}
		}
		if adamID == 0 {
			return fmt.Errorf("--adam-id is required (or include an adamId EQUALS condition in --selector-json)")
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
			result, err := collectAllSelectorPaginated(sel, pageSize, func(s *types.Selector) ([]types.AppAsset, *types.PageDetail, error) {
				return apiClient.AdRejections().FindAssets(adamID, s)
			})
			if err != nil {
				return err
			}
			return printOutput(result)
		}

		result, _, err := apiClient.AdRejections().FindAssets(adamID, sel)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

func init() {
	rootCmd.AddCommand(adRejectionsCmd)

	adRejFindCmd.Flags().String("selector-json", "", "Selector JSON")
	adRejFindCmd.Flags().String("field", "", "Filter field")
	adRejFindCmd.Flags().String("op", "", "Filter operator")
	adRejFindCmd.Flags().String("values", "", "Filter values")
	adRejFindCmd.Flags().Int("limit", 20, "Max results")
	adRejFindCmd.Flags().Int("offset", 0, "Start offset")
	adRejFindCmd.Flags().Bool("all", false, "Fetch all pages")
	adRejectionsCmd.AddCommand(adRejFindCmd)

	adRejGetCmd.Flags().Int64("id", 0, "Product page reason ID")
	adRejGetCmd.Flags().String("product-page-id", "", "Deprecated: use --id (product page reason ID)")
	adRejGetCmd.Flags().MarkDeprecated("product-page-id", "use --id (product page reason ID)")
	adRejectionsCmd.AddCommand(adRejGetCmd)

	adRejFindAssetsCmd.Flags().String("selector-json", "", "Selector JSON")
	adRejFindAssetsCmd.Flags().Int64("adam-id", 0, "App Adam ID")
	adRejFindAssetsCmd.Flags().Int("limit", 20, "Max results")
	adRejFindAssetsCmd.Flags().Int("offset", 0, "Start offset")
	adRejFindAssetsCmd.Flags().Bool("all", false, "Fetch all pages")
	adRejectionsCmd.AddCommand(adRejFindAssetsCmd)
}
