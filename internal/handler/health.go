package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ezjuanify/wallet/internal/model"
)

type healthResponse struct {
	Status model.HealthStatus `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := healthResponse{Status: model.StatusHealthy}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
