package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ezjuanify/wallet/internal/handler/response"
	"github.com/ezjuanify/wallet/internal/model"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := &response.HealthResponse{Status: model.StatusHealthy}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
