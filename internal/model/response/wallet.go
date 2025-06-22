package response

import "github.com/ezjuanify/wallet/internal/model"

type TransactionResponse struct {
	Status          int64         `json:"status"`
	TransactionType model.TxnType `json:"action,omitempty"`
	Wallet          model.Wallet  `json:"wallet"`
	Counterparty    *string       `json:"counterparty,omitempty"`
}

type TransactionQueryResponse struct {
	Status       int64               `json:"status"`
	Criteria     *model.Criteria     `json:"criteria"`
	Transactions []model.Transaction `json:"transactions"`
}
