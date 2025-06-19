package testcase

import (
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/request"
	"github.com/ezjuanify/wallet/internal/validation"
)

type TestCase struct {
	Name           string
	TxnType        model.TransactionType
	InitialWallet  *model.Wallet
	Payload        *request.RequestPayload
	ExpectedWallet *model.Wallet
	ExpectErr      bool
}

func (tc *TestCase) SanitizedPayloadUsername() string {
	sanitized, _ := validation.SanitizeAndValidateUsername(tc.Payload.Username)
	return sanitized
}
