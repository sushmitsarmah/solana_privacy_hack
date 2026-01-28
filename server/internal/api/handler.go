package api

import (
	"encoding/json"
	"net/http"
	"os"

	shadowpay "sol_privacy"
	"sol_privacy/internal/umbra"

	"github.com/go-chi/chi/v5"
)

// Handler holds the ShadowPay client and provides HTTP endpoints
type Handler struct {
	client      *shadowpay.ShadowPay
	umbraClient *umbra.Client
	umbraEnabled bool
}

// NewHandler creates a new API handler
func NewHandler(apiKey string) *Handler {
	h := &Handler{
		client: shadowpay.New(apiKey),
	}

	// Initialize Umbra client if URL is configured
	if umbraURL := os.Getenv("UMBRA_API_URL"); umbraURL != "" {
		h.umbraClient = umbra.NewClient(umbra.Config{
			BaseURL: umbraURL,
		})
		h.umbraEnabled = true
	}

	return h
}

// Routes returns all API routes
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	// Payment routes
	r.Route("/payment", func(r chi.Router) {
		r.Post("/deposit", h.PaymentDeposit)
		r.Post("/withdraw", h.PaymentWithdraw)
		r.Post("/prepare", h.PaymentPrepare)
		r.Post("/authorize", h.PaymentAuthorize)
		r.Post("/verify-access", h.PaymentVerifyAccess)
		r.Post("/settle", h.PaymentSettle)
	})

	// Pool routes
	r.Route("/pool", func(r chi.Router) {
		r.Get("/balance/{wallet}", h.PoolBalance)
		r.Post("/deposit", h.PoolDeposit)
		r.Post("/withdraw", h.PoolWithdraw)
		r.Get("/deposit-address", h.PoolDepositAddress)
	})

	// Token routes
	r.Route("/token", func(r chi.Router) {
		r.Get("/list", h.TokenList)
		r.Post("/add", h.TokenAdd)
		r.Put("/{mint}", h.TokenUpdate)
		r.Delete("/{mint}", h.TokenRemove)
	})

	// Merchant routes
	r.Route("/merchant", func(r chi.Router) {
		r.Get("/earnings", h.MerchantEarnings)
		r.Post("/analytics", h.MerchantAnalytics)
		r.Post("/withdraw", h.MerchantWithdraw)
	})

	// Privacy routes
	r.Route("/privacy", func(r chi.Router) {
		r.Post("/decrypt", h.PrivacyDecrypt)
	})

	// Webhook routes
	r.Route("/webhook", func(r chi.Router) {
		r.Post("/register", h.WebhookRegister)
		r.Get("/config", h.WebhookConfig)
		r.Post("/test", h.WebhookTest)
		r.Get("/logs", h.WebhookLogs)
		r.Get("/stats", h.WebhookStats)
		r.Post("/deactivate", h.WebhookDeactivate)
	})

	// ShadowID routes
	r.Route("/shadowid", func(r chi.Router) {
		r.Post("/auto-register", h.ShadowIDAutoRegister)
		r.Post("/register", h.ShadowIDRegister)
		r.Post("/proof", h.ShadowIDProof)
		r.Get("/root", h.ShadowIDRoot)
		r.Get("/status/{commitment}", h.ShadowIDStatus)
	})

	// Authorization routes
	r.Route("/authorization", func(r chi.Router) {
		r.Post("/authorize", h.AuthorizationAuthorize)
		r.Get("/list/{wallet}", h.AuthorizationList)
		r.Post("/revoke", h.AuthorizationRevoke)
	})

	// Umbra integration routes (only if Umbra is enabled)
	if h.umbraEnabled {
		r.Route("/umbra", func(r chi.Router) {
			r.Post("/stealth-address", h.UmbraStealthAddress)
			r.Post("/deposit", h.UmbraDeposit)
			r.Post("/send", h.UmbraSend)
			r.Post("/withdraw", h.UmbraWithdraw)
			r.Post("/balance", h.UmbraBalance)
			r.Post("/prepare-stealth-payment", h.UmbraPrepareStealthPayment)
		})
	}

	return r
}

// Helper functions for JSON responses
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
