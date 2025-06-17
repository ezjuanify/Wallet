package service

import (
	"context"

	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/validation"
)

type WalletStore interface {
	UpsertWallet(ctx context.Context, username string, amount int64) (*model.Wallet, error)
	FetchWallet(ctx context.Context, username string) (*model.Wallet, error)
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

	currentWallet, err := ws.store.FetchWallet(ctx, username)
	if err != nil {
		return nil, err
	}
	if currentWallet != nil {
		if err := validation.ValidateAmount(currentWallet.Balance + amount); err != nil {
			return nil, err
		}
	}

	updatedWallet, err := ws.store.UpsertWallet(ctx, username, amount)
	if err != nil {
		return nil, err
	}
	return updatedWallet, nil
}
