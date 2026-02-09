package types

// ProductPageDetail represents a custom product page.
type ProductPageDetail struct {
	AdamID           int64  `json:"adamId,omitempty"`
	CreationTime     string `json:"creationTime,omitempty"`
	DeepLink         string `json:"deepLink,omitempty"`
	ID               string `json:"id,omitempty"`
	ModificationTime string `json:"modificationTime,omitempty"`
	Name             string `json:"name,omitempty"`
	State            string `json:"state,omitempty"` // HIDDEN, VISIBLE
}

// ProductPageLocaleDetail represents locale-specific product page details.
type ProductPageLocaleDetail struct {
	AdamID                   int64          `json:"adamId,omitempty"`
	AppName                  string         `json:"appName,omitempty"`
	AppPreviewDeviceWithAssets map[string]any `json:"appPreviewDeviceWithAssets,omitempty"`
	DeviceClasses            any            `json:"deviceClasses,omitempty"`
	Language                 string         `json:"language,omitempty"`
	LanguageCode             string         `json:"languageCode,omitempty"`
	ProductPageID            string         `json:"productPageId,omitempty"`
	PromotionalText          string         `json:"promotionalText,omitempty"`
	ShortDescription         string         `json:"shortDescription,omitempty"`
	SubTitle                 string         `json:"subTitle,omitempty"`
}

// LocaleInfo represents locale information.
type LocaleInfo struct {
	Language string `json:"language,omitempty"`
}

// CountryOrRegion represents a supported country or region.
type CountryOrRegion struct {
	Code        string `json:"code,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// DeviceSizeMapping represents app preview device size mapping.
type DeviceSizeMapping struct {
	DeviceClass string `json:"deviceClass,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}
