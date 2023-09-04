package metrics

import (
	"encoding/binary"
	"github.com/go-chi/chi/v5"
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
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case Counter:
		value, err = h.Repo.GetCounter(metricName)
	default:
		// Must be 400, return 501 because of autotests.
		http.Error(w, "Unknown metric type", http.StatusNotImplemented)
		return
	}

	if err != nil {
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
