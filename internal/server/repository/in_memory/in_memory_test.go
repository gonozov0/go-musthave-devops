package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	filestorage "github.com/gonozov0/go-musthave-devops/internal/server/repository/in_memory/internal/file_storage"
)

func TestLoadMetrics(t *testing.T) {
	fileName := "test_load_metrics.json"
	defer os.Remove(fileName)

	metrics := filestorage.Metrics{
		Gauges: []filestorage.GaugeMetric{
			{
				Name: "TestMetric1",
			},
		},
		Counters: []filestorage.CounterMetric{
			{
				Name: "TestMetric2",
			},
		},
	}
	bytes_, err := json.Marshal(metrics)
	require.NoError(t, err)

	err = os.WriteFile(fileName, bytes_, 0644)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	repo, err := NewInMemoryRepositoryWithFileStorage(ctx, wg, fileName, 0, true)
	require.NoError(t, err)

	// for correct removal of test file
	cancel()
	wg.Wait()

	gaugeMetric := metrics.Gauges[0]
	gaugeValue, err := repo.GetGauge(gaugeMetric.Name)
	require.NoError(t, err)
	require.Equal(t, gaugeMetric.Value, gaugeValue)

	counterMetric := metrics.Counters[0]
	counterValue, err := repo.GetCounter(counterMetric.Name)
	require.NoError(t, err)
	require.Equal(t, counterMetric.Value, counterValue)
}

func TestSaveMetrics(t *testing.T) {
	fileName := "test_trigger_save.json"
	defer os.Remove(fileName)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	repo, err := NewInMemoryRepositoryWithFileStorage(ctx, wg, fileName, 0, false)
	require.NoError(t, err)

	metricValue := int64(42)
	metricID := "TestMetric"
	_, err = repo.UpdateCounter(metricID, metricValue)
	require.NoError(t, err)

	time.Sleep(1 * time.Microsecond)
	cancel()
	wg.Wait()

	bytes_, err := os.ReadFile(fileName)
	require.NoError(t, err)
	require.Equal(
		t,
		fmt.Sprintf(
			"{\"Gauges\":[],\"Counters\":[{\"name\":\"%s\",\"value\":%d}]}",
			metricID,
			metricValue,
		),
		string(bytes_),
		"File content is not equal to expected",
	)
}
