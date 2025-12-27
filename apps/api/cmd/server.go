// Package main is the entry point for the API server.
//
// This package contains the main application entry point and configuration
// for the Drop API server, including environment setup, database connections,
// and HTTP server initialization.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/egeuysall/drop/internal/modules/items"
	"github.com/egeuysall/drop/internal/router"
	"github.com/egeuysall/drop/internal/utils"

	supabase "github.com/egeuysall/drop/internal/supabase"
	generated "github.com/egeuysall/drop/internal/supabase/generated"
	"github.com/joho/godotenv"
)

// main is the entry point of the application.
// It initializes the environment, database connection, and starts the HTTP server.
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Initialize database connection
	dbConn := supabase.Connect()
	defer dbConn.Close()

	// Initialize database queries
	queries := generated.New(dbConn)

	// Initialize utility functions with database queries
	utils.Init(queries)

	// Initialize items repository and service
	itemsRepo := items.NewRepository(queries)
	itemsService := items.NewService(itemsRepo)
	itemsHandler := items.NewHandler(itemsService)

	// Initialize handlers
	handlers := router.Handlers{
		// Add other handlers here as needed
		Items: itemsHandler,
	}

	// Initialize HTTP router with all routes and middleware
	router := router.Router(handlers)

	// Get server port from environment variables
	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT environment variable not set")
	}

	// Validate port is a valid number
	if _, err := strconv.Atoi(portStr); err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}

	// Construct server address
	addr := fmt.Sprintf(":%s", portStr)

	// Log server startup information
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("Environment: %s", getEnvironment())
	log.Printf("Database connection established")

	// Start HTTP server
	log.Printf("Listening on %s...", addr)
	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getEnvironment() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "development"
	}
	return env
}
