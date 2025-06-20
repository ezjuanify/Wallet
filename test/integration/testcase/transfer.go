package testcase

import (
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/request"
	"github.com/ezjuanify/wallet/internal/utils"
)

func AddTransferTestCases() []TestCase {
	return []TestCase{
		{
			Name:    "Integration Test: Successful Transfer - Normal existing account",
			TxnType: model.TypeTransfer,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
				{
					Username:            "MARY",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username:     "juan",
				Amount:       500,
				Counterparty: utils.Ptr("mary"),
			},
			ExpectErr: false,
		},
		{
			Name:    "Integration Test: Successful Transfer - Full transfer",
			TxnType: model.TypeTransfer,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
				{
					Username:            "MARY",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username:     "juan",
				Amount:       1000,
				Counterparty: utils.Ptr("mary"),
			},
			ExpectErr: false,
		},
		{
			Name:    "Integration Test: Fail Transfer - Negative amount",
			TxnType: model.TypeTransfer,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
				{
					Username:            "MARY",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username:     "juan",
				Amount:       -1000,
				Counterparty: utils.Ptr("mary"),
			},
			ExpectErr: true,
		},
		{
			Name:    "Integration Test: Fail Transfer - Insufficient funds from user wallet",
			TxnType: model.TypeTransfer,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
				{
					Username:            "MARY",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username:     "juan",
				Amount:       2000,
				Counterparty: utils.Ptr("mary"),
			},
			ExpectErr: true,
		},
		{
			Name:    "Integration Test: Fail Transfer - Transfer over limit",
			TxnType: model.TypeTransfer,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
				{
					Username:            "MARY",
					Balance:             999999,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username:     "juan",
				Amount:       2000,
				Counterparty: utils.Ptr("mary"),
			},
			ExpectErr: true,
		},
		{
			Name:    "Integration Test: Fail Transfer - User wallet does not exist",
			TxnType: model.TypeTransfer,
			InitialWallets: []model.Wallet{
				{
					Username:            "MARY",
					Balance:             5000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username:     "juan",
				Amount:       2000,
				Counterparty: utils.Ptr("mary"),
			},
			ExpectErr: true,
		},
		{
			Name:    "Integration Test: Fail Transfer - Counterparty wallet does not exist",
			TxnType: model.TypeTransfer,
			InitialWallets: []model.Wallet{
				{
					Username:            "JUAN",
					Balance:             5000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username:     "juan",
				Amount:       2000,
				Counterparty: utils.Ptr("mary"),
			},
			ExpectErr: true,
		},
		{
			Name:           "Integration Test: Fail Transfer - No wallets",
			TxnType:        model.TypeTransfer,
			InitialWallets: nil,
			Payload: &request.RequestPayload{
				Username:     "juan",
				Amount:       2000,
				Counterparty: utils.Ptr("mary"),
			},
			ExpectErr: true,
		},
		{
			Name:    "Integration Test: Fail Transfer - Symbol in username",
			TxnType: model.TypeTransfer,
			InitialWallets: []model.Wallet{
				{
					Username:            "JAN",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
				{
					Username:            "MARY",
					Balance:             999999,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username:     "j@n",
				Amount:       2000,
				Counterparty: utils.Ptr("mary"),
			},
			ExpectErr: true,
		},
		{
			Name:    "Integration Test: Fail Transfer - Symbol in counterparty",
			TxnType: model.TypeTransfer,
			InitialWallets: []model.Wallet{
				{
					Username:            "JAN",
					Balance:             1000,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
				{
					Username:            "MARY",
					Balance:             999999,
					LastDepositAmount:   nil,
					LastDepositUpdated:  nil,
					LastWithdrawAmount:  nil,
					LastWithdrawUpdated: nil,
				},
			},
			Payload: &request.RequestPayload{
				Username:     "juan",
				Amount:       2000,
				Counterparty: utils.Ptr("m@ry"),
			},
			ExpectErr: true,
		},
	}
}
