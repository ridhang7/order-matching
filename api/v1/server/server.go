package server

import (
	"fmt"
	"log"
	"net/http"
	"order-matching/api/v1/database"
	"order-matching/api/v1/routes"
	order_matcher "order-matching/api/v1/services"

	"github.com/gorilla/mux"
)

// Initialize sets up the application
func Initialize() (*mux.Router, error) {
	// Initialize database
	if err := database.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// Set DB in order matcher
	matcher := order_matcher.GetOrderMatcher()
	matcher.SetDB(database.GetDB())

	// Initialize router
	router := mux.NewRouter()

	// Setup routes
	routes.SetupRoutes(router)

	return router, nil
}

// Close cleans up resources
func Close() {
	database.Close()
}

// Run starts the HTTP server
func Run(addr string) error {
	router, err := Initialize()
	if err != nil {
		return err
	}

	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, router)
}
