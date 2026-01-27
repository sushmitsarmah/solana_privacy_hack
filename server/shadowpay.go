package shadowpay

import (
	"context"

	"sol_privacy/internal/authorization"
	"sol_privacy/internal/client"
	"sol_privacy/internal/escrow"
	"sol_privacy/internal/intent"
	"sol_privacy/internal/keys"
	"sol_privacy/internal/merchant"
	"sol_privacy/internal/payment"
	"sol_privacy/internal/pool"
	"sol_privacy/internal/privacy"
	"sol_privacy/internal/receipt"
	"sol_privacy/internal/shadowid"
	"sol_privacy/internal/token"
	"sol_privacy/internal/verify"
	"sol_privacy/internal/webhook"
)

// ShadowPay is the main SDK client for interacting with the ShadowPay API.
type ShadowPay struct {
	client *client.Client

	// Services
	Keys          *keys.Service
	Escrow        *escrow.Service
	Payment       *payment.Service
	Intent        *intent.Service
	Verify        *verify.Service
	Pool          *pool.Service
	ShadowID      *shadowid.Service
	Merchant      *merchant.Service
	Webhook       *webhook.Service
	Privacy       *privacy.Service
	Receipt       *receipt.Service
	Token         *token.Service
	Authorization *authorization.Service
}

// New creates a new ShadowPay SDK client.
// apiKey can be empty if you are only calling public endpoints (like key generation).
func New(apiKey string, opts ...client.Option) *ShadowPay {
	c := client.New(apiKey, opts...)

	// Create a helper function that wraps the client's Do method
	doRequest := func(ctx context.Context, method, path string, body, result interface{}) error {
		req, err := c.NewRequest(ctx, method, path, body)
		if err != nil {
			return err
		}
		return c.Do(req, result)
	}

	return &ShadowPay{
		client:        c,
		Keys:          keys.NewService(doRequest),
		Escrow:        escrow.NewService(doRequest),
		Payment:       payment.NewService(doRequest),
		Intent:        intent.NewService(doRequest),
		Verify:        verify.NewService(doRequest),
		Pool:          pool.NewService(doRequest),
		ShadowID:      shadowid.NewService(doRequest),
		Merchant:      merchant.NewService(doRequest),
		Webhook:       webhook.NewService(doRequest),
		Privacy:       privacy.NewService(doRequest),
		Receipt:       receipt.NewService(doRequest),
		Token:         token.NewService(doRequest),
		Authorization: authorization.NewService(doRequest),
	}
}

// GetAPIKey returns the API key configured for this client.
func (s *ShadowPay) GetAPIKey() string {
	return s.client.GetAPIKey()
}
