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

type WithdrawStore interface {
	WithdrawWallet(ctx context.Context, tx *sql.Tx, username string, amount int64) (*model.Wallet, error)
	FetchWallet(ctx context.Context, username string) (*model.Wallet, error)
}

type WithdrawService struct {
	store WithdrawStore
}

func NewWithdrawService(store WithdrawStore) *WithdrawService {
	logger.Debug("Initializing WithdrawService")
	return &WithdrawService{store: store}
}

func (s *WithdrawService) DoWithdraw(ctx context.Context, tx *sql.Tx, username string, amount int64) (*model.Wallet, error) {
	funcName := "WithdrawService.DoWithdraw"
	logger.Info(fmt.Sprintf("%s - Params received", funcName), zap.String("username", username), zap.Int64("amount", amount))
	username, err := validation.SanitizeAndValidateUsername(username)
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("%s - Username sanitized", funcName), zap.String("username", username))

	if err := validation.ValidateAmount(amount); err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("%s - Amount validated", funcName), zap.Int64("amount", amount))

	currentWallet, err := s.store.FetchWallet(ctx, username)
	if err != nil {
		return nil, err
	}
	if currentWallet == nil {
		return nil, fmt.Errorf("username %s does not have a wallet", username)
	}
	if err := validation.ValidateWalletBalance(currentWallet.Balance - amount); err != nil {
		logger.Error(fmt.Sprintf("%s - Wallet balance validation failed", funcName), zap.String("error", err.Error()))
		return nil, fmt.Errorf("insufficient wallet balance - balance: %d - withdraw: %d - overdraft: %d", currentWallet.Balance, amount, currentWallet.Balance-amount)
	}
	logger.Info(fmt.Sprintf("%s - Wallet balance validated", funcName), []zap.Field{
		zap.Int64("wallet_balance", currentWallet.Balance),
		zap.Int64("payload_amount", amount),
		zap.Int64("resulting_balance", currentWallet.Balance-amount),
	}...)

	updatedWallet, err := s.store.WithdrawWallet(ctx, tx, username, amount)
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("%s - Withdrawn wallet", funcName), zap.Any("wallet", updatedWallet))
	return updatedWallet, nil
}
