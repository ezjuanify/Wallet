package handler

import (
	"fmt"
	"net/http"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model/response"
	"github.com/ezjuanify/wallet/internal/validation"
	"go.uber.org/zap"
)

func (h *WalletHandler) TransactionHandler(w http.ResponseWriter, r *http.Request) {
	fnName := "WalletHandler.TransactionHandler"

	ctx := r.Context()

	aErrs := validation.NewHandlerErrors()

	defer func() {
		FinalizeTransactionResponse(fnName, nil, w, aErrs)
	}()

	queries := r.URL.Query()
	username := queries.Get("username")
	counterparty := queries.Get("counterparty")
	txnType := queries.Get("type")
	limit := queries.Get("limit")

	logger.Info(fmt.Sprintf("%s - Query values", fnName),
		zap.String("username", username),
		zap.String("counterparty", counterparty),
		zap.String("txnType", txnType),
		zap.String("limit", limit),
	)

	transactions, criteria, appErr := h.transactionService.DoFetchTransaction(ctx, username, counterparty, txnType, limit)
	if appErr != nil {
		appErr.Status = http.StatusInternalServerError
		aErrs.AddError(*appErr)
		return
	}
	logger.Info(fmt.Sprintf("%s - Transaction fetched successfully", fnName), zap.Any("transaction", transactions))

	resp := &response.TransactionQueryResponse{
		Status:       http.StatusOK,
		Criteria:     criteria,
		Transactions: transactions,
	}
	logger.Info(fmt.Sprintf("%s - Sending transaction response", fnName), zap.Any("response", resp))
	SendJSONResponse(fnName, w, resp.Status, resp)
}
