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

// DepositRequest represents a request to deposit funds for ZK payments.
type DepositRequest struct {
	WalletAddress string `json:"wallet_address"`
	Amount        int64  `json:"amount"` // Amount in lamports
}

// DepositResponse contains the unsigned transaction for deposit.
type DepositResponse struct {
	UnsignedTxBase64      string `json:"unsigned_tx_base64"`
	RecentBlockhash       string `json:"recent_blockhash"`
	LastValidBlockHeight  int    `json:"last_valid_block_height"`
}

// WithdrawRequest represents a request to withdraw funds from payment account.
type WithdrawRequest struct {
	WalletAddress string `json:"wallet_address"`
	Amount        int64  `json:"amount"` // Amount in lamports
}

// WithdrawResponse contains the unsigned transaction for withdrawal.
type WithdrawResponse struct {
	UnsignedTxBase64      string `json:"unsigned_tx_base64"`
	RecentBlockhash       string `json:"recent_blockhash"`
	LastValidBlockHeight  int    `json:"last_valid_block_height"`
	Message               string `json:"message,omitempty"`
}

// PrepareRequest represents the data needed to prepare a ZK payment.
type PrepareRequest struct {
	ReceiverCommitment string `json:"receiver_commitment"` // Base58 or hex encoded
	Amount             int64  `json:"amount"`
	TokenMint          string `json:"token_mint,omitempty"` // Optional SPL token mint
}

// PrepareResponse contains the data for generating the proof client-side.
type PrepareResponse struct {
	PaymentHash    string `json:"payment_hash"`
	Transaction    string `json:"transaction"`    // Unsigned transaction
	Commitment     string `json:"commitment"`
	Message        string `json:"message,omitempty"`
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

// Deposit creates an unsigned transaction to deposit funds for ZK payments.
// The transaction must be signed by the client and submitted to the network.
func (s *Service) Deposit(ctx context.Context, req DepositRequest) (*DepositResponse, error) {
	var resp DepositResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/v1/payment/deposit", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Withdraw creates an unsigned transaction to withdraw funds from the payment account.
// The transaction must be signed by the client and submitted to the network.
func (s *Service) Withdraw(ctx context.Context, req WithdrawRequest) (*WithdrawResponse, error) {
	var resp WithdrawResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/v1/payment/withdraw", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
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

// AuthorizeRequest represents a request to validate escrow and get an access token.
type AuthorizeRequest struct {
	Commitment string `json:"commitment"`
	Nullifier  string `json:"nullifier"`
	Amount     int64  `json:"amount"`
	Merchant   string `json:"merchant"` // Merchant wallet address
}

// AuthorizeResponse contains the JWT access token for x402 flow.
type AuthorizeResponse struct {
	Success     bool   `json:"success"`
	AccessToken string `json:"access_token"` // JWT token
	ExpiresIn   int    `json:"expires_in"`   // Seconds until expiration
	Message     string `json:"message,omitempty"`
}

// VerifyAccessRequest represents query parameters for access token verification.
type VerifyAccessRequest struct {
	Token string `json:"token"`
}

// VerifyAccessResponse contains the token validation result.
type VerifyAccessResponse struct {
	Valid      bool   `json:"valid"`
	Commitment string `json:"commitment,omitempty"`
	Merchant   string `json:"merchant,omitempty"`
	Amount     int64  `json:"amount,omitempty"`
	ExpiresAt  string `json:"expires_at,omitempty"`
	Message    string `json:"message,omitempty"`
}

// Authorize validates escrow balance and returns an access token for the x402 payment flow.
// The user must have sufficient escrow balance for the payment amount.
func (s *Service) Authorize(ctx context.Context, req AuthorizeRequest) (*AuthorizeResponse, error) {
	var resp AuthorizeResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/v1/payment/authorize", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// VerifyAccess verifies the validity of a JWT access token.
// Used by merchants to validate payment authorization before providing access to resources.
func (s *Service) VerifyAccess(ctx context.Context, token string) (*VerifyAccessResponse, error) {
	var resp VerifyAccessResponse
	req := VerifyAccessRequest{Token: token}
	if err := s.doRequest(ctx, "GET", "/shadowpay/v1/payment/verify-access", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
