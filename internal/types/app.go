package types

// AppInfo represents an iOS app from search results.
type AppInfo struct {
	AdamID             int64  `json:"adamId,omitempty"`
	AppName            string `json:"appName,omitempty"`
	DeveloperName      string `json:"developerName,omitempty"`
	CountryOrRegion    string `json:"countryOrRegion,omitempty"`
}

// EligibilityRecord represents app eligibility for Apple Ads.
type EligibilityRecord struct {
	AdamID    int64  `json:"adamId,omitempty"`
	Eligible  bool   `json:"eligible"`
	MinAge    int    `json:"minAge,omitempty"`
	State     string `json:"state,omitempty"`
	AppName   string `json:"appName,omitempty"`
	SupplySource string `json:"supplySource,omitempty"`
}

// AppDetail represents detailed app information.
type AppDetail struct {
	AdamID            int64  `json:"adamId,omitempty"`
	AppName           string `json:"appName,omitempty"`
	DeveloperName     string `json:"developerName,omitempty"`
	CountryOrRegion   string `json:"countryOrRegion,omitempty"`
	PrimaryGenreID    int64  `json:"primaryGenreId,omitempty"`
	IconURL           string `json:"iconUrl,omitempty"`
}

// LocalizedAppDetail represents locale-specific app information.
type LocalizedAppDetail struct {
	AdamID    int64              `json:"adamId,omitempty"`
	AppName   string             `json:"appName,omitempty"`
	Details   []MediaLocaleDetail `json:"details,omitempty"`
}

// MediaLocaleDetail holds media for a specific locale.
type MediaLocaleDetail struct {
	Language string `json:"language,omitempty"`
}
