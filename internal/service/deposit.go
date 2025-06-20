package service

import (
	"context"
	"database/sql"
	"fmt"

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

func (s *DepositService) DoDeposit(ctx context.Context, tx *sql.Tx, username string, amount int64, isCounterparty bool) (*model.Wallet, error) {
	funcName := "DepositService.DoDeposit"
	logger.Info(fmt.Sprintf("%s - Params received", funcName), zap.String("username", username), zap.Int64("amount", amount))
	username, err := validation.SanitizeAndValidateUsername(username)
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("%s - Username sanitized", funcName), zap.String("username", username))

	if err := validation.ValidateAmount(amount); err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("%s - Amount not breaching upper or lower limit", funcName), zap.Int64("amount", amount))

	currentWallet, err := s.store.FetchWallet(ctx, username)
	if err != nil {
		return nil, err
	}
	if currentWallet == nil {
		if isCounterparty {
			return nil, fmt.Errorf("%s - counterparty %s wallet does not exist", funcName, username)
		}
		logger.Warn(fmt.Sprintf("%s - No wallet found for user", funcName))
	}

	if currentWallet != nil {
		newBalance := currentWallet.Balance + amount
		if err := validation.ValidateWalletBalance(newBalance); err != nil {
			logger.Error(
				fmt.Sprintf("%s - Wallet balance validation failed", funcName),
				zap.Int64("wallet_balance", currentWallet.Balance),
				zap.Int64("amount", amount),
				zap.Int64("resulting_balance", newBalance),
				zap.Error(err),
			)
			return nil, fmt.Errorf("%s exceeds upper limit - wallet balance: %d  - deposit: %d - exceed: %d", funcName, currentWallet.Balance, amount, newBalance)
		}
		logger.Info(
			fmt.Sprintf("%s - Wallet balance validated", funcName),
			zap.Int64("wallet_balance", currentWallet.Balance),
			zap.Int64("amount", amount),
			zap.Int64("resulting_balance", newBalance),
		)
	}

	updatedWallet, err := s.store.UpsertWallet(ctx, tx, username, amount)
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("%s - Upserted wallet", funcName), zap.Any("wallet", updatedWallet))
	return updatedWallet, nil
}
