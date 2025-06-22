package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ezjuanify/wallet/internal/db"
	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model/response"
	"github.com/ezjuanify/wallet/internal/service"
	"github.com/ezjuanify/wallet/internal/validation"
	"go.uber.org/zap"
)

type WalletHandler struct {
	store              *db.Store
	walletService      *service.WalletService
	depositService     *service.DepositService
	withdrawService    *service.WithdrawService
	transactionService *service.TransactionService
}

func NewWalletHandler(
	store *db.Store,
	s *service.WalletService,
	ds *service.DepositService,
	ws *service.WithdrawService,
	ts *service.TransactionService,
) *WalletHandler {
	logger.Debug("Initializing WalletHandler")
	return &WalletHandler{
		store:              store,
		walletService:      s,
		depositService:     ds,
		withdrawService:    ws,
		transactionService: ts,
	}
}

func SendJSONResponse(respName string, w http.ResponseWriter, status int, resp any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error(fmt.Sprintf("%s - Failed to encode response", respName), zap.Error(err))
	}
}

func FinalizeTransactionResponse(fnName string, tx *sql.Tx, w http.ResponseWriter, aErrs *validation.AppErrors) {
	if p := recover(); p != nil {
		if tx != nil {
			tx.Rollback()
		}

		wrappedErr := validation.WalletError{
			Name:      fnName,
			Status:    http.StatusInternalServerError,
			Code:      validation.ERR_PANIC_OCCURED,
			Message:   "Application panic",
			Timestamp: time.Now().UTC(),
			Err:       fmt.Errorf("panic: %v", p),
		}
		aErrs.AddError(wrappedErr)
		aErrs.LogAll()

		panic(p)
	}

	if aErrs.GetErrsCount() > 0 {
		if tx != nil {
			tx.Rollback()
			logger.Warn(fmt.Sprintf("%s - Rolling back due to application errors", fnName))
		}

		aErrs.LogAll()

		first := aErrs.First()
		resp := response.ErrorResponse{
			Status:  first.Status,
			Code:    string(first.Code),
			Message: first.Message,
		}
		SendJSONResponse(fnName, w, first.Status, resp)
		return
	}

	if tx == nil {
		return
	}

	if err := tx.Commit(); err != nil {
		wrappedErr := validation.WalletError{
			Name:      fnName,
			Status:    http.StatusInternalServerError,
			Code:      validation.ERR_TRANSACTION_COMMIT_FAILED,
			Message:   "Failed to commit transaction",
			Timestamp: time.Now().UTC(),
			Err:       err,
		}
		aErrs.AddError(wrappedErr)
		aErrs.LogAll()

		resp := response.ErrorResponse{
			Status:  wrappedErr.Status,
			Code:    string(wrappedErr.Code),
			Message: wrappedErr.Message,
		}
		SendJSONResponse(fnName, w, wrappedErr.Status, resp)
		return
	}
	logger.Info(fmt.Sprintf("%s - Transaction committed", fnName))
}
