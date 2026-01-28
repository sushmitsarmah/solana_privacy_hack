package api

import (
	"encoding/json"
	"net/http"

	"sol_privacy/internal/shadowid"

	"github.com/go-chi/chi/v5"
)

// ShadowIDAutoRegister handles auto-registration via signature
func (h *Handler) ShadowIDAutoRegister(w http.ResponseWriter, r *http.Request) {
	var req shadowid.AutoRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.ShadowID.AutoRegister(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ShadowIDRegister handles commitment registration
func (h *Handler) ShadowIDRegister(w http.ResponseWriter, r *http.Request) {
	var req shadowid.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.ShadowID.Register(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ShadowIDProof handles getting a Merkle proof
func (h *Handler) ShadowIDProof(w http.ResponseWriter, r *http.Request) {
	var req shadowid.ProofRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.ShadowID.GetProof(r.Context(), req.Commitment)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ShadowIDRoot handles getting the Merkle tree root
func (h *Handler) ShadowIDRoot(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.ShadowID.GetRoot(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// ShadowIDStatus handles checking commitment registration status
func (h *Handler) ShadowIDStatus(w http.ResponseWriter, r *http.Request) {
	commitment := chi.URLParam(r, "commitment")
	if commitment == "" {
		respondError(w, http.StatusBadRequest, "commitment required")
		return
	}

	resp, err := h.client.ShadowID.GetStatus(r.Context(), commitment)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}
