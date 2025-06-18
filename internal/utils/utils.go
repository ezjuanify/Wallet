package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ezjuanify/wallet/internal/db"
	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model/request"
	"go.uber.org/zap"
)

func PtrInt64(v int64) *int64 { return &v }

func GetPGConfig() (*db.PGConfig, error) {
	pgconfig := &db.PGConfig{
		Host: "localhost",
		Port: 5432,
		SSL:  "disable",
		DB:   "db_wallet_app",
		User: "db_wallet_app",
		Pass: "db_wallet_app",
	}

	env := func(key string) string { return os.Getenv(key) }

	logger.Debug("Loading env overrides for pgconfig",
		zap.String("PG_HOST", env("PG_HOST")),
		zap.String("PG_PORT", env("PG_PORT")),
		zap.String("PG_SSL", env("PG_SSL")),
		zap.String("PG_DB", env("PG_DB")),
		zap.String("PG_USER", env("PG_USER")),
	)

	if val := env("PG_HOST"); val != "" {
		pgconfig.Host = val
	}

	if val := env("PG_PORT"); val != "" {
		pg_port, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		pgconfig.Port = pg_port
	}

	if val := env("PG_SSL"); val != "" {
		pgconfig.SSL = val
	}

	if val := env("PG_DB"); val != "" {
		pgconfig.DB = val
	}

	if val := env("PG_USER"); val != "" {
		pgconfig.User = val
	}

	if val := env("PG_PASS"); val != "" {
		pgconfig.Pass = val
	}

	logger.Debug("Final pgconfig built",
		zap.String("host", pgconfig.Host),
		zap.Int64("port", pgconfig.Port),
		zap.String("db", pgconfig.DB),
		zap.String("user", pgconfig.User),
		zap.String("ssl", pgconfig.SSL),
	)

	return pgconfig, nil
}

func DecodeRequest(r *http.Request) (*request.RequestPayload, error) {
	var req *request.RequestPayload
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func GenerateTransactionHash(txUser string, txType string, txAmount int64, txCounterparty *string, timestamp string) string {
	var counterparty string
	if txCounterparty != nil {
		counterparty = *txCounterparty
	}
	logger.Debug("Hashing with values",
		zap.String("username", txUser),
		zap.String("type", string(txType)),
		zap.Int64("amount", txAmount),
		zap.String("counterparty", counterparty),
		zap.String("timestamp", timestamp),
	)
	raw := fmt.Sprintf("%s|%s|%d|%s|%s", txUser, txType, txAmount, counterparty, timestamp)
	logger.Debug("Hashing string", zap.String("raw", raw))
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}
