package api

import (
	"encoding/json"
	"net/http"

	"sol_privacy/internal/authorization"

	"github.com/go-chi/chi/v5"
)

// AuthorizationAuthorize handles bot spending authorization
func (h *Handler) AuthorizationAuthorize(w http.ResponseWriter, r *http.Request) {
	var req authorization.AuthorizeSpendingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Authorization.AuthorizeSpending(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// AuthorizationList handles listing authorizations
func (h *Handler) AuthorizationList(w http.ResponseWriter, r *http.Request) {
	wallet := chi.URLParam(r, "wallet")
	if wallet == "" {
		respondError(w, http.StatusBadRequest, "wallet address required")
		return
	}

	resp, err := h.client.Authorization.ListAuthorizations(r.Context(), wallet)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// AuthorizationRevoke handles revoking an authorization
func (h *Handler) AuthorizationRevoke(w http.ResponseWriter, r *http.Request) {
	var req authorization.RevokeAuthorizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Authorization.RevokeAuthorization(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}
