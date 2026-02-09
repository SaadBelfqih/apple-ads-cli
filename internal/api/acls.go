package api

import (
	"encoding/json"
	"fmt"

	"github.com/SaadBelfqih/apple-ads-cli/internal/types"
)

// ACLService handles ACL and user API calls.
type ACLService struct {
	client *Client
}

// ACLs returns an ACLService.
func (c *Client) ACLs() *ACLService {
	return &ACLService{client: c}
}

// List retrieves all user ACLs.
func (s *ACLService) List() ([]types.UserACL, error) {
	body, err := s.client.Get("/acls")
	if err != nil {
		return nil, err
	}

	var resp types.APIListResponse[types.UserACL]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}

// Me retrieves the calling user's details.
func (s *ACLService) Me() (*types.MeDetail, error) {
	body, err := s.client.Get("/me")
	if err != nil {
		return nil, err
	}

	var resp types.APIResponse[types.MeDetail]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return resp.Data, nil
}
