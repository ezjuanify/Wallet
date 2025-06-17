package testcase

import (
	"github.com/ezjuanify/wallet/internal/handler/request"
	"github.com/ezjuanify/wallet/internal/model"
)

type TestCase struct {
	Name           string
	TxnType        model.TransactionType
	InitialWallet  *model.Wallet
	Payload        *request.RequestPayload
	ExpectedWallet *model.Wallet
	ExpectErr      bool
}
