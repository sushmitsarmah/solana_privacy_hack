package api

import (
	"encoding/json"
	"net/http"

	"sol_privacy/internal/webhook"
)

// WebhookRegister handles webhook registration
func (h *Handler) WebhookRegister(w http.ResponseWriter, r *http.Request) {
	var req webhook.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Webhook.Register(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// WebhookConfig handles getting webhook configuration
func (h *Handler) WebhookConfig(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.Webhook.GetConfig(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// WebhookTest handles testing a webhook
func (h *Handler) WebhookTest(w http.ResponseWriter, r *http.Request) {
	var req webhook.TestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Webhook.Test(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// WebhookLogs handles getting webhook logs
func (h *Handler) WebhookLogs(w http.ResponseWriter, r *http.Request) {
	var req webhook.LogsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Try query params if body is empty
		req.WebhookID = r.URL.Query().Get("webhook_id")
		req.Event = r.URL.Query().Get("event")
	}

	resp, err := h.client.Webhook.GetLogs(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// WebhookStats handles getting webhook statistics
func (h *Handler) WebhookStats(w http.ResponseWriter, r *http.Request) {
	resp, err := h.client.Webhook.GetStats(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}

// WebhookDeactivate handles deactivating a webhook
func (h *Handler) WebhookDeactivate(w http.ResponseWriter, r *http.Request) {
	var req webhook.DeactivateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, err := h.client.Webhook.Deactivate(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resp)
}
