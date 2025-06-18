package service

import (
	"context"
	"fmt"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/validation"
	"go.uber.org/zap"
)

type DepositStore interface {
	UpsertWallet(ctx context.Context, username string, amount int64) (*model.Wallet, error)
	FetchWallet(ctx context.Context, username string) (*model.Wallet, error)
}

type DepositService struct {
	store DepositStore
}

func NewDepositService(store DepositStore) *DepositService {
	logger.Debug("Initializing DepositService")
	return &DepositService{store: store}
}

func (s *DepositService) DoDeposit(ctx context.Context, username string, amount int64) (*model.Wallet, error) {
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
	logger.Info(fmt.Sprintf("%s - Amount validated", funcName), zap.Int64("amount", amount))

	currentWallet, err := s.store.FetchWallet(ctx, username)
	if err != nil {
		return nil, err
	}
	if currentWallet != nil {
		logger.Info(fmt.Sprintf("%s - Validating if amount breaches upper limit", funcName), []zap.Field{
			zap.Int64("wallet_balance", currentWallet.Balance),
			zap.Int64("payload_amount", amount),
			zap.Int64("resulting_balance", currentWallet.Balance+amount),
		}...)
		if err := validation.ValidateWalletBalance(currentWallet.Balance + amount); err != nil {
			logger.Error(fmt.Sprintf("%s - Wallet balance validation failed", funcName), zap.String("error", err.Error()))
			return nil, fmt.Errorf("exceeds upper limit - wallet balance: %d  - deposit: %d - exceed: %d", currentWallet.Balance, amount, currentWallet.Balance+amount)
		}
	}

	updatedWallet, err := s.store.UpsertWallet(ctx, username, amount)
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("%s - Upserted wallet", funcName), zap.Any("wallet", updatedWallet))
	return updatedWallet, nil
}
