package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

var geoCmd = &cobra.Command{
	Use:   "geo",
	Short: "Search geolocations",
}

var geoSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for geolocations",
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetString("query")
		countryCode, _ := cmd.Flags().GetString("country-code")
		entity, _ := cmd.Flags().GetString("entity")
		limit, _ := cmd.Flags().GetInt("limit")

		result, err := apiClient.Geo().Search(query, countryCode, entity, limit)
		if err != nil {
			return err
		}
		return printOutput(result)
	},
}

var geoGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get geo location by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		geoID, _ := cmd.Flags().GetString("geo-id")
		result, err := apiClient.Geo().Get(geoID)
		if err != nil {
			return err
		}
		var parsed any
		json.Unmarshal(result, &parsed)
		return printOutput(parsed)
	},
}

func init() {
	rootCmd.AddCommand(geoCmd)

	geoSearchCmd.Flags().String("query", "", "Search query")
	geoSearchCmd.MarkFlagRequired("query")
	geoSearchCmd.Flags().String("country-code", "", "Country code filter")
	geoSearchCmd.Flags().String("entity", "", "Entity type: Country, AdminArea, Locality")
	geoSearchCmd.Flags().Int("limit", 0, "Max results")
	geoCmd.AddCommand(geoSearchCmd)

	geoGetCmd.Flags().String("geo-id", "", "Geo identifier")
	geoGetCmd.MarkFlagRequired("geo-id")
	geoCmd.AddCommand(geoGetCmd)
}
