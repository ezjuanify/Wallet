package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/response"
	"github.com/ezjuanify/wallet/internal/utils"
	"github.com/ezjuanify/wallet/internal/validation"
	"go.uber.org/zap"
)

func (h *WalletHandler) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	fnName := "WalletHandler.WithdrawHandler"

	ctx := r.Context()

	appErrs := validation.NewHandlerErrors()

	tx, err := h.store.BeginTransaction(ctx)
	if err != nil {
		appErrs.AddError(
			validation.WalletError{
				Name:      fnName,
				Status:    http.StatusInternalServerError,
				Code:      validation.ERR_TRANSACTION_START_FAILED,
				Message:   "Failed to start transaction",
				Timestamp: time.Now().UTC(),
				Err:       err,
			},
		)
		return
	}
	logger.Info(fmt.Sprintf("%s - Transaction Started", fnName))

	defer func() {
		FinalizeTransactionResponse(fnName, tx, w, appErrs)
	}()

	payload, err := utils.DecodeRequest(r)
	if err != nil {
		appErrs.AddError(
			validation.WalletError{
				Name:      fnName,
				Status:    http.StatusBadRequest,
				Code:      validation.ERR_INVALID_JSON_BODY,
				Message:   "Failed to decode JSON body",
				Timestamp: time.Now().UTC(),
				Err:       err,
			},
		)
		return
	}
	logger.Info(fmt.Sprintf("%s - Decoded withdraw payload", fnName), zap.Any("payload", payload))

	wallet, appErr := h.withdrawService.DoWithdraw(ctx, tx, payload.Username, payload.Amount)
	if appErr != nil {
		appErr.Status = http.StatusInternalServerError
		appErrs.AddError(*appErr)
		return
	}
	logger.Info(fmt.Sprintf("%s - Withdraw successful", fnName), zap.Any("wallet", wallet))

	transaction, appErr := h.transactionService.LogTransaction(ctx, tx, payload.Username, model.TypeWithdraw, payload.Amount, nil)
	if appErr != nil {
		appErr.Status = http.StatusInternalServerError
		appErrs.AddError(*appErr)
		return
	}
	logger.Info(fmt.Sprintf("%s - Transaction logged", fnName), zap.Any("transaction", transaction))

	resp := &response.TransactionResponse{
		Status:          http.StatusOK,
		TransactionType: model.TypeWithdraw,
		Wallet:          *wallet,
	}
	logger.Info(fmt.Sprintf("%s - Sending withdraw response", fnName), zap.Any("response", resp))
	SendJSONResponse(fnName, w, resp.Status, resp)
}
