package receipt

import (
	"context"
	"fmt"
)

// Service handles receipt operations for transaction verification and history.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new receipt service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// ReceiptBody contains the core receipt data.
type ReceiptBody struct {
	ID            string `json:"id"`
	AmountLamports int64  `json:"amount_lamports"`
	Timestamp     int64  `json:"timestamp"`
	Merchant      string `json:"merchant"`
	Resource      string `json:"resource,omitempty"`
}

// Receipt represents a signed payment receipt.
type Receipt struct {
	Body   ReceiptBody `json:"body"`
	Sig    string      `json:"sig"`    // Ed25519 signature (base58)
	Pubkey string      `json:"pubkey"` // Settler public key (base58)
}

// GetByCommitmentResponse contains the receipt for a specific commitment.
type GetByCommitmentResponse struct {
	Receipt    Receipt `json:"receipt"`
	Commitment string  `json:"commitment"`
	Verified   bool    `json:"verified"`
}

// ListUserReceiptsRequest represents query parameters for listing user receipts.
type ListUserReceiptsRequest struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// ListUserReceiptsResponse contains paginated user receipts.
type ListUserReceiptsResponse struct {
	Receipts   []Receipt `json:"receipts"`
	TotalCount int       `json:"total_count"`
	Limit      int       `json:"limit"`
	Offset     int       `json:"offset"`
}

// TreeMetadata contains receipt Merkle tree information.
type TreeMetadata struct {
	Root        string `json:"root"`
	TreeDepth   int    `json:"tree_depth"`
	LeafCount   int    `json:"leaf_count"`
	LastUpdated string `json:"last_updated"`
}

// GetTreeResponse contains the receipt Merkle tree metadata.
type GetTreeResponse struct {
	WalletAddress string       `json:"wallet_address"`
	TreeMetadata  TreeMetadata `json:"tree_metadata"`
}

// GetByCommitment fetches a receipt by commitment hash.
// Returns the signed receipt with verification status.
func (s *Service) GetByCommitment(ctx context.Context, commitment string) (*GetByCommitmentResponse, error) {
	var resp GetByCommitmentResponse
	path := fmt.Sprintf("/shadowpay/api/receipts/by-commitment?commitment=%s", commitment)
	if err := s.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListUserReceipts retrieves all receipts for a specific user wallet.
// Supports pagination via limit and offset parameters.
func (s *Service) ListUserReceipts(ctx context.Context, walletAddress string, req ListUserReceiptsRequest) (*ListUserReceiptsResponse, error) {
	var resp ListUserReceiptsResponse
	path := fmt.Sprintf("/shadowpay/api/receipts/user/%s", walletAddress)
	if err := s.doRequest(ctx, "GET", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTree retrieves the receipt Merkle tree metadata for a user.
// The tree allows compact verification of receipt authenticity.
func (s *Service) GetTree(ctx context.Context, walletAddress string) (*GetTreeResponse, error) {
	var resp GetTreeResponse
	path := fmt.Sprintf("/shadowpay/api/receipts/tree/%s", walletAddress)
	if err := s.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
