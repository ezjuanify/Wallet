package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ezjuanify/wallet/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	DB *sql.DB
}

type PGConfig struct {
	Host string
	Port int64
	SSL  string
	DB   string
	User string
	Pass string
}

func NewStore(pgconfig *PGConfig) (*Store, error) {
	dsn := fmt.Sprintf("user=%v password=%v host=%v port=%v database=%v sslmode=%v", pgconfig.User, pgconfig.Pass, pgconfig.Host, pgconfig.Port, pgconfig.DB, pgconfig.SSL)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func (s *Store) InsertTransaction(ctx context.Context, txn model.Transaction) error {
	query := `
		INSERT INTO transactions (username, type, amount, counterparty, timestamp, hash) 
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := s.DB.ExecContext(
		ctx,
		query,
		txn.Username,
		txn.Type,
		txn.Amount,
		txn.Counterparty,
		txn.Timestamp,
		txn.Hash,
	)
	return err
}

func (s *Store) FetchWallet(ctx context.Context, username string) (*model.Wallet, error) {
	query := `
		SELECT username, balance, last_deposit_amount, last_deposit_updated, last_withdraw_amount, last_withdraw_updated
		FROM wallets
		WHERE username = $1;
	`

	row := s.DB.QueryRowContext(ctx, query, username)

	var wallet model.Wallet
	err := row.Scan(
		&wallet.Username,
		&wallet.Balance,
		&wallet.LastDepositAmount,
		&wallet.LastDepositUpdated,
		&wallet.LastWithdrawAmount,
		&wallet.LastWithdrawUpdated,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (s *Store) UpsertWallet(ctx context.Context, username string, amount int64) (*model.Wallet, error) {
	query := `
		INSERT INTO wallets (username, balance, last_deposit_amount, last_deposit_updated)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (username)
		DO UPDATE SET 
			balance = wallets.balance + EXCLUDED.balance,
			last_deposit_amount = EXCLUDED.last_deposit_amount,
			last_deposit_updated = now()
		RETURNING username, balance, last_deposit_amount, last_deposit_updated, last_withdraw_amount, last_withdraw_updated;
	`

	var wallet model.Wallet
	err := s.DB.QueryRowContext(
		ctx,
		query,
		username,
		amount,
		amount,
	).Scan(
		&wallet.Username,
		&wallet.Balance,
		&wallet.LastDepositAmount,
		&wallet.LastDepositUpdated,
		&wallet.LastWithdrawAmount,
		&wallet.LastWithdrawUpdated,
	)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}
