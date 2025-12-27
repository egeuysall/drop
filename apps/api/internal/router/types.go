// Package router provides HTTP routing functionality for the API.
//
// This package is responsible for setting up all HTTP routes, middleware,
// and request handling for the Drop API server.
package router

import "net/http"

// Handlers holds all HTTP handlers for the application
type Handlers struct {
	Waitlist      WaitlistHandler
	Auth          AuthHandler
}

// AuthMiddleware holds the auth middleware configuration
type AuthMiddleware struct {
	JWTSecret   string
	SupabaseURL string
}

// WaitlistHandler defines the interface for waitlist HTTP operations
type WaitlistHandler interface {
	Join(w http.ResponseWriter, r *http.Request)
}

// AuthHandler defines the interface for auth HTTP operations
type AuthHandler interface {
	GetMe(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	CompleteOnboarding(w http.ResponseWriter, r *http.Request)
	GetPreferences(w http.ResponseWriter, r *http.Request)
	UpdatePreferences(w http.ResponseWriter, r *http.Request)
}
