package api

import (
	"encoding/json"
	"fmt"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type CreativeService struct {
	client *Client
}

func (c *Client) Creatives() *CreativeService {
	return &CreativeService{client: c}
}

func (s *CreativeService) Create(req *types.CreativeCreate) (*types.Creative, error) {
	body, err := s.client.Post("/creatives", req)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.Creative]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *CreativeService) Get(creativeID int64) (*types.Creative, error) {
	body, err := s.client.Get(fmt.Sprintf("/creatives/%d", creativeID))
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.Creative]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *CreativeService) List(limit, offset int) ([]types.Creative, *types.PageDetail, error) {
	path := "/creatives"
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
	var resp types.APIListResponse[types.Creative]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *CreativeService) Find(selector *types.Selector) ([]types.Creative, *types.PageDetail, error) {
	body, err := s.client.Post("/creatives/find", selector)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.Creative]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}
