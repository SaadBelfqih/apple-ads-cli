package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type GeoService struct {
	client *Client
}

func (c *Client) Geo() *GeoService {
	return &GeoService{client: c}
}

func (s *GeoService) Search(query, countryCode, entity string, limit int) ([]types.SearchEntity, error) {
	path := "/search/geo?query=" + url.QueryEscape(query)
	if countryCode != "" {
		path += "&countrycode=" + url.QueryEscape(countryCode)
	}
	if entity != "" {
		path += "&entity=" + url.QueryEscape(entity)
	}
	if limit > 0 {
		path += fmt.Sprintf("&limit=%d", limit)
	}
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIListResponse[types.SearchEntity]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *GeoService) Get(geoID string) (json.RawMessage, error) {
	path := "/geodata?geoId=" + url.QueryEscape(geoID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	return body, nil
}
