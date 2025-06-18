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

func (h *WalletHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := utils.DecodeRequest(r)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	logger.Debug("Decoded withdraw payload", zap.Any("payload", payload))

	wallet, err := h.WithdrawService.DoWithdraw(ctx, payload.Username, payload.Amount)
	if err != nil {
		logger.Error("Wallet withdraw failed", zap.String("error", err.Error()), zap.String("user", payload.Username))
		http.Error(w, fmt.Sprintf("Withdraw Error: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.transactionService.LogTransaction(ctx, payload.Username, model.TypeWithdraw, payload.Amount, nil); err != nil {
		logger.Warn("Failed to log transaction", zap.String("error", err.Error()))
	}

	resp := &response.TransactionResponse{
		Status:          http.StatusOK,
		TransactionType: model.TypeWithdraw,
		Wallet:          *wallet,
	}
	logger.Debug("Sending withdraw response", zap.Any("response", resp))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("Failed to encode withdraw response", zap.String("error", err.Error()))
	}
}
