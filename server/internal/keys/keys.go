package keys

import (
	"context"
	"fmt"
)

// Service handles API key operations.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new keys service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// GenerateRequest represents the payload to create a new API key.
type GenerateRequest struct {
	WalletAddress  string `json:"wallet_address"`
	TreasuryWallet string `json:"treasury_wallet,omitempty"`
}

// Response represents the API key details.
type Response struct {
	APIKey string `json:"api_key"`
	Wallet string `json:"wallet_address"`
}

// LimitsResponse represents rate limit information.
type LimitsResponse struct {
	Limit     int64 `json:"limit"`
	Remaining int64 `json:"remaining"`
	Reset     int64 `json:"reset"`
}

// Create generates a new API key for a wallet.
func (s *Service) Create(ctx context.Context, req GenerateRequest) (*Response, error) {
	var resp Response
	if err := s.doRequest(ctx, "POST", "/shadowpay/v1/keys/new", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetByWallet retrieves an existing API key for a wallet.
func (s *Service) GetByWallet(ctx context.Context, wallet string) (*Response, error) {
	path := fmt.Sprintf("/shadowpay/v1/keys/by-wallet/%s", wallet)
	var resp Response
	if err := s.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Rotate invalidates the old key and generates a new one.
func (s *Service) Rotate(ctx context.Context) (*Response, error) {
	var resp Response
	if err := s.doRequest(ctx, "POST", "/shadowpay/v1/keys/rotate", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLimits retrieves the current rate limits for the authenticated key.
func (s *Service) GetLimits(ctx context.Context) (*LimitsResponse, error) {
	var resp LimitsResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/v1/keys/limits", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
