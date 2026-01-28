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
func (h *Handler) PaymentPrepare(w http.ResponseWriter, r *http.Request) {
	var req payment.PrepareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Payment.Prepare(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
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
