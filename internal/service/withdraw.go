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

func (s *WithdrawService) DoWithdraw(ctx context.Context, tx *sql.Tx, username string, amount int64) (*model.Wallet, *validation.WalletError) {
	fnName := "WithdrawService.DoWithdraw"
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
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_WALLET_DOES_NOT_EXIST,
			Message:   "No existing wallet found for user",
			Timestamp: time.Now().UTC(),
			Err:       nil,
			Context: []zap.Field{
				zap.String("username", username),
			},
		}
	}

	newBalance := currentWallet.Balance - amount
	if err := validation.ValidateWalletBalance(newBalance); err != nil {
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_WALLET_BALANCE_VALIDATION_FAILED,
			Message:   "Wallet balance would overdraft",
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

	updatedWallet, err := s.store.WithdrawWallet(ctx, tx, username, amount)
	if err != nil {
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_DB_WITHDRAW_FAILED,
			Message:   "Failed to withdraw fromm wallet",
			Timestamp: time.Now().UTC(),
			Err:       err,
			Context: []zap.Field{
				zap.String("username", username),
				zap.Int64("amount", amount),
			},
		}
	}
	logger.Info(fmt.Sprintf("%s - Withdrawn from wallet", fnName), zap.Any("wallet", updatedWallet))
	return updatedWallet, nil
}
