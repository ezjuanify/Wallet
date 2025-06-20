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
	logger.Debug("Initializing TransactionService")
	return &TransactionService{store: store}
}

func (ts *TransactionService) LogTransaction(ctx context.Context, tx *sql.Tx, txnUsername string, txnType model.TxnType, txnAmount int64, txnCounterparty *string) error {
	fnName := "TransactionService.LogTransaction"
	if txnAmount <= 0 {
		return fmt.Errorf("%s - skipping transaction logging for 0 amount", fnName)
	}

	logger.Debug(fmt.Sprintf("%s - Sanitizing username", fnName), zap.String("username", txnUsername))
	txUser, err := validation.SanitizeAndValidateUsername(txnUsername)
	if err != nil {
		return err
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
	logger.Info(fmt.Sprintf("%s - Inserting transaction", fnName), zap.Any("transaction", txn))
	err = ts.store.InsertTransaction(ctx, tx, txn)
	return err
}

func (ts *TransactionService) DoFetchTransaction(ctx context.Context, txnUsername string, txnCounterparty string, txnType string, txnLimit string) ([]model.Transaction, *model.Criteria, error) {
	fnName := "TransactionService.DoFetchTransaction"
	criteriaUsername := validation.SanitizeUsernameWithoutError(txnUsername)
	logger.Debug(fmt.Sprintf("%s - username sanitized", fnName), zap.String("username", criteriaUsername))

	criteriaTxnType := ""
	if model.IsTxnTypeValid(txnType) {
		criteriaTxnType = txnType
	}
	logger.Debug(fmt.Sprintf("%s - txnType valid", fnName), zap.String("txnType", criteriaTxnType))

	criteriaCounterparty := validation.SanitizeUsernameWithoutError(txnCounterparty)
	logger.Debug(fmt.Sprintf("%s - counterparty sanitized", fnName), zap.String("counterparty", criteriaCounterparty))

	criteriaLimit, err := strconv.Atoi(txnLimit)
	if err != nil {
		criteriaLimit = 0
	}
	logger.Debug(fmt.Sprintf("%s - limit converted to int", fnName), zap.Int("limit", criteriaLimit))

	criteria := &model.Criteria{
		Username:     criteriaUsername,
		Counterparty: criteriaCounterparty,
		TxnType:      model.TxnType(criteriaTxnType),
		Limit:        criteriaLimit,
	}
	logger.Info(fmt.Sprintf("%s - criteria", fnName), zap.Any("criteria", criteria))

	transactions, err := ts.store.FetchTransaction(ctx, criteria)
	if err != nil {
		return nil, nil, err
	}
	logger.Info(fmt.Sprintf("%s - transactions result", fnName), zap.Any("transactions", transactions))
	return transactions, criteria, nil
}
