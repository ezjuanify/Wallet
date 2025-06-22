package handler

import (
	"fmt"
	"net/http"

	"github.com/ezjuanify/wallet/internal/logger"
	"github.com/ezjuanify/wallet/internal/model"
	"github.com/ezjuanify/wallet/internal/model/response"
	"go.uber.org/zap"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fnName := "HealthHandler"
	resp := &response.HealthResponse{Status: model.StatusHealthy}
	logger.Info(fmt.Sprintf("%s - Health request received", fnName), zap.Any("response", resp))
	SendJSONResponse(fnName, w, http.StatusOK, resp)
}
