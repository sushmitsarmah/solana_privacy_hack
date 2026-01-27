package merchant

import (
	"context"
)

// Service handles merchant operations including earnings, analytics, and withdrawals.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new merchant service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// TokenEarnings represents earnings for a specific token.
type TokenEarnings struct {
	TokenMint string `json:"token_mint"`
	Symbol    string `json:"symbol,omitempty"`
	Amount    int64  `json:"amount"`
	UsdValue  string `json:"usd_value,omitempty"`
}

// EarningsResponse contains the merchant's total earnings breakdown.
type EarningsResponse struct {
	TotalEarnings     int64           `json:"total_earnings"` // Total in lamports
	TotalUsdValue     string          `json:"total_usd_value,omitempty"`
	TokenBreakdown    []TokenEarnings `json:"token_breakdown"`
	WithdrawableSOL   int64           `json:"withdrawable_sol"`
	PendingSettlement int64           `json:"pending_settlement"`
}

// AnalyticsRequest represents query parameters for analytics.
type AnalyticsRequest struct {
	StartDate string `json:"start_date,omitempty"` // ISO 8601 format
	EndDate   string `json:"end_date,omitempty"`   // ISO 8601 format
	Interval  string `json:"interval,omitempty"`   // "hour", "day", "week", "month"
}

// PaymentStats contains payment statistics for a time period.
type PaymentStats struct {
	Timestamp    string `json:"timestamp"`
	PaymentCount int    `json:"payment_count"`
	TotalAmount  int64  `json:"total_amount"`
	UniqueUsers  int    `json:"unique_users,omitempty"`
}

// AnalyticsResponse contains payment analytics data.
type AnalyticsResponse struct {
	TotalPayments    int            `json:"total_payments"`
	TotalVolume      int64          `json:"total_volume"`
	AveragePayment   int64          `json:"average_payment"`
	UniqueCustomers  int            `json:"unique_customers"`
	TimeSeries       []PaymentStats `json:"time_series"`
	TopResources     []ResourceStat `json:"top_resources,omitempty"`
	SuccessRate      float64        `json:"success_rate"`
	PendingPayments  int            `json:"pending_payments"`
}

// ResourceStat contains statistics for a specific resource.
type ResourceStat struct {
	Resource     string `json:"resource"`
	PaymentCount int    `json:"payment_count"`
	TotalAmount  int64  `json:"total_amount"`
}

// WithdrawRequest represents a request to withdraw merchant earnings.
type WithdrawRequest struct {
	Amount        int64  `json:"amount"`
	Destination   string `json:"destination"`         // Wallet address
	TokenMint     string `json:"token_mint,omitempty"` // Optional for SPL tokens
}

// WithdrawResponse contains the withdrawal transaction details.
type WithdrawResponse struct {
	Success       bool   `json:"success"`
	Transaction   string `json:"transaction"` // Unsigned transaction for signing
	WithdrawalID  string `json:"withdrawal_id"`
	Amount        int64  `json:"amount"`
	Fee           int64  `json:"fee,omitempty"`
	NetAmount     int64  `json:"net_amount"`
	Message       string `json:"message,omitempty"`
}

// DecryptRequest represents a request to decrypt an ElGamal-encrypted amount.
type DecryptRequest struct {
	EncryptedAmount string `json:"encrypted_amount"` // 0x hex 64 bytes
	PrivateKey      string `json:"private_key"`      // 0x hex 32 bytes
}

// DecryptResponse contains the decrypted amount.
type DecryptResponse struct {
	Amount int64 `json:"amount"`
}

// GetEarnings retrieves the merchant's total earnings and token breakdown.
func (s *Service) GetEarnings(ctx context.Context) (*EarningsResponse, error) {
	var resp EarningsResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/api/merchant/earnings", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetAnalytics retrieves payment analytics with optional date filtering.
// Supports filtering by date range and grouping by interval (hour, day, week, month).
func (s *Service) GetAnalytics(ctx context.Context, req AnalyticsRequest) (*AnalyticsResponse, error) {
	var resp AnalyticsResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/api/merchant/analytics", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Withdraw initiates a withdrawal of merchant earnings.
// Returns an unsigned transaction that must be signed and submitted by the merchant.
func (s *Service) Withdraw(ctx context.Context, req WithdrawRequest) (*WithdrawResponse, error) {
	var resp WithdrawResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/merchant/withdraw", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DecryptAmount decrypts an ElGamal-encrypted payment amount.
// Requires the merchant's private key. Used to reveal the actual amount from encrypted payments.
func (s *Service) DecryptAmount(ctx context.Context, req DecryptRequest) (*DecryptResponse, error) {
	var resp DecryptResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/merchant/decrypt", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
