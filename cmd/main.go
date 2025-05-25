package main

import (
	"log"
	"order-matching/api/v1/server"
)

func main() {
	defer server.Close()

	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
