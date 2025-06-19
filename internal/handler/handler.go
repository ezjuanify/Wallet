package handler

import (
	"github.com/ezjuanify/wallet/internal/db"
	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/service"
)

type WalletHandler struct {
	store              *db.Store
	depositService     *service.DepositService
	withdrawService    *service.WithdrawService
	transactionService *service.TransactionService
}

func NewWalletHandler(
	store *db.Store,
	ds *service.DepositService,
	ws *service.WithdrawService,
	ts *service.TransactionService,
) *WalletHandler {
	logger.Debug("Initializing WalletHandler")
	return &WalletHandler{
		store:              store,
		depositService:     ds,
		withdrawService:    ws,
		transactionService: ts,
	}
}
