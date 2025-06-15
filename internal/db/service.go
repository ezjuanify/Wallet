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
		VALUES ($1, $2, $3, $4, $5, $6)
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
