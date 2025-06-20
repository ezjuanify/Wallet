package integration

import (
	"fmt"
	"time"

	"github.com/ezjuanify/wallet/internal/db"
	"github.com/ezjuanify/wallet/internal/model"
)

type DBTestHarness struct {
	store *db.Store
}

func NewDbHarness(store *db.Store) *DBTestHarness {
	return &DBTestHarness{store: store}
}

func (h *DBTestHarness) waitForDB(maxRetries int, interval time.Duration) error {
	var err error

	for i := 0; i < maxRetries; i++ {
		err = h.store.DB.Ping()
		if err == nil {
			fmt.Println("✅ Database is up.")
			return nil
		}
		fmt.Printf("DB not ready yet (%d/%d): %v\n", i+1, maxRetries, err)
		time.Sleep(interval)
	}
	return fmt.Errorf("database not reachable after %d attempts: %v", maxRetries, err)
}

func (h *DBTestHarness) DoTestResetDBState() error {
	query := `
		TRUNCATE TABLE
			wallets,
			transactions
		RESTART IDENTITY 
		CASCADE;
	`

	_, err := h.store.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}
	return err
}

func (h *DBTestHarness) DoTestInsertInitialWallet(wallet *model.Wallet) error {
	query := `
		INSERT INTO wallets (username, balance, last_deposit_amount, last_deposit_updated, last_withdraw_amount, last_withdraw_updated)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := h.store.DB.Exec(
		query,
		wallet.Username,
		wallet.Balance,
		wallet.LastDepositAmount,
		wallet.LastDepositUpdated,
		wallet.LastWithdrawAmount,
		wallet.LastWithdrawUpdated,
	)
	return err
}

func (h *DBTestHarness) DoTestFetchWalletFromDB(username string) (*model.Wallet, error) {
	query := `
		SELECT username, balance, last_deposit_amount, last_deposit_updated, last_withdraw_amount, last_withdraw_updated
		FROM wallets
		WHERE username = $1;
	`

	row := h.store.DB.QueryRow(query, username)

	var w model.Wallet
	err := row.Scan(
		&w.Username,
		&w.Balance,
		&w.LastDepositAmount,
		&w.LastDepositUpdated,
		&w.LastWithdrawAmount,
		&w.LastWithdrawUpdated,
	)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (h *DBTestHarness) DoTestFetchTransaction(username string) (*model.Transaction, error) {
	query := `
		SELECT username, type, amount, counterparty, timestamp, hash
		FROM transactions
		WHERE username = $1
		ORDER BY timestamp DESC
		LIMIT 1;
	`

	row := h.store.DB.QueryRow(query, username)

	var t model.Transaction
	err := row.Scan(
		&t.Username,
		&t.TxnType,
		&t.Amount,
		&t.Counterparty,
		&t.Timestamp,
		&t.Hash,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
