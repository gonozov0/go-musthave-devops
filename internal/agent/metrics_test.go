package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gonozov0/go-musthave-devops/internal/shared"
)

func TestCollectMetrics(t *testing.T) {
	metrics := CollectMetrics()

	assert.NotEmpty(t, metrics)

	for _, metric := range metrics {
		assert.NotEmpty(t, metric.ID)
		if metric.ID == "PollCount" {
			assert.Equal(t, "counter", metric.MType)
			assert.NotNil(t, metric.Delta)
			continue
		}
		assert.Equal(t, "gauge", metric.MType)
		assert.NotNil(t, metric.Value)
	}
}

func TestSendMetrics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	testFloat := 42.0
	testInt := int64(42)

	metrics := []shared.Metric{
		{ID: "TestMetric1", MType: shared.Gauge, Value: &testFloat},
		{ID: "TestMetric2", MType: shared.Counter, Delta: &testInt},
	}

	newMetrics, err := SendMetrics(metrics, server.URL)

	assert.NoError(t, err)
	assert.Empty(t, newMetrics)

}
