package types

// BudgetOrder represents a budget order.
type BudgetOrder struct {
	ID               int64              `json:"id,omitempty"`
	OrgID            int64              `json:"orgId,omitempty"`
	Name             string             `json:"name,omitempty"`
	StartDate        string             `json:"startDate,omitempty"`
	EndDate          string             `json:"endDate,omitempty"`
	Budget           *Money             `json:"budget,omitempty"`
	Status           string             `json:"status,omitempty"` // ACTIVE, COMPLETED, EXHAUSTED, etc.
	OrderNumber      string             `json:"orderNumber,omitempty"`
	SupplySource     string             `json:"supplySource,omitempty"`
	LOCInvoiceDetails *LOCInvoiceDetails `json:"locInvoiceDetails,omitempty"`
}

// BudgetOrderCreate is the request body for creating a budget order.
type BudgetOrderCreate struct {
	Name              string             `json:"name"`
	StartDate         string             `json:"startDate"`
	EndDate           string             `json:"endDate,omitempty"`
	Budget            *Money             `json:"budget"`
	OrderNumber       string             `json:"orderNumber,omitempty"`
	SupplySource      string             `json:"supplySource,omitempty"`
	LOCInvoiceDetails *LOCInvoiceDetails `json:"locInvoiceDetails,omitempty"`
}

// BudgetOrderUpdate is used to update budget order fields.
type BudgetOrderUpdate struct {
	Name              string             `json:"name,omitempty"`
	EndDate           string             `json:"endDate,omitempty"`
	Budget            *Money             `json:"budget,omitempty"`
	LOCInvoiceDetails *LOCInvoiceDetails `json:"locInvoiceDetails,omitempty"`
}
