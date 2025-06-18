package testcase

import (
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/request"
)

func AddDepositTestCase() []TestCase {
	return []TestCase{
		{
			Name:    "Integration Test: Successful Deposit - Empty initial wallet",
			TxnType: model.TransactionType(model.TypeDeposit),
			InitialWallet: &model.Wallet{
				Username:            "JUAN",
				Balance:             0,
				LastDepositAmount:   nil,
				LastDepositUpdated:  nil,
				LastWithdrawAmount:  nil,
				LastWithdrawUpdated: nil,
			},
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   1000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:          "Integration Test: Successful Deposit - Normal",
			TxnType:       model.TransactionType(model.TypeDeposit),
			InitialWallet: nil,
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   500,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:    "Integration Test: Successful Deposit - Username case insensitivity",
			TxnType: model.TransactionType(model.TypeDeposit),
			InitialWallet: &model.Wallet{
				Username: "JUAN",
				Balance:  100,
			},
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   900,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:    "Integration Test: Successful Deposit - Existing wallet value",
			TxnType: model.TransactionType(model.TypeDeposit),
			InitialWallet: &model.Wallet{
				Username:            "JUAN",
				Balance:             1000,
				LastDepositAmount:   nil,
				LastDepositUpdated:  nil,
				LastWithdrawAmount:  nil,
				LastWithdrawUpdated: nil,
			},
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   500,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:          "Integration Test: Successful Deposit - Alphanumeric and underscore username",
			TxnType:       model.TransactionType(model.TypeDeposit),
			InitialWallet: nil,
			Payload: &request.RequestPayload{
				Username: "j_uan_123",
				Amount:   5000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:    "Integration Test: Successful Deposit - Reaches upper limit exactly",
			TxnType: model.TransactionType(model.TypeDeposit),
			InitialWallet: &model.Wallet{
				Username: "JUAN",
				Balance:  888888,
			},
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   111111,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:          "Integration Test: Successful Deposit - Max allowed amount",
			TxnType:       model.TransactionType(model.TypeDeposit),
			InitialWallet: nil,
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   999999,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:          "Integration Test: Successful Deposit - Zero amount",
			TxnType:       model.TransactionType(model.TypeDeposit),
			InitialWallet: nil,
			Payload: &request.RequestPayload{
				Username: "j_123",
				Amount:   0,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:    "Integration Test: Fail Deposit - Negative amount payload",
			TxnType: model.TransactionType(model.TypeDeposit),
			InitialWallet: &model.Wallet{
				Username:            "JUAN",
				Balance:             1000,
				LastDepositAmount:   nil,
				LastDepositUpdated:  nil,
				LastWithdrawAmount:  nil,
				LastWithdrawUpdated: nil,
			},
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   -1000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:          "Integration Test: Fail Deposit - Invalid character username",
			TxnType:       model.TransactionType(model.TypeDeposit),
			InitialWallet: nil,
			Payload: &request.RequestPayload{
				Username: "j@uan",
				Amount:   500,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:          "Integration Test: Fail Deposit - Over limit amount",
			TxnType:       model.TransactionType(model.TypeDeposit),
			InitialWallet: nil,
			Payload: &request.RequestPayload{
				Username: "j_123",
				Amount:   1000000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:    "Integration Test: Successful Deposit - Exceeds limit in wallet",
			TxnType: model.TransactionType(model.TypeDeposit),
			InitialWallet: &model.Wallet{
				Username: "JUAN",
				Balance:  900000,
			},
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   100000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
	}
}
