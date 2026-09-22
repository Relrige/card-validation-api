package main

import (
	"log"
	"net/http"

	"github.com/Relrige/card-validator-api/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/validate", handler.ValidateCardHandler)

	port := ":8080"
	log.Printf("Starting Card Validation Service on port %s", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
