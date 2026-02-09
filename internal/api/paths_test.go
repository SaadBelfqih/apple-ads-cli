package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func newTestClient(t *testing.T, orgID string, rt roundTripFunc) *Client {
	t.Helper()

	return &Client{
		httpClient: &http.Client{Transport: rt, Timeout: 5 * time.Second},
		tokenSrc:   &TokenSource{accessToken: "test-token", expiresAt: time.Now().Add(24 * time.Hour)},
		orgID:      orgID,
	}
}

func okJSON(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestKeywordBulkPaths(t *testing.T) {
	var got []string
	c := newTestClient(t, "123", func(req *http.Request) (*http.Response, error) {
		got = append(got, req.Method+" "+req.URL.Path)
		// Keyword endpoints use list/response envelopes.
		return okJSON(`{"data":[]}`), nil
	})

	if _, err := c.Keywords().Create(1, 2, []types.Keyword{}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := c.Keywords().Delete(1, 2, []int64{10, 11}); err != nil {
		t.Fatalf("delete bulk: %v", err)
	}
	if err := c.Keywords().DeleteOne(1, 2, 10); err != nil {
		t.Fatalf("delete one: %v", err)
	}

	want := []string{
		"POST /api/v5/campaigns/1/adgroups/2/targetingkeywords/bulk",
		"POST /api/v5/campaigns/1/adgroups/2/targetingkeywords/delete/bulk",
		"DELETE /api/v5/campaigns/1/adgroups/2/targetingkeywords/10",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d requests, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("request[%d]=%q, want %q", i, got[i], want[i])
		}
	}
}

func TestNegativeKeywordBulkPaths(t *testing.T) {
	var got []string
	c := newTestClient(t, "123", func(req *http.Request) (*http.Response, error) {
		got = append(got, req.Method+" "+req.URL.Path)
		return okJSON(`{"data":[]}`), nil
	})

	if _, err := c.Negatives().CampaignCreate(99, []types.NegativeKeyword{}); err != nil {
		t.Fatalf("campaign create: %v", err)
	}
	if err := c.Negatives().CampaignDelete(99, []int64{1, 2}); err != nil {
		t.Fatalf("campaign delete bulk: %v", err)
	}
	if _, err := c.Negatives().AdGroupCreate(99, 100, []types.NegativeKeyword{}); err != nil {
		t.Fatalf("ad group create: %v", err)
	}
	if err := c.Negatives().AdGroupDelete(99, 100, []int64{1, 2}); err != nil {
		t.Fatalf("ad group delete bulk: %v", err)
	}

	want := []string{
		"POST /api/v5/campaigns/99/negativekeywords/bulk",
		"POST /api/v5/campaigns/99/negativekeywords/delete/bulk",
		"POST /api/v5/campaigns/99/adgroups/100/negativekeywords/bulk",
		"POST /api/v5/campaigns/99/adgroups/100/negativekeywords/delete/bulk",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d requests, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("request[%d]=%q, want %q", i, got[i], want[i])
		}
	}
}

func TestProductPagesPaths(t *testing.T) {
	var got []string
	c := newTestClient(t, "123", func(req *http.Request) (*http.Response, error) {
		got = append(got, req.Method+" "+req.URL.Path)

		switch {
		case strings.HasSuffix(req.URL.Path, "/product-pages/pp-123"):
			return okJSON(`{"data":{}}`), nil
		case strings.Contains(req.URL.Path, "/creativeappmappings/devices"):
			return okJSON(`[]`), nil
		default:
			return okJSON(`{"data":[]}`), nil
		}
	})

	if _, err := c.ProductPages().List(123456789); err != nil {
		t.Fatalf("list: %v", err)
	}
	if _, err := c.ProductPages().Get("pp-123", 123456789); err != nil {
		t.Fatalf("get: %v", err)
	}
	if _, err := c.ProductPages().Locales("pp-123", 123456789); err != nil {
		t.Fatalf("locales: %v", err)
	}
	if _, err := c.ProductPages().Countries(); err != nil {
		t.Fatalf("countries: %v", err)
	}
	if _, err := c.ProductPages().DeviceSizes(); err != nil {
		t.Fatalf("device sizes: %v", err)
	}

	want := []string{
		"GET /api/v5/apps/123456789/product-pages",
		"GET /api/v5/apps/123456789/product-pages/pp-123",
		"GET /api/v5/apps/123456789/product-pages/pp-123/locale-details",
		"GET /api/v5/countries-or-regions",
		"GET /api/v5/creativeappmappings/devices",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d requests, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("request[%d]=%q, want %q", i, got[i], want[i])
		}
	}
}

func TestAdRejectionPaths(t *testing.T) {
	var got []string
	c := newTestClient(t, "123", func(req *http.Request) (*http.Response, error) {
		got = append(got, req.Method+" "+req.URL.Path)
		if req.Method == http.MethodGet && strings.HasPrefix(req.URL.Path, "/api/v5/product-page-reasons/") {
			return okJSON(`{"data":{}}`), nil
		}
		return okJSON(`{"data":[]}`), nil
	})

	if _, _, err := c.AdRejections().Find(&types.Selector{}); err != nil {
		t.Fatalf("find: %v", err)
	}
	if _, err := c.AdRejections().Get(42); err != nil {
		t.Fatalf("get: %v", err)
	}
	if _, _, err := c.AdRejections().FindAssets(123456789, &types.Selector{}); err != nil {
		t.Fatalf("find assets: %v", err)
	}

	want := []string{
		"POST /api/v5/product-page-reasons/find",
		"GET /api/v5/product-page-reasons/42",
		"POST /api/v5/apps/123456789/assets/find",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d requests, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("request[%d]=%q, want %q", i, got[i], want[i])
		}
	}
}

func TestNoOrgHeaderForACL(t *testing.T) {
	c := newTestClient(t, "", func(req *http.Request) (*http.Response, error) {
		if v := req.Header.Get("X-AP-Context"); v != "" {
			t.Fatalf("expected no X-AP-Context header, got %q", v)
		}
		if v := req.Header.Get("Authorization"); !strings.HasPrefix(v, "Bearer ") {
			t.Fatalf("expected Authorization bearer token, got %q", v)
		}
		return okJSON(`{"data":[]}`), nil
	})

	if _, err := c.ACLs().List(); err != nil {
		t.Fatalf("acls list: %v", err)
	}
}
