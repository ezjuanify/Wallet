package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
)

type mockWithdrawStore struct {
	wallets map[string]model.Wallet
}

func (m *mockWithdrawStore) initializeMockWallet() {
	m.wallets = map[string]model.Wallet{
		"JUAN": {
			Username: "JUAN",
			Balance:  2000,
		},
		"J_U_A_N": {
			Username: "J_U_A_N",
			Balance:  7000,
		},
		"J123": {
			Username: "J123",
			Balance:  5000,
		},
		"J_123": {
			Username: "J_123",
			Balance:  999999,
		},
	}
}

func (m *mockWithdrawStore) WithdrawWallet(ctx context.Context, username string, amount int64) (*model.Wallet, error) {
	currentTimestamp := time.Now().UTC()
	w, ok := m.wallets[username]
	if !ok {
		return nil, fmt.Errorf("Test Withdraw - No wallet found")
	}
	return &model.Wallet{
		Username:            w.Username,
		Balance:             w.Balance - amount,
		LastWithdrawAmount:  &amount,
		LastWithdrawUpdated: &currentTimestamp,
	}, nil
}

func (m *mockWithdrawStore) FetchWallet(ctx context.Context, username string) (*model.Wallet, error) {
	w, ok := m.wallets[username]
	if !ok {
		return nil, nil
	}
	return &model.Wallet{
		Username: w.Username,
		Balance:  w.Balance,
	}, nil
}

func TestDoWithdraw(t *testing.T) {
	logger.InitLogger()
	defer logger.Sync()

	type testCase struct {
		name           string
		username       string
		amount         int64
		expectedWallet *model.Wallet
		expectErr      bool
	}

	tests := []testCase{
		{
			name:     "Successful Withdraw - Normal",
			username: "JUAN",
			amount:   1000,
			expectedWallet: &model.Wallet{
				Username: "JUAN",
				Balance:  1000,
			},
			expectErr: false,
		},
		{
			name:     "Successful Withdraw - Username with lowercase",
			username: "juan",
			amount:   500,
			expectedWallet: &model.Wallet{
				Username: "JUAN",
				Balance:  1500,
			},
			expectErr: false,
		},
		{
			name:     "Successful Withdraw - Username with underscores",
			username: "j_u_a_n",
			amount:   800,
			expectedWallet: &model.Wallet{
				Username: "J_U_A_N",
				Balance:  6200,
			},
			expectErr: false,
		},
		{
			name:     "Successful Withdraw - Remaining balance 1",
			username: "J123",
			amount:   4999,
			expectedWallet: &model.Wallet{
				Username: "J123",
				Balance:  1,
			},
			expectErr: false,
		},
		{
			name:     "Successful Withdraw - Withdraw everything",
			username: "J_123",
			amount:   999999,
			expectedWallet: &model.Wallet{
				Username: "J_123",
				Balance:  0,
			},
			expectErr: false,
		},
		{
			name:     "Successful Withdraw - Username with space padding",
			username: " JUAN ",
			amount:   500,
			expectedWallet: &model.Wallet{
				Username: "JUAN",
				Balance:  1500,
			},
			expectErr: false,
		},

		{
			name:           "Failed Withdraw - Wallet not found",
			username:       "G12345",
			amount:         1000,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:           "Failed Withdraw - Username with illegal special characters",
			username:       "J@123",
			amount:         200,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:           "Failed Withdraw - Overdraft amount",
			username:       "J123",
			amount:         6000,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:           "Failed Withdraw - Negative amount",
			username:       "juan",
			amount:         -200,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:           "Failed Withdraw - Zero amount",
			username:       "juan",
			amount:         0,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:           "Failed Withdraw - Amount exceeding limit",
			username:       "juan",
			amount:         1000000,
			expectedWallet: nil,
			expectErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := &mockWithdrawStore{}
			mock.initializeMockWallet()
			s := &WithdrawService{store: mock}

			actual, err := s.DoWithdraw(context.Background(), test.username, test.amount)

			if test.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}

			if !test.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !test.expectErr && actual != nil && test.expectedWallet.Username != actual.Username {
				t.Errorf("expected username %s but got %s instead", test.expectedWallet.Username, actual.Username)
			}

			if !test.expectErr && actual != nil && test.expectedWallet.Balance != actual.Balance {
				t.Errorf("expected balance %d but got %d instead", test.expectedWallet.Balance, actual.Balance)
			}

			if !test.expectErr && actual != nil && test.amount != *actual.LastWithdrawAmount {
				t.Errorf("expected lastDepositAmount %d but got %d instead", test.amount, actual.LastWithdrawAmount)
			}

			if !test.expectErr && actual != nil && test.amount > 0 && actual.LastWithdrawUpdated == nil {
				t.Errorf("expected lastDepositUpdated to contain timestamp but got nil instead")
			}
		})
	}
}
