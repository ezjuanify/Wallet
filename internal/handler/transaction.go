package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model/response"
	"go.uber.org/zap"
)

func (wh *WalletHandler) TransactionHandler(w http.ResponseWriter, r *http.Request) {
	fnName := "WalletHandler.TransactionHandler"
	ctx := r.Context()

	queries := r.URL.Query()
	username := queries.Get("username")
	counterparty := queries.Get("counterparty")
	txnType := queries.Get("type")
	limit := queries.Get("limit")

	logger.Info(fmt.Sprintf("%s - Queried values", fnName),
		zap.String("username", username),
		zap.String("counterparty", counterparty),
		zap.String("txnType", txnType),
		zap.String("limit", limit),
	)

	transactions, criteria, err := wh.transactionService.DoFetchTransaction(ctx, username, counterparty, txnType, limit)
	if err != nil {
		logger.Error("Failed to fetch transaction", zap.Error(err))
		http.Error(w, fmt.Sprintf("Failed to fetch transaction: %s", err), http.StatusBadRequest)
		return
	}

	resp := &response.TransactionQueryResponse{
		Status:       http.StatusOK,
		Criteria:     criteria,
		Transactions: transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("Failed to encode transaction query response", zap.Error(err))
	}
}
