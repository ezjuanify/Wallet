package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
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

func (cfg *PGConfig) RedactedDSN() string {
	return fmt.Sprintf("user=%v password=**** host=%v port=%v database=%v sslmode=%v",
		cfg.User, cfg.Host, cfg.Port, cfg.DB, cfg.SSL)
}

func (cfg *PGConfig) Redacted() []zap.Field {
	return []zap.Field{
		zap.String("Host", cfg.Host),
		zap.Int64("Port", cfg.Port),
		zap.String("SSL", cfg.SSL),
		zap.String("DB", cfg.DB),
		zap.String("User", cfg.User),
	}
}

func NewStore(pgconfig *PGConfig) (*Store, error) {
	logger.Debug("Initializing DB object with config",
		zap.String("user", pgconfig.User),
		zap.String("host", pgconfig.Host),
		zap.Int64("port", pgconfig.Port),
		zap.String("database", pgconfig.DB),
		zap.String("sslmode", pgconfig.SSL),
	)

	dsn := fmt.Sprintf("user=%v password=%v host=%v port=%v database=%v sslmode=%v", pgconfig.User, pgconfig.Pass, pgconfig.Host, pgconfig.Port, pgconfig.DB, pgconfig.SSL)

	logger.Debug("DSN string", zap.String("dsn", pgconfig.RedactedDSN()))

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func (s *Store) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return s.DB.BeginTx(ctx, nil)
}

func (s *Store) FetchTransaction(ctx context.Context, criteria *model.Criteria) ([]model.Transaction, error) {
	fnName := "DBStore.FetchTransaction"
	logger.Debug(fmt.Sprintf("%s - parameters", fnName), zap.Any("criteria", criteria))
	var (
		query      strings.Builder
		args       []interface{}
		conditions []string
		argPos     = 1
	)

	query.WriteString("SELECT id, username, type, amount, counterparty, timestamp, hash FROM transactions")

	if criteria.Username != "" {
		conditions = append(conditions, fmt.Sprintf("username = $%d", argPos))
		args = append(args, criteria.Username)
		argPos++
	}
	if criteria.Counterparty != "" {
		conditions = append(conditions, fmt.Sprintf("counterparty = $%d", argPos))
		args = append(args, criteria.Counterparty)
		argPos++
	}
	if criteria.TxnType != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argPos))
		args = append(args, criteria.TxnType)
		argPos++
	}
	if len(conditions) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(conditions, " AND "))
	}

	query.WriteString(" ORDER BY timestamp DESC")

	if criteria.Limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT $%d ", argPos))
		args = append(args, criteria.Limit)
		argPos++
	}

	logger.Info(fmt.Sprintf("%s - Query built", fnName), zap.String("query", query.String()))

	transactions := []model.Transaction{}

	rows, err := s.DB.QueryContext(
		ctx,
		query.String(),
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var txn model.Transaction
		err := rows.Scan(
			&txn.ID,
			&txn.Username,
			&txn.TxnType,
			&txn.Amount,
			&txn.Counterparty,
			&txn.Timestamp,
			&txn.Hash,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, txn)
	}

	logger.Debug(fmt.Sprintf("%s - Transactions found", fnName), zap.Any("transactions", transactions), zap.Int("count", len(transactions)))

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *Store) InsertTransaction(ctx context.Context, tx *sql.Tx, txn model.Transaction) error {
	fnName := "DBStore.InsertTransaction"
	logger.Debug(fmt.Sprintf("%s - parameters", fnName), zap.Any("transaction", txn))
	query := `
		INSERT INTO transactions (username, type, amount, counterparty, timestamp, hash) 
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	logger.Debug("InsertTransaction - query", zap.String("query", query))

	_, err := tx.ExecContext(
		ctx,
		query,
		txn.Username,
		txn.TxnType,
		txn.Amount,
		txn.Counterparty,
		txn.Timestamp,
		txn.Hash,
	)
	return err
}

func (s *Store) FetchWallet(ctx context.Context, username string) (*model.Wallet, error) {
	fnName := "DBStore.FetchWallet"
	logger.Debug(fmt.Sprintf("%s - parameters", fnName), zap.String("username", username))
	query := `
		SELECT username, balance, last_deposit_amount, last_deposit_updated, last_withdraw_amount, last_withdraw_updated
		FROM wallets
		WHERE username = $1;
	`
	logger.Debug(fmt.Sprintf("%s - query", fnName), zap.String("query", query))

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
			logger.Debug(fmt.Sprintf("%s - No wallet found for username", fnName), zap.String("username", username))
			return nil, nil
		}
		return nil, err
	}
	logger.Debug(fmt.Sprintf("%s - after query execution", fnName), zap.Any("wallet", wallet))
	return &wallet, nil
}

func (s *Store) FetchAllWallet(ctx context.Context) ([]model.Wallet, error) {
	fnName := "DBStore.FetchAllWallet"
	logger.Debug(fmt.Sprintf("%s - no params to receive", fnName))
	query := `
		SELECT username, balance, last_deposit_amount, last_deposit_updated, last_withdraw_amount, last_withdraw_updated
		FROM wallets;
	`
	logger.Debug(fmt.Sprintf("%s - query", fnName), zap.String("query", query))

	var wallets []model.Wallet

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var wallet model.Wallet
		err := rows.Scan(
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

		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func (s *Store) UpsertWallet(ctx context.Context, tx *sql.Tx, username string, amount int64) (*model.Wallet, error) {
	fnName := "DBStore.UpsertWallet"
	logger.Debug(fmt.Sprintf("%s - parameters", fnName), zap.String("username", username), zap.Int64("amount", amount))
	query := `
		INSERT INTO wallets (username, balance, last_deposit_amount, last_deposit_updated)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (username)
		DO UPDATE SET 
		balance              = wallets.balance + EXCLUDED.balance,
		last_deposit_amount  = EXCLUDED.last_deposit_amount,
		last_deposit_updated = now()
		RETURNING username, balance, last_deposit_amount, last_deposit_updated, last_withdraw_amount, last_withdraw_updated;
	`
	logger.Debug(fmt.Sprintf("%s - query", fnName), zap.String("query", query))

	var wallet model.Wallet
	err := tx.QueryRowContext(
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
	logger.Debug(fmt.Sprintf("%s - after query execution", fnName), zap.Any("wallet", wallet))
	return &wallet, nil
}

func (s *Store) WithdrawWallet(ctx context.Context, tx *sql.Tx, username string, amount int64) (*model.Wallet, error) {
	fnName := "DBStore.WithdrawWallet"
	logger.Debug(fmt.Sprintf("%s - parameters", fnName), zap.String("username", username), zap.Int64("amount", amount))
	query := `
		UPDATE wallets
		SET
			balance               = balance - $1,
			last_withdraw_amount  = $1,
			last_withdraw_updated = now()
		WHERE
			username = $2
		AND balance >= $1
		RETURNING username, balance, last_deposit_amount, last_deposit_updated, last_withdraw_amount, last_withdraw_updated;
	`
	logger.Debug(fmt.Sprintf("%s - query", fnName), zap.String("query", query))

	var wallet model.Wallet
	err := tx.QueryRowContext(
		ctx,
		query,
		amount,
		username,
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
	logger.Debug(fmt.Sprintf("%s - after query execution", fnName), zap.Any("wallet", wallet))
	return &wallet, nil
}
