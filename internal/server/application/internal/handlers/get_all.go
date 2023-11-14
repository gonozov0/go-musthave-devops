package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) GetAllMetrics(w http.ResponseWriter, r *http.Request) {
	gaugeMetrics, err := h.repo.GetAllGauges()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	counterMetrics, err := h.repo.GetAllCounters()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, "<html><body><h1>Metrics</h1>")
	fmt.Fprint(w, "<h2>Gauges</h2><ul>")
	for _, metric := range gaugeMetrics {
		fmt.Fprintf(w, "<li>%s: %v</li>", metric.Name, metric.Value)
	}
	fmt.Fprint(w, "</ul>")

	fmt.Fprint(w, "<h2>Counters</h2><ul>")
	for _, metric := range counterMetrics {
		fmt.Fprintf(w, "<li>%s: %v</li>", metric.Name, metric.Value)
	}
	fmt.Fprint(w, "</ul></body></html>")
}
