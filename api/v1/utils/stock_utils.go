package utils

import (
	"order-matching/api/v1/models"
)

// GetAllStockSymbols returns all available stock symbols
func GetAllStockSymbols() []models.StockSymbol {
	return []models.StockSymbol{
		models.StockNXTECH,
		models.StockQNTUM,
		models.StockCYBEX,
		models.StockSOLRX,
		models.StockFUSON,
		models.StockGENUM,
		models.StockMEDIX,
		models.StockAITHN,
		models.StockNRLNK,
		models.StockCOGNT,
	}
}
