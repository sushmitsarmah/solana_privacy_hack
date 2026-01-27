package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sol_privacy/internal/errors"
)

const (
	DefaultBaseURL = "https://shadow.radr.fun"
	UserAgent      = "shadowpay-go-client/1.0.0"
)

// Client is the main entry point for the ShadowPay API.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiKey     string
	userAgent  string
}

// Option allows for functional configuration of the Client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(rawURL string) Option {
	return func(c *Client) {
		if parsed, err := url.Parse(strings.TrimRight(rawURL, "/")); err == nil {
			c.baseURL = parsed
		}
	}
}

// WithHTTPClient overrides the default http.Client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// New creates a new ShadowPay API client.
// apiKey can be empty if you are only calling public endpoints (like key generation).
func New(apiKey string, opts ...Option) *Client {
	base, _ := url.Parse(DefaultBaseURL)
	c := &Client{
		baseURL:    base,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     apiKey,
		userAgent:  UserAgent,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// NewRequest creates an authenticated HTTP request.
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, fmt.Errorf("failed to encode body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}

	return req, nil
}

// Do executes the HTTP request and decodes the response.
func (c *Client) Do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for API errors
	if resp.StatusCode >= 400 {
		return c.handleError(resp)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) handleError(resp *http.Response) error {
	var apiErr errors.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		// Fallback if JSON decoding fails
		return fmt.Errorf("api error (status %d): %s", resp.StatusCode, resp.Status)
	}
	apiErr.StatusCode = resp.StatusCode
	return &apiErr
}

// GetAPIKey returns the API key configured for this client.
func (c *Client) GetAPIKey() string {
	return c.apiKey
}
