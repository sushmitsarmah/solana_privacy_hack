package intent

import (
	"context"
)

// Service handles payment intent operations.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new intent service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// CreateRequest represents a request to create a payment intent.
type CreateRequest struct {
	Amount    int64  `json:"amount"`
	Recipient string `json:"recipient"`
	Reference string `json:"reference"` // Often a unique ID
}

// Response represents a payment intent response.
type Response struct {
	IntentID     string `json:"intent_id"`
	ClientSecret string `json:"client_secret"`
	Status       string `json:"status"`
}

// VerifyRequest represents a request to verify a payment intent.
type VerifyRequest struct {
	IntentID string `json:"intent_id"`
}

// VerifyResponse represents a payment verification response.
type VerifyResponse struct {
	Status    string `json:"status"`
	Verified  bool   `json:"verified"`
	Timestamp int64  `json:"timestamp"`
}

// Create creates a new standard payment intent.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*Response, error) {
	var resp Response
	if err := s.doRequest(ctx, "POST", "/shadowpay/v1/pay/intent", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Verify checks the status of a payment intent.
func (s *Service) Verify(ctx context.Context, intentID string) (*VerifyResponse, error) {
	req := VerifyRequest{IntentID: intentID}
	var resp VerifyResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/v1/pay/verify", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPublicKey retrieves the server's public key for payment verification.
func (s *Service) GetPublicKey(ctx context.Context) (string, error) {
	type keyResponse struct {
		PublicKey string `json:"public_key"`
	}

	var resp keyResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/v1/pay/pubkey", nil, &resp); err != nil {
		return "", err
	}
	return resp.PublicKey, nil
}
