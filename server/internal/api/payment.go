package api

import (
	"encoding/json"
	"net/http"

	"sol_privacy/internal/payment"
)

// PaymentDeposit handles deposit to payment account
func (h *Handler) PaymentDeposit(w http.ResponseWriter, r *http.Request) {
	var req payment.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Payment.Deposit(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// PaymentWithdraw handles withdrawal from payment account
func (h *Handler) PaymentWithdraw(w http.ResponseWriter, r *http.Request) {
	var req payment.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Payment.Withdraw(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// PaymentPrepare handles ZK payment preparation
// Enhanced to optionally generate stealth addresses via Umbra
func (h *Handler) PaymentPrepare(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ReceiverCommitment string `json:"receiver_commitment"`
		Amount             int64  `json:"amount"`
		TokenMint          string `json:"token_mint,omitempty"`
		GenerateStealth    bool   `json:"generate_stealth,omitempty"`
		RecipientPublicKey string `json:"recipient_public_key,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	receiverCommitment := req.ReceiverCommitment

	// If stealth address generation is requested and Umbra is enabled
	if req.GenerateStealth && h.umbraEnabled && req.RecipientPublicKey != "" {
		stealthResp, err := h.umbraClient.GenerateStealthAddress(r.Context(), req.RecipientPublicKey)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to generate stealth address: "+err.Error())
			return
		}
		receiverCommitment = stealthResp.Data.EphemeralPublicKey
	}

	// Prepare payment with ShadowPay
	prepareResp, err := h.client.Payment.Prepare(r.Context(), payment.PrepareRequest{
		ReceiverCommitment: receiverCommitment,
		Amount:             req.Amount,
		TokenMint:          req.TokenMint,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// If stealth was generated, include it in the response
	if req.GenerateStealth && h.umbraEnabled && req.RecipientPublicKey != "" {
		stealthResp, _ := h.umbraClient.GenerateStealthAddress(r.Context(), req.RecipientPublicKey)
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"payment_hash": prepareResp.PaymentHash,
			"commitment":   prepareResp.Commitment,
			"message":      prepareResp.Message,
			"transaction":  prepareResp.Transaction,
			"stealth_address": map[string]string{
				"ephemeral_public_key":  stealthResp.Data.EphemeralPublicKey,
				"ephemeral_private_key": stealthResp.Data.EphemeralPrivateKey,
				"recipient_public_key":  stealthResp.Data.RecipientPublicKey,
			},
		})
		return
	}

	respondJSON(w, http.StatusOK, prepareResp)
}

// PaymentAuthorize handles payment authorization
func (h *Handler) PaymentAuthorize(w http.ResponseWriter, r *http.Request) {
	var req payment.AuthorizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Payment.Authorize(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// PaymentVerifyAccess handles access token verification
func (h *Handler) PaymentVerifyAccess(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Payment.VerifyAccess(r.Context(), req.AccessToken)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// PaymentSettle handles payment settlement
func (h *Handler) PaymentSettle(w http.ResponseWriter, r *http.Request) {
	var req payment.SettleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Payment.Settle(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}
