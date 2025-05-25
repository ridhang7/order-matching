package order_matcher

import (
	"database/sql"
	"fmt"
	"order-matching/api/v1/models"
	"sync"
)

// OrderMatcher handles the order matching logic
type OrderMatcher struct {
	mu         sync.Mutex
	db         *sql.DB
	BuyOrders  []models.Order // Sorted by price (desc) and time (asc)
	SellOrders []models.Order // Sorted by price (asc) and time (asc)
}

var (
	instance *OrderMatcher
	once     sync.Once
)

// GetOrderMatcher returns the singleton instance of OrderMatcher
func GetOrderMatcher() *OrderMatcher {
	once.Do(func() {
		instance = &OrderMatcher{
			BuyOrders:  make([]models.Order, 0),
			SellOrders: make([]models.Order, 0),
		}
	})
	return instance
}

// SetDB sets the database connection for the order matcher
func (m *OrderMatcher) SetDB(db *sql.DB) {
	m.db = db
}

// ProcessOrder processes a new order and attempts to match it
func (m *OrderMatcher) ProcessOrder(order *models.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Begin transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Process order based on type
	var matchingOrders []models.Order
	if order.Type == models.OrderTypeBuy {
		matchingOrders, err = m.getMatchingSellOrders(tx, order)
	} else {
		matchingOrders, err = m.getMatchingBuyOrders(tx, order)
	}
	if err != nil {
		return fmt.Errorf("failed to get matching orders: %v", err)
	}

	// For market orders with no matches, cancel immediately
	if order.Category == models.OrderCategoryMarket && len(matchingOrders) == 0 {
		order.Status = models.OrderStatusCancelled
		if err := m.updateOrder(tx, order); err != nil {
			return fmt.Errorf("failed to cancel market order: %v", err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %v", err)
		}
		return nil
	}

	// Match orders
	for _, matchingOrder := range matchingOrders {
		if order.FilledQuantity >= order.Quantity {
			break
		}

		// Calculate trade quantity
		remainingQuantity := order.Quantity - order.FilledQuantity
		matchingRemaining := matchingOrder.Quantity - matchingOrder.FilledQuantity
		tradeQuantity := min(remainingQuantity, matchingRemaining)

		// Determine trade price:
		// - For limit/limit matches: use the resting order's price
		// - For market/limit matches: use the limit order's price
		var tradePrice float64
		if order.Category == models.OrderCategoryLimit && matchingOrder.Category == models.OrderCategoryLimit {
			// Both are limit orders, use the resting (matching) order's price
			tradePrice = matchingOrder.Price
		} else {
			// At least one is a market order, use the limit order's price
			if order.Category == models.OrderCategoryLimit {
				tradePrice = order.Price
			} else {
				tradePrice = matchingOrder.Price
			}
		}

		// Create trade
		var buyOrderID, sellOrderID uint
		if order.Type == models.OrderTypeBuy {
			buyOrderID = order.ID
			sellOrderID = matchingOrder.ID
		} else {
			buyOrderID = matchingOrder.ID
			sellOrderID = order.ID
		}

		trade := &models.Trade{
			BuyOrderID:  buyOrderID,
			SellOrderID: sellOrderID,
			StockSymbol: order.StockSymbol,
			Quantity:    tradeQuantity,
			Price:       tradePrice,
		}

		// Create trade record
		_, err = tx.Exec(`
			INSERT INTO trades (buy_order_id, sell_order_id, stock_symbol,
							quantity, price, executed_at)
			VALUES (?, ?, ?, ?, ?, NOW())`,
			trade.BuyOrderID, trade.SellOrderID, trade.StockSymbol,
			trade.Quantity, trade.Price)
		if err != nil {
			return fmt.Errorf("failed to create trade: %v", err)
		}

		// Update matching order
		matchingOrder.FilledQuantity += tradeQuantity
		if matchingOrder.FilledQuantity >= matchingOrder.Quantity {
			matchingOrder.Status = models.OrderStatusMatched
		} else {
			matchingOrder.Status = models.OrderStatusPartiallyFilled
		}

		// Update matching order in database
		if err := m.updateOrder(tx, &matchingOrder); err != nil {
			return fmt.Errorf("failed to update matching order: %v", err)
		}

		// Update current order
		order.FilledQuantity += tradeQuantity
		if order.FilledQuantity >= order.Quantity {
			order.Status = models.OrderStatusMatched
		} else {
			order.Status = models.OrderStatusPartiallyFilled
		}
	}

	// Handle remaining quantity for market orders
	if order.Category == models.OrderCategoryMarket && order.FilledQuantity < order.Quantity {
		order.Status = models.OrderStatusCancelled
	} else if order.FilledQuantity < order.Quantity {
		// Only add limit orders to the order book
		if order.Type == models.OrderTypeBuy {
			m.BuyOrders = append(m.BuyOrders, *order)
		} else {
			m.SellOrders = append(m.SellOrders, *order)
		}
	}

	// Update order in database
	if err := m.updateOrder(tx, order); err != nil {
		return fmt.Errorf("failed to update order: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// CancelOrder cancels a pending order
func (m *OrderMatcher) CancelOrder(order *models.Order) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Begin transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Update order status
	order.Status = models.OrderStatusCancelled

	// Update in database
	if err := m.updateOrder(tx, order); err != nil {
		return fmt.Errorf("failed to update order: %v", err)
	}

	// Remove from appropriate order list
	var orders *[]models.Order
	if order.Type == models.OrderTypeBuy {
		orders = &m.BuyOrders
	} else {
		orders = &m.SellOrders
	}

	for i, o := range *orders {
		if o.ID == order.ID {
			*orders = append((*orders)[:i], (*orders)[i+1:]...)
			break
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Helper functions

func (m *OrderMatcher) getMatchingSellOrders(tx *sql.Tx, buyOrder *models.Order) ([]models.Order, error) {
	var query string
	var args []interface{}

	// Base query with common conditions
	query = `
		SELECT id, type, category, stock_symbol, quantity, filled_quantity,
		       price, status, user_id, created_at, updated_at
		FROM orders
		WHERE type = 'SELL'
		  AND stock_symbol = ?
		  AND status IN ('PENDING', 'PARTIALLY_FILLED')`
	args = append(args, buyOrder.StockSymbol)

	// Add price condition for limit orders
	if buyOrder.Category == models.OrderCategoryLimit {
		query += ` AND price <= ?`
		args = append(args, buyOrder.Price)
	}

	// Add order by clause for price-time priority
	query += ` ORDER BY price ASC, created_at ASC`

	// Execute query
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.Type, &order.Category, &order.StockSymbol,
			&order.Quantity, &order.FilledQuantity, &order.Price,
			&order.Status, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (m *OrderMatcher) getMatchingBuyOrders(tx *sql.Tx, sellOrder *models.Order) ([]models.Order, error) {
	var query string
	var args []interface{}

	// Base query with common conditions
	query = `
		SELECT id, type, category, stock_symbol, quantity, filled_quantity,
		       price, status, user_id, created_at, updated_at
		FROM orders
		WHERE type = 'BUY'
		  AND stock_symbol = ?
		  AND status IN ('PENDING', 'PARTIALLY_FILLED')`
	args = append(args, sellOrder.StockSymbol)

	// Add price condition for limit orders
	if sellOrder.Category == models.OrderCategoryLimit {
		query += ` AND price >= ?`
		args = append(args, sellOrder.Price)
	}

	// Add order by clause for price-time priority
	query += ` ORDER BY price DESC, created_at ASC`

	// Execute query
	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.Type, &order.Category, &order.StockSymbol,
			&order.Quantity, &order.FilledQuantity, &order.Price,
			&order.Status, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (m *OrderMatcher) updateOrder(tx *sql.Tx, order *models.Order) error {
	_, err := tx.Exec(`
		UPDATE orders 
		SET filled_quantity = ?, status = ?, updated_at = NOW()
		WHERE id = ?`,
		order.FilledQuantity, order.Status, order.ID)
	return err
}

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}
