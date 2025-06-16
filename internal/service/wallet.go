package service

import (
	"context"

	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/validation"
)

type WalletStore interface {
	UpsertWallet(ctx context.Context, username string, amount int64) (*model.Wallet, error)
}

type WalletService struct {
	store WalletStore
}

func NewWalletService(store WalletStore) *WalletService {
	return &WalletService{store: store}
}

func (ws *WalletService) DoDeposit(ctx context.Context, username string, amount int64) (*model.Wallet, error) {
	username, err := validation.SanitizeAndValidateUsername(username)
	if err != nil {
		return nil, err
	}

	if err := validation.ValidateAmount(amount); err != nil {
		return nil, err
	}

	wallet, err := ws.store.UpsertWallet(ctx, username, amount)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
