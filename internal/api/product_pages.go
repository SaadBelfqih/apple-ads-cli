package api

import (
	"encoding/json"
	"fmt"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type ProductPageService struct {
	client *Client
}

func (c *Client) ProductPages() *ProductPageService {
	return &ProductPageService{client: c}
}

func (s *ProductPageService) List(adamID int64) ([]types.ProductPageDetail, error) {
	path := fmt.Sprintf("/apps/%d/product-pages", adamID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIListResponse[types.ProductPageDetail]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *ProductPageService) Get(productPageID string, adamID int64) (*types.ProductPageDetail, error) {
	path := fmt.Sprintf("/apps/%d/product-pages/%s", adamID, productPageID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.ProductPageDetail]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *ProductPageService) Locales(productPageID string, adamID int64) ([]types.ProductPageLocaleDetail, error) {
	path := fmt.Sprintf("/apps/%d/product-pages/%s/locale-details", adamID, productPageID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIListResponse[types.ProductPageLocaleDetail]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *ProductPageService) Countries() ([]types.CountryOrRegion, error) {
	body, err := s.client.Get("/countries-or-regions")
	if err != nil {
		return nil, err
	}
	var resp types.APIListResponse[types.CountryOrRegion]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *ProductPageService) DeviceSizes() (json.RawMessage, error) {
	body, err := s.client.Get("/creativeappmappings/devices")
	if err != nil {
		return nil, err
	}
	return body, nil
}
