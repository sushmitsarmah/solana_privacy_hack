package api

import (
	"encoding/json"
	"net/http"

	"sol_privacy/internal/merchant"
)

// MerchantEarnings handles getting merchant earnings
func (h *Handler) MerchantEarnings(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.Merchant.GetEarnings(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// MerchantAnalytics handles getting merchant analytics
func (h *Handler) MerchantAnalytics(w http.ResponseWriter, r *http.Request) {
	var req merchant.AnalyticsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Merchant.GetAnalytics(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// MerchantWithdraw handles merchant earnings withdrawal
func (h *Handler) MerchantWithdraw(w http.ResponseWriter, r *http.Request) {
	var req merchant.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Merchant.Withdraw(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}
