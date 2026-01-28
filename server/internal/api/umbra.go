package api

import (
	"encoding/json"
	"net/http"

	"sol_privacy/internal/payment"
	"sol_privacy/internal/umbra"
)

// UmbraStealthAddress generates a stealth address for anonymous payments.
func (h *Handler) UmbraStealthAddress(w http.ResponseWriter, r *http.Request) {
	if !h.umbraEnabled {
		respondError(w, http.StatusNotImplemented, "Umbra integration is not enabled")
		return
	}

	var req struct {
		RecipientPublicKey string `json:"recipient_public_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RecipientPublicKey == "" {
		respondError(w, http.StatusBadRequest, "Missing required field: recipient_public_key")
		return
	}

	// Forward request to Umbra server
	stealthReq := umbra.StealthAddressRequest{
		RecipientPublicKey: req.RecipientPublicKey,
	}

	resp, err := h.umbraClient.GenerateStealthAddress(r.Context(), stealthReq.RecipientPublicKey)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"ephemeral_public_key":  resp.Data.EphemeralPublicKey,
			"ephemeral_private_key": resp.Data.EphemeralPrivateKey,
			"recipient_public_key":  resp.Data.RecipientPublicKey,
		},
		"message": resp.Message,
	})
}

// UmbraDeposit deposits SOL into the Umbra privacy pool.
func (h *Handler) UmbraDeposit(w http.ResponseWriter, r *http.Request) {
	if !h.umbraEnabled {
		respondError(w, http.StatusNotImplemented, "Umbra integration is not enabled")
		return
	}

	var req struct {
		PrivateKey         string  `json:"private_key"`
		Amount             float64 `json:"amount"`
		DestinationAddress string  `json:"destination_address,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PrivateKey == "" || req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "Missing required fields: private_key, amount")
		return
	}

	// Forward request to Umbra server
	depositReq := umbra.DepositRequest{
		PrivateKey:         req.PrivateKey,
		Amount:             req.Amount,
		DestinationAddress: req.DestinationAddress,
	}

	resp, err := h.umbraClient.Deposit(r.Context(), depositReq)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"signature":           resp.Data.Signature,
			"amount":              resp.Data.Amount,
			"amount_lamports":     resp.Data.AmountLamports,
			"destination_address": resp.Data.DestinationAddress,
			"public_key":          resp.Data.PublicKey,
			"explorer_url":        resp.Data.ExplorerURL,
		},
		"message": resp.Message,
	})
}

// UmbraSend performs an anonymous/confidential transfer.
func (h *Handler) UmbraSend(w http.ResponseWriter, r *http.Request) {
	if !h.umbraEnabled {
		respondError(w, http.StatusNotImplemented, "Umbra integration is not enabled")
		return
	}

	var req struct {
		PrivateKey       string  `json:"private_key"`
		RecipientAddress string  `json:"recipient_address"`
		Amount           float64 `json:"amount"`
		Mint             string  `json:"mint,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PrivateKey == "" || req.RecipientAddress == "" || req.Amount <= 0 {
		respondError(w, http.StatusBadRequest, "Missing required fields: private_key, recipient_address, amount")
		return
	}

	// Forward request to Umbra server
	sendReq := umbra.SendRequest{
		PrivateKey:       req.PrivateKey,
		RecipientAddress: req.RecipientAddress,
		Amount:           req.Amount,
		Mint:             req.Mint,
	}

	resp, err := h.umbraClient.Send(r.Context(), sendReq)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"signature":        resp.Data.Signature,
			"amount":           resp.Data.Amount,
			"amount_lamports":  resp.Data.AmountLamports,
			"recipient_address": resp.Data.RecipientAddress,
			"sender_public_key": resp.Data.SenderPublicKey,
			"token_mint":       resp.Data.TokenMint,
			"explorer_url":     resp.Data.ExplorerURL,
		},
		"message": resp.Message,
	})
}

// UmbraWithdraw withdraws funds from the Umbra privacy pool.
func (h *Handler) UmbraWithdraw(w http.ResponseWriter, r *http.Request) {
	if !h.umbraEnabled {
		respondError(w, http.StatusNotImplemented, "Umbra integration is not enabled")
		return
	}

	var req struct {
		PrivateKey          string `json:"private_key"`
		CommitmentIndex     int64  `json:"commitment_index"`
		GenerationIndex     int64  `json:"generation_index"`
		DepositTime         int64  `json:"deposit_time"`
		RelayerPublicKey    string `json:"relayer_public_key,omitempty"`
		Mint                string `json:"mint,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PrivateKey == "" || req.CommitmentIndex < 0 || req.GenerationIndex < 0 || req.DepositTime <= 0 {
		respondError(w, http.StatusBadRequest, "Missing required fields: private_key, commitment_index, generation_index, deposit_time")
		return
	}

	// Forward request to Umbra server
	withdrawReq := umbra.WithdrawRequest{
		PrivateKey:       req.PrivateKey,
		CommitmentIndex:  req.CommitmentIndex,
		GenerationIndex:  req.GenerationIndex,
		DepositTime:      req.DepositTime,
		RelayerPublicKey: req.RelayerPublicKey,
		Mint:             req.Mint,
	}

	resp, err := h.umbraClient.Withdraw(r.Context(), withdrawReq)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"signature":          resp.Data.Signature,
			"destination_address": resp.Data.DestinationAddress,
			"token_mint":         resp.Data.TokenMint,
			"claim_artifacts":    resp.Data.ClaimArtifacts,
			"explorer_url":       resp.Data.ExplorerURL,
		},
		"message": resp.Message,
		"privacy": map[string]interface{}{
			"note": "Withdrawal is completely anonymous. Link between deposit and withdrawal is cryptographically unprovable.",
			"zero_knowledge": "Uses ZK-SNARK proofs to prove deposit ownership without revealing identity",
		},
	})
}

