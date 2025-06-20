package testcase

import (
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/request"
	"github.com/ezjuanify/wallet/internal/validation"
)

type TestCase struct {
	Name                       string
	TxnType                    model.TxnType
	Payload                    *request.RequestPayload
	InitialWallets             []model.Wallet
	ExpectedWallet             *model.Wallet
	ExpectedCounterpartyWallet *model.Wallet
	ExpectErr                  bool
}

func (tc *TestCase) BuildExpectedWallet(username string, isCounterparty bool) {
	username, _ = validation.SanitizeAndValidateUsername(username)

	var expected *model.Wallet

	for i := range tc.InitialWallets {
		if tc.InitialWallets[i].Username == username {
			expected = &tc.InitialWallets[i]
			break
		}
	}

	if expected == nil {
		expected = &model.Wallet{
			Username: username,
			Balance:  0,
		}
	}

	switch tc.TxnType {
	case model.TypeDeposit:
		expected.Balance += tc.Payload.Amount
		expected.LastDepositAmount = &tc.Payload.Amount
	case model.TypeWithdraw:
		expected.Balance -= tc.Payload.Amount
		expected.LastWithdrawAmount = &tc.Payload.Amount
	case model.TypeTransfer:
		if isCounterparty {
			expected.Balance += tc.Payload.Amount
			expected.LastDepositAmount = &tc.Payload.Amount
		} else {
			expected.Balance -= tc.Payload.Amount
			expected.LastWithdrawAmount = &tc.Payload.Amount
		}
	}

	if isCounterparty {
		tc.ExpectedCounterpartyWallet = expected
		return
	}
	tc.ExpectedWallet = expected
}
