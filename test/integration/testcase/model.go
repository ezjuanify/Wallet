package testcase

import (
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/request"
)

type TestCase struct {
	Name           string
	TxnType        model.TransactionType
	InitialWallet  *model.Wallet
	Payload        *request.RequestPayload
	ExpectedWallet *model.Wallet
	ExpectErr      bool
}
