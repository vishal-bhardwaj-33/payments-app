package main

import (
	"log"
	"payments-app/internal/services/payments"
)

func main() {
	if err := payments.StartServer(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
