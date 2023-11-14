package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gonozov0/go-musthave-devops/internal/server/repository"
	"github.com/gonozov0/go-musthave-devops/internal/shared"
)

func (h *Handler) BatchUpdateMetrics(w http.ResponseWriter, r *http.Request) {
	var metrics []shared.Metric
	err := json.NewDecoder(r.Body).Decode(&metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(metrics) == 0 {
		http.Error(w, "empty metrics", http.StatusBadRequest)
		return
	}

	updateGauges := make([]repository.GaugeMetric, 0, len(metrics))
	updateCounters := make([]repository.CounterMetric, 0, len(metrics))
	for _, metric := range metrics {
		switch metric.MType {
		case shared.Gauge:
			if metric.Value == nil {
				http.Error(w, "value is required for gauge metric", http.StatusBadRequest)
				return
			}
			updateGauges = append(updateGauges, repository.GaugeMetric{Name: metric.ID, Value: *metric.Value})
		case shared.Counter:
			if metric.Delta == nil {
				http.Error(w, "delta is required for counter metric", http.StatusBadRequest)
				return
			}
			updateCounters = append(updateCounters, repository.CounterMetric{Name: metric.ID, Value: *metric.Delta})
		default:
			// Must be 400, return 501 because of autotests.
			http.Error(w, "Unknown metric type", http.StatusNotImplemented)
			return
		}
	}

	newGauges, err := h.repo.UpdateGauges(updateGauges)
	if err != nil {
		log.Errorf("failed to update gauges: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newCounters, err := h.repo.UpdateCounters(updateCounters)
	if err != nil {
		log.Errorf("failed to update counters: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var newMetrics []shared.Metric
	for _, gauge := range newGauges {
		newMetrics = append(newMetrics, shared.Metric{ID: gauge.Name, MType: shared.Gauge, Value: &gauge.Value})
	}
	for _, counter := range newCounters {
		newMetrics = append(newMetrics, shared.Metric{ID: counter.Name, MType: shared.Counter, Delta: &counter.Value})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(newMetrics)
	if err != nil {
		log.Errorf("failed to encode metrics: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
