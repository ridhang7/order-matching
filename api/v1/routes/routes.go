package routes

import (
	"order-matching/api/v1/controllers/orders"
	"order-matching/api/v1/controllers/trades"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *mux.Router) {
	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Orders routes
	api.HandleFunc("/orders", orders.CreateOrder).Methods("POST")
	api.HandleFunc("/orders", orders.GetAllOrders).Methods("GET")
	api.HandleFunc("/orders/{id:[0-9]+}", orders.GetOrder).Methods("GET")
	api.HandleFunc("/orders/{id:[0-9]+}/cancel", orders.CancelOrder).Methods("POST")
	api.HandleFunc("/orders/stock/{symbol}", orders.GetOrdersByStock).Methods("GET")

	// Trades routes
	api.HandleFunc("/trades", trades.GetAllTrades).Methods("GET")
	api.HandleFunc("/trades/{id:[0-9]+}", trades.GetTradeByID).Methods("GET")
}
