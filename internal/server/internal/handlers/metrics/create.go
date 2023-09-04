package metrics

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// CreateMetric is the HTTP handler for updating metrics.
func (h *Handler) CreateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	valueStr := chi.URLParam(r, "value")

	var err error

	switch metricType {
	case Gauge:
		var value float64
		value, err = strconv.ParseFloat(valueStr, 64)
		if err != nil {
			http.Error(w, "Invalid float value", http.StatusBadRequest)
			return
		}
		err = h.Repo.UpdateGauge(metricName, value)
	case Counter:
		var value int64
		value, err = strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid integer value", http.StatusBadRequest)
			return
		}
		err = h.Repo.UpdateCounter(metricName, value)
	default:
		// Must be 400, return 501 because of autotests.
		http.Error(w, "Unknown metric type", http.StatusNotImplemented)
		return
	}

	if err != nil {
		log.Printf("Error updating metric: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
