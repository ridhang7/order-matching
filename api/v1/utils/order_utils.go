package utils

import (
	"order-matching/api/v1/models"
)

// GetOrderPrice returns the order price, or current market price for market orders
func GetOrderPrice(order *models.Order) float64 {
	if order.Category == models.OrderCategoryMarket {
		return order.Stock.CurrentPrice
	}
	return order.Price
}

// IsOrderActive returns true if the order is still active (pending or partially filled)
func IsOrderActive(order *models.Order) bool {
	return order.Status == models.OrderStatusPending || order.Status == models.OrderStatusPartiallyFilled
}

// GetRemainingQuantity returns the unfilled quantity of the order
func GetRemainingQuantity(order *models.Order) uint {
	return order.Quantity - order.FilledQuantity
}

// UpdateOrderStatus updates the order status based on filled quantity
func UpdateOrderStatus(order *models.Order) {
	if order.FilledQuantity == 0 {
		order.Status = models.OrderStatusPending
	} else if order.FilledQuantity < order.Quantity {
		order.Status = models.OrderStatusPartiallyFilled
	} else {
		order.Status = models.OrderStatusMatched
	}
}
