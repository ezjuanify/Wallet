package handler

import "github.com/ezjuanify/wallet/internal/service"

type WalletHandler struct {
	walletService      *service.WalletService
	transactionService *service.TransactionService
}

func NewWalletHandler(ws *service.WalletService, ts *service.TransactionService) *WalletHandler {
	return &WalletHandler{walletService: ws, transactionService: ts}
}
