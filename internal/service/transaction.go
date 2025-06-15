package service

import (
	"context"
	"time"

	"github.com/ezjuanify/wallet/internal/db"
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/utils"
)

type TransactionService struct {
	store *db.Store
}

func NewTransactionService(store *db.Store) *TransactionService {
	return &TransactionService{store: store}
}

func (ts *TransactionService) LogTransaction(ctx context.Context, txUser string, txType model.TransactionType, txAmount int64, txCounterparty *string) error {
	timestamp := time.Now().UTC()
	hash := utils.GenerateTransactionHash(txUser, string(txType), txAmount, txCounterparty, timestamp.Format(time.RFC3339))

	txn := model.Transaction{
		Username:     txUser,
		Type:         string(txType),
		Amount:       txAmount,
		Counterparty: txCounterparty,
		Timestamp:    timestamp,
		Hash:         hash,
	}
	err := ts.store.InsertTransaction(ctx, txn)
	return err
}
