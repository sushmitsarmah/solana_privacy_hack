package authorization

import (
	"context"
	"fmt"
)

// Service handles automated payment authorization for bots and services.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new authorization service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// AuthorizeSpendingRequest represents a request to authorize bot/service spending.
type AuthorizeSpendingRequest struct {
	UserWallet        string `json:"user_wallet"`
	AuthorizedService string `json:"authorized_service"`
	MaxAmountPerTx    string `json:"max_amount_per_tx"`    // In SOL (string format)
	MaxDailySpend     string `json:"max_daily_spend"`      // In SOL (string format)
	ValidUntil        int64  `json:"valid_until"`          // Unix timestamp
	UserSignature     string `json:"user_signature"`       // Base58 encoded signature
}

// AuthorizeSpendingResponse contains the authorization confirmation.
type AuthorizeSpendingResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	AuthorizationID int    `json:"authorization_id"`
}

// RevokeAuthorizationRequest represents a request to revoke bot/service spending authorization.
type RevokeAuthorizationRequest struct {
	UserWallet        string `json:"user_wallet"`
	AuthorizedService string `json:"authorized_service"`
	UserSignature     string `json:"user_signature"` // Base58 encoded signature
}

// RevokeAuthorizationResponse contains the revocation confirmation.
type RevokeAuthorizationResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	AuthorizationID int    `json:"authorization_id"`
}

// Authorization represents a spending authorization for a bot/service.
type Authorization struct {
	ID                int    `json:"id"`
	UserWallet        string `json:"user_wallet"`
	AuthorizedService string `json:"authorized_service"`
	MaxAmountPerTx    int64  `json:"max_amount_per_tx"`    // In lamports
	MaxDailySpend     int64  `json:"max_daily_spend"`      // In lamports
	SpentToday        int64  `json:"spent_today"`          // In lamports
	LastResetDate     string `json:"last_reset_date"`      // Date string
	ValidUntil        int64  `json:"valid_until"`          // Unix timestamp
	Revoked           bool   `json:"revoked"`
	CreatedAt         int64  `json:"created_at"`           // Unix timestamp
}

// ListAuthorizationsResponse contains the list of user authorizations.
type ListAuthorizationsResponse struct {
	Authorizations []Authorization `json:"authorizations"`
}

// AuthorizeSpending registers a bot/service to spend from user's escrow automatically.
// Includes per-transaction and daily limits with expiration.
// User must sign the authorization message to prove ownership.
func (s *Service) AuthorizeSpending(ctx context.Context, req AuthorizeSpendingRequest) (*AuthorizeSpendingResponse, error) {
	var resp AuthorizeSpendingResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/authorize-spending", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RevokeAuthorization revokes a bot/service's permission to spend from user's escrow.
// User must sign the revocation message to prove ownership.
func (s *Service) RevokeAuthorization(ctx context.Context, req RevokeAuthorizationRequest) (*RevokeAuthorizationResponse, error) {
	var resp RevokeAuthorizationResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/revoke-authorization", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListAuthorizations retrieves all active spending authorizations for a user wallet.
// Shows per-transaction limits, daily spending caps, current usage, and expiration.
func (s *Service) ListAuthorizations(ctx context.Context, walletAddress string) (*ListAuthorizationsResponse, error) {
	var resp ListAuthorizationsResponse
	path := fmt.Sprintf("/shadowpay/api/my-authorizations/%s", walletAddress)
	if err := s.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
