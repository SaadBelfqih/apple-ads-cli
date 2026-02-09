package types

// AdGroup represents an Apple Ads ad group.
type AdGroup struct {
	ID                    int64                  `json:"id,omitempty"`
	CampaignID            int64                  `json:"campaignId,omitempty"`
	OrgID                 int64                  `json:"orgId,omitempty"`
	Name                  string                 `json:"name,omitempty"`
	Status                string                 `json:"status,omitempty"`
	ServingStatus         string                 `json:"servingStatus,omitempty"`
	DisplayStatus         string                 `json:"displayStatus,omitempty"`
	DefaultBidAmount      *Money                 `json:"defaultBidAmount,omitempty"`
	CpaGoal               *Money                 `json:"cpaGoal,omitempty"`
	AutomatedKeywordsOptIn bool                  `json:"automatedKeywordsOptIn,omitempty"`
	StartTime             string                 `json:"startTime,omitempty"`
	EndTime               string                 `json:"endTime,omitempty"`
	ModificationTime      string                 `json:"modificationTime,omitempty"`
	Deleted               bool                   `json:"deleted,omitempty"`
	TargetingDimensions   *TargetingDimensions   `json:"targetingDimensions,omitempty"`
	ServingStateReasons   []string               `json:"servingStateReasons,omitempty"`
}

// AdGroupCreate is the request body for creating an ad group.
type AdGroupCreate struct {
	Name                   string               `json:"name"`
	DefaultBidAmount       *Money               `json:"defaultBidAmount"`
	CpaGoal                *Money               `json:"cpaGoal,omitempty"`
	AutomatedKeywordsOptIn bool                 `json:"automatedKeywordsOptIn"`
	StartTime              string               `json:"startTime,omitempty"`
	EndTime                string               `json:"endTime,omitempty"`
	Status                 string               `json:"status,omitempty"`
	TargetingDimensions    *TargetingDimensions `json:"targetingDimensions,omitempty"`
}

// AdGroupUpdate is used to update ad group fields.
type AdGroupUpdate struct {
	Name                   string               `json:"name,omitempty"`
	DefaultBidAmount       *Money               `json:"defaultBidAmount,omitempty"`
	CpaGoal                *Money               `json:"cpaGoal,omitempty"`
	AutomatedKeywordsOptIn *bool                `json:"automatedKeywordsOptIn,omitempty"`
	Status                 string               `json:"status,omitempty"`
	StartTime              string               `json:"startTime,omitempty"`
	EndTime                string               `json:"endTime,omitempty"`
	TargetingDimensions    *TargetingDimensions `json:"targetingDimensions,omitempty"`
}

// TargetingDimensions holds all targeting criteria for an ad group.
type TargetingDimensions struct {
	Age            *AgeCriteria            `json:"age,omitempty"`
	Gender         *GenderCriteria         `json:"gender,omitempty"`
	DeviceClass    *DeviceClassCriteria    `json:"deviceClass,omitempty"`
	Daypart        *DaypartCriteria        `json:"daypart,omitempty"`
	AdminArea      *LocationCriteria       `json:"adminArea,omitempty"`
	Locality       *LocationCriteria       `json:"locality,omitempty"`
	Country        *LocationCriteria       `json:"country,omitempty"`
	AppDownloaders *AppDownloadersCriteria `json:"appDownloaders,omitempty"`
	AppCategories  *AppCategoriesCriteria  `json:"appCategories,omitempty"`
}

type AgeCriteria struct {
	Included []AgeRange `json:"included,omitempty"`
}

type AgeRange struct {
	MinAge int `json:"minAge,omitempty"`
	MaxAge int `json:"maxAge,omitempty"`
}

type GenderCriteria struct {
	Included []string `json:"included,omitempty"` // M, F
}

type DeviceClassCriteria struct {
	Included []string `json:"included,omitempty"` // IPHONE, IPAD
}

type DaypartCriteria struct {
	UserTime *DaypartDetail `json:"userTime,omitempty"`
}

type DaypartDetail struct {
	Included []int `json:"included,omitempty"` // 0-167 (hour of week)
}

type LocationCriteria struct {
	Included []string `json:"included,omitempty"`
}

type AppDownloadersCriteria struct {
	Included []string `json:"included,omitempty"`
	Excluded []string `json:"excluded,omitempty"`
}

type AppCategoriesCriteria struct {
	Included []int64 `json:"included,omitempty"`
}
