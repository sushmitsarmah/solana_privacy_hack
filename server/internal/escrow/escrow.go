package escrow

import (
	"context"
	"fmt"

	"sol_privacy/internal/types"
)

// Service handles escrow operations.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new escrow service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// BalanceResponse represents the balance of a user's escrow account.
type BalanceResponse struct {
	WalletAddress string `json:"wallet_address"`
	Balance       int64  `json:"balance"` // In lamports
	Mint          string `json:"mint,omitempty"`
}

// TransactionRequest represents a request to generate an unsigned escrow transaction.
type TransactionRequest struct {
	WalletAddress string `json:"wallet_address"`
	Amount        int64  `json:"amount"` // In lamports or smallest token unit
	Mint          string `json:"mint,omitempty"`
}

// GetBalance retrieves the SOL escrow balance for a wallet.
func (s *Service) GetBalance(ctx context.Context, wallet string) (*BalanceResponse, error) {
	path := fmt.Sprintf("/shadowpay/api/escrow/balance/%s", wallet)
	var resp BalanceResponse
	if err := s.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTokenBalance retrieves the SPL token escrow balance for a wallet.
func (s *Service) GetTokenBalance(ctx context.Context, wallet, mint string) (*BalanceResponse, error) {
	path := fmt.Sprintf("/shadowpay/api/escrow/balance-token/%s/%s", wallet, mint)
	var resp BalanceResponse
	if err := s.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Deposit creates an unsigned transaction to deposit SOL into escrow.
func (s *Service) Deposit(ctx context.Context, req TransactionRequest) (*types.UnsignedTxResponse, error) {
	var resp types.UnsignedTxResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/escrow/deposit", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Withdraw creates an unsigned transaction to withdraw SOL from escrow.
func (s *Service) Withdraw(ctx context.Context, req TransactionRequest) (*types.UnsignedTxResponse, error) {
	var resp types.UnsignedTxResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/escrow/withdraw", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// WithdrawToken creates an unsigned transaction to withdraw SPL tokens from escrow.
func (s *Service) WithdrawToken(ctx context.Context, req TransactionRequest) (*types.UnsignedTxResponse, error) {
	var resp types.UnsignedTxResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/escrow/withdraw-tokens", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
