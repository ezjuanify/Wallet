package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/ezjuanify/wallet/internal/db"
	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/utils"
	"github.com/ezjuanify/wallet/internal/validation"
	"go.uber.org/zap"
)

type TransactionService struct {
	store *db.Store
}

func NewTransactionService(store *db.Store) *TransactionService {
	logger.Info("Initializing TransactionService")
	return &TransactionService{store: store}
}

func (ts *TransactionService) LogTransaction(ctx context.Context, tx *sql.Tx, txnUsername string, txnType model.TxnType, txnAmount int64, txnCounterparty *string) (*model.Transaction, *validation.WalletError) {
	fnName := "TransactionService.LogTransaction"
	if txnAmount <= 0 {
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_ZERO_AMOUNT,
			Message:   "Skip logging transaction due to zero amount",
			Timestamp: time.Now().UTC(),
			Err:       nil,
			Context: []zap.Field{
				zap.Int64("amount", txnAmount),
			},
		}
	}
	logger.Info(fmt.Sprintf("%s - Amount valid", fnName), zap.Int64("amount", txnAmount))

	txUser, err := validation.SanitizeAndValidateUsername(txnUsername)
	if err != nil {
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_SANITIZE_USERNAME_FAILED,
			Message:   "Failed to sanitize username",
			Timestamp: time.Now().UTC(),
			Err:       err,
			Context: []zap.Field{
				zap.String("username", txUser),
			},
		}
	}
	logger.Info(fmt.Sprintf("%s - Username sanitized", fnName), zap.String("username", txUser))

	timestamp := time.Now().UTC()
	hash := utils.GenerateTransactionHash(txUser, txnType, txnAmount, txnCounterparty, timestamp.Format(time.RFC3339))
	logger.Info(fmt.Sprintf("%s - Generated hash", fnName), zap.String("hash", hash))

	txn := model.Transaction{
		Username:     txUser,
		TxnType:      txnType,
		Amount:       txnAmount,
		Counterparty: txnCounterparty,
		Timestamp:    timestamp,
		Hash:         hash,
	}
	err = ts.store.InsertTransaction(ctx, tx, txn)
	if err != nil {
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_LOG_TRANSACTION_FAILED,
			Message:   "Failed to log transaction",
			Timestamp: time.Now().UTC(),
			Context: []zap.Field{
				zap.Any("transaction", txn),
			},
		}
	}
	logger.Info(fmt.Sprintf("%s - Transaction logged successfully", fnName), zap.Any("transaction", txn))
	return &txn, nil
}

func (ts *TransactionService) DoFetchTransaction(ctx context.Context, txnUsername string, txnCounterparty string, txnType string, txnLimit string) ([]model.Transaction, *model.Criteria, *validation.WalletError) {
	fnName := "TransactionService.DoFetchTransaction"
	queryUsername := validation.SanitizeUsernameWithoutError(txnUsername)
	logger.Info(fmt.Sprintf("%s - Username sanitized", fnName), zap.String("username", queryUsername))

	queryTxnType := ""
	if model.IsTxnTypeValid(txnType) {
		queryTxnType = txnType
	}
	logger.Info(fmt.Sprintf("%s - Transaction type valid", fnName), zap.String("txnType", queryTxnType))

	queryCounterparty := validation.SanitizeUsernameWithoutError(txnCounterparty)
	logger.Info(fmt.Sprintf("%s - Counterparty sanitized", fnName), zap.String("counterparty", queryCounterparty))

	queryLimit, err := strconv.Atoi(txnLimit)
	if err != nil {
		queryLimit = 0
	}
	logger.Info(fmt.Sprintf("%s - Limit converted to int", fnName), zap.Int("limit", queryLimit))

	query := &model.Criteria{
		Username:     queryUsername,
		Counterparty: queryCounterparty,
		TxnType:      model.TxnType(queryTxnType),
		Limit:        queryLimit,
	}
	logger.Info(fmt.Sprintf("%s - query", fnName), zap.Any("query", query))

	transactions, err := ts.store.FetchTransaction(ctx, query)
	if err != nil {
		return nil, nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_FETCH_TRANSACTION_FAILED,
			Message:   "Failed to fetch transaction",
			Timestamp: time.Now().UTC(),
			Err:       err,
			Context: []zap.Field{
				zap.Any("query", query),
			},
		}
	}
	logger.Info(fmt.Sprintf("%s - Transaction fetched successfully", fnName), zap.Any("transactions", transactions))
	return transactions, query, nil
}
