package response

import "github.com/ezjuanify/wallet/internal/model"

type TransactionResponse struct {
	Status          int64                 `json:"status"`
	TransactionType model.TransactionType `json:"action"`
	Wallet          model.Wallet          `json:"wallet"`
	Counterparty    *string               `json:"counterparty,omitempty"`
}
