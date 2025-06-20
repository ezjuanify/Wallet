package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/response"
	"github.com/ezjuanify/wallet/internal/utils"
	"go.uber.org/zap"
)

func (h *WalletHandler) DepositHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	tx, err := h.store.BeginTransaction(ctx)
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				logger.Error("Failed to commit transaction", zap.String("error", err.Error()))
			}
		}
	}()

	payload, err := utils.DecodeRequest(r)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	logger.Debug("Decoded deposit payload", zap.Any("payload", payload))

	wallet, err := h.depositService.DoDeposit(ctx, tx, payload.Username, payload.Amount, false)
	if err != nil {
		logger.Error("Wallet deposit failed", zap.String("error", err.Error()), zap.String("user", payload.Username))
		http.Error(w, fmt.Sprintf("Deposit Error: %v", err), http.StatusInternalServerError)
		return
	}

	if err = h.transactionService.LogTransaction(ctx, tx, payload.Username, model.TypeDeposit, payload.Amount, nil); err != nil {
		logger.Warn("Failed to log transaction", zap.String("error", err.Error()))
	}

	resp := &response.TransactionResponse{
		Status:          http.StatusOK,
		TransactionType: model.TypeDeposit,
		Wallet:          *wallet,
	}
	logger.Debug("Sending deposit response", zap.Any("response", resp))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("Failed to encode deposit response", zap.String("error", err.Error()))
	}
}
