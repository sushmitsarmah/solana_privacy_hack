package payment

import (
	"context"
)

// Service handles ZK payment operations.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new payment service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// PrepareRequest represents the data needed to prepare a ZK payment.
type PrepareRequest struct {
	Amount           int64  `json:"amount"`
	Recipient        string `json:"recipient"`
	Resource         string `json:"resource"`
	SenderCommitment string `json:"sender_commitment"`
}

// PrepareResponse contains the data for generating the proof client-side.
type PrepareResponse struct {
	PaymentHash string `json:"payment_hash"`
	// Add other fields as per specific API response for prepare
}

// SettleRequest represents the request to settle a ZK payment via the relayer.
type SettleRequest struct {
	X402Version         int          `json:"x402Version"`
	PaymentHeader       string       `json:"paymentHeader"` // Base64 encoded JSON
	Resource            string       `json:"resource"`
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

// SettleResponse represents the result of the settlement.
type SettleResponse struct {
	Success bool   `json:"success"`
	TxSig   string `json:"tx_sig"`
	Message string `json:"message"`
}

// Prepare initiates the ZK payment flow.
func (s *Service) Prepare(ctx context.Context, req PrepareRequest) (*PrepareResponse, error) {
	var resp PrepareResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/v1/payment/prepare", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Settle submits a ZK proof to the relayer for settlement.
func (s *Service) Settle(ctx context.Context, req SettleRequest) (*SettleResponse, error) {
	var resp SettleResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/v1/payment/settle", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
