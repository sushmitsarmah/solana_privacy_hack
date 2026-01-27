package shadowid

import (
	"context"
	"fmt"
)

// Service handles anonymous identity operations using Merkle tree-based commitments.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new ShadowID service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// AutoRegisterRequest represents a request to register a wallet via signature (production-recommended).
type AutoRegisterRequest struct {
	WalletAddress string `json:"wallet_address"`
	Signature     string `json:"signature"` // Base58 encoded signature
	Message       string `json:"message"`
}

// AutoRegisterResponse contains the registration result.
type AutoRegisterResponse struct {
	Success    bool   `json:"success"`
	Commitment string `json:"commitment"`
	LeafIndex  int    `json:"leaf_index,omitempty"`
	Message    string `json:"message,omitempty"`
}

// RegisterRequest represents a request to register a Poseidon hash commitment in the tree.
type RegisterRequest struct {
	Commitment string `json:"commitment"` // Poseidon hash commitment
}

// RegisterResponse contains the registration confirmation.
type RegisterResponse struct {
	Success   bool   `json:"success"`
	LeafIndex int    `json:"leaf_index"`
	TxHash    string `json:"tx_hash,omitempty"`
	Message   string `json:"message,omitempty"`
}

// ProofRequest represents a request to retrieve a Merkle proof for a commitment.
type ProofRequest struct {
	Commitment string `json:"commitment"`
}

// ProofResponse contains the Merkle proof for the commitment.
type ProofResponse struct {
	Commitment string   `json:"commitment"`
	LeafIndex  int      `json:"leaf_index"`
	Proof      []string `json:"proof"` // Array of sibling hashes
	Root       string   `json:"root"`
}

// RootResponse contains the current Merkle tree root.
type RootResponse struct {
	Root      string `json:"root"`
	TreeDepth int    `json:"tree_depth"`
	LeafCount int    `json:"leaf_count"`
}

// StatusResponse contains the registration status of a commitment.
type StatusResponse struct {
	Commitment string `json:"commitment"`
	Registered bool   `json:"registered"`
	LeafIndex  int    `json:"leaf_index,omitempty"`
}

// AutoRegister registers a wallet via signature (production-recommended method).
// User must sign a message with their wallet to prove ownership.
func (s *Service) AutoRegister(ctx context.Context, req AutoRegisterRequest) (*AutoRegisterResponse, error) {
	var resp AutoRegisterResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/shadowid/auto-register", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Register adds a Poseidon hash commitment to the Merkle tree.
// Returns the leaf index where the commitment was inserted.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	var resp RegisterResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/shadowid/register", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetProof retrieves the Merkle proof for a given commitment.
// The proof allows anonymous verification of membership in the identity set.
func (s *Service) GetProof(ctx context.Context, commitment string) (*ProofResponse, error) {
	var resp ProofResponse
	req := ProofRequest{Commitment: commitment}
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/shadowid/proof", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetRoot fetches the current Merkle tree root.
// The root is used to verify proofs and represents the current state of all registered identities.
func (s *Service) GetRoot(ctx context.Context) (*RootResponse, error) {
	var resp RootResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/api/shadowid/root", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetStatus checks if a commitment is registered in the tree.
func (s *Service) GetStatus(ctx context.Context, commitment string) (*StatusResponse, error) {
	var resp StatusResponse
	path := fmt.Sprintf("/shadowpay/shadowid/v1/id/status/%s", commitment)
	if err := s.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
