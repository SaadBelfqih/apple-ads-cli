package types

// Money represents a currency amount.
type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// Selector is used with POST /find endpoints to filter results.
type Selector struct {
	Conditions []*Condition `json:"conditions,omitempty"`
	Fields     []string     `json:"fields,omitempty"`
	OrderBy    []*Sorting   `json:"orderBy,omitempty"`
	Pagination *Pagination  `json:"pagination,omitempty"`
}

// Condition represents a single filter condition in a Selector.
type Condition struct {
	Field    string   `json:"field"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

// Sorting specifies a sort field and direction.
type Sorting struct {
	Field     string `json:"field"`
	SortOrder string `json:"sortOrder"` // ASCENDING or DESCENDING
}

// Pagination controls offset-based pagination.
type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// PageDetail is returned in API responses to indicate pagination state.
type PageDetail struct {
	TotalResults int `json:"totalResults"`
	StartIndex   int `json:"startIndex"`
	ItemsPerPage int `json:"itemsPerPage"`
}

// LOCInvoiceDetails holds invoice/billing info for budget orders.
type LOCInvoiceDetails struct {
	BillingContactEmail string `json:"billingContactEmail,omitempty"`
	BuyerEmail          string `json:"buyerEmail,omitempty"`
	BuyerName           string `json:"buyerName,omitempty"`
	ClientName          string `json:"clientName,omitempty"`
	OrderNumber         string `json:"orderNumber,omitempty"`
}
