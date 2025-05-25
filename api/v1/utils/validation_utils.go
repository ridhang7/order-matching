package utils

import (
	"errors"
	"order-matching/api/v1/models"
)

var (
	// Order-related errors
	ErrInvalidOrderType     = errors.New("invalid order type")
	ErrInvalidOrderCategory = errors.New("invalid order category")
	ErrInvalidOrderStatus   = errors.New("invalid order status")
	ErrInvalidPrice         = errors.New("price is required and must be greater than 0 for limit order")
	ErrInvalidQuantity      = errors.New("quantity must be greater than 0")

	// Stock-related errors
	ErrInvalidStockSymbol = errors.New("invalid stock symbol")
	ErrInvalidStockName   = errors.New("stock name cannot be empty")
	ErrInvalidStockPrice  = errors.New("price must be greater than 0")
	ErrInvalidPriceRange  = errors.New("day high cannot be less than day low")
	ErrInvalidVolume      = errors.New("volume cannot be negative")
	ErrInvalidMarketCap   = errors.New("market cap must be greater than 0")
	ErrInvalidSector      = errors.New("sector cannot be empty")

	// Trade-related errors
	ErrInvalidTradeOrders = errors.New("both buy and sell order IDs are required")
	ErrSameOrderTrade     = errors.New("buy and sell order IDs cannot be the same")
)

// isValidStockSymbol checks if a given stock symbol is valid
func isValidStockSymbol(symbol models.StockSymbol) bool {
	validSymbols := GetAllStockSymbols()
	for _, s := range validSymbols {
		if s == symbol {
			return true
		}
	}
	return false
}

// ValidateOrder performs validation on the order
func ValidateOrder(order *models.Order) error {
	// Validate order type
	switch order.Type {
	case models.OrderTypeBuy, models.OrderTypeSell:
		// Valid
	default:
		return ErrInvalidOrderType
	}

	// Validate order category
	switch order.Category {
	case models.OrderCategoryLimit, models.OrderCategoryMarket:
		// Valid
	default:
		return ErrInvalidOrderCategory
	}

	// Validate stock symbol
	if !isValidStockSymbol(order.StockSymbol) {
		return ErrInvalidStockSymbol
	}

	// Validate price for limit orders
	if order.Category == models.OrderCategoryLimit && order.Price <= 0 {
		return ErrInvalidPrice
	}

	// Validate quantity
	if order.Quantity <= 0 {
		return ErrInvalidQuantity
	}

	return nil
}

// ValidateStock performs validation on the stock data
func ValidateStock(stock *models.Stock) error {
	if !isValidStockSymbol(stock.Symbol) {
		return ErrInvalidStockSymbol
	}

	if stock.Name == "" {
		return ErrInvalidStockName
	}

	if stock.CurrentPrice <= 0 {
		return ErrInvalidStockPrice
	}

	if stock.DayHigh < stock.DayLow {
		return ErrInvalidPriceRange
	}

	if stock.Volume < 0 {
		return ErrInvalidVolume
	}

	if stock.MarketCap <= 0 {
		return ErrInvalidMarketCap
	}

	if stock.Sector == "" {
		return ErrInvalidSector
	}

	return nil
}

// ValidateTrade performs validation on the trade
func ValidateTrade(trade *models.Trade) error {
	if !isValidStockSymbol(trade.StockSymbol) {
		return ErrInvalidStockSymbol
	}

	if trade.Quantity == 0 {
		return ErrInvalidQuantity
	}

	if trade.Price <= 0 {
		return ErrInvalidStockPrice
	}

	if trade.BuyOrderID == 0 || trade.SellOrderID == 0 {
		return ErrInvalidTradeOrders
	}

	if trade.BuyOrderID == trade.SellOrderID {
		return ErrSameOrderTrade
	}

	return nil
}

// ValidateOrderStatus checks if the order status is valid
func ValidateOrderStatus(status models.OrderStatus) error {
	switch status {
	case models.OrderStatusPending,
		models.OrderStatusPartiallyFilled,
		models.OrderStatusMatched,
		models.OrderStatusCancelled:
		return nil
	default:
		return ErrInvalidOrderStatus
	}
}
