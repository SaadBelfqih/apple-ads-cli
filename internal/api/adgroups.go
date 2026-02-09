package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type AdGroupService struct {
	client *Client
}

func (c *Client) AdGroups() *AdGroupService {
	return &AdGroupService{client: c}
}

func (s *AdGroupService) Create(campaignID int64, req *types.AdGroupCreate) (*types.AdGroup, error) {
	body, err := s.client.Post(fmt.Sprintf("/campaigns/%d/adgroups", campaignID), req)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.AdGroup]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *AdGroupService) Get(campaignID, adGroupID int64, fields string) (*types.AdGroup, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d", campaignID, adGroupID)
	if fields != "" {
		path += "?fields=" + url.QueryEscape(fields)
	}
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.AdGroup]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *AdGroupService) List(campaignID int64, limit, offset int, fields string) ([]types.AdGroup, *types.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups", campaignID)
	sep := "?"
	if fields != "" {
		path += sep + "fields=" + url.QueryEscape(fields)
		sep = "&"
	}
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
	var resp types.APIListResponse[types.AdGroup]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *AdGroupService) Find(campaignID int64, selector *types.Selector) ([]types.AdGroup, *types.PageDetail, error) {
	body, err := s.client.Post(fmt.Sprintf("/campaigns/%d/adgroups/find", campaignID), selector)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.AdGroup]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *AdGroupService) FindAll(selector *types.Selector) ([]types.AdGroup, *types.PageDetail, error) {
	body, err := s.client.Post("/adgroups/find", selector)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.AdGroup]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *AdGroupService) Update(campaignID, adGroupID int64, req *types.AdGroupUpdate) (*types.AdGroup, error) {
	body, err := s.client.Put(fmt.Sprintf("/campaigns/%d/adgroups/%d", campaignID, adGroupID), req)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.AdGroup]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *AdGroupService) Delete(campaignID, adGroupID int64) error {
	_, err := s.client.Delete(fmt.Sprintf("/campaigns/%d/adgroups/%d", campaignID, adGroupID))
	return err
}
