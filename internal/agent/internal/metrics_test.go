package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectMetrics(t *testing.T) {
	metrics := CollectMetrics()

	assert.NotEmpty(t, metrics)

	for _, metric := range metrics {
		assert.NotEmpty(t, metric.Name)
		assert.Equal(t, "gauge", metric.Type)
		assert.NotNil(t, metric.Value)
	}
}

func TestSendMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	metrics := []Metric{
		{Name: "TestMetric1", Type: "gauge", Value: 42.0},
		{Name: "TestMetric2", Type: "gauge", Value: 43.0},
	}

	newMetrics, err := SendMetrics(metrics, server.URL)

	assert.NoError(t, err)
	assert.Empty(t, newMetrics)
}
