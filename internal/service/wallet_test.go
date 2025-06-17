package service

import (
	"context"
	"testing"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
)

type mockStore struct{}

func (m *mockStore) UpsertWallet(ctx context.Context, username string, amount int64) (*model.Wallet, error) {
	return &model.Wallet{
		Username: username,
		Balance:  amount,
	}, nil
}

func (m *mockStore) FetchWallet(ctx context.Context, username string) (*model.Wallet, error) {
	return nil, nil
}

func (m *mockStore) WithdrawWallet(ctx context.Context, username string, amount int64) (*model.Wallet, error) {
	return &model.Wallet{
		Username: username,
		Balance:  amount,
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
			name:     "Successful deposit",
			username: "JUAN",
			amount:   1000,
			expectedWallet: &model.Wallet{
				Username: "JUAN",
				Balance:  1000,
			},
			expectErr: false,
		},
		{
			name:     "Successful alphanumeric username",
			username: "JUAN123",
			amount:   500,
			expectedWallet: &model.Wallet{
				Username: "JUAN123",
				Balance:  500,
			},
			expectErr: false,
		},
		{
			name:     "Valid username with lower case and underscore",
			username: "j_ua_n",
			amount:   750,
			expectedWallet: &model.Wallet{
				Username: "J_UA_N",
				Balance:  750,
			},
			expectErr: false,
		},
		{
			name:           "Invalid username with special character",
			username:       "J@123",
			amount:         200,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:           "Invalid amount with negative number",
			username:       "juan",
			amount:         -200,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:           "Invalid number with zero amount",
			username:       "juan",
			amount:         0,
			expectedWallet: nil,
			expectErr:      true,
		},
		{
			name:     "Invalid number at 999999",
			username: "juan",
			amount:   999999,
			expectedWallet: &model.Wallet{
				Username: "JUAN",
				Balance:  999999,
			},
			expectErr: false,
		},
		{
			name:           "Invalid number at one million",
			username:       "juan",
			amount:         1000000,
			expectedWallet: nil,
			expectErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock := &mockStore{}
			ws := &WalletService{store: mock}
			actual, err := ws.DoDeposit(context.Background(), test.username, test.amount)

			if test.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}

			if !test.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !test.expectErr && actual != nil && actual.Username != test.expectedWallet.Username {
				t.Errorf("expected username %s but got %s instead", test.expectedWallet.Username, actual.Username)
			}

			if !test.expectErr && actual != nil && actual.Balance != test.expectedWallet.Balance {
				t.Errorf("expected balance %d but got %d instead", test.expectedWallet.Balance, actual.Balance)
			}
		})
	}
}
