package integration

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ezjuanify/wallet/internal/db"
	"github.com/ezjuanify/wallet/internal/handler"
	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/service"
)

const (
	TEST_WALLET_HOST = "localhost"
	TEST_WALLET_PORT = ":8081"
)

var (
	dbTestHarness *DBTestHarness
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	defer logger.Sync()

	startTestDB()

	pgconfig := &db.PGConfig{
		Host: PG_HOST,
		Port: PG_PORT,
		SSL:  PG_SSL,
		DB:   PG_DB,
		User: PG_USER,
		Pass: PG_PASS,
	}
	store, err := db.NewStore(pgconfig)
	if err != nil {
		log.Fatalf("Failed to connect to test DB: %v\n", err)
	}

	ds := service.NewDepositService(store)
	ws := service.NewWithdrawService(store)
	ts := service.NewTransactionService(store)
	wh := handler.NewWalletHandler(store, ds, ws, ts)
	dbTestHarness = NewDbHarness(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/deposit", wh.DepositHandler)
	mux.HandleFunc("/withdraw", wh.WithdrawHandler)
	mux.HandleFunc("/transfer", wh.TransferHandler)

	go func() {
		log.Printf("Integration server starting on :%s\n", TEST_WALLET_PORT)
		if err := http.ListenAndServe(TEST_WALLET_PORT, mux); err != nil {
			log.Fatalf("Failed to start integration server: %v", err)
		}
	}()

	if err := dbTestHarness.waitForDB(10, 2*time.Second); err != nil {
		log.Fatalf("DB timeout: %v", err)
	}

	if err := waitForHTTPServerReady(10, 2*time.Second, TEST_WALLET_HOST, TEST_WALLET_PORT); err != nil {
		log.Fatalf("API timeout: %v", err)
	}

	code := m.Run()
	stopTestDB()
	os.Exit(code)
}
