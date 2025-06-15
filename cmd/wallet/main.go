package main

import (
	"log"
	"net/http"

	"github.com/ezjuanify/wallet/internal/db"
	"github.com/ezjuanify/wallet/internal/handler"
	"github.com/ezjuanify/wallet/internal/service"
)

func main() {
	pgconfig := &db.PGConfig{
		Host: "localhost",
		Port: 5432,
		SSL:  "disable",
		DB:   "db_wallet_app",
		User: "db_wallet_app",
		Pass: "db_wallet_app",
	}
	store, err := db.NewStore(pgconfig)
	if err != nil {
		log.Fatalf("Failed to initialize DB Store: %v", err)
	}

	ws := service.NewWalletService(store)
	ts := service.NewTransactionService(store)
	wh := handler.NewWalletHandler(ws, ts)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/deposit", wh.DepositResponse)

	serve := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server starting on :8080")
	if err := serve.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
