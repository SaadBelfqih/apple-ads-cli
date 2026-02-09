package api

import (
	"encoding/json"
	"fmt"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type ReportService struct {
	client *Client
}

func (c *Client) Reports() *ReportService {
	return &ReportService{client: c}
}

func (s *ReportService) Campaigns(req *types.ReportingRequest) (json.RawMessage, error) {
	return s.client.Post("/reports/campaigns", req)
}

func (s *ReportService) AdGroups(campaignID int64, req *types.ReportingRequest) (json.RawMessage, error) {
	return s.client.Post(fmt.Sprintf("/reports/campaigns/%d/adgroups", campaignID), req)
}

func (s *ReportService) Keywords(campaignID int64, adGroupID *int64, req *types.ReportingRequest) (json.RawMessage, error) {
	var path string
	if adGroupID != nil {
		path = fmt.Sprintf("/reports/campaigns/%d/adgroups/%d/keywords", campaignID, *adGroupID)
	} else {
		path = fmt.Sprintf("/reports/campaigns/%d/keywords", campaignID)
	}
	return s.client.Post(path, req)
}

func (s *ReportService) SearchTerms(campaignID int64, adGroupID *int64, req *types.ReportingRequest) (json.RawMessage, error) {
	var path string
	if adGroupID != nil {
		path = fmt.Sprintf("/reports/campaigns/%d/adgroups/%d/searchterms", campaignID, *adGroupID)
	} else {
		path = fmt.Sprintf("/reports/campaigns/%d/searchterms", campaignID)
	}
	return s.client.Post(path, req)
}

func (s *ReportService) Ads(campaignID int64, req *types.ReportingRequest) (json.RawMessage, error) {
	return s.client.Post(fmt.Sprintf("/reports/campaigns/%d/ads", campaignID), req)
}
