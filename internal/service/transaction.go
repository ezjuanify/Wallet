package service

import (
	"context"
	"fmt"
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

func (ts *TransactionService) LogTransaction(ctx context.Context, txUser string, txType model.TransactionType, txAmount int64, txCounterparty *string) error {
	if txAmount <= 0 {
		return fmt.Errorf("skipping transaction logging for 0 amount")
	}

	logger.Debug("TransactionService.LogTransaction - Sanitizing username", zap.String("username", txUser))
	txUser, err := validation.SanitizeAndValidateUsername(txUser)
	if err != nil {
		return err
	}
	logger.Info("TransactionService.LogTransaction - Username sanitized", zap.String("username", txUser))

	timestamp := time.Now().UTC()
	hash := utils.GenerateTransactionHash(txUser, string(txType), txAmount, txCounterparty, timestamp.Format(time.RFC3339))
	logger.Info("TransactionService.LogTransaction - Generated hash", zap.String("hash", hash))

	txn := model.Transaction{
		Username:     txUser,
		Type:         string(txType),
		Amount:       txAmount,
		Counterparty: txCounterparty,
		Timestamp:    timestamp,
		Hash:         hash,
	}
	logger.Info("TransactionService.LogTransaction - Inserting transaction", zap.Any("transaction", txn))
	err = ts.store.InsertTransaction(ctx, txn)
	return err
}
