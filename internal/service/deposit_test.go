package service

import (
	"context"
	"testing"
	"time"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
)

type mockDepositStore struct {
	wallets map[string]model.Wallet
}

func (m *mockDepositStore) initializeMockWallet() {
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

func (m *mockDepositStore) UpsertWallet(ctx context.Context, username string, amount int64) (*model.Wallet, error) {
	currentTimestamp := time.Now().UTC()
	w, ok := m.wallets[username]
	if !ok {
		return &model.Wallet{
			Username:           username,
			Balance:            amount,
			LastDepositAmount:  &amount,
			LastDepositUpdated: &currentTimestamp,
		}, nil
	}
	return &model.Wallet{
		Username:           w.Username,
		Balance:            w.Balance + amount,
		LastDepositAmount:  &amount,
		LastDepositUpdated: &currentTimestamp,
	}, nil
}

func (m *mockDepositStore) FetchWallet(ctx context.Context, username string) (*model.Wallet, error) {
	return &model.Wallet{
		Username: m.wallets[username].Username,
		Balance:  m.wallets[username].Balance,
	}, nil
}

func TestDoDeposit(t *testing.T) {
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
			name:     "Successful Deposit - Existing wallet",
			username: "JUAN",
			amount:   1000,
			expectedWallet: &model.Wallet{
				Username: "JUAN",
				Balance:  3000,
			},
			expectErr: false,
		},
		{
			name:     "Successful Deposit - Existing wallet with lowercase and underscore username",
			username: "j_u_a_n",
			amount:   750,
			expectedWallet: &model.Wallet{
				Username: "J_U_A_N",
				Balance:  7750,
			},
			expectErr: false,
		},
		{
			name:     "Successful Deposit - Existing wallet with username padded with spaces",
			username: " j123 ",
			amount:   1000,
			expectedWallet: &model.Wallet{
				Username: "J123",
				Balance:  6000,
			},
			expectErr: false,
		},
		{
			name:     "Successful Deposit - New wallet",
			username: "JUAN123",
			amount:   500,
			expectedWallet: &model.Wallet{
				Username: "JUAN123",
				Balance:  500,
			},
			expectErr: false,
		},
		{
			name:     "Successful Deposit - New wallet with lowercase and underscore username",
			username: "__j__123",
			amount:   500,
			expectedWallet: &model.Wallet{
				Username: "__J__123",
				Balance:  500,
			},
			expectErr: false,
		},
		{
			name:     "Successful Deposit - New wallet with username padded with spaces",
			username: " __juan__ ",
			amount:   500,
			expectedWallet: &model.Wallet{
				Username: "__JUAN__",
				Balance:  500,
			},
			expectErr: false,
		},
		{
			name:     "Successful Deposit - Amount at upper limit",
			username: "_J_",
			amount:   999999,
			expectedWallet: &model.Wallet{
				Username: "_J_",
				Balance:  999999,
			},
			expectErr: false,
		},
		{
			name:           "Failed Deposit - Invalid username with special character",
			username:       "J@123",
			amount:         200,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:           "Failed Deposit - Invalid amount with negative number",
			username:       "juan",
			amount:         -200,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:           "Failed Deposit - Invalid number with zero amount",
			username:       "juan",
			amount:         0,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:           "Failed Deposit - Breach wallet amount on existing wallet",
			username:       "J_123",
			amount:         1,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:           "Failed Deposit - Breach wallet amount on ne wallets",
			username:       "_J_test_",
			amount:         1000000,
			expectedWallet: nil,
			expectErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := &mockDepositStore{}
			mock.initializeMockWallet()
			s := &DepositService{store: mock}
			actual, err := s.DoDeposit(context.Background(), test.username, test.amount)

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

			if !test.expectErr && actual != nil && test.amount != *actual.LastDepositAmount {
				t.Errorf("expected lastDepositAmount %d but got %d instead", test.amount, actual.LastDepositAmount)
			}

			if !test.expectErr && actual != nil && test.amount > 0 && actual.LastDepositUpdated == nil {
				t.Errorf("expected lastDepositUpdated to contain timestamp but got nil instead")
			}
		})
	}
}
