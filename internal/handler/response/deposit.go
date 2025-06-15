package response

import "github.com/ezjuanify/wallet/internal/model"

type DepositResponse struct {
	Status          int64                 `json:"status"`
	TransactionType model.TransactionType `json:"action"`
	Wallet          model.Wallet          `json:"wallet"`
}
