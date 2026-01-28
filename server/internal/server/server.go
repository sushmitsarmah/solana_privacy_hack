package server

import (
	"log"
	"net/http"
	"time"

	"sol_privacy/internal/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Config holds server configuration
type Config struct {
	APIKey string
	Port   string
}

// Run starts the HTTP server
func Run(cfg Config) error {
	if cfg.APIKey == "" {
		log.Fatal("API key is required")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.1:*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize API handlers
	apiHandler := api.NewHandler(cfg.APIKey)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"shadowpay-api"}`))
	})

	// Mount API routes
	r.Mount("/api", apiHandler.Routes())

	// Start server
	log.Printf("🚀 ShadowPay API Server starting on port %s", cfg.Port)
	log.Printf("📊 Health check: http://localhost:%s/health", cfg.Port)
	log.Printf("🔌 API endpoint: http://localhost:%s/api", cfg.Port)
	log.Printf("📖 Example: curl http://localhost:%s/api/pool/balance/<wallet>", cfg.Port)

	return http.ListenAndServe(":"+cfg.Port, r)
}
