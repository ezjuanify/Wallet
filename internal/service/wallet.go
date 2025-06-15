package service

import (
	"github.com/ezjuanify/wallet/internal/db"
)

type WalletService struct {
	store *db.Store
}

func NewWalletService(store *db.Store) *WalletService {
	return &WalletService{store: store}
}
