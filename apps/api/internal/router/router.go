package router

import (
	"time"

	"github.com/egeuysall/drop/internal/handlers"
	appmid "github.com/egeuysall/drop/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func Router() *chi.Mux {
	r := chi.NewRouter()

	// Initialize handlers
	waitlistHandler := handlers.NewWaitlistHandler()

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
	r.Get("/", handlers.HandleRoot)
	r.Get("/health", handlers.HandleHealth)

	// Protected API v1 routes - require authentication
	r.Route("/v1", func(r chi.Router) {
		r.Use(appmid.RequireAuth())

		// Routes
		r.Post("/waitlist/join", waitlistHandler.Join)
	})

	return r
}
