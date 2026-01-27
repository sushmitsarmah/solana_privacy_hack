package token

import (
	"context"
	"fmt"
)

// Service handles SPL token management operations.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new token service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// Token represents an SPL token configuration.
type Token struct {
	Mint     string `json:"mint"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Enabled  bool   `json:"enabled"`
}

// ListSupportedResponse contains the list of supported SPL tokens.
type ListSupportedResponse struct {
	Tokens []Token `json:"tokens"`
}

// AddRequest represents a request to add a new SPL token.
type AddRequest struct {
	Mint     string `json:"mint"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Enabled  bool   `json:"enabled"`
}

// AddResponse contains the result of adding a token.
type AddResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UpdateRequest represents a request to update token configuration.
type UpdateRequest struct {
	Enabled  *bool   `json:"enabled,omitempty"`
	Symbol   *string `json:"symbol,omitempty"`
	Decimals *int    `json:"decimals,omitempty"`
}

// UpdateResponse contains the result of updating a token.
type UpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// RemoveResponse contains the result of removing a token.
type RemoveResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ListSupported retrieves all currently supported SPL tokens for payments.
func (s *Service) ListSupported(ctx context.Context) (*ListSupportedResponse, error) {
	var resp ListSupportedResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/api/tokens/supported", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Add adds a new SPL token to the supported tokens list.
// Requires admin authentication via API key.
func (s *Service) Add(ctx context.Context, req AddRequest) (*AddResponse, error) {
	var resp AddResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/tokens/add", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Update modifies the configuration of an existing SPL token.
// Requires admin authentication via API key.
func (s *Service) Update(ctx context.Context, mint string, req UpdateRequest) (*UpdateResponse, error) {
	var resp UpdateResponse
	path := fmt.Sprintf("/shadowpay/api/tokens/update/%s", mint)
	if err := s.doRequest(ctx, "PATCH", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Remove disables an SPL token from the supported tokens list.
// Requires admin authentication via API key.
func (s *Service) Remove(ctx context.Context, mint string) (*RemoveResponse, error) {
	var resp RemoveResponse
	path := fmt.Sprintf("/shadowpay/api/tokens/remove/%s", mint)
	if err := s.doRequest(ctx, "DELETE", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
