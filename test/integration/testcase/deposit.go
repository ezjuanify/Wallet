package testcase

import (
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/request"
)

func AddDepositTestCases() []TestCase {
	return []TestCase{
		{
			Name:    "Integration Test: Successful Deposit - Empty initial wallet",
			TxnType: model.TypeDeposit,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             0,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   1000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:           "Integration Test: Successful Deposit - Normal",
			TxnType:        model.TypeDeposit,
			InitialWallets: nil,
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   500,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:    "Integration Test: Successful Deposit - Username case insensitivity",
			TxnType: model.TypeDeposit,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             100,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
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
			TxnType: model.TypeDeposit,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   500,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:           "Integration Test: Successful Deposit - Alphanumeric and underscore username",
			TxnType:        model.TypeDeposit,
			InitialWallets: nil,
			Payload: &request.RequestPayload{
				Username: "j_uan_123",
				Amount:   5000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:    "Integration Test: Successful Deposit - Reaches upper limit exactly",
			TxnType: model.TypeDeposit,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             888888,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   111111,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:           "Integration Test: Successful Deposit - Max allowed amount",
			TxnType:        model.TypeDeposit,
			InitialWallets: nil,
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   999999,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:           "Integration Test: Successful Deposit - Username with spaces",
			TxnType:        model.TypeDeposit,
			InitialWallets: nil,
			Payload: &request.RequestPayload{
				Username: " _ju_an_ ",
				Amount:   999999,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:           "Integration Test: Fail Deposit - Zero amount",
			TxnType:        model.TypeDeposit,
			InitialWallets: nil,
			Payload: &request.RequestPayload{
				Username: "j_123",
				Amount:   0,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:    "Integration Test: Fail Deposit - Negative payload amount",
			TxnType: model.TypeDeposit,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   -1000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:           "Integration Test: Fail Deposit - Invalid character username",
			TxnType:        model.TypeDeposit,
			InitialWallets: nil,
			Payload: &request.RequestPayload{
				Username: "j@uan",
				Amount:   500,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:           "Integration Test: Fail Deposit - Over limit amount",
			TxnType:        model.TypeDeposit,
			InitialWallets: nil,
			Payload: &request.RequestPayload{
				Username: "j_123",
				Amount:   1000000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:    "Integration Test: Fail Deposit - Exceeds limit in wallet",
			TxnType: model.TypeDeposit,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             900000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: "juan",
				Amount:   100000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:    "Integration Test: Fail Deposit - Spaces in username",
			TxnType: model.TypeWithdraw,
			InitialWallets: []model.Wallet{
				{
					Username:            "J_123",
					Balance:             5000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: "j_ 1 2 3",
				Amount:   500,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
	}
}
