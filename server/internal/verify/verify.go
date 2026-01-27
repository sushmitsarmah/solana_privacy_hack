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

// SupportedResponse contains the supported x402 payment methods.
type SupportedResponse struct {
	X402Version int      `json:"x402Version"`
	Schemes     []Scheme `json:"schemes"`
}

// Scheme represents a supported payment scheme.
type Scheme struct {
	Scheme      string `json:"scheme"`      // e.g., "zkproof"
	Network     string `json:"network"`     // e.g., "solana-mainnet"
	Description string `json:"description"`
}

// VerifyRequest represents a full x402 verify request.
type VerifyRequest struct {
	X402Version         int          `json:"x402Version"`
	PaymentHeader       string       `json:"paymentHeader"` // Base64 encoded JSON
	PaymentRequirements Requirements `json:"paymentRequirements"`
}

// Requirements details the constraints for the payment.
type Requirements struct {
	Scheme            string `json:"scheme"`             // e.g., "zkproof"
	Network           string `json:"network"`            // e.g., "solana-mainnet"
	MaxAmountRequired string `json:"maxAmountRequired"`  // In SOL (string format)
	Resource          string `json:"resource"`
	Description       string `json:"description"`
	MimeType          string `json:"mimeType"`
	PayTo             string `json:"payTo"`
	MaxTimeoutSeconds int    `json:"maxTimeoutSeconds"`
}

// VerifyResponse represents the x402 verification result.
type VerifyResponse struct {
	IsValid       bool   `json:"isValid"`
	InvalidReason string `json:"invalidReason,omitempty"`
	PaymentToken  string `json:"paymentToken,omitempty"` // Token for settlement
}

// SettleRequest represents a request to execute on-chain payment settlement.
type SettleRequest struct {
	X402Version         int          `json:"x402Version"`
	PaymentHeader       string       `json:"paymentHeader"` // Base64 encoded JSON
	PaymentRequirements Requirements `json:"paymentRequirements"`
	Resource            string       `json:"resource"`
	Metadata            string       `json:"metadata,omitempty"`
}

// SettleResponse represents the settlement result.
type SettleResponse struct {
	Success     bool   `json:"success"`
	TxSignature string `json:"tx_signature"`
	NetworkID   string `json:"network_id"`
	Message     string `json:"message,omitempty"`
}

// PremiumResponse represents the demo paywall resource response.
type PremiumResponse struct {
	Content    string `json:"content,omitempty"`    // HTML content if paid
	StatusCode int    `json:"status_code"`          // 200 or 402
	Message    string `json:"message,omitempty"`
}

// X402 verifies a payment token or proof (simplified version).
func (s *Service) X402(ctx context.Context, token string) (*Response, error) {
	req := Request{Token: token}
	var resp Response
	if err := s.doRequest(ctx, "POST", "/shadowpay/verify", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSupported retrieves the supported x402 payment methods.
func (s *Service) GetSupported(ctx context.Context) (*SupportedResponse, error) {
	var resp SupportedResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/supported", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Verify validates a zero-knowledge proof payment per x402 standard.
// Returns a payment token that can be used for settlement.
func (s *Service) Verify(ctx context.Context, req VerifyRequest) (*VerifyResponse, error) {
	var resp VerifyResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/verify", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Settle executes on-chain payment settlement per x402 protocol.
// Can be used in both manual and automated (relayer) modes.
func (s *Service) Settle(ctx context.Context, req SettleRequest) (*SettleResponse, error) {
	var resp SettleResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/settle", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPremium retrieves the demo paywalled resource (0.001 SOL).
// Returns 402 Payment Required if unpaid, or premium content if paid.
func (s *Service) GetPremium(ctx context.Context) (*PremiumResponse, error) {
	var resp PremiumResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/premium", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
