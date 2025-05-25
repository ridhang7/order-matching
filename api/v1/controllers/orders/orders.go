package orders

import (
	"encoding/json"
	"net/http"
	"order-matching/api/v1/database"
	"order-matching/api/v1/models"
	order_matcher "order-matching/api/v1/services"
	"strconv"

	"github.com/gorilla/mux"
)

// OrderRequest represents the request body for creating an order
type OrderRequest struct {
	Type        models.OrderType     `json:"type"`
	Category    models.OrderCategory `json:"category"`
	StockSymbol models.StockSymbol   `json:"stock_symbol"`
	Quantity    uint                 `json:"quantity"`
	Price       float64              `json:"price"`
	UserID      uint                 `json:"user_id"`
}

// OrderResponse represents the response for order-related endpoints
type OrderResponse struct {
	BuyOrders  []models.Order `json:"buy_orders"`
	SellOrders []models.Order `json:"sell_orders"`
}

// CreateOrder handles the creation of a new order
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate stock exists
	if _, err := models.GetStockBySymbol(database.GetDB(), req.StockSymbol); err != nil {
		http.Error(w, "Invalid stock symbol", http.StatusBadRequest)
		return
	}

	// Create order
	order := &models.Order{
		Type:        req.Type,
		Category:    req.Category,
		StockSymbol: req.StockSymbol,
		Quantity:    req.Quantity,
		Price:       req.Price,
		Status:      models.OrderStatusPending,
		UserID:      req.UserID, // TODO: Get from auth context
	}

	// Save order to database
	if err := models.CreateOrder(database.GetDB(), order); err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Process order through matching engine
	matcher := order_matcher.GetOrderMatcher()
	if err := matcher.ProcessOrder(order); err != nil {
		http.Error(w, "Failed to process order", http.StatusInternalServerError)
		return
	}

	// Reload order with stock data
	order, err := models.GetOrderByID(database.GetDB(), order.ID)
	if err != nil {
		http.Error(w, "Failed to load order details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// GetOrdersByStock retrieves all orders for a specific stock
func GetOrdersByStock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	symbol := models.StockSymbol(vars["symbol"])

	// Validate stock exists
	_, err := models.GetStockBySymbol(database.GetDB(), symbol)
	if err != nil {
		http.Error(w, "Invalid stock symbol", http.StatusBadRequest)
		return
	}

	// Get all orders for the stock
	orders, err := models.GetOrdersByStock(database.GetDB(), symbol)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	// Split into buy and sell orders
	var buyOrders, sellOrders []models.Order
	for _, order := range orders {
		if order.Type == models.OrderTypeBuy {
			buyOrders = append(buyOrders, order)
		} else {
			sellOrders = append(sellOrders, order)
		}
	}

	response := OrderResponse{
		BuyOrders:  buyOrders,
		SellOrders: sellOrders,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAllOrders retrieves all orders
func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := models.GetAllOrders(database.GetDB())
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetOrder retrieves a specific order by ID
func GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := models.GetOrderByID(database.GetDB(), uint(id))
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// CancelOrder cancels a specific order
func CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := models.GetOrderByID(database.GetDB(), uint(id))
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Cancel order through matching engine
	matcher := order_matcher.GetOrderMatcher()
	if err := matcher.CancelOrder(order); err != nil {
		http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
