package types

// Creative represents an org-level creative.
type Creative struct {
	ID               int64  `json:"id,omitempty"`
	OrgID            int64  `json:"orgId,omitempty"`
	AdamID           int64  `json:"adamId,omitempty"`
	Name             string `json:"name,omitempty"`
	Type             string `json:"type,omitempty"`
	State            string `json:"state,omitempty"`
	CreationTime     string `json:"creationTime,omitempty"`
	ModificationTime string `json:"modificationTime,omitempty"`
	Deleted          bool   `json:"deleted,omitempty"`

	ProductPageID         string                      `json:"productPageId,omitempty"`
	CustomProductPageCreative *CustomProductPageCreative `json:"customProductPageCreative,omitempty"`
}

// CustomProductPageCreative contains localized creative details.
type CustomProductPageCreative struct {
	Localizations []CreativeLocalization `json:"localizations,omitempty"`
}

// CreativeLocalization represents a locale-specific creative config.
type CreativeLocalization struct {
	Language string `json:"language,omitempty"`
}

// CreativeCreate is the request body for creating a creative.
type CreativeCreate struct {
	AdamID        int64  `json:"adamId"`
	Name          string `json:"name,omitempty"`
	ProductPageID string `json:"productPageId,omitempty"`
}
