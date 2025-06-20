package testcase

import (
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/request"
)

func AddWithdrawTestCases() []TestCase {
	return []TestCase{
		{
			Name:    "Integration Test: Successful Withdraw - Normal empty account",
			TxnType: model.TransactionType(model.TypeWithdraw),
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
				Amount:   1000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:    "Integration Test: Successful Withdraw - Remainder balance with lowercase alphanumeric and underscore characters",
			TxnType: model.TransactionType(model.TypeWithdraw),
			InitialWallets: []model.Wallet{
				{
					Username:            "_JUAN123_",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: "_juan123_",
				Amount:   500,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:    "Integration Test: Successful Withdraw - Username with space padding",
			TxnType: model.TransactionType(model.TypeWithdraw),
			InitialWallets: []model.Wallet{
				{
					Username:            "_J_U_L5_",
					Balance:             999999,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: " _j_u_l5_ ",
				Amount:   999999,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      false,
		},
		{
			Name:    "Integration Test: Fail Withdraw - Overdraft wallet balance",
			TxnType: model.TransactionType(model.TypeWithdraw),
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
				Amount:   1500,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:           "Integration Test: Fail Withdraw - No wallet found",
			TxnType:        model.TransactionType(model.TypeWithdraw),
			InitialWallets: nil,
			Payload: &request.RequestPayload{
				Username: "_test_J_",
				Amount:   5000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:    "Integration Test: Fail Withdraw - Zero amount",
			TxnType: model.TransactionType(model.TypeWithdraw),
			InitialWallets: []model.Wallet{
				{
					Username:            "J_123",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: "j_123",
				Amount:   0,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:    "Integration Test: Fail Withdraw - Negative payload amount",
			TxnType: model.TransactionType(model.TypeWithdraw),
			InitialWallets: []model.Wallet{
				{
					Username:            "J_123",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: "j_123",
				Amount:   -1000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:    "Integration Test: Fail Deposit - Invalid character username",
			TxnType: model.TransactionType(model.TypeWithdraw),
			InitialWallets: []model.Wallet{
				{
					Username:            "J_123",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: "j@uan",
				Amount:   500,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:    "Integration Test: Fail Deposit - Overdraft amount",
			TxnType: model.TransactionType(model.TypeWithdraw),
			InitialWallets: []model.Wallet{
				{
					Username:            "J_123",
					Balance:             1,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username: "j_123",
				Amount:   1000000,
			},
			ExpectedWallet: &model.Wallet{},
			ExpectErr:      true,
		},
		{
			Name:    "Integration Test: Fail Deposit - Spaces in username",
			TxnType: model.TransactionType(model.TypeWithdraw),
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
