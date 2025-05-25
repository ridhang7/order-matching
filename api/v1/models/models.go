package models

import (
	"database/sql"
	"time"
)

// Stock represents a stock in the system
type Stock struct {
	Symbol       StockSymbol
	Name         string
	Description  string
	CurrentPrice float64
	DayHigh      float64
	DayLow       float64
	Volume       int64
	MarketCap    float64
	Sector       string
	LastUpdated  time.Time
}

// Order represents a trading order
type Order struct {
	ID             uint
	Type           OrderType
	Category       OrderCategory
	StockSymbol    StockSymbol
	Quantity       uint
	FilledQuantity uint
	Price          float64
	Status         OrderStatus
	UserID         uint
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Stock          *Stock
}

// Trade represents a matched trade between two orders
type Trade struct {
	ID          uint
	BuyOrderID  uint
	SellOrderID uint
	StockSymbol StockSymbol
	Quantity    uint
	Price       float64
	ExecutedAt  time.Time
	BuyOrder    *Order
	SellOrder   *Order
	Stock       *Stock
}

// GetStockBySymbol retrieves a stock by its symbol
func GetStockBySymbol(db *sql.DB, symbol StockSymbol) (*Stock, error) {
	stock := &Stock{}
	err := db.QueryRow(`
		SELECT symbol, name, description, current_price, day_high, day_low, 
		       volume, market_cap, sector, last_updated 
		FROM stocks 
		WHERE symbol = ?`, symbol).Scan(
		&stock.Symbol, &stock.Name, &stock.Description,
		&stock.CurrentPrice, &stock.DayHigh, &stock.DayLow,
		&stock.Volume, &stock.MarketCap, &stock.Sector, &stock.LastUpdated)
	if err != nil {
		return nil, err
	}
	return stock, nil
}

