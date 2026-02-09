package types

// APIResponse is the generic envelope returned by the Apple Ads API.
type APIResponse[T any] struct {
	Data       *T          `json:"data,omitempty"`
	Pagination *PageDetail `json:"pagination,omitempty"`
	Error      *APIError   `json:"error,omitempty"`
}

// APIListResponse is the generic envelope for list/find endpoints.
type APIListResponse[T any] struct {
	Data       []T         `json:"data,omitempty"`
	Pagination *PageDetail `json:"pagination,omitempty"`
	Error      *APIError   `json:"error,omitempty"`
}

// APIError represents an error returned by the API.
type APIError struct {
	Errors []APIErrorDetail `json:"errors,omitempty"`
}

// APIErrorDetail is a single error entry.
type APIErrorDetail struct {
	MessageCode string `json:"messageCode"`
	Message     string `json:"message"`
	Field       string `json:"field,omitempty"`
}
