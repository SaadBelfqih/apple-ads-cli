package api

import (
	"encoding/json"
	"fmt"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

type BudgetOrderService struct {
	client *Client
}

func (c *Client) BudgetOrders() *BudgetOrderService {
	return &BudgetOrderService{client: c}
}

func (s *BudgetOrderService) Create(req *types.BudgetOrderCreate) (*types.BudgetOrder, error) {
	body, err := s.client.Post("/budgetorders", req)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.BudgetOrder]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *BudgetOrderService) Get(orderID int64) (*types.BudgetOrder, error) {
	path := fmt.Sprintf("/budgetorders/%d", orderID)
	body, err := s.client.Get(path)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.BudgetOrder]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

func (s *BudgetOrderService) List(limit, offset int) ([]types.BudgetOrder, *types.PageDetail, error) {
	path := "/budgetorders"
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
	var resp types.APIListResponse[types.BudgetOrder]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, resp.Pagination, nil
}

func (s *BudgetOrderService) Update(orderID int64, req *types.BudgetOrderUpdate) (*types.BudgetOrder, error) {
	path := fmt.Sprintf("/budgetorders/%d", orderID)
	body, err := s.client.Put(path, req)
	if err != nil {
		return nil, err
	}
	var resp types.APIResponse[types.BudgetOrder]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}