// GetOrderByID retrieves an order by its ID
func GetOrderByID(db *sql.DB, id uint) (*Order, error) {
	order := &Order{}
	err := db.QueryRow(`
		SELECT o.id, o.type, o.category, o.stock_symbol, o.quantity, 
		       o.filled_quantity, o.price, o.status, o.user_id, 
		       o.created_at, o.updated_at
		FROM orders o
		WHERE o.id = ?`, id).Scan(
		&order.ID, &order.Type, &order.Category, &order.StockSymbol,
		&order.Quantity, &order.FilledQuantity, &order.Price,
		&order.Status, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Load associated stock
	stock, err := GetStockBySymbol(db, order.StockSymbol)
	if err != nil {
		return nil, err
	}
	order.Stock = stock

	return order, nil
}

// CreateOrder creates a new order in the database
func CreateOrder(db *sql.DB, order *Order) error {
	result, err := db.Exec(`
		INSERT INTO orders (type, category, stock_symbol, quantity, 
		                   filled_quantity, price, status, user_id, 
		                   created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		order.Type, order.Category, order.StockSymbol,
		order.Quantity, order.FilledQuantity, order.Price,
		order.Status, order.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	order.ID = uint(id)
	return nil
}

// UpdateOrder updates an existing order in the database
func UpdateOrder(db *sql.DB, order *Order) error {
	_, err := db.Exec(`
		UPDATE orders 
		SET filled_quantity = ?, status = ?, updated_at = NOW()
		WHERE id = ?`,
		order.FilledQuantity, order.Status, order.ID)
	return err
}

// CreateTrade creates a new trade in the database
func CreateTrade(db *sql.DB, trade *Trade) error {
	_, err := db.Exec(`
		INSERT INTO trades (buy_order_id, sell_order_id, stock_symbol,
		                   quantity, price, executed_at)
		VALUES (?, ?, ?, ?, ?, NOW())`,
		trade.BuyOrderID, trade.SellOrderID, trade.StockSymbol,
		trade.Quantity, trade.Price)
	return err
}

// GetTradeByID retrieves a trade by its ID with all associated data
func GetTradeByID(db *sql.DB, id uint) (*Trade, error) {
	trade := &Trade{}
	err := db.QueryRow(`
		SELECT id, buy_order_id, sell_order_id, stock_symbol,
		       quantity, price, executed_at
		FROM trades
		WHERE id = ?`, id).Scan(
		&trade.ID, &trade.BuyOrderID, &trade.SellOrderID,
		&trade.StockSymbol, &trade.Quantity, &trade.Price,
		&trade.ExecutedAt)
	if err != nil {
		return nil, err
	}

	// Load associated orders and stock
	buyOrder, err := GetOrderByID(db, trade.BuyOrderID)
	if err != nil {
		return nil, err
	}
	trade.BuyOrder = buyOrder

	sellOrder, err := GetOrderByID(db, trade.SellOrderID)
	if err != nil {
		return nil, err
	}
	trade.SellOrder = sellOrder

	stock, err := GetStockBySymbol(db, trade.StockSymbol)
	if err != nil {
		return nil, err
	}
	trade.Stock = stock

	return trade, nil
}

// GetOrdersByUserID retrieves all orders for a specific user
func GetOrdersByUserID(db *sql.DB, userID uint) ([]Order, error) {
	rows, err := db.Query(`
		SELECT o.id, o.type, o.category, o.stock_symbol, o.quantity, 
		       o.filled_quantity, o.price, o.status, o.user_id, 
		       o.created_at, o.updated_at
		FROM orders o
		WHERE o.user_id = ?
		ORDER BY o.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.Type, &order.Category, &order.StockSymbol,
			&order.Quantity, &order.FilledQuantity, &order.Price,
			&order.Status, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Load associated stock
		stock, err := GetStockBySymbol(db, order.StockSymbol)
		if err != nil {
			return nil, err
		}
		order.Stock = stock

		orders = append(orders, order)
	}
	return orders, nil
}

// GetOrdersByStock retrieves all orders for a specific stock
func GetOrdersByStock(db *sql.DB, symbol StockSymbol) ([]Order, error) {
	rows, err := db.Query(`
		SELECT o.id, o.type, o.category, o.stock_symbol, o.quantity, 
		       o.filled_quantity, o.price, o.status, o.user_id, 
		       o.created_at, o.updated_at
		FROM orders o
		WHERE o.stock_symbol = ?
		ORDER BY o.created_at DESC`, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.Type, &order.Category, &order.StockSymbol,
			&order.Quantity, &order.FilledQuantity, &order.Price,
			&order.Status, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Load associated stock
		stock, err := GetStockBySymbol(db, order.StockSymbol)
		if err != nil {
			return nil, err
		}
		order.Stock = stock

		orders = append(orders, order)
	}
	return orders, nil
}

// GetAllOrders retrieves all orders from the database
func GetAllOrders(db *sql.DB) ([]Order, error) {
	rows, err := db.Query(`
		SELECT o.id, o.type, o.category, o.stock_symbol, o.quantity, 
		       o.filled_quantity, o.price, o.status, o.user_id, 
		       o.created_at, o.updated_at
		FROM orders o
		ORDER BY o.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.Type, &order.Category, &order.StockSymbol,
			&order.Quantity, &order.FilledQuantity, &order.Price,
			&order.Status, &order.UserID, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Load associated stock
		stock, err := GetStockBySymbol(db, order.StockSymbol)
		if err != nil {
			return nil, err
		}
		order.Stock = stock

		orders = append(orders, order)
	}
	return orders, nil
}

// GetAllTrades retrieves all trades from the database
func GetAllTrades(db *sql.DB) ([]Trade, error) {
	rows, err := db.Query(`
		SELECT id, buy_order_id, sell_order_id, stock_symbol,
		       quantity, price, executed_at
		FROM trades
		ORDER BY executed_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []Trade
	for rows.Next() {
		var trade Trade
		err := rows.Scan(
			&trade.ID, &trade.BuyOrderID, &trade.SellOrderID,
			&trade.StockSymbol, &trade.Quantity, &trade.Price,
			&trade.ExecutedAt)
		if err != nil {
			return nil, err
		}

		// Load associated orders and stock
		buyOrder, err := GetOrderByID(db, trade.BuyOrderID)
		if err != nil {
			return nil, err
		}
		trade.BuyOrder = buyOrder

		sellOrder, err := GetOrderByID(db, trade.SellOrderID)
		if err != nil {
			return nil, err
		}
		trade.SellOrder = sellOrder

		stock, err := GetStockBySymbol(db, trade.StockSymbol)
		if err != nil {
			return nil, err
		}
		trade.Stock = stock

		trades = append(trades, trade)
	}
	return trades, nil
}
