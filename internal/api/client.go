package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SaadBelfqih/apple-ads-cli/internal/config"
	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

const (
	baseURL       = "https://api.searchads.apple.com/api/v5"
	maxRetries    = 4
	baseRetryWait = 2 * time.Second
)

// Client is the Apple Ads API HTTP client.
type Client struct {
	httpClient *http.Client
	tokenSrc   *TokenSource
	orgID      string
	verbose    bool
}

// NewClient creates a new API client from config.
func NewClient(cfg *config.Config) (*Client, error) {
	ts, err := NewTokenSource(cfg)
	if err != nil {
		return nil, fmt.Errorf("init auth: %w", err)
	}

	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		tokenSrc:   ts,
		orgID:      cfg.OrgID,
	}, nil
}

// SetVerbose enables verbose logging.
func (c *Client) SetVerbose(v bool) {
	c.verbose = v
}

// SetOrgID overrides the org ID from config.
func (c *Client) SetOrgID(id string) {
	c.orgID = id
}

// Get performs a GET request.
func (c *Client) Get(path string) ([]byte, error) {
	return c.do("GET", path, nil)
}

// Post performs a POST request with a JSON body.
func (c *Client) Post(path string, body any) ([]byte, error) {
	return c.doJSON("POST", path, body)
}

// Put performs a PUT request with a JSON body.
func (c *Client) Put(path string, body any) ([]byte, error) {
	return c.doJSON("PUT", path, body)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) ([]byte, error) {
	return c.do("DELETE", path, nil)
}

// DeleteWithBody performs a POST-as-DELETE (for bulk delete endpoints).
func (c *Client) DeleteWithBody(path string, body any) ([]byte, error) {
	return c.doJSON("POST", path, body)
}

func (c *Client) doJSON(method, path string, body any) ([]byte, error) {
	var buf *bytes.Buffer
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		buf = bytes.NewBuffer(data)
	}

	var bodyReader io.Reader
	var bodyBytes []byte
	if buf != nil {
		bodyBytes = buf.Bytes()
		bodyReader = bytes.NewReader(bodyBytes)
	}

	return c.doWithRetry(method, path, bodyReader, bodyBytes)
}

func (c *Client) do(method, path string, body io.Reader) ([]byte, error) {
	return c.doWithRetry(method, path, body, nil)
}

func (c *Client) doWithRetry(method, path string, body io.Reader, bodyBytes []byte) ([]byte, error) {
	var lastErr error
	var wait time.Duration
	didAuthRefresh := false

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			if c.verbose {
				fmt.Printf("Retrying in %v (attempt %d/%d)...\n", wait, attempt, maxRetries)
			}
			if wait > 0 {
				time.Sleep(wait)
			}

			// Reset body for retry
			if bodyBytes != nil {
				body = bytes.NewReader(bodyBytes)
			}
		}

		result, retry, retryAfter, err := c.doOnce(method, path, body)
		if err == nil {
			return result, nil
		}
		lastErr = err

		// Some auth failures can be resolved by forcing a token refresh, but don't loop forever.
		if apiErr, ok := err.(*APIError); ok && apiErr.StatusCode == http.StatusUnauthorized && !didAuthRefresh {
			didAuthRefresh = true
			if c.verbose {
				fmt.Printf("401 Unauthorized; refreshing token and retrying once...\n")
			}
			if c.tokenSrc != nil {
				c.tokenSrc.Invalidate()
			}
			wait = 0
			continue
		}

		if !retry {
			return nil, err
		}

		wait = retryWait(attempt, retryAfter)
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func retryWait(attempt int, retryAfter time.Duration) time.Duration {
	// attempt is 0-based (0 is the first request). After each failed attempt,
	// compute how long to wait before the next attempt.
	wait := baseRetryWait * time.Duration(1<<attempt)
	if wait > 16*time.Second {
		wait = 16 * time.Second
	}
	if retryAfter > wait {
		wait = retryAfter
	}
	return wait
}

func parseRetryAfter(headerValue string) time.Duration {
	v := strings.TrimSpace(headerValue)
	if v == "" {
		return 0
	}
	// RFC 9110 allows either delta-seconds or an HTTP-date.
	if secs, err := strconv.Atoi(v); err == nil {
		if secs <= 0 {
			return 0
		}
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		d := time.Until(t)
		if d < 0 {
			return 0
		}
		return d
	}
	return 0
}

func isSafeToRetry(method, path string) bool {
	switch method {
	case http.MethodGet, http.MethodPut, http.MethodDelete:
		return true
	case http.MethodPost:
		// Read-only selectors and reports are safe to retry.
		if strings.HasSuffix(path, "/find") {
			return true
		}
		if strings.HasPrefix(path, "/reports/") {
			return true
		}
		// Bulk delete endpoints use POST but are idempotent with respect to creation.
		if strings.Contains(path, "/delete/bulk") {
			return true
		}
		return false
	default:
		return false
	}
}

func (c *Client) doOnce(method, path string, body io.Reader) ([]byte, bool, time.Duration, error) {
	url := baseURL + path

	token, err := c.tokenSrc.Token()
	if err != nil {
		return nil, false, 0, fmt.Errorf("get token: %w", err)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, false, 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if c.orgID != "" {
		req.Header.Set("X-AP-Context", "orgId="+c.orgID)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.verbose {
		fmt.Printf("%s %s\n", method, url)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, isSafeToRetry(method, path), 0, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, isSafeToRetry(method, path), 0, fmt.Errorf("read response: %w", err)
	}

	if c.verbose {
		fmt.Printf("Status: %d\n", resp.StatusCode)
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: resp.StatusCode, RawBody: string(respBody)}
		var errResp types.APIError
		if json.Unmarshal(respBody, &errResp) == nil {
			apiErr.Errors = errResp.Errors
		}

		// Always retry on 429. Respect Retry-After when present.
		if resp.StatusCode == http.StatusTooManyRequests {
			apiErr.RetryAfter = parseRetryAfter(resp.Header.Get("Retry-After"))
			return nil, true, apiErr.RetryAfter, apiErr
		}

		// Retry on 5xx only when it's safe to do so (avoid duplicating POST creates).
		if resp.StatusCode >= 500 && isSafeToRetry(method, path) {
			apiErr.RetryAfter = parseRetryAfter(resp.Header.Get("Retry-After"))
			return nil, true, apiErr.RetryAfter, apiErr
		}
		return nil, false, 0, apiErr
	}

	return respBody, false, 0, nil
}
