package types

// CustomReportRequest is the request body for creating impression share reports.
type CustomReportRequest struct {
	Name             string        `json:"name,omitempty"`
	StartTime        string        `json:"startTime"`
	EndTime          string        `json:"endTime"`
	Granularity      string        `json:"granularity,omitempty"`
	GroupBy          []string      `json:"groupBy,omitempty"`
	Selector         *SovSelector  `json:"selector,omitempty"`
}

// SovSelector is the selector for impression share reports.
type SovSelector struct {
	Conditions []*SovCondition `json:"conditions,omitempty"`
	OrderBy    []*Sorting      `json:"orderBy,omitempty"`
	Pagination *Pagination     `json:"pagination,omitempty"`
}

// SovCondition is a filter condition for impression share reports.
type SovCondition struct {
	Field    string   `json:"field"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

// CustomReportResponse represents a custom (impression share) report.
type CustomReportResponse struct {
	ID               int64  `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	State            string `json:"state,omitempty"` // QUEUED, RUNNING, COMPLETED, FAILED
	StartTime        string `json:"startTime,omitempty"`
	EndTime          string `json:"endTime,omitempty"`
	Granularity      string `json:"granularity,omitempty"`
	DateRange        string `json:"dateRange,omitempty"`
	CreationTime     string `json:"creationTime,omitempty"`
	ModificationTime string `json:"modificationTime,omitempty"`
	ReportRows       []any  `json:"reportRows,omitempty"`
}
