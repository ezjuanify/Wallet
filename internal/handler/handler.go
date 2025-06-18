package handler

import (
	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/service"
)

type WalletHandler struct {
	DepositService     *service.DepositService
	WithdrawService    *service.WithdrawService
	transactionService *service.TransactionService
}

func NewWalletHandler(
	ds *service.DepositService,
	ws *service.WithdrawService,
	ts *service.TransactionService,
) *WalletHandler {
	logger.Debug("Initializing WalletHandler")
	return &WalletHandler{
		DepositService:     ds,
		WithdrawService:    ws,
		transactionService: ts,
	}
}
