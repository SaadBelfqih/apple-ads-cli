package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type AppService struct {
	client *Client
}

func (c *Client) Apps() *AppService {
	return &AppService{client: c}
}

func (s *AppService) Search(query string, returnOwnedApps bool, limit, offset int) ([]types.AppInfo, *types.PageDetail, error) {
	path := "/search/apps?query=" + url.QueryEscape(query)
	if returnOwnedApps {
		path += "&returnOwnedApps=true"
	}
	if limit > 0 {
		path += fmt.Sprintf("&limit=%d", limit)
	}
	if offset > 0 {
		path += fmt.Sprintf("&offset=%d", offset)
	}
	body, err := s.client.Get(path)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.AppInfo]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *AppService) Eligibility(selector *types.Selector) ([]types.EligibilityRecord, error) {
	body, err := s.client.Post("/app-eligibility/find", selector)
	if err != nil {
		return nil, err
	}
	var resp types.APIListResponse[types.EligibilityRecord]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *AppService) Details(adamID int64) (*types.AppDetail, error) {
	path := fmt.Sprintf("/apps/%d", adamID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.AppDetail]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *AppService) LocalizedDetails(adamID int64) (*types.LocalizedAppDetail, error) {
	path := fmt.Sprintf("/apps/%d/localized-details", adamID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.LocalizedAppDetail]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}
