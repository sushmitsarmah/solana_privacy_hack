package verify

import (
	"context"
)

// Service handles X402 verification operations.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new verify service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// Request represents a request to verify an X402 token.
type Request struct {
	Token string `json:"token"`
}

// Response represents an X402 verification response.
type Response struct {
	Valid     bool   `json:"valid"`
	ExpiresAt int64  `json:"expires_at"`
	Scope     string `json:"scope"`
}

// X402 verifies a payment token or proof.
func (s *Service) X402(ctx context.Context, token string) (*Response, error) {
	req := Request{Token: token}
	var resp Response
	if err := s.doRequest(ctx, "POST", "/shadowpay/verify", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
