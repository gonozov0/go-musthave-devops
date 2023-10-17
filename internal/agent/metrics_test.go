package agent

import (
	"net/http/httptest"
	"testing"

	"github.com/gonozov0/go-musthave-devops/internal/server/application"
	repository "github.com/gonozov0/go-musthave-devops/internal/server/repository/in_memory"

	"github.com/stretchr/testify/require"

	"github.com/gonozov0/go-musthave-devops/internal/shared"
)

func TestCollectMetrics(t *testing.T) {
	metrics := CollectMetrics()

	require.NotEmpty(t, metrics)

	for _, metric := range metrics {
		require.NotEmpty(t, metric.ID)
		if metric.ID == "PollCount" {
			require.Equal(t, "counter", metric.MType)
			require.NotNil(t, metric.Delta)
			continue
		}
		require.Equal(t, "gauge", metric.MType)
		require.NotNil(t, metric.Value)
	}
}

func TestSendMetrics(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	router := application.NewRouter(repo)
	server := httptest.NewServer(router)
	defer server.Close()

	testFloat := 42.0
	testInt := int64(42)

	metrics := []shared.Metric{
		{ID: "TestMetric1", MType: shared.Gauge, Value: &testFloat},
		{ID: "TestMetric2", MType: shared.Counter, Delta: &testInt},
	}

	newMetrics, err := SendMetrics(metrics, server.URL)

	require.NoError(t, err)
	require.Empty(t, newMetrics)
}
