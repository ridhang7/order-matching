package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err)
		return
	}

	// Print environment variables (without password)
	fmt.Printf("DB_HOST: %s\n", os.Getenv("DB_HOST"))
	fmt.Printf("DB_PORT: %s\n", os.Getenv("DB_PORT"))
	fmt.Printf("DB_USER: %s\n", os.Getenv("DB_USER"))
	fmt.Printf("DB_NAME: %s\n", os.Getenv("DB_NAME"))

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	fmt.Printf("Attempting to connect with DSN (password hidden): %s\n",
		fmt.Sprintf("%s:***@tcp(%s:%s)/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME")))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error opening connection: %v\n", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Printf("Error connecting to the database: %v\n", err)
		return
	}
	fmt.Println("Successfully connected to MySQL database!")
}
