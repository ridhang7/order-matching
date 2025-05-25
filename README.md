# Order Matching System

## Dependencies
- github.com/go-sql-driver/mysql v1.7.1
- github.com/gorilla/mux v1.8.1
- github.com/joho/godotenv v1.5.1

## Database Schema
```sql
-- Stocks table
CREATE TABLE stocks (
    symbol VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    current_price DECIMAL(10,2) NOT NULL,
    day_high DECIMAL(10,2) NOT NULL,
    day_low DECIMAL(10,2) NOT NULL,
    volume BIGINT NOT NULL,
    market_cap D6ECIMAL(15,2) NOT NULL,
    sector VARCHAR(50) NOT NULL,
    last_updated TIMESTAMP NOT NULL
);

-- Orders table
CREATE TABLE orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    type ENUM('BUY', 'SELL') NOT NULL,
    category ENUM('LIMIT', 'MARKET') NOT NULL,
    stock_symbol VARCHAR(10) NOT NULL,
    quantity INT UNSIGNED NOT NULL,
    filled_quantity INT UNSIGNED DEFAULT 0,
    price DECIMAL(10,2) NOT NULL,
    status ENUM('PENDING', 'PARTIALLY_FILLED', 'MATCHED', 'CANCELLED') DEFAULT 'PENDING',
    user_id BIGINT UNSIGNED NOT NULL,
    FOREIGN KEY (stock_symbol) REFERENCES stocks(symbol)
);

-- Trades table
CREATE TABLE trades (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    buy_order_id BIGINT UNSIGNED NOT NULL,
    sell_order_id BIGINT UNSIGNED NOT NULL,
    stock_symbol VARCHAR(10) NOT NULL,
    quantity INT UNSIGNED NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    executed_at TIMESTAMP NOT NULL,
    FOREIGN KEY (buy_order_id) REFERENCES orders(id),
    FOREIGN KEY (sell_order_id) REFERENCES orders(id),
    FOREIGN KEY (stock_symbol) REFERENCES stocks(symbol)
);
```

## Running the Application

1. Set up environment variables in `.env`:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=order_matching
```

2. Run the application:
```bash
go run cmd/main.go
```

The server will start on port 8080.

## API Endpoints

### Orders
- `POST /api/v1/orders` - Create a new order
- `GET /api/v1/orders` - Get all orders
- `GET /api/v1/orders/{id}` - Get order by ID
- `POST /api/v1/orders/{id}/cancel` - Cancel an order
- `GET /api/v1/orders/stock/{symbol}` - Get orders by stock symbol

### Trades
- `GET /api/v1/trades` - Get all trades
- `GET /api/v1/trades/{id}` - Get trade by ID

The system supports the following stock symbols:
- NXTECH (Nexus Technologies)
- QNTUM (Quantum Dynamics)
- CYBEX (Cyber Matrix Systems)
- SOLRX (Solar Matrix Energy)
- FUSON (Fusion Power Corp)
- GENUM (Genome Solutions)
- MEDIX (Medical Innovations X)
- AITHN (AI Think Networks)
- NRLNK (Neural Link Systems)
- COGNT (Cognitive Tech Labs)
