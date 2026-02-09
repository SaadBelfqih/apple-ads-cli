package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

var productPagesCmd = &cobra.Command{
	Use:   "product-pages",
	Short: "Manage custom product pages",
}

var ppListCmd = &cobra.Command{
	Use:   "list",
	Short: "List product pages for an app",
	RunE: func(cmd *cobra.Command, args []string) error {
		adamID, _ := cmd.Flags().GetInt64("adam-id")
		result, err := apiClient.ProductPages().List(adamID)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var ppGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a product page by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")
		adamID, _ := cmd.Flags().GetInt64("adam-id")
		result, err := apiClient.ProductPages().Get(id, adamID)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var ppLocalesCmd = &cobra.Command{
	Use:   "locales",
	Short: "Get product page locale details",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, _ := cmd.Flags().GetString("id")
		adamID, _ := cmd.Flags().GetInt64("adam-id")
		result, err := apiClient.ProductPages().Locales(id, adamID)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var ppCountriesCmd = &cobra.Command{
	Use:   "countries",
	Short: "List supported countries and regions",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := apiClient.ProductPages().Countries()
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var ppDeviceSizesCmd = &cobra.Command{
	Use:   "device-sizes",
	Short: "Get app preview device size mapping",
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := apiClient.ProductPages().DeviceSizes()
		if err != nil {
			return err
		}
		// Raw JSON response
		var parsed any
		json.Unmarshal(result, &parsed)
		return printOutput(parsed)
	},
}

func init() {
	rootCmd.AddCommand(productPagesCmd)

	ppListCmd.Flags().Int64("adam-id", 0, "App Adam ID")
	ppListCmd.MarkFlagRequired("adam-id")
	productPagesCmd.AddCommand(ppListCmd)

	ppGetCmd.Flags().String("id", "", "Product page ID")
	ppGetCmd.MarkFlagRequired("id")
	ppGetCmd.Flags().Int64("adam-id", 0, "App Adam ID")
	ppGetCmd.MarkFlagRequired("adam-id")
	productPagesCmd.AddCommand(ppGetCmd)

	ppLocalesCmd.Flags().String("id", "", "Product page ID")
	ppLocalesCmd.MarkFlagRequired("id")
	ppLocalesCmd.Flags().Int64("adam-id", 0, "App Adam ID")
	ppLocalesCmd.MarkFlagRequired("adam-id")
	productPagesCmd.AddCommand(ppLocalesCmd)

	productPagesCmd.AddCommand(ppCountriesCmd)
	productPagesCmd.AddCommand(ppDeviceSizesCmd)
}
