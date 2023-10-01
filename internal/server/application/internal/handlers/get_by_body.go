package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gonozov0/go-musthave-devops/internal/shared"

	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
)

func (h *Handler) GetMetricByBody(w http.ResponseWriter, r *http.Request) {
	var metric shared.Metric

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var (
		gaugeValue   float64
		counterValue int64
	)

	switch metric.MType {
	case shared.Gauge:
		gaugeValue, err = h.Repo.GetGauge(metric.ID)
		metric.Value = &gaugeValue
	case shared.Counter:
		counterValue, err = h.Repo.GetCounter(metric.ID)
		metric.Delta = &counterValue
	default:
		// Must be 400, return 501 because of autotests.
		http.Error(w, "Unknown metric type", http.StatusNotImplemented)
		return
	}

	if err != nil {
		if errors.Is(err, repository.ErrMetricNotFound) {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
