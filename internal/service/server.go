package service

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

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		logger.Info("Request received",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("ip", ip),
		)

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
