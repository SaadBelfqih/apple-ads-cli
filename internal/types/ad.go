package types

// Ad represents an ad in an ad group.
type Ad struct {
	ID               int64  `json:"id,omitempty"`
	OrgID            int64  `json:"orgId,omitempty"`
	CampaignID       int64  `json:"campaignId,omitempty"`
	AdGroupID        int64  `json:"adGroupId,omitempty"`
	Name             string `json:"name,omitempty"`
	CreativeID       int64  `json:"creativeId,omitempty"`
	Status           string `json:"status,omitempty"`
	ServingStatus    string `json:"servingStatus,omitempty"`
	CreativeType     string `json:"creativeType,omitempty"`
	Deleted          bool   `json:"deleted,omitempty"`
	ModificationTime string `json:"modificationTime,omitempty"`
	CreationTime     string `json:"creationTime,omitempty"`
}

// AdCreate is the request body for creating an ad.
type AdCreate struct {
	Name       string `json:"name,omitempty"`
	CreativeID int64  `json:"creativeId"`
	Status     string `json:"status,omitempty"`
}

// AdUpdate is used to update ad fields.
type AdUpdate struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}
