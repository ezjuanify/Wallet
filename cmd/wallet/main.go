package main

import (
	"log"
	"net/http"

	"github.com/ezjuanify/wallet/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)

	serve := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server starting on :8080")
	if err := serve.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
