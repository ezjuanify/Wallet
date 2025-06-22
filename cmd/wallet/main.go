package main

import (
	"time"

	"github.com/ezjuanify/wallet/internal/appserv"
	"github.com/ezjuanify/wallet/internal/db"
	"github.com/ezjuanify/wallet/internal/handler"
	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/service"
	"github.com/ezjuanify/wallet/internal/utils"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.Sync()

	start := time.Now()
	defer func() {
		logger.Info("Startup complete", zap.Duration("took", time.Since(start)))
	}()

	pgconfig, err := utils.GetPGConfig()
	if err != nil {
		logger.Warn("Failed to get DB config, falling back to default config", zap.String("error", err.Error()))
		logger.Debug("using default DB config", zap.Any("pgconfig", pgconfig.Redacted()))
	}
	logger.Info("Successfully fetched DB config")

	dbFields := []zap.Field{
		zap.String("host", pgconfig.Host),
		zap.Int64("port", pgconfig.Port),
		zap.String("database", pgconfig.DB),
		zap.String("ssl", pgconfig.SSL),
	}
	store, err := db.NewStore(pgconfig)
	if err != nil {
		logger.Fatal("Failed to establish connection with DB", zap.String("error", err.Error()))
	}
	logger.Info("Successfully connected to DB", dbFields...)

	s := service.NewWalletService(store)
	ds := service.NewDepositService(store)
	ws := service.NewWithdrawService(store)
	ts := service.NewTransactionService(store)
	wh := handler.NewWalletHandler(store, s, ds, ws, ts)
	logger.Info("All services initialized")

	ap := appserv.NewAppServer()
	logger.Debug("Attaching HealthHandler")
	ap.Mux.HandleFunc(appserv.HEALTH, handler.HealthHandler)
	logger.Debug("Attaching DepositHandler")
	ap.Mux.HandleFunc(appserv.DEPOSIT, wh.DepositHandler)
	logger.Debug("Attaching WithdrawHandler")
	ap.Mux.HandleFunc(appserv.WITHDRAW, wh.WithdrawHandler)
	logger.Debug("Attaching TransferHandler")
	ap.Mux.HandleFunc(appserv.TRANSFER, wh.TransferHandler)
	logger.Debug("Attaching TransactionHandler")
	ap.Mux.HandleFunc(appserv.TRANSACTION, wh.TransactionHandler)
	logger.Debug("Attaching BalanceHandler")
	ap.Mux.HandleFunc(appserv.BALANCE, wh.BalanceHandler)
	logger.Info("All API handlers attached")

	if err := ap.GetEnvPort(); err != nil {
		logger.Warn("Failed to get app port from env variables, falling back to default port")
		logger.Debug("Using default port due to env fallback", zap.Int("app_port", ap.Port))
	}
	logger.Info("Server port set", zap.Int("port", ap.Port))

	logger.Info("Starting Wallet API", zap.Int("port", ap.Port))
	if err := ap.StartServer(); err != nil {
		logger.Fatal("Wallet API failed to start", zap.String("error", err.Error()))
	}
}