// UmbraBalance retrieves the encrypted balance for a token.
func (h *Handler) UmbraBalance(w http.ResponseWriter, r *http.Request) {
	if !h.umbraEnabled {
		respondError(w, http.StatusNotImplemented, "Umbra integration is not enabled")
		return
	}

	var req struct {
		PrivateKey string `json:"private_key"`
		Mint       string `json:"mint,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PrivateKey == "" {
		respondError(w, http.StatusBadRequest, "Missing required field: private_key")
		return
	}

	// Forward request to Umbra server
	balanceReq := umbra.BalanceRequest{
		PrivateKey: req.PrivateKey,
		Mint:       req.Mint,
	}

	resp, err := h.umbraClient.GetBalance(r.Context(), balanceReq)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"balance":     resp.Data.Balance,
			"balance_sol": resp.Data.BalanceSOL,
			"mint":        resp.Data.Mint,
			"public_key":  resp.Data.PublicKey,
		},
		"message": resp.Message,
	})
}

// UmbraPrepareStealthPayment is the key integration point:
// Generates a stealth address and prepares a ZK payment in one call.
func (h *Handler) UmbraPrepareStealthPayment(w http.ResponseWriter, r *http.Request) {
	if !h.umbraEnabled {
		respondError(w, http.StatusNotImplemented, "Umbra integration is not enabled")
		return
	}

	var req struct {
		RecipientPublicKey string  `json:"recipient_public_key"`
		Amount             float64 `json:"amount"`
		TokenMint          string  `json:"token_mint,omitempty"`
		PrivateKey         string  `json:"private_key"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RecipientPublicKey == "" || req.Amount <= 0 || req.PrivateKey == "" {
		respondError(w, http.StatusBadRequest, "Missing required fields: recipient_public_key, amount, private_key")
		return
	}

	// Step 1: Generate stealth address via Umbra
	stealthResp, err := h.umbraClient.GenerateStealthAddress(r.Context(), req.RecipientPublicKey)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate stealth address: "+err.Error())
		return
	}

	// Step 2: Deposit to Umbra pool using the stealth address as destination
	depositResp, err := h.umbraClient.Deposit(r.Context(), umbra.DepositRequest{
		PrivateKey:         req.PrivateKey,
		Amount:             req.Amount,
		DestinationAddress: stealthResp.Data.EphemeralPublicKey,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to deposit to Umbra pool: "+err.Error())
		return
	}

	// Step 3: Prepare ShadowPay payment with the stealth address as commitment
	prepareResp, err := h.client.Payment.Prepare(r.Context(), payment.PrepareRequest{
		ReceiverCommitment: stealthResp.Data.EphemeralPublicKey,
		Amount:             int64(req.Amount * 1e9), // Convert to lamports
		TokenMint:          req.TokenMint,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to prepare payment: "+err.Error())
		return
	}

	// Return combined response
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"stealth_address": map[string]interface{}{
			"ephemeral_public_key":  stealthResp.Data.EphemeralPublicKey,
			"ephemeral_private_key": stealthResp.Data.EphemeralPrivateKey,
			"recipient_public_key":  stealthResp.Data.RecipientPublicKey,
		},
		"deposit": map[string]interface{}{
			"signature":           depositResp.Data.Signature,
			"amount":              depositResp.Data.Amount,
			"amount_lamports":     depositResp.Data.AmountLamports,
			"destination_address": depositResp.Data.DestinationAddress,
			"explorer_url":        depositResp.Data.ExplorerURL,
		},
		"payment": map[string]interface{}{
			"payment_hash": prepareResp.PaymentHash,
			"commitment":   prepareResp.Commitment,
			"message":      prepareResp.Message,
		},
		"message": "Stealth payment prepared successfully",
	})
}
