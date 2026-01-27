package privacy

import (
	"context"
)

// Service handles ElGamal encryption operations on the BN254 curve.
type Service struct {
	doRequest func(ctx context.Context, method, path string, body, result interface{}) error
}

// NewService creates a new privacy service.
func NewService(doRequest func(ctx context.Context, method, path string, body, result interface{}) error) *Service {
	return &Service{
		doRequest: doRequest,
	}
}

// KeygenResponse contains a newly generated ElGamal keypair on BN254.
type KeygenResponse struct {
	PublicKey  string `json:"public_key"`  // 0x hex encoded
	PrivateKey string `json:"private_key"` // 0x hex 32 bytes - store securely!
}

// DecryptRequest represents a request to decrypt an ElGamal ciphertext.
type DecryptRequest struct {
	Ciphertext string `json:"ciphertext"`  // 0x hex 64 bytes
	PrivateKey string `json:"private_key"` // 0x hex 32 bytes
}

// DecryptResponse contains the decrypted plaintext amount.
type DecryptResponse struct {
	Amount int64 `json:"amount"`
}

// GenerateKeypair generates a new ElGamal keypair on the BN254 curve.
// The private key should be stored securely by the client for decrypting payments.
func (s *Service) GenerateKeypair(ctx context.Context) (*KeygenResponse, error) {
	var resp KeygenResponse
	if err := s.doRequest(ctx, "GET", "/shadowpay/api/privacy/keygen", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Decrypt decrypts an ElGamal-encrypted ciphertext using the provided private key.
// Used to reveal the actual payment amount from encrypted transactions.
func (s *Service) Decrypt(ctx context.Context, req DecryptRequest) (*DecryptResponse, error) {
	var resp DecryptResponse
	if err := s.doRequest(ctx, "POST", "/shadowpay/api/privacy/decrypt", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
