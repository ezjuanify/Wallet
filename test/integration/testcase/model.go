package testcase

import (
	"github.com/ezjuanify/wallet/internal/model"
)

type RequestPayload struct {
	Username     string  `json:"username"`
	Amount       int64   `json:"amount"`
	Counterparty *string `json:"counterparty,omitempty"`
}

type TestCase struct {
	Name           string
	TxnType        model.TransactionType
	InitialWallet  *model.Wallet
	Payload        *RequestPayload
	ExpectedWallet *model.Wallet
	ExpectErr      bool
}
