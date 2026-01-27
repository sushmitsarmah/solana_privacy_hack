package pool

import (
	"context"
	"fmt"
)

// Service handles privacy pool operations for mixing funds across users.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new pool service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// BalanceResponse contains the user's pool balance and escrow status.
type BalanceResponse struct {
	WalletAddress string `json:"wallet_address"`
	Balance       int64  `json:"balance"`
	MinDeposit    int64  `json:"min_deposit"`
}

// DepositRequest represents a request to deposit SOL into the pool.
// Creates an unsigned transaction for the client to sign.
type DepositRequest struct {
	WalletAddress string `json:"wallet_address"`
	Amount        int64  `json:"amount"` // Must be at least 0.01 SOL (10000000 lamports)
}

// DepositResponse contains the unsigned transaction for deposit.
type DepositResponse struct {
	Transaction string `json:"transaction"` // Unsigned serialized transaction
	Message     string `json:"message,omitempty"`
}

// WithdrawRequest represents a request to withdraw SOL from the pool.
// Incurs a 0.2% fee to discourage savings use.
type WithdrawRequest struct {
	WalletAddress string `json:"wallet_address"`
	Amount        int64  `json:"amount"`
}

// WithdrawResponse contains the withdrawal transaction details.
type WithdrawResponse struct {
	Transaction string `json:"transaction"` // Unsigned serialized transaction
	NetAmount   int64  `json:"net_amount"`  // Amount after 0.2% fee
	Fee         int64  `json:"fee"`
	Message     string `json:"message,omitempty"`
}

// DepositAddressResponse contains the pool PDA address.
type DepositAddressResponse struct {
	DepositAddress string `json:"deposit_address"`
	Network        string `json:"network"`
}

// GetBalance retrieves the user's available pool balance and escrow status.
func (s *Service) GetBalance(ctx context.Context, walletAddress string) (*BalanceResponse, error) {
	var resp BalanceResponse
	path := fmt.Sprintf("/shadowpay/api/pool/balance/%s", walletAddress)
	if err := s.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Deposit creates an unsigned transaction to deposit SOL into the privacy pool.
// Minimum deposit is 0.01 SOL (10000000 lamports).
// Funds are mixed with other users for maximum privacy.
func (s *Service) Deposit(ctx context.Context, req DepositRequest) (*DepositResponse, error) {
	var resp DepositResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/pool/deposit", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Withdraw withdraws SOL from the pool with a 0.2% fee.
// Fee is applied to discourage using the pool as a savings account.
func (s *Service) Withdraw(ctx context.Context, req WithdrawRequest) (*WithdrawResponse, error) {
	var resp WithdrawResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/pool/withdraw", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDepositAddress obtains the pool PDA address for reference.
func (s *Service) GetDepositAddress(ctx context.Context) (*DepositAddressResponse, error) {
	var resp DepositAddressResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/api/pool/deposit-address", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
