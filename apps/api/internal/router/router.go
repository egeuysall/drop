package router

import (
	"encoding/json"
	"net/http"
	"time"

	appmid "github.com/egeuysall/drop/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware configuration
	// These middleware are applied to all incoming requests
	r.Use(
		middleware.Recoverer,
		middleware.RealIP,
		middleware.Timeout(3*time.Second),
		middleware.NoCache,
		middleware.Compress(5),
		httprate.LimitByIP(30, time.Minute),
		appmid.SetContentType(),
		appmid.Cors(),
	)

	// Public routes - accessible without authentication
	r.Get("/", handleRoot)
	r.Get("/health", handleHealth)

	// Protected API v1 routes - require authentication
	r.Route("/v1", func(r chi.Router) {
		r.Use(appmid.RequireAuth())

		// Routes
		// TODO: Implement routes
	})

	return r
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Drop API"}
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	statusResponse := map[string]string{"status": "Healthy"}
	json.NewEncoder(w).Encode(statusResponse)
}
