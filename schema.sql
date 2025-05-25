-- Create database
CREATE DATABASE IF NOT EXISTS order_matching;
USE order_matching;

-- Create stocks table
CREATE TABLE IF NOT EXISTS stocks (
    symbol VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    current_price DECIMAL(10,2) NOT NULL,
    day_high DECIMAL(10,2) NOT NULL,
    day_low DECIMAL(10,2) NOT NULL,
    volume BIGINT NOT NULL,
    market_cap DECIMAL(15,2) NOT NULL,
    sector VARCHAR(50) NOT NULL,
    last_updated TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_sector (sector),
    INDEX idx_price (current_price)
);

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    type ENUM('BUY', 'SELL') NOT NULL,
    category ENUM('LIMIT', 'MARKET') NOT NULL,
    stock_symbol VARCHAR(10) NOT NULL,
    quantity INT UNSIGNED NOT NULL,
    filled_quantity INT UNSIGNED DEFAULT 0,
    price DECIMAL(10,2) NOT NULL,
    status ENUM('PENDING', 'PARTIALLY_FILLED', 'MATCHED', 'CANCELLED') DEFAULT 'PENDING',
    user_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (stock_symbol) REFERENCES stocks(symbol),
    INDEX idx_type_status (type, status),
    INDEX idx_stock_status (stock_symbol, status),
    INDEX idx_user (user_id)
);

-- Create trades table
CREATE TABLE IF NOT EXISTS trades (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    buy_order_id BIGINT UNSIGNED NOT NULL,
    sell_order_id BIGINT UNSIGNED NOT NULL,
    stock_symbol VARCHAR(10) NOT NULL,
    quantity INT UNSIGNED NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    executed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (buy_order_id) REFERENCES orders(id),
    FOREIGN KEY (sell_order_id) REFERENCES orders(id),
    FOREIGN KEY (stock_symbol) REFERENCES stocks(symbol),
    INDEX idx_stock_time (stock_symbol, executed_at),
    INDEX idx_orders (buy_order_id, sell_order_id)
);

-- Insert initial stock data
INSERT INTO stocks (symbol, name, description, current_price, day_high, day_low, volume, market_cap, sector, last_updated) VALUES
('NXTECH', 'Nexus Technologies', 'Advanced technology solutions provider', 150.00, 155.00, 145.00, 1000000, 15000000000.00, 'Technology', NOW()),
('QNTUM', 'Quantum Dynamics', 'Quantum computing research and development', 200.00, 210.00, 195.00, 800000, 20000000000.00, 'Technology', NOW()),
('CYBEX', 'Cyber Matrix Systems', 'Cybersecurity solutions provider', 175.00, 180.00, 170.00, 1200000, 17500000000.00, 'Technology', NOW()),
('SOLRX', 'Solar Matrix Energy', 'Renewable energy solutions', 125.00, 130.00, 120.00, 1500000, 12500000000.00, 'Energy', NOW()),
('FUSON', 'Fusion Power Corp', 'Nuclear fusion research and development', 300.00, 310.00, 290.00, 600000, 30000000000.00, 'Energy', NOW()),
('GENUM', 'Genome Solutions', 'Genetic research and biotechnology', 250.00, 260.00, 240.00, 700000, 25000000000.00, 'Healthcare', NOW()),
('MEDIX', 'Medical Innovations X', 'Medical device manufacturer', 180.00, 185.00, 175.00, 900000, 18000000000.00, 'Healthcare', NOW()),
('AITHN', 'AI Think Networks', 'Artificial intelligence solutions', 220.00, 225.00, 215.00, 1100000, 22000000000.00, 'Technology', NOW()),
('NRLNK', 'Neural Link Systems', 'Brain-computer interface technology', 275.00, 280.00, 270.00, 500000, 27500000000.00, 'Technology', NOW()),
('COGNT', 'Cognitive Tech Labs', 'Cognitive computing solutions', 190.00, 195.00, 185.00, 1000000, 19000000000.00, 'Technology', NOW());

-- Create stored procedure for updating order status
DELIMITER //

DROP PROCEDURE IF EXISTS update_order_status //
CREATE PROCEDURE update_order_status(IN order_id BIGINT UNSIGNED)
BEGIN
    DECLARE total_qty INT UNSIGNED;
    DECLARE filled_qty INT UNSIGNED;
    
    -- Get order quantities
    SELECT quantity, filled_quantity 
    INTO total_qty, filled_qty
    FROM orders 
    WHERE id = order_id;
    
    -- Update status based on filled quantity
    IF filled_qty = 0 THEN
        UPDATE orders SET status = 'PENDING' WHERE id = order_id;
    ELSEIF filled_qty < total_qty THEN
        UPDATE orders SET status = 'PARTIALLY_FILLED' WHERE id = order_id;
    ELSE
        UPDATE orders SET status = 'MATCHED' WHERE id = order_id;
    END IF;
END //

DROP PROCEDURE IF EXISTS match_market_order //
CREATE PROCEDURE match_market_order(IN order_id BIGINT UNSIGNED)
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE match_id BIGINT UNSIGNED;
    DECLARE match_qty INT UNSIGNED;
    DECLARE match_price DECIMAL(10,2);
    DECLARE remaining_qty INT UNSIGNED;
    DECLARE order_type ENUM('BUY', 'SELL');
    DECLARE order_symbol VARCHAR(10);
    
    DECLARE match_cursor CURSOR FOR
        SELECT id, quantity - filled_quantity, price
        FROM orders
        WHERE stock_symbol = order_symbol
        AND type = CASE order_type WHEN 'BUY' THEN 'SELL' ELSE 'BUY' END
        AND status IN ('PENDING', 'PARTIALLY_FILLED')
        ORDER BY CASE order_type 
            WHEN 'BUY' THEN price END ASC,
            CASE order_type WHEN 'SELL' THEN price END DESC,
            created_at ASC;
    
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
    
    -- Get order details
    SELECT type, stock_symbol, quantity - filled_quantity
    INTO order_type, order_symbol, remaining_qty
    FROM orders
    WHERE id = order_id AND category = 'MARKET' AND status != 'MATCHED';
    
    -- Start transaction
    START TRANSACTION;
    
    OPEN match_cursor;
    match_loop: LOOP
        FETCH match_cursor INTO match_id, match_qty, match_price;
        IF done OR remaining_qty = 0 THEN
            LEAVE match_loop;
        END IF;
        
        -- Calculate trade quantity
        SET match_qty = LEAST(remaining_qty, match_qty);
        
        -- Create trade record
        INSERT INTO trades (
            buy_order_id, 
            sell_order_id, 
            stock_symbol,
            quantity,
            price,
            executed_at
        ) VALUES (
            CASE order_type WHEN 'BUY' THEN order_id ELSE match_id END,
            CASE order_type WHEN 'SELL' THEN order_id ELSE match_id END,
            order_symbol,
            match_qty,
            match_price,
            NOW()
        );
        
        -- Update quantities
        UPDATE orders 
        SET filled_quantity = filled_quantity + match_qty
        WHERE id IN (order_id, match_id);
        
        -- Update remaining quantity
        SET remaining_qty = remaining_qty - match_qty;
    END LOOP;
    
    CLOSE match_cursor;
    
    -- Cancel remaining quantity for market orders
    IF remaining_qty > 0 THEN
        UPDATE orders SET status = 'CANCELLED' WHERE id = order_id;
    END IF;
    
    -- Update status for matched orders
    CALL update_order_status(order_id);
    
    COMMIT;
END //

DELIMITER ; 