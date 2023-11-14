package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gonozov0/go-musthave-devops/internal/server/repository"

	"github.com/gonozov0/go-musthave-devops/internal/shared"

	"github.com/go-chi/chi/v5"
)

// getGaugeValue returns the string value of the gauge metric.
func (h *Handler) getGaugeValue(metricName string) (string, error) {
	value, err := h.repo.GetGauge(metricName)
	if err != nil {
		return "", err
	}
	valueStr := strconv.FormatFloat(value, 'f', -1, 64)
	return valueStr, nil
}

// getCounterValue returns the string value of the counter metric.
func (h *Handler) getCounterValue(metricName string) (string, error) {
	value, err := h.repo.GetCounter(metricName)
	if err != nil {
		return "", err
	}
	valueStr := strconv.FormatInt(value, 10)
	return valueStr, nil
}

// GetMetricByURL is the HTTP handler for getting metrics.
func (h *Handler) GetMetricByURL(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	var err error
	var value string

	switch metricType {
	case shared.Gauge:
		value, err = h.getGaugeValue(metricName)
	case shared.Counter:
		value, err = h.getCounterValue(metricName)
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

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
