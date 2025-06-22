package response

import "github.com/ezjuanify/wallet/internal/model"

type TransactionResponse struct {
	Status          int           `json:"status"`
	TransactionType model.TxnType `json:"action"`
	Wallet          model.Wallet  `json:"wallet"`
	Counterparty    *string       `json:"counterparty,omitempty"`
}

type TransactionQueryResponse struct {
	Status       int                 `json:"status"`
	Criteria     *model.Criteria     `json:"criteria"`
	Transactions []model.Transaction `json:"transactions"`
}

type WalletResponse struct {
	Status  int            `json:"status"`
	Message *string        `json:"message,omitempty"`
	Wallet  *model.Wallet  `json:"wallet,omitempty"`
	Wallets []model.Wallet `json:"wallets,omitempty"`
}
