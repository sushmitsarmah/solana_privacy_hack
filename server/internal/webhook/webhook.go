package webhook

import (
	"context"
)

// Service handles webhook registration and management for real-time payment notifications.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new webhook service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// RegisterRequest represents a request to register a webhook URL for events.
type RegisterRequest struct {
	URL    string   `json:"url"`               // HTTPS URL to receive webhook notifications
	Events []string `json:"events"`            // Supported: "payment.received", "payment.settled", "payment.failed"
	Secret string   `json:"secret,omitempty"`  // Optional HMAC secret for signature verification
}

// RegisterResponse contains the webhook registration confirmation.
type RegisterResponse struct {
	Success   bool   `json:"success"`
	WebhookID string `json:"webhook_id"`
	URL       string `json:"url"`
	Events    []string `json:"events"`
	CreatedAt string `json:"created_at"`
	Message   string `json:"message,omitempty"`
}

// ConfigResponse contains the merchant's webhook configuration.
type ConfigResponse struct {
	WebhookID string   `json:"webhook_id"`
	URL       string   `json:"url"`
	Events    []string `json:"events"`
	Active    bool     `json:"active"`
	Secret    string   `json:"secret,omitempty"` // Masked in response
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at,omitempty"`
}

// TestRequest represents a request to send a test notification.
type TestRequest struct {
	WebhookID string `json:"webhook_id,omitempty"`
	Event     string `json:"event,omitempty"` // Optional: specify event type
}

// TestResponse contains the test notification result.
type TestResponse struct {
	Success      bool   `json:"success"`
	StatusCode   int    `json:"status_code"`
	ResponseTime int    `json:"response_time_ms"`
	Message      string `json:"message,omitempty"`
	Error        string `json:"error,omitempty"`
}

// LogEntry represents a single webhook delivery attempt.
type LogEntry struct {
	ID           string `json:"id"`
	WebhookID    string `json:"webhook_id"`
	Event        string `json:"event"`
	StatusCode   int    `json:"status_code"`
	ResponseTime int    `json:"response_time_ms"`
	Success      bool   `json:"success"`
	Attempt      int    `json:"attempt"`
	Timestamp    string `json:"timestamp"`
	Error        string `json:"error,omitempty"`
	PayloadID    string `json:"payload_id,omitempty"`
}

// LogsRequest represents query parameters for webhook logs.
type LogsRequest struct {
	WebhookID string `json:"webhook_id,omitempty"`
	Event     string `json:"event,omitempty"`
	Success   *bool  `json:"success,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Offset    int    `json:"offset,omitempty"`
}

// LogsResponse contains paginated webhook delivery history.
type LogsResponse struct {
	Logs       []LogEntry `json:"logs"`
	TotalCount int        `json:"total_count"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}

// StatsResponse contains webhook delivery metrics.
type StatsResponse struct {
	WebhookID        string  `json:"webhook_id"`
	TotalDeliveries  int     `json:"total_deliveries"`
	SuccessfulDeliveries int `json:"successful_deliveries"`
	FailedDeliveries int     `json:"failed_deliveries"`
	SuccessRate      float64 `json:"success_rate"`
	AverageResponseTime int  `json:"average_response_time_ms"`
	LastDelivery     string  `json:"last_delivery,omitempty"`
	LastSuccess      string  `json:"last_success,omitempty"`
	LastFailure      string  `json:"last_failure,omitempty"`
}

// DeactivateRequest represents a request to disable a webhook.
type DeactivateRequest struct {
	WebhookID string `json:"webhook_id"`
}

// DeactivateResponse contains the deactivation confirmation.
type DeactivateResponse struct {
	Success   bool   `json:"success"`
	WebhookID string `json:"webhook_id"`
	Message   string `json:"message,omitempty"`
}

// Register registers a webhook URL to receive payment event notifications.
// Supported events: "payment.received", "payment.settled", "payment.failed"
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	var resp RegisterResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/webhooks/register", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetConfig retrieves the merchant's current webhook configuration.
func (s *Service) GetConfig(ctx context.Context) (*ConfigResponse, error) {
	var resp ConfigResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/api/webhooks/config", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Test sends a test notification to the registered webhook URL.
// Useful for verifying webhook endpoint functionality.
func (s *Service) Test(ctx context.Context, req TestRequest) (*TestResponse, error) {
	var resp TestResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/webhooks/test", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLogs retrieves paginated webhook delivery history.
// Supports filtering by event type, success status, and pagination.
func (s *Service) GetLogs(ctx context.Context, req LogsRequest) (*LogsResponse, error) {
	var resp LogsResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/api/webhooks/logs", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetStats retrieves webhook delivery metrics and success rates.
func (s *Service) GetStats(ctx context.Context) (*StatsResponse, error) {
	var resp StatsResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/api/webhooks/stats", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Deactivate disables a registered webhook.
// The webhook will stop receiving event notifications.
func (s *Service) Deactivate(ctx context.Context, req DeactivateRequest) (*DeactivateResponse, error) {
	var resp DeactivateResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/webhooks/deactivate", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
