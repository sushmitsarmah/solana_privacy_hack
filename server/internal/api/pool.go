package api

import (
	"encoding/json"
	"net/http"

	"sol_privacy/internal/pool"

	"github.com/go-chi/chi/v5"
)

// PoolBalance handles pool balance check
func (h *Handler) PoolBalance(w http.ResponseWriter, r *http.Request) {
	wallet := chi.URLParam(r, "wallet")
	if wallet == "" {
		respondError(w, http.StatusBadRequest, "wallet address required")
		return
	}

	resp, err := h.client.Pool.GetBalance(r.Context(), wallet)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// PoolDeposit handles pool deposit
func (h *Handler) PoolDeposit(w http.ResponseWriter, r *http.Request) {
	var req pool.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Pool.Deposit(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// PoolWithdraw handles pool withdrawal
func (h *Handler) PoolWithdraw(w http.ResponseWriter, r *http.Request) {
	var req pool.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Pool.Withdraw(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// PoolDepositAddress handles getting pool deposit address
func (h *Handler) PoolDepositAddress(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.Pool.GetDepositAddress(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}
