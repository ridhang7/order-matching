package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

// Initialize sets up the database connection
func Initialize() error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error connecting to database: %v", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
}
