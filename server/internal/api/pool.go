package api

import (
	"encoding/json"
	"net/http"

	"sol_privacy/internal/pool"
	"sol_privacy/internal/umbra"

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
// Enhanced to optionally deposit to Umbra privacy pool as well
func (h *Handler) PoolDeposit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WalletAddress      string  `json:"wallet_address"`
		Amount             int64   `json:"amount"`
		PrivateKey         string  `json:"private_key,omitempty"`
		DepositToUmbra     bool    `json:"deposit_to_umbra,omitempty"`
		UmbraDestination   string  `json:"umbra_destination,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Deposit to ShadowPay pool
	poolResp, err := h.client.Pool.Deposit(r.Context(), pool.DepositRequest{
		WalletAddress: req.WalletAddress,
		Amount:        req.Amount,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"transaction": poolResp.Transaction,
		"message":     poolResp.Message,
	}

	// Optionally also deposit to Umbra
	if req.DepositToUmbra && h.umbraEnabled && req.PrivateKey != "" {
		// Convert lamports to SOL for Umbra
		amountInSOL := float64(req.Amount) / 1e9

		umbraResp, err := h.umbraClient.Deposit(r.Context(), umbra.DepositRequest{
			PrivateKey:         req.PrivateKey,
			Amount:             amountInSOL,
			DestinationAddress: req.UmbraDestination,
		})
		if err != nil {
			// Don't fail the whole request if Umbra deposit fails
			response["umbra_deposit"] = map[string]interface{}{
				"error": err.Error(),
			}
		} else {
			response["umbra_deposit"] = map[string]interface{}{
				"signature":           umbraResp.Data.Signature,
				"amount":              umbraResp.Data.Amount,
				"destination_address": umbraResp.Data.DestinationAddress,
				"explorer_url":        umbraResp.Data.ExplorerURL,
				"success":             true,
			}
		}
	}

	respondJSON(w, http.StatusOK, response)
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
