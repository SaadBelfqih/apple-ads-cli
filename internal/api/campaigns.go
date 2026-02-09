package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

// CampaignService handles campaign API calls.
type CampaignService struct {
	client *Client
}

// Campaigns returns a CampaignService.
func (c *Client) Campaigns() *CampaignService {
	return &CampaignService{client: c}
}

// Create creates a new campaign.
func (s *CampaignService) Create(req *types.CampaignCreate) (*types.Campaign, error) {
	body, err := s.client.Post("/campaigns", req)
	if err != nil {
		return nil, err
	}

	var resp types.APIResponse[types.Campaign]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

// Get retrieves a single campaign by ID.
func (s *CampaignService) Get(campaignID int64, fields string) (*types.Campaign, error) {
	path := fmt.Sprintf("/campaigns/%d", campaignID)
	if fields != "" {
		path += "?fields=" + url.QueryEscape(fields)
	}

	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}

	var resp types.APIResponse[types.Campaign]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

// List retrieves all campaigns.
func (s *CampaignService) List(limit, offset int, fields string) ([]types.Campaign, *types.PageDetail, error) {
	path := "/campaigns"
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

	var resp types.APIListResponse[types.Campaign]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

// Find searches campaigns using a selector.
func (s *CampaignService) Find(selector *types.Selector) ([]types.Campaign, *types.PageDetail, error) {
	body, err := s.client.Post("/campaigns/find", selector)
	if err != nil {
		return nil, nil, err
	}

	var resp types.APIListResponse[types.Campaign]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

// Update updates a campaign. Uses the {"campaign":{...}} envelope.
func (s *CampaignService) Update(campaignID int64, req *types.CampaignUpdate) (*types.Campaign, error) {
	envelope := &types.UpdateCampaignRequest{Campaign: req}
	body, err := s.client.Put(fmt.Sprintf("/campaigns/%d", campaignID), envelope)
	if err != nil {
		return nil, err
	}

	var resp types.APIResponse[types.Campaign]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

// Delete deletes a campaign.
func (s *CampaignService) Delete(campaignID int64) error {
	_, err := s.client.Delete(fmt.Sprintf("/campaigns/%d", campaignID))
	return err
}
