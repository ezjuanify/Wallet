package response

import (
	"github.com/ezjuanify/wallet/internal/model"
)

type HealthResponse struct {
	Status model.HealthStatus `json:"status"`
}
