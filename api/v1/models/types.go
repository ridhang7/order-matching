package models

// OrderType represents the type of order (BUY/SELL)
type OrderType string

const (
	OrderTypeBuy  OrderType = "BUY"
	OrderTypeSell OrderType = "SELL"
)

// OrderCategory represents the category of order (LIMIT/MARKET)
type OrderCategory string

const (
	OrderCategoryLimit  OrderCategory = "LIMIT"
	OrderCategoryMarket OrderCategory = "MARKET"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending         OrderStatus = "PENDING"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusMatched         OrderStatus = "MATCHED"
	OrderStatusCancelled       OrderStatus = "CANCELLED"
)

// StockSymbol represents valid stock symbols
type StockSymbol string

const (
	StockNXTECH StockSymbol = "NXTECH"
	StockQNTUM  StockSymbol = "QNTUM"
	StockCYBEX  StockSymbol = "CYBEX"
	StockSOLRX  StockSymbol = "SOLRX"
	StockFUSON  StockSymbol = "FUSON"
	StockGENUM  StockSymbol = "GENUM"
	StockMEDIX  StockSymbol = "MEDIX"
	StockAITHN  StockSymbol = "AITHN"
	StockNRLNK  StockSymbol = "NRLNK"
	StockCOGNT  StockSymbol = "COGNT"
)
