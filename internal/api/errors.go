package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

// APIError wraps Apple Ads API error responses.
type APIError struct {
	StatusCode int
	Errors     []types.APIErrorDetail
	RawBody    string
	RetryAfter time.Duration
}

func (e *APIError) Error() string {
	retryInfo := ""
	if e.RetryAfter > 0 {
		retryInfo = fmt.Sprintf(" (retry after %s)", e.RetryAfter)
	}

	if len(e.Errors) == 0 {
		return fmt.Sprintf("API error %d%s: %s", e.StatusCode, retryInfo, e.RawBody)
	}

	var msgs []string
	for _, err := range e.Errors {
		msg := err.Message
		if err.Field != "" {
			msg = fmt.Sprintf("%s (field: %s)", msg, err.Field)
		}
		msgs = append(msgs, msg)
	}
	return fmt.Sprintf("API error %d%s: %s", e.StatusCode, retryInfo, strings.Join(msgs, "; "))
}
