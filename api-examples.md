# Order Matching System API Documentation

## Base URL
All endpoints are prefixed with `/api/v1`

## Order Endpoints

### 1. Create Order
Creates a new order (buy or sell).

```bash
# Create a LIMIT BUY order
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "type": "BUY",
    "category": "LIMIT",
    "stock_symbol": "AAPL",
    "quantity": 100,
    "price": 150.50,
    "user_id": 1
  }'

# Create a MARKET SELL order
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "type": "SELL",
    "category": "MARKET",
    "stock_symbol": "AAPL",
    "quantity": 50,
    "user_id": 2
  }'
```

Example Success Response:
```json
{
  "message": "Order processed successfully",
  "order": {
    "id": 1,
    "type": "BUY",
    "category": "LIMIT",
    "stock_symbol": "AAPL",
    "quantity": 100,
    "price": 150.50,
    "user_id": 1,
    "status": "PENDING",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
  }
}
```

### 2. Get All Orders
Retrieves all orders in the system.

```bash
curl -X GET http://localhost:8080/api/v1/orders
```

Example Response:
```json
{
  "orders": [
    {
      "id": 1,
      "type": "BUY",
      "category": "LIMIT",
      "stock_symbol": "AAPL",
      "quantity": 100,
      "price": 150.50,
      "user_id": 1,
      "status": "PENDING",
      "created_at": "2024-03-20T10:00:00Z",
      "updated_at": "2024-03-20T10:00:00Z"
    }
  ]
}
```

### 3. Get Order by ID
Retrieves a specific order by its ID.

```bash
curl -X GET http://localhost:8080/api/v1/orders/1
```

Example Response:
```json
{
  "order": {
    "id": 1,
    "type": "BUY",
    "category": "LIMIT",
    "stock_symbol": "AAPL",
    "quantity": 100,
    "price": 150.50,
    "user_id": 1,
    "status": "PENDING",
    "created_at": "2024-03-20T10:00:00Z",
    "updated_at": "2024-03-20T10:00:00Z"
  }
}
```

### 4. Cancel Order
Cancels a pending order.

```bash
curl -X PUT http://localhost:8080/api/v1/orders/1/cancel
```

Example Response:
```json
{
  "message": "Order cancelled successfully",
  "order": {
    "id": 1,
    "status": "CANCELLED",
    "updated_at": "2024-03-20T10:05:00Z"
  }
}
```

### 5. Get Order Book
Retrieves the current order book for a specific stock.

```bash
curl -X GET "http://localhost:8080/api/v1/orders/book?symbol=AAPL"
```

Example Response:
```json
{
  "order_book": {
    "buy_orders": [
      {
        "id": 1,
        "type": "BUY",
        "price": 150.50,
        "quantity": 100,
        "status": "PENDING"
      }
    ],
    "sell_orders": [
      {
        "id": 2,
        "type": "SELL",
        "price": 151.00,
        "quantity": 50,
        "status": "PENDING"
      }
    ]
  }
}
```

## Trade Endpoints

### 1. Get All Trades
Retrieves all trades in the system.

```bash
curl -X GET "http://localhost:8080/api/v1/trades?symbol=AAPL"
```

Example Response:
```json
{
  "trades": [
    {
      "id": 1,
      "buy_order_id": 1,
      "sell_order_id": 2,
      "stock_symbol": "AAPL",
      "quantity": 50,
      "price": 150.75,
      "executed_at": "2024-03-20T10:01:00Z"
    }
  ]
}
```

### 2. Get Trade by ID
Retrieves a specific trade by its ID.

```bash
curl -X GET "http://localhost:8080/api/v1/trades/1"
```

Example Response:
```json
{
  "trade": {
    "id": 1,
    "buy_order_id": 1,
    "sell_order_id": 2,
    "stock_symbol": "AAPL",
    "quantity": 50,
    "price": 150.75,
    "executed_at": "2024-03-20T10:01:00Z"
  }
}
```

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request",
  "details": "Field validation failed"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found",
  "details": "Order/Trade with ID {id} not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "details": "Failed to process request"
}
``` 