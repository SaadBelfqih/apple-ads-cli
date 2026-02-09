package types

// SearchEntity represents a geo location search result.
type SearchEntity struct {
	ID          string `json:"id,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Entity      string `json:"entity,omitempty"` // Country, AdminArea, Locality
	CountryCode string `json:"countryCode,omitempty"`
}

// GeoRequest represents a request to get geo data.
type GeoRequest struct {
	ID string `json:"id,omitempty"`
}
