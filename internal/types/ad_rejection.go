package types

// ProductPageReason represents an ad creative rejection reason based on a product page.
type ProductPageReason struct {
	AdamID           int64   `json:"adamId,omitempty"`
	AppPreviewDevice string  `json:"appPreviewDevice,omitempty"`
	AssetGenID       string  `json:"assetGenId,omitempty"`
	Comment          string  `json:"comment,omitempty"`
	CountryOrRegion  string  `json:"countryOrRegion,omitempty"`
	ID               int64   `json:"id,omitempty"`
	LanguageCode     string  `json:"languageCode,omitempty"`
	ProductPageID    *string `json:"productPageId,omitempty"`
	ReasonCode       string  `json:"reasonCode,omitempty"`
	ReasonLevel      string  `json:"reasonLevel,omitempty"`
	ReasonType       string  `json:"reasonType,omitempty"`
	SupplySource     string  `json:"supplySource,omitempty"`
}

// AppAsset represents an app asset.
type AppAsset struct {
	AdamID           int64  `json:"adamId,omitempty"`
	AppPreviewDevice string `json:"appPreviewDevice,omitempty"`
	AssetGenID       string `json:"assetGenId,omitempty"`
	AssetType        string `json:"assetType,omitempty"` // APP_PREVIEW, SCREENSHOT
	AssetURL         string `json:"assetURL,omitempty"`
	AssetVideoURL    string `json:"assetVideoUrl,omitempty"`
	Deleted          bool   `json:"deleted,omitempty"`
	Orientation      string `json:"orientation,omitempty"` // LANDSCAPE, PORTRAIT, UNKNOWN
	SourceHeight     int32  `json:"sourceHeight,omitempty"`
	SourceWidth      int32  `json:"sourceWidth,omitempty"`
}
