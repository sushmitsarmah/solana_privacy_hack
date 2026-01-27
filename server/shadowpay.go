package shadowpay

import (
	"context"

	"sol_privacy/internal/client"
	"sol_privacy/internal/escrow"
	"sol_privacy/internal/intent"
	"sol_privacy/internal/keys"
	"sol_privacy/internal/payment"
	"sol_privacy/internal/verify"
)

// ShadowPay is the main SDK client for interacting with the ShadowPay API.
type ShadowPay struct {
	client *client.Client

	// Services
	Keys    *keys.Service
	Escrow  *escrow.Service
	Payment *payment.Service
	Intent  *intent.Service
	Verify  *verify.Service
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
		client:  c,
		Keys:    keys.NewService(doRequest),
		Escrow:  escrow.NewService(doRequest),
		Payment: payment.NewService(doRequest),
		Intent:  intent.NewService(doRequest),
		Verify:  verify.NewService(doRequest),
	}
}

// GetAPIKey returns the API key configured for this client.
func (s *ShadowPay) GetAPIKey() string {
	return s.client.GetAPIKey()
}
