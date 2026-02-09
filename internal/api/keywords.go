package api

import (
	"encoding/json"
	"fmt"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type KeywordService struct {
	client *Client
}

func (c *Client) Keywords() *KeywordService {
	return &KeywordService{client: c}
}

func (s *KeywordService) Create(campaignID, adGroupID int64, keywords []types.Keyword) ([]types.Keyword, error) {
	// Bulk create endpoint per Apple Ads API docs.
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/bulk", campaignID, adGroupID)
	body, err := s.client.Post(path, keywords)
	if err != nil {
		return nil, err
	}
	var resp types.APIListResponse[types.Keyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *KeywordService) Get(campaignID, adGroupID, keywordID int64) (*types.Keyword, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/%d", campaignID, adGroupID, keywordID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.Keyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *KeywordService) List(campaignID, adGroupID int64, limit, offset int) ([]types.Keyword, *types.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords", campaignID, adGroupID)
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
	var resp types.APIListResponse[types.Keyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *KeywordService) Find(campaignID, adGroupID int64, selector *types.Selector) ([]types.Keyword, *types.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/find", campaignID, adGroupID)
	body, err := s.client.Post(path, selector)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.Keyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *KeywordService) FindCampaign(campaignID int64, selector *types.Selector) ([]types.Keyword, *types.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/targetingkeywords/find", campaignID)
	body, err := s.client.Post(path, selector)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.Keyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *KeywordService) Update(campaignID, adGroupID int64, keywords []types.Keyword) ([]types.Keyword, error) {
	// Bulk update endpoint per Apple Ads API docs.
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/bulk", campaignID, adGroupID)
	body, err := s.client.Put(path, keywords)
	if err != nil {
		return nil, err
	}
	var resp types.APIListResponse[types.Keyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

// Delete bulk-deletes keywords by IDs.
func (s *KeywordService) Delete(campaignID, adGroupID int64, keywordIDs []int64) error {
	// Bulk delete endpoint per Apple Ads API docs.
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/delete/bulk", campaignID, adGroupID)
	_, err := s.client.DeleteWithBody(path, keywordIDs)
	return err
}

// DeleteOne deletes a single keyword.
func (s *KeywordService) DeleteOne(campaignID, adGroupID, keywordID int64) error {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/%d", campaignID, adGroupID, keywordID)
	_, err := s.client.Delete(path)
	return err
}
