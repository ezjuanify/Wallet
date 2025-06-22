package appserv

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ezjuanify/wallet/internal/logger"
	"go.uber.org/zap"
)

var (
	defaultPort = 8080
)

type AppServer struct {
	Mux  *http.ServeMux
	Port int
}

const (
	DEPOSIT        = "/deposit"
	WITHDRAW       = "/withdraw"
	TRANSFER       = "/transfer"
	HEALTH         = "/health"
	TRANSACTION    = "/transactions"
	BALANCE        = "/balance"
	ADMIN_BALANCES = "/admin/balances"
)

var POSTEndpoint = map[string]struct{}{
	DEPOSIT:  {},
	WITHDRAW: {},
	TRANSFER: {},
}

var GETEndpoint = map[string]struct{}{
	TRANSACTION:    {},
	HEALTH:         {},
	BALANCE:        {},
	ADMIN_BALANCES: {},
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		logger.Info("Request received",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("ip", ip),
		)

		switch r.Method {
		case http.MethodPost:
			if _, ok := POSTEndpoint[r.URL.Path]; !ok {
				logger.Warn(fmt.Sprintf("No %s method for %s endpoint", r.Method, r.URL.Path))
				http.Error(w, "Invalid POST endpoint", http.StatusNotFound)
				return
			}
		case http.MethodGet:
			if _, ok := GETEndpoint[r.URL.Path]; !ok {
				logger.Warn(fmt.Sprintf("No %s method for %s endpoint", r.Method, r.URL.Path))
				http.Error(w, "Invalid GET endpoint", http.StatusNotFound)
				return
			}
		default:
			logger.Error(fmt.Sprintf("%s method not allowed - %s", r.Method, r.URL.Path))
			http.Error(w, fmt.Sprintf("Request method %s not allowed", r.Method), http.StatusNotFound)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)

		logger.Info("Request completed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

func NewAppServer() *AppServer {
	logger.Debug("Initializing App Server")
	return &AppServer{
		Mux:  http.NewServeMux(),
		Port: defaultPort,
	}
}

func (s *AppServer) GetEnvPort() error {
	env := func(key string) string { return os.Getenv(key) }

	logger.Debug("Attempting to fetch env port", zap.String("APP_PORT", env("APP_PORT")))
	if val := env("APP_PORT"); val != "" {
		app_port, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		s.Port = app_port
	}
	return nil
}

func (s *AppServer) StartServer() error {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), requestLogger(s.Mux)); err != nil {
		return err
	}
	return nil
}
