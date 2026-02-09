package types

// ReportingRequest is the request body for standard report endpoints.
type ReportingRequest struct {
	StartTime        string    `json:"startTime"`
	EndTime          string    `json:"endTime"`
	Granularity      string    `json:"granularity,omitempty"` // HOURLY, DAILY, WEEKLY, MONTHLY
	GroupBy          []string  `json:"groupBy,omitempty"`     // e.g., ["countryOrRegion"]
	TimeZone         string    `json:"timeZone,omitempty"`
	ReturnGrandTotals bool    `json:"returnGrandTotals,omitempty"`
	ReturnRowTotals  bool     `json:"returnRowTotals,omitempty"`
	ReturnRecordsWithNoMetrics bool `json:"returnRecordsWithNoMetrics,omitempty"`
	Selector         *Selector `json:"selector,omitempty"`
}

// ReportingResponse represents a standard report response.
type ReportingResponse struct {
	ReportingDataResponse *ReportingDataResponse `json:"reportingDataResponse,omitempty"`
}

// ReportingDataResponse contains the report data.
type ReportingDataResponse struct {
	Row        []ReportRow `json:"row,omitempty"`
	GrandTotals *ReportRow `json:"grandTotals,omitempty"`
}

// ReportRow represents a single row in a report.
type ReportRow struct {
	Other    map[string]any `json:"other,omitempty"`
	Total    *SpendRow      `json:"total,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Insights *Insights      `json:"insights,omitempty"`
	Granularity []GranularityRow `json:"granularity,omitempty"`
}

// GranularityRow represents a time-bucketed row.
type GranularityRow struct {
	Date  string    `json:"date,omitempty"`
	Total *SpendRow `json:"total,omitempty"`
}

// SpendRow holds metric fields.
type SpendRow struct {
	Impressions     int64   `json:"impressions"`
	Taps            int64   `json:"taps"`
	Installs        int64   `json:"installs"`
	NewDownloads    int64   `json:"newDownloads"`
	Redownloads     int64   `json:"redownloads"`
	LatOnInstalls   int64   `json:"latOnInstalls"`
	LatOffInstalls  int64   `json:"latOffInstalls"`
	TTR             float64 `json:"ttr"`
	ConversionRate  float64 `json:"conversionRate"`
	AvgCPA          *Money  `json:"avgCPA,omitempty"`
	AvgCPT          *Money  `json:"avgCPT,omitempty"`
	LocalSpend      *Money  `json:"localSpend,omitempty"`
}

// Insights holds bid recommendation data.
type Insights struct {
	BidRecommendation *BidRecommendation `json:"bidRecommendation,omitempty"`
}

type BidRecommendation struct {
	SuggestedBidAmount *Money `json:"suggestedBidAmount,omitempty"`
}
