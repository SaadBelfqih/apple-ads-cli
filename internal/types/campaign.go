package types

// Campaign represents an Apple Ads campaign.
type Campaign struct {
	ID                        int64    `json:"id,omitempty"`
	OrgID                     int64    `json:"orgId,omitempty"`
	Name                      string   `json:"name,omitempty"`
	BudgetAmount              *Money   `json:"budgetAmount,omitempty"`
	DailyBudgetAmount         *Money   `json:"dailyBudgetAmount,omitempty"`
	AdamID                    int64    `json:"adamId,omitempty"`
	CountriesOrRegions        []string `json:"countriesOrRegions,omitempty"`
	Status                    string   `json:"status,omitempty"`
	ServingStatus             string   `json:"servingStatus,omitempty"`
	DisplayStatus             string   `json:"displayStatus,omitempty"`
	SupplySources             []string `json:"supplySources,omitempty"` // APPSTORE_SEARCH_RESULTS, APPSTORE_SEARCH_TAB
	AdChannelType             string   `json:"adChannelType,omitempty"`
	BudgetOrders              []int64  `json:"budgetOrders,omitempty"`
	StartTime                 string   `json:"startTime,omitempty"`
	EndTime                   string   `json:"endTime,omitempty"`
	ServingStateReasons       []string `json:"servingStateReasons,omitempty"`
	ModificationTime          string   `json:"modificationTime,omitempty"`
	Deleted                   bool     `json:"deleted,omitempty"`
	CountryOrRegionServingStateReasons map[string][]string `json:"countryOrRegionServingStateReasons,omitempty"`
	LOCEnabled                bool     `json:"locEnabled,omitempty"`
}

// CampaignCreate is the request body for creating a campaign.
type CampaignCreate struct {
	Name               string   `json:"name"`
	BudgetAmount       *Money   `json:"budgetAmount,omitempty"`
	DailyBudgetAmount  *Money   `json:"dailyBudgetAmount,omitempty"`
	AdamID             int64    `json:"adamId"`
	CountriesOrRegions []string `json:"countriesOrRegions"`
	Status             string   `json:"status,omitempty"` // ENABLED or PAUSED
	SupplySources      []string `json:"supplySources,omitempty"`
	AdChannelType      string   `json:"adChannelType,omitempty"`
	LOCInvoiceDetails  *LOCInvoiceDetails `json:"locInvoiceDetails,omitempty"`
}

// CampaignUpdate is used to update campaign fields.
type CampaignUpdate struct {
	Name               string   `json:"name,omitempty"`
	BudgetAmount       *Money   `json:"budgetAmount,omitempty"`
	DailyBudgetAmount  *Money   `json:"dailyBudgetAmount,omitempty"`
	Status             string   `json:"status,omitempty"`
	CountriesOrRegions []string `json:"countriesOrRegions,omitempty"`
	LOCInvoiceDetails  *LOCInvoiceDetails `json:"locInvoiceDetails,omitempty"`
}

// UpdateCampaignRequest wraps the update in the required envelope.
// Apple's API requires PUT /campaigns/{id} to use {"campaign":{...}}.
type UpdateCampaignRequest struct {
	Campaign *CampaignUpdate `json:"campaign"`
}
