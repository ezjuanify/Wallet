package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ezjuanify/wallet/internal/db"
	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/validation"
	"go.uber.org/zap"
)

type WalletService struct {
	store *db.Store
}

func NewWalletService(store *db.Store) *WalletService {
	return &WalletService{store: store}
}

func (s *WalletService) DoFetchWallet(ctx context.Context, username string) (*model.Wallet, *validation.WalletError) {
	fnName := "WalletService.DoFetchWallet"
	logger.Info(fmt.Sprintf("%s - Params received", fnName), zap.String("username", username))

	username = validation.SanitizeUsernameWithoutError(username)
	logger.Info(fmt.Sprintf("%s - Username sanitized", fnName), zap.String("username", username))

	wallet, err := s.store.FetchWallet(ctx, username)
	if err != nil {
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_FETCH_WALLET_FAILED,
			Message:   "Error while fetching wallet",
			Timestamp: time.Now().UTC(),
			Err:       err,
			Context: []zap.Field{
				zap.String("username", username),
			},
		}
	}

	if wallet == nil {
		return nil, &validation.WalletError{
			Name:      fnName,
			Code:      validation.ERR_WALLET_DOES_NOT_EXIST,
			Message:   "User does not have an existing wallet",
			Timestamp: time.Now().UTC(),
			Err:       nil,
			Context: []zap.Field{
				zap.String("username", username),
			},
		}
	}

	return wallet, nil
}
