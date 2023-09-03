package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gonozov0/go-musthave-devops/internal/server/internal/storage"
)

// UpdateMetricsHandler is a struct that holds the repository to update metrics.
type UpdateMetricsHandler struct {
	Repo storage.Repository
}

// NewUpdateMetricsHandler constructs a new UpdateMetricsHandler.
func NewUpdateMetricsHandler(repo storage.Repository) *UpdateMetricsHandler {
	return &UpdateMetricsHandler{
		Repo: repo,
	}
}

// UpdateMetrics is the HTTP handler for updating metrics.
func (h *UpdateMetricsHandler) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	metricType := parts[2]
	metricName := parts[3]
	valueStr := parts[4]

	var err error

	switch metricType {
	case storage.Gauge:
		var value float64
		value, err = strconv.ParseFloat(valueStr, 64)
		if err != nil {
			http.Error(w, "Invalid float value", http.StatusBadRequest)
			return
		}
		err = h.Repo.UpdateGauge(metricName, value)
	case storage.Counter:
		var value int64
		value, err = strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid integer value", http.StatusBadRequest)
			return
		}
		err = h.Repo.UpdateCounter(metricName, value)
	default:
		http.Error(w, "Unknown metric type", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Error updating metric: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
