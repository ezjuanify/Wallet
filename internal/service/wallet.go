package service

import (
	"context"
	"fmt"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/validation"
	"go.uber.org/zap"
)

type WalletStore interface {
	UpsertWallet(ctx context.Context, username string, amount int64) (*model.Wallet, error)
	FetchWallet(ctx context.Context, username string) (*model.Wallet, error)
	WithdrawWallet(ctx context.Context, username string, amount int64) (*model.Wallet, error)
}

type WalletService struct {
	store WalletStore
}

func NewWalletService(store WalletStore) *WalletService {
	logger.Debug("Initializing WalletService")
	return &WalletService{store: store}
}

func (ws *WalletService) DoDeposit(ctx context.Context, username string, amount int64) (*model.Wallet, error) {
	logger.Info("WalletService.DoDeposit - Params received", zap.String("username", username), zap.Int64("amount", amount))
	username, err := validation.SanitizeAndValidateUsername(username)
	if err != nil {
		return nil, err
	}
	logger.Info("WalletService.DoDeposit - Username sanitized", zap.String("username", username))

	if err := validation.ValidateAmount(amount); err != nil {
		return nil, err
	}
	logger.Info("WalletService.DoDeposit - Amount validated", zap.Int64("amount", amount))

	currentWallet, err := ws.store.FetchWallet(ctx, username)
	if err != nil {
		return nil, err
	}
	if currentWallet != nil {
		logger.Info("WalletService.DoDeposit - Validating if amount breaches upper limit", []zap.Field{
			zap.Int64("wallet_balance", currentWallet.Balance),
			zap.Int64("payload_amount", amount),
			zap.Int64("resulting_balance", currentWallet.Balance+amount),
		}...)
		if err := validation.ValidateWalletBalance(currentWallet.Balance + amount); err != nil {
			logger.Error("wallet balance validation failed", zap.String("error", err.Error()))
			return nil, fmt.Errorf("exceeds upper limit - wallet balance: %d  - deposit: %d - exceed: %d", currentWallet.Balance, amount, currentWallet.Balance+amount)
		}
	}

	updatedWallet, err := ws.store.UpsertWallet(ctx, username, amount)
	if err != nil {
		return nil, err
	}
	logger.Info("WalletService.DoDeposit - Upserted wallet", zap.Any("wallet", updatedWallet))
	return updatedWallet, nil
}

func (ws *WalletService) DoWithdraw(ctx context.Context, username string, amount int64) (*model.Wallet, error) {
	logger.Info("WalletService.DoWithdraw - Params received", zap.String("username", username), zap.Int64("amount", amount))
	username, err := validation.SanitizeAndValidateUsername(username)
	if err != nil {
		return nil, err
	}
	logger.Info("WalletService.DoWithdraw - Username sanitized", zap.String("username", username))

	if err := validation.ValidateAmount(amount); err != nil {
		return nil, err
	}
	logger.Info("WalletService.DoWithdraw - Amount validated", zap.Int64("amount", amount))

	currentWallet, err := ws.store.FetchWallet(ctx, username)
	if err != nil {
		return nil, err
	}
	if currentWallet == nil {
		return nil, fmt.Errorf("username %s does not have a wallet", username)
	}
	logger.Info("WalletService.DoWithdraw - Validating if amount breaches lower limit", []zap.Field{
		zap.Int64("wallet_balance", currentWallet.Balance),
		zap.Int64("payload_amount", amount),
		zap.Int64("resulting_balance", currentWallet.Balance-amount),
	}...)
	if err := validation.ValidateWalletBalance(currentWallet.Balance - amount); err != nil {
		logger.Error("wallet balance validation failed", zap.String("error", err.Error()))
		return nil, fmt.Errorf("insufficient wallet balance - balance: %d - withdraw: %d - overdraft: %d", currentWallet.Balance, amount, currentWallet.Balance-amount)
	}

	updatedWallet, err := ws.store.WithdrawWallet(ctx, username, amount)
	if err != nil {
		return nil, err
	}
	logger.Info("WalletService.DoWithdraw - Upserted wallet", zap.Any("wallet", updatedWallet))
	return updatedWallet, nil
}
