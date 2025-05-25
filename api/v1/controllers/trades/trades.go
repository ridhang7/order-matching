package trades

import (
	"encoding/json"
	"net/http"
	"order-matching/api/v1/database"
	"order-matching/api/v1/models"
	"strconv"

	"github.com/gorilla/mux"
)

// GetAllTrades retrieves all trades
func GetAllTrades(w http.ResponseWriter, r *http.Request) {
	// Get trades from database
	trades, err := models.GetAllTrades(database.GetDB())
	if err != nil {
		http.Error(w, "Failed to fetch trades", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trades)
}

// GetTradeByID retrieves a specific trade by ID
func GetTradeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid trade ID", http.StatusBadRequest)
		return
	}

	trade, err := models.GetTradeByID(database.GetDB(), uint(id))
	if err != nil {
		http.Error(w, "Trade not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trade)
}
