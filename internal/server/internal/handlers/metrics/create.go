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
	metricValue := chi.URLParam(r, "metricValue")

	var err error

	switch metricType {
	case Gauge:
		var value float64
		value, err = strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Invalid float metricValue", http.StatusBadRequest)
			return
		}
		err = h.Repo.CreateGauge(metricName, value)
	case Counter:
		var value int64
		value, err = strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid integer metricValue", http.StatusBadRequest)
			return
		}
		err = h.Repo.CreateCounter(metricName, value)
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
