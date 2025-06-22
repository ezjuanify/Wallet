package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/validation"
	"go.uber.org/zap"
)

type DepositStore interface {
	UpsertWallet(ctx context.Context, tx *sql.Tx, username string, amount int64) (*model.Wallet, error)
	FetchWallet(ctx context.Context, username string) (*model.Wallet, error)
}

type DepositService struct {
	store DepositStore
}

func NewDepositService(store DepositStore) *DepositService {
	logger.Debug("Initializing DepositService")
	return &DepositService{store: store}
}

func (s *DepositService) DoDeposit(ctx context.Context, tx *sql.Tx, username string, amount int64, isCounterparty bool) (*model.Wallet, *validation.WalletError) {
	fnName := "DepositService.DoDeposit"
	logger.Info(fmt.Sprintf("%s - Params received", fnName), zap.String("username", username), zap.Int64("amount", amount))
	username, err := validation.SanitizeAndValidateUsername(username)
	if err != nil {
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_SANITIZE_USERNAME_FAILED,
			Message:   "Failed to sanitize username",
			Timestamp: time.Now().UTC(),
			Err:       err,
			Context: []zap.Field{
				zap.String("username", username),
			},
		}
	}
	logger.Info(fmt.Sprintf("%s - Username sanitized", fnName), zap.String("username", username))

	if err := validation.ValidateAmount(amount); err != nil {
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_AMOUNT_VALIDATION_FAILED,
			Message:   "Amount validation failed",
			Timestamp: time.Now().UTC(),
			Err:       err,
			Context: []zap.Field{
				zap.Int64("amount", amount),
			},
		}
	}
	logger.Info(fmt.Sprintf("%s - Amount validated", fnName), zap.Int64("amount", amount))

	currentWallet, err := s.store.FetchWallet(ctx, username)
	if err != nil {
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_FETCH_WALLET_FAILED,
			Message:   "Failed to fetch wallet",
			Timestamp: time.Now().UTC(),
			Err:       err,
			Context:   nil,
		}
	}
	if currentWallet == nil {
		if isCounterparty {
			return nil, &validation.WalletError{
				Name:      fnName,
				Code:      validation.ERR_WALLET_DOES_NOT_EXIST,
				Message:   "Counterparty wallet does not exist",
				Timestamp: time.Now().UTC(),
				Err:       nil,
				Context: []zap.Field{
					zap.String("counterparty", username),
				},
			}
		}
		logger.Warn(fmt.Sprintf("%s - No wallet found for user", fnName))
	}

	if currentWallet != nil {
		newBalance := currentWallet.Balance + amount
		if err := validation.ValidateWalletBalance(newBalance); err != nil {
			return nil, &validation.WalletError{
				Name:      fnName,
				Code:      validation.ERR_WALLET_BALANCE_VALIDATION_FAILED,
				Message:   "Wallet balance would exceed limit",
				Timestamp: time.Now().UTC(),
				Err:       err,
				Context: []zap.Field{
					zap.String("username", username),
					zap.Int64("balance", currentWallet.Balance),
					zap.Int64("amount", amount),
					zap.Int64("resulting", newBalance),
				},
			}
		}
		logger.Info(
			fmt.Sprintf("%s - Wallet balance validated", fnName),
			zap.Int64("wallet_balance", currentWallet.Balance),
			zap.Int64("amount", amount),
			zap.Int64("resulting_balance", newBalance),
		)
	}

	updatedWallet, err := s.store.UpsertWallet(ctx, tx, username, amount)
	if err != nil {
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_DB_UPSERT_FAILED,
			Message:   "Failed to upsert wallet",
			Timestamp: time.Now().UTC(),
			Err:       err,
			Context: []zap.Field{
				zap.String("username", username),
				zap.Int64("amount", amount),
			},
		}
	}
	logger.Info(fmt.Sprintf("%s - Upserted wallet", fnName), zap.Any("wallet", updatedWallet))
	return updatedWallet, nil
}
