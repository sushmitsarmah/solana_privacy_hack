// Package umbra provides a client for the Umbra privacy server.
package umbra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Config holds configuration for the Umbra client.
type Config struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Client is a client for the Umbra server API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Umbra client.
func NewClient(config Config) *Client {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	return &Client{
		baseURL:    config.BaseURL,
		httpClient: config.HTTPClient,
	}
}

// StealthAddressRequest represents a request to generate a stealth address.
type StealthAddressRequest struct {
	RecipientPublicKey string `json:"recipientPublicKey"`
}

// StealthAddressResponse represents a stealth address generation response.
type StealthAddressResponse struct {
	Success bool `json:"success"`
	Data    struct {
		EphemeralPublicKey  string `json:"ephemeralPublicKey"`
		EphemeralPrivateKey string `json:"ephemeralPrivateKey"`
		RecipientPublicKey  string `json:"recipientPublicKey"`
	} `json:"data"`
	Message string `json:"message"`
}

// GenerateStealthAddress generates a stealth address for a recipient.
func (c *Client) GenerateStealthAddress(ctx context.Context, recipientPublicKey string) (*StealthAddressResponse, error) {
	req := StealthAddressRequest{
		RecipientPublicKey: recipientPublicKey,
	}

	var resp StealthAddressResponse
	if err := c.post(ctx, "/api/umbra/stealth-address", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to generate stealth address: %w", err)
	}

	return &resp, nil
}

// DepositRequest represents a deposit into the Umbra privacy pool.
type DepositRequest struct {
	PrivateKey         string  `json:"privateKey"`
	Amount             float64 `json:"amount"`
	DestinationAddress string  `json:"destinationAddress,omitempty"`
}

// DepositResponse represents a deposit response.
type DepositResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Signature         string `json:"signature"`
		Amount            float64 `json:"amount"`
		AmountLamports    int64  `json:"amountLamports"`
		DestinationAddress string `json:"destinationAddress"`
		PublicKey         string `json:"publicKey"`
		ExplorerURL       string `json:"explorerUrl"`
	} `json:"data"`
	Message string `json:"message"`
}

// Deposit deposits SOL into the Umbra privacy pool.
func (c *Client) Deposit(ctx context.Context, req DepositRequest) (*DepositResponse, error) {
	var resp DepositResponse
	if err := c.post(ctx, "/api/umbra/deposit", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to deposit to Umbra pool: %w", err)
	}

	return &resp, nil
}

// SendRequest represents an anonymous transfer request.
type SendRequest struct {
	PrivateKey       string  `json:"privateKey"`
	RecipientAddress string  `json:"recipientAddress"`
	Amount           float64 `json:"amount"`
	Mint             string  `json:"mint,omitempty"`
}

// SendResponse represents an anonymous transfer response.
type SendResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Signature       string `json:"signature"`
		Amount          float64 `json:"amount"`
		AmountLamports  int64  `json:"amountLamports"`
		RecipientAddress string `json:"recipientAddress"`
		SenderPublicKey string `json:"senderPublicKey"`
		TokenMint       string `json:"tokenMint"`
		ExplorerURL     string `json:"explorerUrl"`
	} `json:"data"`
	Message string `json:"message"`
}

// Send performs an anonymous/confidential transfer.
func (c *Client) Send(ctx context.Context, req SendRequest) (*SendResponse, error) {
	var resp SendResponse
	if err := c.post(ctx, "/api/umbra/send", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to send anonymous transfer: %w", err)
	}

	return &resp, nil
}

// WithdrawRequest represents a withdrawal from the privacy pool.
type WithdrawRequest struct {
	PrivateKey          string `json:"privateKey"`
	CommitmentIndex     int64  `json:"commitmentIndex"`
	GenerationIndex     int64  `json:"generationIndex"`
	DepositTime         int64  `json:"depositTime"`
	RelayerPublicKey    string `json:"relayerPublicKey,omitempty"`
	Mint                string `json:"mint,omitempty"`
}

// WithdrawResponse represents a withdrawal response.
type WithdrawResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Signature         string `json:"signature"`
		DestinationAddress string `json:"destinationAddress"`
		TokenMint         string `json:"tokenMint"`
		ClaimArtifacts    map[string]interface{} `json:"claimArtifacts"`
		ExplorerURL       string `json:"explorerUrl"`
	} `json:"data"`
	Message string `json:"message"`
}

// Withdraw withdraws funds from the Umbra privacy pool.
func (c *Client) Withdraw(ctx context.Context, req WithdrawRequest) (*WithdrawResponse, error) {
	var resp WithdrawResponse
	if err := c.post(ctx, "/api/umbra/withdraw", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to withdraw from Umbra pool: %w", err)
	}

	return &resp, nil
}

// BalanceRequest represents a balance check request.
type BalanceRequest struct {
	PrivateKey string `json:"privateKey"`
	Mint       string `json:"mint,omitempty"`
}

// BalanceResponse represents a balance check response.
type BalanceResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Balance    string `json:"balance"`
		BalanceSOL float64 `json:"balanceSOL"`
		Mint       string `json:"mint"`
		PublicKey  string `json:"publicKey"`
	} `json:"data"`
	Message string `json:"message"`
}

// GetBalance retrieves the encrypted balance for a token.
func (c *Client) GetBalance(ctx context.Context, req BalanceRequest) (*BalanceResponse, error) {
	var resp BalanceResponse
	if err := c.post(ctx, "/api/umbra/balance", req, &resp); err != nil {
		return nil, fmt.Errorf("failed to get Umbra balance: %w", err)
	}

	return &resp, nil
}

// helper method to make POST requests
func (c *Client) post(ctx context.Context, path string, body interface{}, result interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errorResp struct {
			Error   string `json:"error"`
			Details string `json:"details"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return fmt.Errorf("umbra API error: %s - %s", errorResp.Error, errorResp.Details)
		}
		return fmt.Errorf("umbra API error: status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
