package api

import (
	"encoding/json"
	"net/http"

	"sol_privacy/internal/privacy"
)

// PrivacyDecrypt handles decrypting an amount
func (h *Handler) PrivacyDecrypt(w http.ResponseWriter, r *http.Request) {
	var req privacy.DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Privacy.Decrypt(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}
