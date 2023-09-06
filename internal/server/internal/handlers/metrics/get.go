package metrics

import (
	"encoding/binary"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/gonozov0/go-musthave-devops/internal/server/internal/storage"
	"net/http"
)

// GetMetric is the HTTP handler for getting metrics.
func (h *Handler) GetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	var err error
	var value interface{}

	switch metricType {
	case Gauge:
		value, err = h.Repo.GetGauge(metricName)
	case Counter:
		value, err = h.Repo.GetCounter(metricName)
	default:
		// Must be 400, return 501 because of autotests.
		http.Error(w, "Unknown metric type", http.StatusNotImplemented)
		return
	}

	if err != nil {
		if errors.Is(err, storage.MetricNotFoundError) {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = binary.Write(w, binary.LittleEndian, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
