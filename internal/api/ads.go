package api

import (
	"encoding/json"
	"fmt"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type AdService struct {
	client *Client
}

func (c *Client) Ads() *AdService {
	return &AdService{client: c}
}

func (s *AdService) Create(campaignID, adGroupID int64, req *types.AdCreate) (*types.Ad, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/ads", campaignID, adGroupID)
	body, err := s.client.Post(path, req)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.Ad]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *AdService) Get(campaignID, adGroupID, adID int64) (*types.Ad, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/ads/%d", campaignID, adGroupID, adID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.Ad]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *AdService) List(campaignID, adGroupID int64, limit, offset int) ([]types.Ad, *types.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/ads", campaignID, adGroupID)
	sep := "?"
	if limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, limit)
		sep = "&"
	}
	if offset > 0 {
		path += fmt.Sprintf("%soffset=%d", sep, offset)
	}
	body, err := s.client.Get(path)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.Ad]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *AdService) Find(campaignID, adGroupID int64, selector *types.Selector) ([]types.Ad, *types.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/ads/find", campaignID, adGroupID)
	body, err := s.client.Post(path, selector)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.Ad]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *AdService) FindAll(selector *types.Selector) ([]types.Ad, *types.PageDetail, error) {
	body, err := s.client.Post("/ads/find", selector)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.Ad]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *AdService) Update(campaignID, adGroupID, adID int64, req *types.AdUpdate) (*types.Ad, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/ads/%d", campaignID, adGroupID, adID)
	body, err := s.client.Put(path, req)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.Ad]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *AdService) Delete(campaignID, adGroupID, adID int64) error {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/ads/%d", campaignID, adGroupID, adID)
	_, err := s.client.Delete(path)
	return err
}
