package api

import (
	"encoding/json"
	"fmt"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type ImpressionShareService struct {
	client *Client
}

func (c *Client) ImpressionShare() *ImpressionShareService {
	return &ImpressionShareService{client: c}
}

func (s *ImpressionShareService) Create(req *types.CustomReportRequest) (*types.CustomReportResponse, error) {
	body, err := s.client.Post("/custom-reports", req)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.CustomReportResponse]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *ImpressionShareService) Get(reportID int64) (*types.CustomReportResponse, error) {
	path := fmt.Sprintf("/custom-reports/%d", reportID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.CustomReportResponse]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *ImpressionShareService) List(limit, offset int) ([]types.CustomReportResponse, *types.PageDetail, error) {
	path := "/custom-reports"
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
	var resp types.APIListResponse[types.CustomReportResponse]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}
