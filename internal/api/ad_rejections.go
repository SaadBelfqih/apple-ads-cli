package api

import (
	"encoding/json"
	"fmt"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type AdRejectionService struct {
	client *Client
}

func (c *Client) AdRejections() *AdRejectionService {
	return &AdRejectionService{client: c}
}

func (s *AdRejectionService) Find(selector *types.Selector) ([]types.ProductPageReason, *types.PageDetail, error) {
	body, err := s.client.Post("/product-page-reasons/find", selector)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.ProductPageReason]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *AdRejectionService) Get(productPageReasonID int64) (*types.ProductPageReason, error) {
	path := fmt.Sprintf("/product-page-reasons/%d", productPageReasonID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.ProductPageReason]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *AdRejectionService) FindAssets(adamID int64, selector *types.Selector) ([]types.AppAsset, *types.PageDetail, error) {
	body, err := s.client.Post(fmt.Sprintf("/apps/%d/assets/find", adamID), selector)
	if err != nil {
		return nil, nil, err
	}
	var resp types.APIListResponse[types.AppAsset]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}
