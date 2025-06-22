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

func (h *WalletHandler) TransferHandler(w http.ResponseWriter, r *http.Request) {
	fnName := "WalletHandler.TransferHandler"

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
	logger.Debug("Decoded transfer payload", zap.Any("payload", payload))

	wallet, appErr := h.withdrawService.DoWithdraw(ctx, tx, payload.Username, payload.Amount)
	if appErr != nil {
		appErr.Status = http.StatusInternalServerError
		appErrs.AddError(*appErr)
		return
	}
	logger.Info(fmt.Sprintf("%s - Transfer out successful", fnName), zap.Any("wallet", wallet))

	_, appErr = h.depositService.DoDeposit(ctx, tx, *payload.Counterparty, payload.Amount, true)
	if appErr != nil {
		appErr.Status = http.StatusInternalServerError
		appErrs.AddError(*appErr)
		return
	}
	logger.Info(fmt.Sprintf("%s - Transfer in successful", fnName), zap.Any("wallet", wallet))

	username, err := validation.SanitizeAndValidateUsername(payload.Username)
	if err != nil {
		appErrs.AddError(
			validation.WalletError{
				Name:      fnName,
				Status:    http.StatusInternalServerError,
				Code:      validation.ERR_SANITIZE_USERNAME_FAILED,
				Message:   "Failed to sanitize username",
				Timestamp: time.Now().UTC(),
				Err:       err,
			},
		)
		return
	}
	logger.Info(fmt.Sprintf("%s - Username sanitized", fnName), zap.Any("username", username))

	counterparty, err := validation.SanitizeAndValidateUsername(*payload.Counterparty)
	if err != nil {
		appErrs.AddError(
			validation.WalletError{
				Name:      fnName,
				Status:    http.StatusInternalServerError,
				Code:      validation.ERR_SANITIZE_USERNAME_FAILED,
				Message:   "Failed to sanitize counterparty",
				Timestamp: time.Now().UTC(),
				Err:       err,
			},
		)
		return
	}
	logger.Info(fmt.Sprintf("%s - Counterparty sanitized", fnName), zap.Any("counterparty", counterparty))

	outTransaction, appErr := h.transactionService.LogTransaction(ctx, tx, username, model.TypeTransferOut, payload.Amount, &counterparty)
	if appErr != nil {
		appErr.Status = http.StatusInternalServerError
		appErrs.AddError(*appErr)
		return
	}
	logger.Info(fmt.Sprintf("%s - Transfer out transaction logged successfully", fnName), zap.Any("outTransaction", outTransaction))

	inTransaction, appErr := h.transactionService.LogTransaction(ctx, tx, counterparty, model.TypeTransferIn, payload.Amount, &username)
	if appErr != nil {
		appErr.Status = http.StatusInternalServerError
		appErrs.AddError(*appErr)
		return
	}
	logger.Info(fmt.Sprintf("%s - Transfer in transaction logged successfully", fnName), zap.Any("inTransaction", inTransaction))

	resp := &response.TransactionResponse{
		Status:          http.StatusOK,
		TransactionType: model.TypeTransfer,
		Wallet:          *wallet,
		Counterparty:    &counterparty,
	}
	logger.Info(fmt.Sprintf("%s - Sending transfer response", fnName), zap.Any("response", resp))
	SendJSONResponse(fnName, w, int(resp.Status), resp)
}
