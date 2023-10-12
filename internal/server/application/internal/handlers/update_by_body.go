package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gonozov0/go-musthave-devops/internal/shared"

	log "github.com/sirupsen/logrus"
)

func (h *Handler) UpdateMetricByBody(w http.ResponseWriter, r *http.Request) {
	var metric shared.Metric

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var (
		newValue float64
		newDelta int64
	)

	switch metric.MType {
	case shared.Gauge:
		if metric.Value == nil {
			http.Error(w, "Invalid metric value for type Gauge", http.StatusBadRequest)
			return
		}
		newValue, err = h.Repo.UpdateGauge(metric.ID, *metric.Value)
		metric.Value = &newValue
	case shared.Counter:
		if metric.Delta == nil {
			http.Error(w, "Invalid metric delta for type Counter", http.StatusBadRequest)
			return
		}
		newDelta, err = h.Repo.UpdateCounter(metric.ID, *metric.Delta)
		metric.Delta = &newDelta
	default:
		// Must be 400, return 501 because of autotests.
		http.Error(w, "Unknown metric type", http.StatusNotImplemented)
		return
	}

	if err != nil {
		log.Errorf("Error updating metric: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(metric)
	if err != nil {
		log.Errorf("Error encoding metric: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
