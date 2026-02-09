package api

import (
	"encoding/json"
	"fmt"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type NegativeKeywordService struct {
	client *Client
}

func (c *Client) Negatives() *NegativeKeywordService {
	return &NegativeKeywordService{client: c}
}

// Campaign-level negative keywords

func (s *NegativeKeywordService) CampaignCreate(campaignID int64, keywords []types.NegativeKeyword) ([]types.NegativeKeyword, error) {
	// Bulk create endpoint per Apple Ads API docs.
	path := fmt.Sprintf("/campaigns/%d/negativekeywords/bulk", campaignID)
	body, err := s.client.Post(path, keywords)
	if err != nil {
		return nil, err
	}
	var resp types.APIListResponse[types.NegativeKeyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *NegativeKeywordService) CampaignGet(campaignID, keywordID int64) (*types.NegativeKeyword, error) {
	path := fmt.Sprintf("/campaigns/%d/negativekeywords/%d", campaignID, keywordID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.NegativeKeyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *NegativeKeywordService) CampaignList(campaignID int64, limit, offset int) ([]types.NegativeKeyword, *types.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/negativekeywords", campaignID)
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
	var resp types.APIListResponse[types.NegativeKeyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *NegativeKeywordService) CampaignFind(campaignID int64, selector *types.Selector) ([]types.NegativeKeyword, *types.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/negativekeywords/find", campaignID)
	body, err := s.client.Post(path, selector)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.NegativeKeyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *NegativeKeywordService) CampaignUpdate(campaignID int64, keywords []types.NegativeKeyword) ([]types.NegativeKeyword, error) {
	// Bulk update endpoint per Apple Ads API docs.
	path := fmt.Sprintf("/campaigns/%d/negativekeywords/bulk", campaignID)
	body, err := s.client.Put(path, keywords)
	if err != nil {
		return nil, err
	}
	var resp types.APIListResponse[types.NegativeKeyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *NegativeKeywordService) CampaignDelete(campaignID int64, keywordIDs []int64) error {
	// Bulk delete endpoint per Apple Ads API docs.
	path := fmt.Sprintf("/campaigns/%d/negativekeywords/delete/bulk", campaignID)
	_, err := s.client.DeleteWithBody(path, keywordIDs)
	return err
}

// Ad group-level negative keywords

func (s *NegativeKeywordService) AdGroupCreate(campaignID, adGroupID int64, keywords []types.NegativeKeyword) ([]types.NegativeKeyword, error) {
	// Bulk create endpoint per Apple Ads API docs.
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/negativekeywords/bulk", campaignID, adGroupID)
	body, err := s.client.Post(path, keywords)
	if err != nil {
		return nil, err
	}
	var resp types.APIListResponse[types.NegativeKeyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *NegativeKeywordService) AdGroupGet(campaignID, adGroupID, keywordID int64) (*types.NegativeKeyword, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/negativekeywords/%d", campaignID, adGroupID, keywordID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.NegativeKeyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *NegativeKeywordService) AdGroupList(campaignID, adGroupID int64, limit, offset int) ([]types.NegativeKeyword, *types.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/negativekeywords", campaignID, adGroupID)
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
	var resp types.APIListResponse[types.NegativeKeyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *NegativeKeywordService) AdGroupFind(campaignID, adGroupID int64, selector *types.Selector) ([]types.NegativeKeyword, *types.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/negativekeywords/find", campaignID, adGroupID)
	body, err := s.client.Post(path, selector)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.NegativeKeyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *NegativeKeywordService) AdGroupUpdate(campaignID, adGroupID int64, keywords []types.NegativeKeyword) ([]types.NegativeKeyword, error) {
	// Bulk update endpoint per Apple Ads API docs.
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/negativekeywords/bulk", campaignID, adGroupID)
	body, err := s.client.Put(path, keywords)
	if err != nil {
		return nil, err
	}
	var resp types.APIListResponse[types.NegativeKeyword]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *NegativeKeywordService) AdGroupDelete(campaignID, adGroupID int64, keywordIDs []int64) error {
	// Bulk delete endpoint per Apple Ads API docs.
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/negativekeywords/delete/bulk", campaignID, adGroupID)
	_, err := s.client.DeleteWithBody(path, keywordIDs)
	return err
}
