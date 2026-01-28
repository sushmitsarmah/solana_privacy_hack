package api

import (
	"encoding/json"
	"net/http"

	"sol_privacy/internal/token"

	"github.com/go-chi/chi/v5"
)

// TokenList handles listing supported tokens
func (h *Handler) TokenList(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.Token.ListSupported(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// TokenAdd handles adding a new token
func (h *Handler) TokenAdd(w http.ResponseWriter, r *http.Request) {
	var req token.AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Token.Add(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// TokenUpdate handles updating a token
func (h *Handler) TokenUpdate(w http.ResponseWriter, r *http.Request) {
	mint := chi.URLParam(r, "mint")
	if mint == "" {
		respondError(w, http.StatusBadRequest, "mint address required")
		return
	}

	var req token.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Token.Update(r.Context(), mint, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// TokenRemove handles removing a token
func (h *Handler) TokenRemove(w http.ResponseWriter, r *http.Request) {
	mint := chi.URLParam(r, "mint")
	if mint == "" {
		respondError(w, http.StatusBadRequest, "mint address required")
		return
	}

	resp, err := h.client.Token.Remove(r.Context(), mint)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}
